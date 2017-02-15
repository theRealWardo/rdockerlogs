package host

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/jeromer/syslogparser/rfc3164"
)

const DOCKER_TAG = "docker/"

type Host struct {
	Ip            string
	User          string
	IdentityFile  string
	InspectFormat string
	containers    map[string]string
}

func (h *Host) cmd(cmd string) *exec.Cmd {
	return exec.Command("ssh", fmt.Sprintf("%s@%s", h.User, h.Ip), "-i", h.IdentityFile, cmd)
}

func (h *Host) resolveContainer(id string) string {
	if value, ok := h.containers[id]; ok {
		return value
	}
	cmd := h.cmd(fmt.Sprintf("docker inspect -f '%s' %s", h.InspectFormat, id))
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Run()
	l, _ := stdout.ReadString('\n')
	lines := strings.Split(l, "\n")
	l = lines[0]
	if h.containers == nil {
		h.containers = make(map[string]string)
	}
	if len(l) > 0 {
		h.containers[id] = l
		return l
	}
	h.containers[id] = id
	return id
}

func (h *Host) WriteLines(out chan string, data string) {
	lines := strings.Split(data, "\n")
	for _, line := range lines {
		if len(line) > 0 {
			// Hack in the priority...
			p := rfc3164.NewParser([]byte("<34>" + line))
			err := p.Parse()
			if err != nil {
				out <- fmt.Sprintf("%s E - %s", h.Ip, err)
				out <- fmt.Sprintf("%s   - %s", h.Ip, line)
				continue
			}
			d := p.Dump()
			tag := d["tag"].(string)
			if strings.Index(tag, DOCKER_TAG) == 0 {
				tag = h.resolveContainer(tag[len(DOCKER_TAG):])
			}
			out <- fmt.Sprintf("%s [%s] - %s", h.Ip, tag, d["content"])
		}
	}
}

func (h *Host) Logs(out chan string) {
	cmd := h.cmd("tail -f /var/log/syslog")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	for {
		l, _ := stdout.ReadString('\n')
		h.WriteLines(out, l)
		l, _ = stderr.ReadString('\n')
		h.WriteLines(out, l)
	}
}
