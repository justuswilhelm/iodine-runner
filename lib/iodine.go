package lib

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/thejerf/suture"
)

// CreateIodine creates a valid Iodine configuration
func createIodine(supervisor *suture.Supervisor) *iodine {
	result := &iodine{
		Domain:     "t1.justus.pw",
		Device:     "utunX",
		supervisor: supervisor,
	}

	if os.Getenv("IODINE_PASS") == "" {
		log.Fatalf("Expecting IODINE_PASS to be in environment.")
	}

	if os.Getenv("IODINE_NAMESERVER") == "" {
		result.NameServer = os.Getenv("IODINE_NAMESERVER")
	}
	return result
}

func (i *iodine) configure() {
	stderr, err := i.cmd.StderrPipe()
	if err != nil {
		log.Panicf("Could not connect iodine stderr: %+v", err)
	}
	i.stderr = bufio.NewReader(stderr)
	if err := i.cmd.Start(); err != nil {
		log.Panicf("Could not start iodine: %+v", err)
	}
}

func (i *iodine) log() bool {
	stderr, err := ioutil.ReadAll(i.stderr)
	if err != nil {
		log.Panicf("Error reading from stderr: %+v", err)
	}
	log.Printf("Iodine stderr output: %s", stderr)
	return len(stderr) > 0

}

// Start starts the iodine process
func (i *iodine) Serve() {
	log.Printf("Serving iodine")
	i.cmd = exec.Command(
		"iodine",
		"-f",
		"-d", i.Device,
		"-I1",
	)
	if i.NameServer != "" {
		i.cmd.Args = append(i.cmd.Args, i.NameServer)
	}
	i.cmd.Args = append(i.cmd.Args, i.Domain)
	i.configure()

	readStderr(i.stderr, "Iodine", func(arg string) {
		if strings.Contains(arg, "Connection setup complete") {
			i.supervisor.Add(createSocksProxy(i.supervisor, i))
		}
	})
	err := i.cmd.Wait()
	if err != nil {
		log.Panicf("Iodine did not start successfully: %+v", err)
	}
}

func (i *iodine) Stop() {
	proc := i.cmd.Process
	if err := proc.Kill(); err != nil {
		log.Panic(err)
	}
}
