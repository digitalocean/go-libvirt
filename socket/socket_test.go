package socket

import (
	"bytes"
	"testing"

	"github.com/digitalocean/go-libvirt/internal/constants"
)

var testHeader = []byte{
	0x20, 0x00, 0x80, 0x86, // program
	0x00, 0x00, 0x00, 0x01, // version
	0x00, 0x00, 0x00, 0x01, // procedure
	0x00, 0x00, 0x00, 0x00, // type
	0x00, 0x00, 0x00, 0x00, // serial
	0x00, 0x00, 0x00, 0x00, // status
}

func TestPktLen(t *testing.T) {
	data := []byte{0x00, 0x00, 0x00, 0xa} // uint32:10
	r := bytes.NewBuffer(data)

	expected := uint32(10)
	actual, err := pktlen(r)
	if err != nil {
		t.Error(err)
	}

	if expected != actual {
		t.Errorf("expected packet length %q, got %q", expected, actual)
	}
}

func TestExtractHeader(t *testing.T) {
	r := bytes.NewBuffer(testHeader)
	h, err := extractHeader(r)
	if err != nil {
		t.Error(err)
	}

	if h.Program != constants.Program {
		t.Errorf("expected Program %q, got %q", constants.Program, h.Program)
	}

	if h.Version != constants.ProtocolVersion {
		t.Errorf("expected version %q, got %q", constants.ProtocolVersion, h.Version)
	}

	if h.Procedure != constants.ProcConnectOpen {
		t.Errorf("expected procedure %q, got %q", constants.ProcConnectOpen, h.Procedure)
	}

	if h.Type != Call {
		t.Errorf("expected type %q, got %q", Call, h.Type)
	}

	if h.Status != StatusOK {
		t.Errorf("expected status %q, got %q", StatusOK, h.Status)
	}
}
