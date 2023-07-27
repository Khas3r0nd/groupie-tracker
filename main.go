package main

import (
	"log"
	"net/http"

	"groupie-tracker-filters/internal"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))
	mux.HandleFunc("/", internal.Home)
	mux.HandleFunc("/artists/", internal.Concert)
	mux.HandleFunc("/search",internal.SearchHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Printf("Localhost on http://localhost:8080")
	err := server.ListenAndServe()
	log.Fatal(err)
}
