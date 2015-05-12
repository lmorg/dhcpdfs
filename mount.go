package main

import (
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
	"log"
)

func Mount() error {
	vfs := pathfs.NewPathNodeFs(
		&fs{FileSystem: pathfs.NewDefaultFileSystem()}, nil,
	)
	server, _, err := nodefs.MountRoot(f_mount_point, vfs.Root(), nil)
	if err != nil {
		log.Fatalf("Mount failed (%s): %v\n", f_mount_point, err)
	}
	server.Serve()

	return nil
}
