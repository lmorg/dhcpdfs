package main

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

func ReadFile(filename string) (content *string) {
	content = new(string)
	var (
		reader *bufio.Reader
		err    error
	)

	fi, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		return
	}
	defer fi.Close()

	reader = bufio.NewReader(fi)

	for {
		b, _, err := reader.ReadLine()
		if err != nil {
			if err.Error() != "EOF" {
				log.Println(err)
			}
			return
		}

		b = bytes.TrimSpace(b)
		if len(b) == 0 {
			continue
		}
		if b[0] != '#' {
			*content += string(b)
		}
	}

	return
}

func ParseFile(content *string, only_get_active bool) (leases *[]Leases) {
	leases = new([]Leases)
	rx, _ := regexp.Compile(`(?ms)lease ([.:0-9a-f]+).*?\{(.*?)\}`)
	matched := rx.FindAllStringSubmatch(*content, -1)

	for _, block := range matched {
		items := strings.Split(block[2], ";")
		lease := new(Leases)
		lease.IP = block[1]

		for _, val := range items {
			ParseLease(strings.TrimSpace(val), lease)
		}

		if !only_get_active || lease.Active {
			*leases = append(*leases, *lease)
		}
	}

	return
}

func ParseLease(value string, lease *Leases) {
	var err error
	rx, _ := regexp.Compile(`"(.*?)"`)
	split := strings.Split(value, " ")
	switch split[0] {
	default:
		return
	case "starts":
		lease.Starts, err = time.Parse(DATE_TIME_FORMAT, split[2]+" "+split[3])
	case "ends":
		lease.Ends, err = time.Parse(DATE_TIME_FORMAT, split[2]+" "+split[3])
	case "cltt":
		lease.CLTT, err = time.Parse(DATE_TIME_FORMAT, split[2]+" "+split[3])
	case "tstp":
		lease.TSTP, err = time.Parse(DATE_TIME_FORMAT, split[2]+" "+split[3])
	case "binding":
		if split[2] == "active" {
			lease.Active = true
		}
	case "hardware":
		lease.MAC = split[2]
	case "uid":
		lease.UID = rx.FindStringSubmatch(value)[1]
	case "client-hostname":
		lease.Hostname = rx.FindStringSubmatch(value)[1]
	case "set":
		if split[1] == "vendor-class-identifier" {
			lease.VendorClass = rx.FindStringSubmatch(value)[1]
		}
	}

	if err != nil {
		log.Println(err)
	}
}
