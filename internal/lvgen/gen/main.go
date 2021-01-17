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
	"path/filepath"
	"strings"

	"github.com/digitalocean/go-libvirt/internal/lvgen"
)

var protoPaths = [...]string{
	"src/remote/remote_protocol.x",
	"src/remote/qemu_protocol.x",
}

func main() {
	lvPath := os.Getenv("LIBVIRT_SOURCE")
	if lvPath == "" {
		fmt.Println("set $LIBVIRT_SOURCE to point to the root of the libvirt sources and retry")
		os.Exit(1)
	}
	fmt.Println("protocol file processing")
	for _, p := range protoPaths {
		protoPath := filepath.Join(lvPath, p)
		fmt.Println("  processing", p)
		err := processProto(protoPath)
		if err != nil {
			fmt.Println("go-libvirt code generator failed:", err)
			os.Exit(1)
		}
	}
}

func processProto(lvFile string) error {
	rdr, err := os.Open(lvFile)
	if err != nil {
		fmt.Printf("failed to open protocol file at %v: %v\n", lvFile, err)
		os.Exit(1)
	}
	defer rdr.Close()

	// extract the base filename, without extension, for the generator to use.
	name := strings.TrimSuffix(filepath.Base(lvFile), filepath.Ext(lvFile))

	return lvgen.Generate(name, rdr)
}
