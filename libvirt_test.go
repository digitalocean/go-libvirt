// Copyright 2018 The go-libvirt Authors.
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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/digitalocean/go-libvirt/libvirttest"
)

func TestConnect(t *testing.T) {
	conn := libvirttest.New()
	l := New(conn)

	err := l.Connect()
	if err != nil {
		t.Error(err)
	}
	defer l.Disconnect()
}

func TestLibvirt_ConnectToURI(t *testing.T) {
	type args struct {
		uri ConnectURI
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "connect to test:///default socket",
			args: args{
				uri: TestDefault,
			},
			wantErr: false,
		},
		{
			name: "connect to qemu:///session socket",
			args: args{
				uri: QEMUSession,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn := libvirttest.New()
			libvirtConn := New(conn)

			if err := libvirtConn.ConnectToURI(tt.args.uri); (err != nil) != tt.wantErr {
				t.Errorf("ConnectToURI() error = %v, wantErr %v", err, tt.wantErr)
			}
			libvirtConn.Disconnect()
		})
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

func TestDisconnectCleanup(t *testing.T) {
	conn := libvirttest.New()
	l := New(conn)

	err := l.Connect()
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}

	// register a callback
	ch := make(chan response)
	l.register(1, ch)

	err = l.Disconnect()
	if err != nil {
		t.Errorf("disconnect failed: %v", err)
	}

	// now make sure that the Disconnect deregistered the callback
	// do it via a select just to avoid the 10 minute unit test timeout on
	// failure
	select {
	case <-ch:
		t.Log("callback channel appropriately closed")
	case <-time.After(10 * time.Second):
		t.Errorf("callback channel not closed")
	}
}

func TestLostConnectionCleanup(t *testing.T) {
	conn := libvirttest.New()
	l := New(conn)

	err := l.Connect()
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}

	// register a callback
	ch := make(chan response)
	l.register(1, ch)

	// forcibly just close the fake libvirt tester side of the pipe to lose
	// the connection.
	// This should be detected and cleaned up without having to call Disconnect.
	conn.Test.Close()

	// now make sure that deregistered the callback
	// do it via a select just to avoid the 10 minute unit test timeout on
	// failure
	select {
	case <-ch:
		t.Log("callback channel appropriately closed")
	case <-time.After(10 * time.Second):
		t.Fatalf("callback channel not closed")
	}
}

func TestMigrate(t *testing.T) {
	conn := libvirttest.New()
	l := New(conn)

	err := l.Connect()
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer l.Disconnect()

	flags := MigrateLive |
		MigratePeer2peer |
		MigratePersistDest |
		MigrateChangeProtection |
		MigrateAbortOnError |
		MigrateAutoConverge |
		MigrateNonSharedDisk

	dom, err := l.DomainLookupByName("test")
	if err != nil {
		t.Fatalf("failed to lookup domain: %v", err)
	}
	dconnuri := []string{"qemu+tcp://foo/system"}
	if _, err := l.DomainMigratePerform3Params(dom, dconnuri,
		[]TypedParam{}, []byte{}, flags); err != nil {
		t.Fatalf("unexpected live migration error: %v", err)
	}
}

func TestMigrateSetMaxSpeed(t *testing.T) {
	conn := libvirttest.New()
	l := New(conn)

	err := l.Connect()
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer l.Disconnect()

	dom, err := l.DomainLookupByName("test")
	if err != nil {
		t.Fatalf("failed to lookup domain: %v", err)
	}

	if err := l.DomainMigrateSetMaxSpeed(dom, 100, 0); err != nil {
		t.Fatalf("unexpected error setting max speed for migrate: %v", err)
	}
}

