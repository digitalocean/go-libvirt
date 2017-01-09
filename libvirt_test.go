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
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/digitalocean/go-libvirt/internal/constants"
	"github.com/digitalocean/go-libvirt/libvirttest"
)

func TestConnect(t *testing.T) {
	conn := libvirttest.New()
	l := New(conn)

	err := l.Connect()
	if err != nil {
		t.Error(err)
	}
}

func TestDisconnect(t *testing.T) {
	conn := libvirttest.New()
	l := New(conn)

	err := l.Disconnect()
	if err != nil {
		t.Error(err)
	}
}

func TestMigrate(t *testing.T) {
	conn := libvirttest.New()
	l := New(conn)

	var flags MigrateFlags
	flags = MigrateFlagLive |
		MigrateFlagPeerToPeer |
		MigrateFlagPersistDestination |
		MigrateFlagChangeProtection |
		MigrateFlagAbortOnError |
		MigrateFlagAutoConverge |
		MigrateFlagNonSharedDisk

	if err := l.Migrate("test", "qemu+tcp://foo/system", flags); err != nil {
		t.Fatalf("unexpected live migration error: %v", err)
	}
}

func TestMigrateInvalidDest(t *testing.T) {
	conn := libvirttest.New()
	l := New(conn)

	var flags MigrateFlags
	flags = MigrateFlagLive |
		MigrateFlagPeerToPeer |
		MigrateFlagPersistDestination |
		MigrateFlagChangeProtection |
		MigrateFlagAbortOnError |
		MigrateFlagAutoConverge |
		MigrateFlagNonSharedDisk

	dest := ":$'"
	if err := l.Migrate("test", dest, flags); err == nil {
		t.Fatalf("expected invalid dest uri %q to fail", dest)
	}
}

func TestMigrateSetMaxSpeed(t *testing.T) {
	conn := libvirttest.New()
	l := New(conn)

	if err := l.MigrateSetMaxSpeed("test", 100); err != nil {
		t.Fatalf("unexpected error setting max speed for migrate: %v", err)
	}
}

func TestDomains(t *testing.T) {
	conn := libvirttest.New()
	l := New(conn)

	domains, err := l.Domains()
	if err != nil {
		t.Error(err)
	}

	wantLen := 2
	gotLen := len(domains)
	if gotLen != wantLen {
		t.Errorf("expected %d domains to be returned, got %d", wantLen, gotLen)
	}

	for i, d := range domains {
		wantID := i + 1
		if d.ID != wantID {
			t.Errorf("expected domain ID %q, got %q", wantID, d.ID)
		}

		wantName := fmt.Sprintf("aaaaaaa-%d", i+1)
		if d.Name != wantName {
			t.Errorf("expected domain name %q, got %q", wantName, d.Name)
		}
	}
}

func TestDomainState(t *testing.T) {
	conn := libvirttest.New()
	l := New(conn)

	wantState := DomainState(DomainStateRunning)
	gotState, err := l.DomainState("test")
	if err != nil {
		t.Error(err)
	}

	if gotState != wantState {
		t.Errorf("expected domain state %d, got %d", wantState, gotState)
	}
}

func TestEvents(t *testing.T) {
	conn := libvirttest.New()
	l := New(conn)
	done := make(chan struct{})

	stream, err := l.Events("test")
	if err != nil {
		t.Error(err)
	}

	go func() {
		var e DomainEvent
		select {
		case e = <-stream:
		case <-time.After(time.Second * 5):
			t.Error("expected event, received timeout")
		}

		result := struct {
			Device string `json:"device"`
			Len    int    `json:"len"`
			Offset int    `json:"offset"`
			Speed  int    `json:"speed"`
			Type   string `json:"type"`
		}{}

		if err := json.Unmarshal(e.Details, &result); err != nil {
			t.Error(err)
		}

		expected := "drive-ide0-0-0"
		if result.Device != expected {
			t.Errorf("expected device %q, got %q", expected, result.Device)
		}

		done <- struct{}{}
	}()

	// send an event to the listener goroutine
	conn.Test.Write(append(testEventHeader, testEvent...))

	// wait for completion
	<-done
}

