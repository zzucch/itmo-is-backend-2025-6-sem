package deploy

import (
	"encoding/json"
	"log"
	"net/http"
	"sync/atomic"
)

func Setup() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/deploy", HandleDeploy)

	return mux
}

var deployInProgress atomic.Bool

func HandleDeploy(responseWriter http.ResponseWriter, request *http.Request) {
	if !deployInProgress.CompareAndSwap(false, true) {
		http.Error(
			responseWriter,
			"deployment in progress",
			http.StatusConflict,
		)

		return
	}
	defer deployInProgress.Store(false)

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

	if err := Deploy(); err != nil {
		log.Print(err)

		http.Error(
			responseWriter,
			err.Error(),
			http.StatusInternalServerError,
		)
	}
}
