package main

import (
	"fmt"
	"log"
	"net/http"
)

const port = 8092

func main() {
	// configure the songs directory name and port
	const songsDir = "songs"
	const filmsDir = "./assets"

	// add a handler for the song files
	http.Handle("/", addHeaders(http.FileServer(http.Dir(filmsDir))))
	http.HandleFunc("/ping", HelloServer)
	fmt.Printf("Starting server on %v\n", port)
	//log.Printf("Serving %s on HTTP port: %v\n", filmsDir, port)

	// serve and log errors
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

// addHeaders will act as middleware to give us CORS support
func addHeaders(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		h.ServeHTTP(w, r)
	}
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}