func TestRun(t *testing.T) {
	conn := libvirttest.New()
	l := New(conn)

	res, err := l.Run("test", []byte(`{"query-version"}`))
	if err != nil {
		t.Error(err)
	}

	type version struct {
		Return struct {
			Package string `json:"package"`
			QEMU    struct {
				Major int `json:"major"`
				Micro int `json:"micro"`
				Minor int `json:"minor"`
			} `json:"qemu"`
		} `json:"return"`
	}

	var v version
	err = json.Unmarshal(res, &v)
	if err != nil {
		t.Error(err)
	}

	expected := 2
	if v.Return.QEMU.Major != expected {
		t.Errorf("expected qemu major version %d, got %d", expected, v.Return.QEMU.Major)
	}
}

func TestRunFail(t *testing.T) {
	conn := libvirttest.New()
	conn.Fail = true
	l := New(conn)

	_, err := l.Run("test", []byte(`{"drive-foo"}`))
	if err == nil {
		t.Error("expected qemu error")
	}
}

func TestStoragePool(t *testing.T) {
	conn := libvirttest.New()
	l := New(conn)

	wantName := "default"
	pool, err := l.StoragePool(wantName)
	if err != nil {
		t.Error(err)
	}

	gotName := pool.Name
	if gotName != wantName {
		t.Errorf("expected name %q, got %q", wantName, gotName)
	}

	// bb30a11c-0846-4827-8bba-3e6b5cf1b65f
	wantUUID := [constants.UUIDSize]byte{
		0xbb, 0x30, 0xa1, 0x1c, 0x08, 0x46, 0x48, 0x27,
		0x8b, 0xba, 0x3e, 0x6b, 0x5c, 0xf1, 0xb6, 0x5f,
	}
	gotUUID := pool.UUID
	if gotUUID != wantUUID {
		t.Errorf("expected UUID %q, got %q", wantUUID, gotUUID)
	}
}

func TestStoragePools(t *testing.T) {
	conn := libvirttest.New()
	l := New(conn)

	pools, err := l.StoragePools(StoragePoolsFlagActive)
	if err != nil {
		t.Error(err)
	}

	wantLen := 1
	gotLen := len(pools)
	if gotLen != wantLen {
		t.Errorf("expected %d storage pool, got %d", wantLen, gotLen)
	}

	wantName := "default"
	gotName := pools[0].Name
	if gotName != wantName {
		t.Errorf("expected name %q, got %q", wantName, gotName)
	}

	// bb30a11c-0846-4827-8bba-3e6b5cf1b65f
	wantUUID := [constants.UUIDSize]byte{
		0xbb, 0x30, 0xa1, 0x1c, 0x08, 0x46, 0x48, 0x27,
		0x8b, 0xba, 0x3e, 0x6b, 0x5c, 0xf1, 0xb6, 0x5f,
	}
	gotUUID := pools[0].UUID
	if gotUUID != wantUUID {
		t.Errorf("expected UUID %q, got %q", wantUUID, gotUUID)
	}
}

func TestUndefine(t *testing.T) {
	conn := libvirttest.New()
	l := New(conn)

	var flags UndefineFlags
	if err := l.Undefine("test", flags); err != nil {
		t.Fatalf("unexpected undefine error: %v", err)
	}
}

func TestDestroy(t *testing.T) {
	conn := libvirttest.New()
	l := New(conn)

	var flags DestroyFlags
	if err := l.Destroy("test", flags); err != nil {
		t.Fatalf("unexpected destroy error: %v", err)
	}
}

func TestVersion(t *testing.T) {
	conn := libvirttest.New()
	l := New(conn)

	version, err := l.Version()
	if err != nil {
		t.Error(err)
	}

	expected := "1.3.4"
	if version != expected {
		t.Errorf("expected version %q, got %q", expected, version)
	}
}

func TestDefineXML(t *testing.T) {
	conn := libvirttest.New()
	l := New(conn)

	var flags DomainDefineXMLFlags
	var buf []byte
	if err := l.DefineXML(buf, flags); err != nil {
		t.Fatalf("unexpected define error: %v", err)
	}
}
