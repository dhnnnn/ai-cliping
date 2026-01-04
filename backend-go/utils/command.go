package utils

import (
	"log"
	"os/exec"
)

func RunCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Println("Command error:", err)
		log.Println(string(output))
		return err
	}

	log.Println(string(output))
	return nil
}
