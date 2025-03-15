package server

import (
	"fmt"
	"net/http"
	"time"
)

func (controllers *controllers) notificationsHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	http.ServeFile(w, r, "./static/notifications.html")
}

func (controllers *controllers) notificationsSSEHandler(
	w http.ResponseWriter,
	_ *http.Request,
) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	for {
		select {
		case <-time.After(200 * time.Millisecond):
			message := fmt.Sprintf(
				"data: Update at %s\n\n",
				time.Now().Format(time.TimeOnly),
			)
			_, err := w.Write([]byte(message))
			if err != nil {
				return
			}

			flusher.Flush()
		}
	}
}
