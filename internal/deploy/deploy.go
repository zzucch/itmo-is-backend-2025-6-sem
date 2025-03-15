package deploy

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

func Setup() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/deploy", Deploy)

	return mux
}

func Deploy(responseWriter http.ResponseWriter, request *http.Request) {
	const scriptPath = "./build/deploy.sh"

	if err := killProcessOnPort("3000"); err != nil {
		log.Print(err)

		http.Error(
			responseWriter,
			err.Error(),
			http.StatusInternalServerError,
		)

		return
	}

	cmd := exec.Command("bash", scriptPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Print(err)

		http.Error(
			responseWriter,
			err.Error(),
			http.StatusInternalServerError,
		)

		return
	}

	responseWriter.Write(output)
}

func killProcessOnPort(port string) error {
	cmd := exec.Command("lsof", "-i", ":"+port, "-t")

	output, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(err.Error(), "exit status 1") {
			log.Printf("no process found on port=%s\n", port)
			return nil
		}

		return fmt.Errorf("failed to find process on port=%s: %v", port, err)
	}

	pids := strings.Fields(string(output))
	if len(pids) == 0 {
		log.Printf("no process found on port=%s\n", port)
		return nil
	}

	for _, pid := range pids {
		killCmd := exec.Command("kill", "-9", pid)

		if err := killCmd.Run(); err != nil {
			return fmt.Errorf("failed to kill process pid=%s: %v", pid, err)
		}

		log.Printf("killed process pid=%s port=%s\n", pid, port)
	}

	return nil
}
