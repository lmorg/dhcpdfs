package main

import (
	"flag"
	"github.com/hanwen/go-fuse/fuse"
	"os"
	"time"
)

type Leases struct {
	IP          string
	MAC         string
	Starts      time.Time
	Ends        time.Time
	TSTP        time.Time
	CLTT        time.Time
	Active      bool
	UID         string
	Hostname    string
	VendorClass string
	DirEntry    fuse.DirEntry
}

const DATE_TIME_FORMAT = "2006/01/02 15:04:05"

const (
	COPYRIGHT = "© 2015 Laurence Morgan\n"
	VERSION   = "0.1\n"
)

// these will be pushed out to flags eventually,
// hence the weird naming convention.
var (
	f_mount_point string
	f_lease_file  string
)

func main() {
	Flags()
	Mount()
}

func Flags() {
	flag.StringVar(&f_mount_point, "m", "", "mount point")
	flag.StringVar(&f_lease_file, "l", "", "lease file")

	flag.Parse()

	if f_mount_point == "" || f_lease_file == "" {
		flag.Usage()
		os.Exit(1)
	}
}
