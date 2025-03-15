package deploy

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/is-web-y26/m3302-milovatskiy/internal/config"
)

func Deploy() error {
	const deployScriptPath = "./build/deploy.sh"
	const runScriptPath = "./build/run.sh"

	if err := killProcessOnPort(config.GetDefaultConfig().Port); err != nil {
		return err
	}

	cmd := exec.Command("bash", deployScriptPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	log.Print(output)

	cmd = exec.Command("bash", runScriptPath)
	if err := cmd.Start(); err != nil {
		return err
	}

	go func(cmd *exec.Cmd) {
		if err := cmd.Wait(); err != nil {
			log.Printf("error waiting: err=%v", err)
		}

		log.Print("waited")
	}(cmd)

	return nil
}

func killProcessOnPort(port string) error {
	cmd := exec.Command("lsof", "-i", ":"+port, "-t")

	output, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(err.Error(), "exit status 1") {
			log.Printf("no process found on port=%s", port)
			return nil
		}

		return fmt.Errorf("failed to find process on port=%s: %v", port, err)
	}

	pids := strings.Fields(string(output))
	if len(pids) == 0 {
		log.Printf("no process found on port=%s", port)
		return nil
	}

	for _, pid := range pids {
		killCmd := exec.Command("kill", "-9", pid)

		if err := killCmd.Run(); err != nil {
			return fmt.Errorf("failed to kill process pid=%s: %v", pid, err)
		}

		log.Printf("killed process pid=%s port=%s", pid, port)
	}

	return nil
}
