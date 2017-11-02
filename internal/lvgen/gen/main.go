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

package main

import (
	"fmt"
	"os"
	"path"

	"github.com/digitalocean/go-libvirt/internal/lvgen"
)

// TODO: make these an argument
const lvPath = "../../../libvirt"
const protoPath = "src/remote/remote_protocol.x"

func main() {
	fmt.Println("Generating golang bindings for libvirt")
	lvFile := path.Join(lvPath, protoPath)
	rdr, err := os.Open(lvFile)
	if err != nil {
		fmt.Printf("failed to open protocol file at %v: %v\n", lvFile, err)
		os.Exit(1)
	}
	defer rdr.Close()

	if err = lvgen.Generate(rdr); err != nil {
		fmt.Println("go-libvirt code generator failed:", err)
		os.Exit(1)
	}
}
