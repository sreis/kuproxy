package haproxy

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// initial haproxy configuration file
const template string = `
global
    daemon
    pidfile /var/run/haproxy.pid
    stats socket /var/run/haproxy.sock mode 600 level admin

defaults
    mode tcp

# we need an entry like this to ensure haproxy is running
listen nginx *:80
    mode http
    balance roundrobin
`

// Create temporary haproxy config file. Don't forget to call:
//     defer os.Remove(config.Name())
// after calling this function.
func makeConfig() (config *os.File, err error) {

	if config, err = ioutil.TempFile("", "haproxy"); err != nil {
		return nil, err
	}

	if _, err = config.Write([]byte(template)); err != nil {
		return nil, err
	}

	return config, nil
}

// Start haproxy instance. haproxy binary must be in $PATH.
// Uses pid file option to allow configuration reload.
func Start() error {
	binary, err := exec.LookPath("haproxy")
	if err != nil {
		return err
	}

	config, err := makeConfig()
	if err != nil {
		return nil
	}
	defer os.Remove(config.Name())

	log.Printf("Starting haproxy.")

	cmd := exec.Command(binary, "-f", config.Name())
	err = cmd.Start()

	// wait for haproxy to startup
	time.Sleep(200 * time.Millisecond)

	return err
}

// Generate configuration file and reload haproxy.
func Reload() error {
	log.Println("Reloading haproxy configuration.")

	binary, err := exec.LookPath("haproxy")
	if err != nil {
		return err
	}

	config, err := makeConfig()
	if err != nil {
		return nil
	}
	defer os.Remove(config.Name())

	cmd := exec.Command(binary, "-f", config.Name(), "-p", "/var/run/haproxy.pid", "-st", "$(cat /var/run/haproxy.pid)")
	err = cmd.Start()

	return err
}

// We need to stop haproxy the old fashioned way. :)
func Stop() error {
	log.Println("Stopping haproxy.")

	b, err := ioutil.ReadFile("/var/run/haproxy.pid")
	if err != nil {
		return err
	}

	pids := strings.Split(string(b), "\n")
	for _, pid := range pids {
		if i, err := strconv.Atoi(pid); err == nil {
			syscall.Kill(i, 9)
		}
	}
	return nil
}

// Helper function to send an acl command to haproxy.
func acl(command string) (status string, err error) {
	conn, err := net.Dial("unix", "/var/run/haproxy.sock")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	fmt.Fprintf(conn, command)

	status, err = bufio.NewReader(conn).ReadString('\n')

	return status, nil

}

func ShowStat() error {
	status, err := acl("show stat\r\n")
	if err != nil {
		return err
	}

	log.Printf("> show stat: %s", status)
	return nil
}
