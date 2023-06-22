package libvirt

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/digitalocean/go-libvirt/socket"
	"github.com/digitalocean/go-libvirt/socket/dialers"
)

// ConnectToURI returns a new, connected client instance using the appropriate
// dialer for the given libvirt URI.
func ConnectToURI(uri *url.URL) (*Libvirt, error) {
	dialer, err := dialerForURI(uri)
	if err != nil {
		return nil, err
	}

	lv := NewWithDialer(dialer)

	if err := lv.ConnectToURI(RemoteURI(uri)); err != nil {
		return nil, fmt.Errorf("failed to connect to libvirt: %w", err)
	}

	return lv, nil
}

// RemoteURI returns the libvirtd URI corresponding to a given client URI.
// The client URI contains details of the connection method, but once connected
// to libvirtd, all connections are local. So e.g. the client may want to
// connect to qemu+tcp://example.com/system but once the socket is established
// it will ask the remote libvirtd for qemu:///system.
func RemoteURI(uri *url.URL) ConnectURI {
	remoteURI := (&url.URL{
		Scheme: strings.Split(uri.Scheme, "+")[0],
		Path:   uri.Path,
	}).String()
	if name := uri.Query().Get("name"); name != "" {
		remoteURI = name
	}
	return ConnectURI(remoteURI)
}

func dialerForURI(uri *url.URL) (socket.Dialer, error) {
	transport := "unix"
	if scheme := strings.SplitN(uri.Scheme, "+", 2); len(scheme) > 1 {
		transport = scheme[1]
	} else if uri.Host != "" {
		transport = "tls"
	}

	switch transport {
	case "unix":
		options := []dialers.LocalOption{}
		if s := uri.Query().Get("socket"); s != "" {
			options = append(options, dialers.WithSocket(s))
		}
		if err := checkModeOption(uri); err != nil {
			return nil, err
		}
		return dialers.NewLocal(options...), nil
	case "tcp":
		options := []dialers.RemoteOption{}
		if port := uri.Port(); port != "" {
			options = append(options, dialers.UsePort(port))
		}
		return dialers.NewRemote(uri.Hostname(), options...), nil
	default:
		return nil, fmt.Errorf("unsupported libvirt transport %s", transport)
	}
}

func checkModeOption(uri *url.URL) error {
	mode := uri.Query().Get("mode")
	switch strings.ToLower(mode) {
	case "":
	case "legacy", "auto":
	case "direct":
		return errors.New("cannot connect in direct mode")
	default:
		return fmt.Errorf("invalid ssh mode %v", mode)
	}
	return nil
}
