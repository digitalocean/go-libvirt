// Copyright 2017 The go-libvirt Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package lvgen contains the instructions for regenerating the libvirt
// bindings. We do this by parsing the remote_protocol.x file included in the
// libvirt sources. Bindings will be generated if you run `go generate` in this
// directory.
package lvgen

// Before running `go generate`:
// 1) Make sure goyacc is installed from golang.org/x/tools (you can use this
//    command: `go get golang.org/x/tools/...`)
// 2) Set the environment variable LIBVIRT_SOURCE to point to the top level
//    directory containing the version of libvirt for which you want to generate
//    bindings.

//go:generate goyacc sunrpc.y
//go:generate go run gen/main.go
