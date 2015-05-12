// A Go mirror of libfuse's hello.c

package main

import (
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
	"log"
	"strings"
	"time"
)

type fs struct {
	pathfs.FileSystem
}

func (me *fs) GetAttr(path string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	attr := new(fuse.Attr)
	now := uint64(time.Now().Unix())
	attr.Atime = now
	attr.Ctime = now
	attr.Mtime = now

	isDir := func() {
		attr.Mode = fuse.S_IFDIR | 0555
		attr.Size = 0
	}

	isFile := func(size uint64) {
		attr.Mode = fuse.S_IFREG | 0444
		attr.Size = size
	}

	// root mount point
	if path == "" {
		isDir()
		return attr, fuse.OK
	}

	if !strings.Contains(path, "/") {

		// tier 1 directory?
		if DirHierarchy[path].Leases != nil {
			isDir()
			return attr, fuse.OK

		}

		// tier 1 file?
		if FileHierarchy[path].Size != 0 {
			isFile(FileHierarchy[path].Size)
			return attr, fuse.OK

		}
		return nil, fuse.ENOENT
	}

	split := strings.Split(path, "/")

	// is tier 2
	if len(split) == 2 &&
		DirHierarchy[split[0]].Leases != nil &&
		DirHierarchy[split[0]].Rx.MatchString(split[1]) {

		isDir()
		return attr, fuse.OK
	}

	// is tier 3
	if len(split) == 3 &&
		DirHierarchy[split[0]].Leases != nil &&
		DirHierarchy[split[0]].Rx.MatchString(split[1]) {

		leases, err := DirHierarchy[split[0]].Leases()
		if err != nil {
			return nil, fuse.ENOENT
		}

		// file not found (we always expect an IP)
		if leases[split[1]].IP == "" {
			return nil, fuse.ENOENT
		}

		// file not found
		file := leases[split[1]].Property(split[2])
		if file == nil {
			return nil, fuse.ENOENT
		}

		isFile(file.Size)
		return attr, fuse.OK
	}

	// file not found
	return nil, fuse.ENOENT
}

func (me *fs) OpenDir(path string, context *fuse.Context) (dir []fuse.DirEntry, code fuse.Status) {
	dir, err := Directory(path)
	if err != nil {
		log.Println("Error opening directory:", err)
		return nil, fuse.ENOENT
	}

	return dir, fuse.OK
}

func (me *fs) Open(path string, flags uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	if flags&fuse.O_ANYWRITE != 0 {
		return nil, fuse.EPERM
	}

	// root files
	if FileHierarchy[path].Size != 0 {
		return nodefs.NewDataFile(FileHierarchy[path].Content), fuse.OK
	}

	// the only other files are in tier 3,
	// so exit if anything other than 3rd tier
	split := strings.Split(path, "/")
	if len(split) != 3 {
		return nil, fuse.ENOENT
	}

	if DirHierarchy[split[0]].Leases != nil &&
		DirHierarchy[split[0]].Rx.MatchString(split[1]) {

		leases, err := DirHierarchy[split[0]].Leases()
		if err != nil {
			return nil, fuse.ENOENT
		}

		if leases[split[1]].IP != "" &&
			LeaseHierarchy[split[2]].Size != 0 {

			f := leases[split[1]].Property(split[2])
			if f == nil {
				return nil, fuse.ENOENT
			}

			return nodefs.NewDataFile(f.Content), fuse.OK

		}
	}

	return nil, fuse.ENOENT
}
