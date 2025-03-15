package deploy

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"sync/atomic"
)

func Setup() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/deploy", Deploy)

	return mux
}

var deployingInProgress atomic.Bool

func Deploy(responseWriter http.ResponseWriter, request *http.Request) {
	if !deployingInProgress.CompareAndSwap(false, true) {
		http.Error(
			responseWriter,
			"deployment in progress",
			http.StatusConflict,
		)

		return
	}
	defer deployingInProgress.Store(false)

	const deployScriptPath = "./build/deploy.sh"
	const runScriptPath = "./build/run.sh"

	requesterInfo := map[string]any{
		"remote_addr":    request.RemoteAddr,
		"method":         request.Method,
		"url":            request.URL.String(),
		"host":           request.Host,
		"proto":          request.Proto,
		"headers":        request.Header,
		"content_length": request.ContentLength,
		"referer":        request.Referer(),
		"user_agent":     request.UserAgent(),
	}

	requesterInfoJSON, err := json.Marshal(requesterInfo)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("deploy requested: %s", string(requesterInfoJSON))

	if err := killProcessOnPort("3000"); err != nil {
		log.Print(err)

		http.Error(
			responseWriter,
			err.Error(),
			http.StatusInternalServerError,
		)

		return
	}

	cmd := exec.Command("bash", deployScriptPath)
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

	cmd = exec.Command("bash", runScriptPath)
	if err := cmd.Start(); err != nil {
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
