# dhcpdfs
DHCP as a file system (requires FUSE)

This was cobbled together over a couple of afternoons to teach myself FUSE. While I do run
it on a home server, lacking features I'd expect a production-ready fork would include:
* caching the dhcpd lease file (or better yet, the objects created from it)
* support reading from multiple lease files
* better handle duplicated entries on the lease file (eg update rather than overwrite)
* and I'd make the command line usage more user friendly

I may well add these features myself - if and when I get time, or there becomes demand for it.

### Installation guide (assumes existing Go environment):

    GOFUSE=github.com/hanwen/go-fuse
    go get $GOFUSE
    cd $GOPATH/src/$GOFUSE
    bash all.bash
    
    go get github.com/lmorg/dhcpdfs
    go install github.com/lmorg/dhcpdfs

(dhcpdfs was written in Go but uses cgo, which may affect your ability to cross-compile)
    
### Example usage:

    mkdir ~/dhcpdfs
    $GOPATH/bin/dhcpdfs -m ~/dhcpdfs -l /var/db/dhcpd/dhcpd.leases
