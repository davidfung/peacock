package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/davidfung/peacock/probe/alive"
	"github.com/davidfung/peacock/probe/motd"
)

const chanSize = 100

func addTS(msg string) string {
	msg = fmt.Sprintf("[%s] %s", time.Now().Format("060102 150405"), msg)
	return msg
}

func streamHandler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ch := make(chan string, chanSize)

	alive.Load(ch)
	motd.Load(ch)

	for {
		select {
		case <-r.Context().Done(): // client disconnected
			log.Println("[i] client disconnected")
			return
		case msg := <-ch:
			fmt.Fprintf(w, "data: %s\n\n", addTS(msg))
			flusher.Flush()
		}
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("[i] new client connected")
	fmt.Fprint(w, `
<!DOCTYPE html>
<html>
<body>
<script>
const evtSource = new EventSource('/stream');
evtSource.onmessage = (event) => {
  const line = document.createElement('div');
  line.textContent = event.data;
  document.getElementById('output').appendChild(line);
  window.scrollTo(0, document.body.scrollHeight);
};
</script>
<div id="output">Welcome to Peacock!</div>
</body>
</html>
`)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/stream", streamHandler)
	svr := http.Server{
		Addr:         ":8000",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}
	log.Println("[i] listening on :8000")
	log.Fatal(svr.ListenAndServe())
}
