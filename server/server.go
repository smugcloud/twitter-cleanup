package server

import (
	"log"
	"net/http"
)

// HTTP server that will expose an API to query health
// and provide details on the last run

type Server struct{}

func (s *Server) HttpServer() {
	helloHandler := func(w http.ResponseWriter, req *http.Request) {
	}

	http.HandleFunc("/hello", helloHandler)
	log.Fatal(http.ListenAndServe(":9000", nil))
}
