package haproxy

import (
	"log"
	"os/exec"
)

var cmd *exec.Cmd

// Start haproxy instance. haproxy binary must be in $PATH.
func Start() error {
	binary, err := exec.LookPath("haproxy")
	if err != nil {
		return err
	}

	log.Println("Starting haproxy.")
	cmd = exec.Command(binary)
	err = cmd.Start()

	return err
}

// Stop haproxy instance.
func Stop() error {
	log.Println("Stopping haproxy.")
	err := cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}
