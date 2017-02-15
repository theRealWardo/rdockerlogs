package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/therealwardo/rdockerlogs/host"
)

var hostStr = flag.String("hosts", "", "Hosts to tail logs on.")
var userStr = flag.String("user", "", "The username to login to the host as.")
var inspectFormat = flag.String("format", "{{.Name}}", "The format to pass to docker inspect, used for the first block of each log line.")
var rancherFile = flag.String("rancher", "", "The Rancher config file to use if talking to Rancher.")
var rancherService = flag.String("service", "", "The Rancher service to tail logs for.")
var identityFile = flag.String("i", "", "The identity file to use with SSH.")
var configFile = flag.String("config", "", "A config file to use instead of flags.")

type ConfigData struct {
	hosts         []string
	user          string
	inspectFormat string
	rancherFile   string
	// rancherConfigData *rancher.ConfigData
	rancherService string
	identityFile   string
}

func buildConfigData() *ConfigData {
	if len(*configFile) > 0 {
		if len(*hostStr) > 0 {
			log.Fatal("-hosts cannot be specified with -config")
		}
		if len(*userStr) > 0 {
			log.Fatal("-user cannot be specified with -config")
		}
		if len(*inspectFormat) > 0 {
			log.Fatal("-format cannot be specified with -config")
		}
		if len(*rancherFile) > 0 {
			log.Fatal("-rancher cannot be specified with -config")
		}
		if len(*rancherService) > 0 {
			log.Fatal("-service cannot be specified with -config")
		}
		if len(*identityFile) > 0 {
			log.Fatal("-i cannot be specified with -config")
		}
		// TODO: parse config files.
		return &ConfigData{}
	}

	config := &ConfigData{}

	if len(*hostStr) > 0 {
		config.hosts = strings.Split(*hostStr, ",")
	}
	if len(*userStr) > 0 {
		config.user = *userStr
	}
	if len(*inspectFormat) > 0 {
		config.inspectFormat = *inspectFormat
	}
	if len(*rancherFile) > 0 {
		fmt.Printf("loading rancher config... %s\n", *rancherFile)
		// config.rancherConfigData = rancher.LoadConfig()
	}
	if len(*rancherService) > 0 {
		config.rancherService = *rancherService
	}
	if len(*identityFile) > 0 {
		config.identityFile = *identityFile
	}
	return config
}

func resolveHosts(config *ConfigData) []*host.Host {
	hosts := make([]*host.Host, 0)
	longestIp := 0
	for _, ip := range config.hosts {
		if len(ip) > longestIp {
			longestIp = len(ip)
		}
	}
	for _, ip := range config.hosts {
		hosts = append(hosts, &host.Host{
			Ip:            ip,
			User:          config.user,
			IdentityFile:  config.identityFile,
			InspectFormat: config.inspectFormat,
		})
	}
	return hosts
}

func main() {
	flag.Parse()

	config := buildConfigData()
	// Prepare the output channel.
	output := make(chan string)
	// Close connections channel.
	close := make(chan bool)
	// Listen for SIGINT to tell us to close down.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			close <- true
		}
	}()
	// Resolve hosts to connect to.
	hosts := resolveHosts(config)
	for _, h := range hosts {
		// Pull logs from the host.
		go h.Logs(output)
	}
	for {
		select {
		case data := <-output:
			fmt.Println(data)
		case <-close:
			fmt.Println("\n\nSIGINT received, shutting down...")
			return
		}
	}
}
