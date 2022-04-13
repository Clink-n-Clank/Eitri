package handler

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Clink-n-Clank/Eitri/internal/log"
)

// Install go install something by path
func Install(path ...string) error {
	for _, p := range path {
		if !strings.Contains(p, "@") {
			p += "@latest"
		}

		log.PrintInfo(fmt.Sprintf("go install %s", p))

		cmd := exec.Command("go", "install", p)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}
