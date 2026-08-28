package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func streamHandler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done(): // client disconnected
			log.Println("[i] client disconnected")
			return
		case t := <-ticker.C:
			fmt.Fprintf(w, "data: current time is %s\n\n", t.Format(time.RFC3339))
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
<div id="output">Waiting...</div>
</body>
</html>
`)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/stream", streamHandler)
	log.Println("[i] listening on :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
