// Copyright 2016 The go-libvirt Authors.
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

package libvirt

import (
	"bytes"

	"github.com/davecgh/go-xdr/xdr2"
	"github.com/vtolstov/go-libvirt/internal/constants"
)

// Secret represents a secret managed by the libvirt daemon.
type Secret struct {
	UUID      [constants.UUIDSize]byte
	UsageType SecretUsageType
	UsageID   string
}

// SecretUsageType specifies the usage for a libvirt secret.
type SecretUsageType uint32

const (
	// SecretUsageTypeNone specifies no usage.
	SecretUsageTypeNone SecretUsageType = iota
	// SecretUsageTypeVolume specifies a volume secret.
	SecretUsageTypeVolume
	// SecretUsageTypeCeph specifies secrets for ceph devices.
	SecretUsageTypeCeph
	// SecretUsageTypeISCSI specifies secrets for ISCSI devices.
	SecretUsageTypeISCSI
)

// Secrets returns all secrets managed by the libvirt daemon.
func (l *Libvirt) Secrets() ([]Secret, error) {
	req := struct {
		NeedResults uint32
		Flags       uint32
	}{
		NeedResults: 1,
		Flags:       0, // unused per libvirt source, callers should pass 0
	}

	buf, err := encode(&req)
	if err != nil {
		return nil, err
	}

	resp, err := l.request(constants.ProcConnectListAllSecrets, constants.ProgramRemote, &buf)
	if err != nil {
		return nil, err
	}

	r := <-resp
	if r.Status != StatusOK {
		return nil, decodeError(r.Payload)
	}

	result := struct {
		Secrets []Secret
		Count   uint32
	}{}

	dec := xdr.NewDecoder(bytes.NewReader(r.Payload))
	_, err = dec.Decode(&result)
	if err != nil {
		return nil, err
	}

	return result.Secrets, nil
}
