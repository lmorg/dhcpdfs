package main

import (
	"errors"
	"github.com/hanwen/go-fuse/fuse"
	"regexp"
	"strings"
)

type Dir struct {
	Rx     *regexp.Regexp
	Leases func() (map[string]Leases, error)
}

var DirHierarchy map[string]Dir

func init() {
	var rx_compile = func(s string) (rx *regexp.Regexp) {
		rx, _ = regexp.Compile(s)
		return
	}

	DirHierarchy = make(map[string]Dir)

	DirHierarchy = map[string]Dir{
		"active-ip": Dir{
			Rx:     rx_compile(`^([0-9]+\.){3}[0-9]+$`),
			Leases: PwdActiveIPs,
		},
		"active-mac": Dir{
			Rx:     rx_compile(`^([0-9a-f]{2}:){5}[0-9a-f]{2}$`),
			Leases: PwdActiveMACs,
		},
		"all-mac": Dir{
			Rx:     rx_compile(`^([0-9a-f]{2}:){5}[0-9a-f]{2}$`),
			Leases: Empty, // not yet implimented
		},
	}

}

func Empty() (dir map[string]Leases, err error) {
	return
}

func Directory(path string) (dir []fuse.DirEntry, err error) {
	const NOT_FOUND = "File or path not found"

	if path == "" {

		// list root directories
		for name, _ := range DirHierarchy {
			dir = append(dir, fuse.DirEntry{
				Name: name,
				Mode: fuse.S_IFDIR,
			})

		}

		// list root files
		for name, _ := range FileHierarchy {
			dir = append(dir, fuse.DirEntry{
				Name: name,
				Mode: fuse.S_IFREG,
			})
		}

		return
	}

	// if 0 tiers or more than 2, then we can safely
	// assume that something is amiss.
	split := strings.Split(path, "/")
	if len(split) == 0 || len(split) > 2 {
		err = errors.New(NOT_FOUND)
		return
	}

	// tier 1 not found. Don't both decending
	if DirHierarchy[split[0]].Leases == nil {
		err = errors.New(NOT_FOUND)
		return
	}

	// establish deeper nesting
	leases, err := DirHierarchy[split[0]].Leases()
	if err != nil {
		return
	}

	// list tier 2 directories
	if len(split) == 1 {
		for name, _ := range leases {
			dir = append(dir, leases[name].DirEntry)
		}
		return
	}

	// tier 2 not found
	if leases[split[1]].IP == "" {
		err = errors.New(NOT_FOUND)
		return
	}

	// list tier 3 files
	for name, _ := range LeaseHierarchy {
		dir = append(dir, fuse.DirEntry{
			Name: name,
			Mode: fuse.S_IFREG,
		})
	}

	return
}

func PwdActiveIPs() (lmap map[string]Leases, err error) {
	lmap = make(map[string]Leases)
	content := ReadFile(f_lease_file)
	leases := ParseFile(content, true)

	for i := 0; i < len(*leases); i++ {
		(*leases)[i].DirEntry = fuse.DirEntry{
			Name: (*leases)[i].IP,
			Mode: fuse.S_IFDIR,
		}
		lmap[(*leases)[i].IP] = (*leases)[i]
	}

	return
}

func PwdActiveMACs() (lmap map[string]Leases, err error) {
	lmap = make(map[string]Leases)
	content := ReadFile(f_lease_file)
	leases := ParseFile(content, true)

	for i := 0; i < len(*leases); i++ {
		(*leases)[i].DirEntry = fuse.DirEntry{
			Name: (*leases)[i].MAC,
			Mode: fuse.S_IFDIR,
		}
		lmap[(*leases)[i].MAC] = (*leases)[i]
	}

	return
}
