package lib

import (
	"fmt"
	"golang.org/x/net/proxy"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"time"
)

// curl -v --socks5-hostname 127.0.0.1:5000 http://httpbin.org/ip

const (
	testURL       = "https://httpbin.org/ip"
	checkInterval = 30
)

func createNetworkConfig(parent *SocksProxy) *NetworkConfig {
	host := "127.0.0.1"
	port := 5000
	proxy, err := proxy.SOCKS5(
		"tcp", fmt.Sprintf("%s:%d", host, port), nil, proxy.Direct)
	if err != nil {
		log.Panicf("Error when creating proxy: %+v", err)
	}
	transport := &http.Transport{Dial: proxy.Dial}
	return &NetworkConfig{
		NetworkName: "Wi-Fi",
		proxy:       &proxy,
		SocksHost:   host,
		SocksPort:   port,
		client: &http.Client{
			Transport: transport,
		},
		parent: parent,
	}
}

func (n *NetworkConfig) testProxy() error {
	log.Printf("Performing test ping to %s", testURL)
	start := time.Now()
	res, err := n.client.Get(testURL)
	if err != nil {
		return fmt.Errorf(
			"Error when retriving %s via SOCKS proxy: %+v", testURL, err,
		)
	}
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("Error when reading from %s: %+v", testURL, err)
	}

	if err := res.Body.Close(); err != nil {
		return fmt.Errorf("Error when closing HTTP connection: %+v", err)
	}
	log.Printf(
		"Socks proxy works. Server response after %.2f s: %s",
		time.Now().Sub(start).Seconds(),
		content,
	)
	return nil
}

func (n *NetworkConfig) networkSetup() error {
	log.Printf("Performing networksetup")
	cmd := exec.Command(
		"networksetup",
		"-setsocksfirewallproxy", n.NetworkName,
		n.SocksHost, fmt.Sprintf("%d", n.SocksPort),
	)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Error when running networksetup: %+v", err)
	}
	return nil
}

func (n *NetworkConfig) networkSetupOff() error {
	cmd := exec.Command(
		"networksetup",
		"-setsocksfirewallproxystate", n.NetworkName, "off",
	)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Error when running networksetup: %+v", err)
	}
	return nil
}

// Serve network configuration
func (n *NetworkConfig) Serve() {
	log.Printf("Serving NetworkConfig")
	if err := n.testProxy(); err != nil {
		log.Panicf("Error when testing Proxy: %+v", err)
	}
	if err := n.networkSetup(); err != nil {
		log.Panicf("Error when performing network setup: %+v", err)
	}

	for {
		time.Sleep(checkInterval * time.Second)
		if err := n.testProxy(); err != nil {
			log.Panicf("Error when testing Proxy: %+v", err)
		}
		log.Printf("Checking again after %d seconds", checkInterval)
	}
}

// Stop network configuration
func (n *NetworkConfig) Stop() {
	n.parent.Stop()
	log.Print("Tearing down Socks proxy network setup")
	if err := n.networkSetupOff(); err != nil {
		log.Panicf("Error when disabling socks proxy: %+v", err)
	}
}
