package lib

import (
	"bufio"
	"io"
	"log"
	"os/exec"
	"strings"

	"github.com/thejerf/suture"
)

//SSH_FLAGS="-o CompressionLevel=9 -C -v"
//SSH_FWD_FLAGS="-D 5000 -N"
//ADDRESS="root@192.168.1.99"

//ssh $SSH_FWD_FLAGS $SSH_FLAGS $ADDRESS

func createSocksProxy(supervisor *suture.Supervisor, parent *iodine) *SocksProxy {
	return &SocksProxy{
		SocksPort:  "5000",
		Address:    "root@192.168.1.99",
		supervisor: supervisor,
		parent:     parent,
	}
}

func (s *SocksProxy) configure() {
	s.cmd = exec.Command(
		"ssh",
		"-D", s.SocksPort,
		"-N",
		"-v",
		"-o", "CompressionLevel=9",
		"-C",
		s.Address,
	)
	stderr, err := s.cmd.StderrPipe()
	if err != nil {
		log.Panicf("Error when creating stderr pipe for socks proxy: %+v", err)
	}
	stdin, err := s.cmd.StdinPipe()
	if err != nil {
		log.Panicf("Error when creating stdin pipe for socks proxy: %+v", err)
	}
	if err := s.cmd.Start(); err != nil {
		log.Panicf("Error when starting socks proxy: %+v", err)
	}

	if _, err := io.WriteString(stdin, "yes\n"); err != nil {
		log.Panicf("Error when writing yes to socks proxy stdin: %+v", err)
	}
	s.stderr = bufio.NewReader(stderr)
}

// Serve SocksProxy
func (s *SocksProxy) Serve() {
	log.Printf("Serving Socks Proxy")
	s.configure()
	go readStderr(s.stderr, "Socks Proxy", func(arg string) {
		if strings.Contains(arg, "Entering interactive session.") {
			s.supervisor.Add(createNetworkConfig(s))
		}
	})
	log.Panicf("Error when running socks proxy: %+v", s.cmd.Wait())
}

// Stop SocksProxy
func (s *SocksProxy) Stop() {
	err := s.cmd.Process.Kill()
	if err != nil {
		log.Panicf("Error when killing socks proxy: %+v", err)
	}
	s.parent.Stop()
}
