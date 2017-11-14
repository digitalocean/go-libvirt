package lvgen

// This file contains the instructions for regenerating the libvirt bindings.
// We do this by parsing the remote_protocol.x file included in the libvirt
// sources. Bindings will be generated if you run `go generate` in this
// directory.

// Before running `go generate`:
// 1) Make sure goyacc is installed from golang.org/x/tools (you can use this
//    command: `go get golang.org/x/tools/...`)
// 2) Set the environment variable LIBVIRT_SOURCE to point to the top level
//    directory containing the version of libvirt for which you want to generate
//    bindings.

//go:generate goyacc sunrpc.y
//go:generate go run gen/main.go
