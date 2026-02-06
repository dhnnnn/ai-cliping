package utils

import (
	"log"
	"os/exec"
)

func RunCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()

	// Always log output for debugging
	if len(output) > 0 {
		log.Println(string(output))
	}

	if err != nil {
		log.Printf("Command failed: %v", err)
		return err
	}

	return nil
}
