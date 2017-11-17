libvirt [![GoDoc](http://godoc.org/github.com/digitalocean/go-libvirt?status.svg)](http://godoc.org/github.com/digitalocean/go-libvirt) [![Build Status](https://travis-ci.org/digitalocean/go-libvirt.svg?branch=master)](https://travis-ci.org/digitalocean/go-libvirt) [![Report Card](https://goreportcard.com/badge/github.com/digitalocean/go-libvirt)](https://goreportcard.com/report/github.com/digitalocean/go-libvirt)
====

Package `libvirt` provides a pure Go interface for interacting with Libvirt.

Rather than using Libvirt's C bindings, this package makes use of
Libvirt's RPC interface, as documented [here](https://libvirt.org/internals/rpc.html).
Connections to the libvirt server may be local, or remote. RPC packets are encoded
using the XDR standard as defined by [RFC 4506](https://tools.ietf.org/html/rfc4506.html).

Libvirt's RPC interface is quite extensive, and changes from one version to the next, so
this project uses a code generator to build the go bindings. The code generator should
be run whenever you want to build go-libvirt for a new version of libvirt. To do this,
you'll need to set an environment variable `LIBVIRT_SOURCE` to the directory containing
the untarred libvirt sources, and then run `go generate ./...` from the go-libvirt directory.
The code generator consumes [src/remote/remote_protocol.x](https://github.com/libvirt/libvirt/blob/master/src/remote/remote_protocol.x)
and produces go bindings for all the remote procedures defined there.

[Pull requests are welcome](https://github.com/digitalocean/go-libvirt/blob/master/CONTRIBUTING.md)!

Feel free to join us in [`#go-qemu` on freenode](https://webchat.freenode.net/)
if you'd like to discuss the project.

Warning
-------

The libvirt project strongly recommends *against* talking to the RPC interface
directly. They consider it to be a private implementation detail with the
possibility of being entirely rearchitected in the future.

While these package are reasonably well-tested and have seen some use inside of
DigitalOcean, there may be subtle bugs which could cause the packages to act
in unexpected ways.  Use at your own risk!

In addition, the API is not considered stable at this time.  If you would like
to include package `libvirt` in a project, we highly recommend vendoring it into
your project.

Example
-------

```go
package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/digitalocean/go-libvirt"
)

func main() {
	//c, err := net.DialTimeout("tcp", "127.0.0.1:16509", 2*time.Second)
	//c, err := net.DialTimeout("tcp", "192.168.1.12:16509", 2*time.Second)
	c, err := net.DialTimeout("unix", "/var/run/libvirt/libvirt-sock", 2*time.Second)
	if err != nil {
		log.Fatalf("failed to dial libvirt: %v", err)
	}

	l := libvirt.New(c)
	if err := l.Connect(); err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	v, err := l.Version()
	if err != nil {
		log.Fatalf("failed to retrieve libvirt version: %v", err)
	}
	fmt.Println("Version:", v)

	domains, err := l.Domains()
	if err != nil {
		log.Fatalf("failed to retrieve domains: %v", err)
	}

	fmt.Println("ID\tName\t\tUUID")
	fmt.Printf("--------------------------------------------------------\n")
	for _, d := range domains {
		fmt.Printf("%d\t%s\t%x\n", d.ID, d.Name, d.UUID)
	}

	if err := l.Disconnect(); err != nil {
		log.Fatal("failed to disconnect: %v", err)
	}
}

```

```
Version: 1.3.4
ID	Name		UUID
--------------------------------------------------------
1	Test-1		dc329f87d4de47198cfd2e21c6105b01
2	Test-2		dc229f87d4de47198cfd2e21c6105b01
```
