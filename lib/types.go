package lib

import (
	"bufio"
	"net/http"
	"os/exec"

	"github.com/thejerf/suture"
	"golang.org/x/net/proxy"
)

// Iodine stores information needed to run iodine
type iodine struct {
	Domain     string
	Device     string
	NameServer string

	user       string
	stderr     *bufio.Reader
	cmd        *exec.Cmd
	supervisor *suture.Supervisor
}

// NetworkConfig stores OS X network configuration
type NetworkConfig struct {
	SocksHost   string
	SocksPort   int
	NetworkName string

	proxy  *proxy.Dialer
	client *http.Client
	parent *SocksProxy
}

// SocksProxy stores socks proxy configuration
type SocksProxy struct {
	SocksPort string
	Address   string

	cmd        *exec.Cmd
	stderr     *bufio.Reader
	supervisor *suture.Supervisor
	parent     *iodine
}