func TestDomains(t *testing.T) {
	conn := libvirttest.New()
	l := New(conn)

	err := l.Connect()
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer l.Disconnect()

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
		if d.ID != int32(wantID) {
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

	err := l.Connect()
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer l.Disconnect()

	wantState := DomainState(DomainRunning)
	gotState, err := l.DomainState("test")
	if err != nil {
		t.Error(err)
	}

	if gotState != wantState {
		t.Errorf("expected domain state %d, got %d", wantState, gotState)
	}
}

func TestDomainMemoryStats(t *testing.T) {
	conn := libvirttest.New()
	l := New(conn)

	err := l.Connect()
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer l.Disconnect()

	wantDomainMemoryStats := []DomainMemoryStat{
		{
			Tag: 6,
			Val: 1048576,
		},
		{
			Tag: 7,
			Val: 91272,
		},
	}

	d, err := l.lookup("test")
	if err != nil {
		t.Error(err)
	}

	gotDomainMemoryStats, err := l.DomainMemoryStats(d, 8, 0)
	if err != nil {
		t.Error(err)
	}

	if len(gotDomainMemoryStats) == 0 {
		t.Error("No memory stats returned!")
	}

	for i := range wantDomainMemoryStats {
		if wantDomainMemoryStats[i] != gotDomainMemoryStats[i] {
			t.Errorf("expected domain memory stat %v, got %v", wantDomainMemoryStats[i], gotDomainMemoryStats[i])
		}
	}
}

func TestEvents(t *testing.T) {
	conn := libvirttest.New()
	l := New(conn)

	err := l.Connect()
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer l.Disconnect()

	done := make(chan error)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream, err := l.SubscribeQEMUEvents(ctx, "test")
	if err != nil {
		t.Error(err)
	}

	go func() {
		var e DomainEvent

		select {
		case e = <-stream:
		case <-time.After(time.Second * 5):
			done <- errors.New("expected event, received timeout")
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

		done <- nil
	}()

	// send an event to the listener goroutine
	conn.Test.Write(append(testEventHeader, testEvent...))

	// wait for completion
	if err := <-done; err != nil {
		t.Fatal(err)
	}
}

func TestRun(t *testing.T) {
	conn := libvirttest.New()
	l := New(conn)

	err := l.Connect()
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer l.Disconnect()

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

	err := l.Connect()
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer l.Disconnect()

	_, err = l.Run("test", []byte(`{"drive-foo"}`))
	if err == nil {
		t.Error("expected qemu error")
	}
}

func TestSecrets(t *testing.T) {
	conn := libvirttest.New()
	l := New(conn)

	err := l.Connect()
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer l.Disconnect()

	secrets, err := l.Secrets()
	if err != nil {
		t.Fatal(err)
	}

	wantLen := 1
	gotLen := len(secrets)
	if gotLen != wantLen {
		t.Fatalf("expected %d secrets, got %d", wantLen, gotLen)
	}

	s := secrets[0]
	wantType := int32(SecretUsageTypeVolume)
	if s.UsageType != wantType {
		t.Errorf("expected usage type %d, got %d", wantType, s.UsageType)
	}

	wantID := "/tmp"
	if s.UsageID != wantID {
		t.Errorf("expected usage id %q, got %q", wantID, s.UsageID)
	}

	// 19fdc2f2-fa64-46f3-bacf-42a8aafca6dd
	wantUUID := [UUIDBuflen]byte{
		0x19, 0xfd, 0xc2, 0xf2, 0xfa, 0x64, 0x46, 0xf3,
		0xba, 0xcf, 0x42, 0xa8, 0xaa, 0xfc, 0xa6, 0xdd,
	}
	if s.UUID != wantUUID {
		t.Errorf("expected UUID %q, got %q", wantUUID, s.UUID)
	}
}

func TestStoragePool(t *testing.T) {
	conn := libvirttest.New()
	l := New(conn)

	err := l.Connect()
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer l.Disconnect()

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
	wantUUID := [UUIDBuflen]byte{
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

	err := l.Connect()
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer l.Disconnect()

	pools, err := l.StoragePools(ConnectListStoragePoolsTransient)
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
	wantUUID := [UUIDBuflen]byte{
		0xbb, 0x30, 0xa1, 0x1c, 0x08, 0x46, 0x48, 0x27,
		0x8b, 0xba, 0x3e, 0x6b, 0x5c, 0xf1, 0xb6, 0x5f,
	}
	gotUUID := pools[0].UUID
	if gotUUID != wantUUID {
		t.Errorf("expected UUID %q, got %q", wantUUID, gotUUID)
	}
}

func TestStoragePoolRefresh(t *testing.T) {
	conn := libvirttest.New()
	l := New(conn)

	err := l.Connect()
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer l.Disconnect()

	pool, err := l.StoragePool("default")
	if err != nil {
		t.Error(err)
	}
	err = l.StoragePoolRefresh(pool, 0)
	if err != nil {
		t.Error(err)
	}
}

func TestUndefine(t *testing.T) {
	conn := libvirttest.New()
	l := New(conn)

	err := l.Connect()
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer l.Disconnect()

	var flags DomainUndefineFlagsValues
	if err := l.Undefine("test", flags); err != nil {
		t.Fatalf("unexpected undefine error: %v", err)
	}
}

func TestDestroy(t *testing.T) {
	conn := libvirttest.New()
	l := New(conn)

	err := l.Connect()
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer l.Disconnect()

	var flags DomainDestroyFlagsValues
	if err := l.Destroy("test", flags); err != nil {
		t.Fatalf("unexpected destroy error: %v", err)
	}
}

func TestVersion(t *testing.T) {
	conn := libvirttest.New()
	l := New(conn)

	err := l.Connect()
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer l.Disconnect()

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

	err := l.Connect()
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer l.Disconnect()

	var flags DomainDefineFlags
	var buf []byte
	if err := l.DefineXML(buf, flags); err != nil {
		t.Fatalf("unexpected define error: %v", err)
	}
}

func TestDomainCreateWithFlags(t *testing.T) {
	conn := libvirttest.New()
	l := New(conn)

	err := l.Connect()
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer l.Disconnect()

	d, err := l.lookup("test")
	if err != nil {
		t.Fatalf("failed to lookup domain: %v", err)
	}
	var flags uint32
	if _, err := l.DomainCreateWithFlags(d, flags); err != nil {
		t.Fatalf("unexpected create error: %v", err)
	}
}

func TestShutdown(t *testing.T) {
	conn := libvirttest.New()
	l := New(conn)

	err := l.Connect()
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer l.Disconnect()

	var flags DomainShutdownFlagValues
	if err := l.Shutdown("test", flags); err != nil {
		t.Fatalf("unexpected shutdown error: %v", err)
	}
}

func TestReboot(t *testing.T) {
	conn := libvirttest.New()
	l := New(conn)

	err := l.Connect()
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer l.Disconnect()

	var flags DomainRebootFlagValues
	if err := l.Reboot("test", flags); err != nil {
		t.Fatalf("unexpected reboot error: %v", err)
	}
}

func TestReset(t *testing.T) {
	conn := libvirttest.New()
	l := New(conn)

	err := l.Connect()
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer l.Disconnect()

	if err := l.Reset("test"); err != nil {
		t.Fatalf("unexpected reset error: %v", err)
	}
}

func TestSetBlockIOTune(t *testing.T) {
	conn := libvirttest.New()
	l := New(conn)

	err := l.Connect()
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer l.Disconnect()

	if err := l.SetBlockIOTune("test", "vda", BlockLimit{"write_bytes_sec", 5000000}); err != nil {
		t.Fatalf("unexpected SetBlockIOTune error: %v", err)
	}
}

func TestGetBlockIOTune(t *testing.T) {
	conn := libvirttest.New()
	l := New(conn)

	err := l.Connect()
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer l.Disconnect()

	limits, err := l.GetBlockIOTune("do-test", "vda")
	if err != nil {
		t.Fatalf("unexpected GetBlockIOTune error: %v", err)
	}

	lim := BlockLimit{"write_bytes_sec", 500000}
	if limits[2] != lim {
		t.Fatalf("unexpected result in limits list: %v", limits[2])
	}
}

// ConnectGetAllDomainStats is a good test, since its return values include an
// array of TypedParams embedded in a struct.
func TestConnectGetAllDomainStats(t *testing.T) {
	conn := libvirttest.New()
	l := New(conn)

	err := l.Connect()
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer l.Disconnect()

	doms, _, err := l.ConnectListAllDomains(1, 0)
	if err != nil {
		t.Fatalf("unexpected error listing all domains")
	}

	stats, err := l.ConnectGetAllDomainStats(doms, uint32(DomainStatsState), 0)
	if err != nil {
		t.Fatalf("unexpected error getting domain stats")
	}

	// Verify that the stats records we got back look like this:
	// [{state.state {1 1}} {state.reason {1 2}}]
	// [{state.state {1 1}} {state.reason {1 1}}]
	if len(stats) != 2 {
		t.Fatalf("expected 2 domain stats records, got %d", len(stats))
	}

	tp1 := NewTypedParamValueInt(1)
	if stats[1].Params[0].Field != "state.state" {
		t.Fatalf("unexpected domain statistic %v", stats[1].Params[0].Field)
	}
	if stats[1].Params[1].Field != "state.reason" {
		t.Fatalf("unexpected domain statistic %v", stats[1].Params[1].Field)
	}
	if stats[1].Params[0].Value != *tp1 {
		t.Fatalf("unexpected typed param value %v", stats[1].Params[0].Value)
	}
	if stats[1].Params[1].Value != *tp1 {
		t.Fatalf("unexpected typed param value %v", stats[1].Params[1].Value)
	}
}
