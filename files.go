package main

type File struct {
	Size    uint64
	Content []byte
}

var FileHierarchy map[string]File
var LeaseHierarchy map[string]File

func init() {
	FileHierarchy = make(map[string]File)
	LeaseHierarchy = make(map[string]File)

	FileHierarchy = map[string]File{
		"VERSION": File{
			Size:    uint64(len(VERSION)),
			Content: []byte(VERSION),
		},
		"COPYRIGHT": File{
			Size:    uint64(len(COPYRIGHT)),
			Content: []byte(COPYRIGHT),
		},
		"README": File{
			Size:    uint64(len(README)),
			Content: []byte(README),
		},
	}

	// initialise the map.
	// a bit of a hack, but since lease hierarchy is
	// only used to determine which file names are
	// valid, any non-zero size will be considered
	// valid.
	LeaseHierarchy = map[string]File{
		"ip":           File{Size: 1},
		"mac":          File{Size: 1},
		"active":       File{Size: 1},
		"uid":          File{Size: 1},
		"hostname":     File{Size: 1},
		"vendor-class": File{Size: 1},
	}
}

func (l Leases) Property(p string) (f *File) {
	bytes := func(s string) []byte {
		return []byte(s + "\n")
	}

	f = new(File)
	switch p {
	default:
		return nil
	case "ip":
		f.Content = bytes(l.IP)
	case "mac":
		f.Content = bytes(l.MAC)
	case "active":
		if l.Active {
			f.Content = bytes("1")
		} else {
			f.Content = bytes("0")
		}
	case "uid":
		f.Content = bytes(l.UID)
	case "hostname":
		f.Content = bytes(l.Hostname)
	case "vendor-class":
		f.Content = bytes(l.VendorClass)
	}

	f.Size = uint64(len(f.Content))
	return
}
