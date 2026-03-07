package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const port = ":8080"
	mux := http.NewServeMux()
	apiConf := apiConfig{}
	fsHandler := http.StripPrefix("/app/", apiConf.middlewareMetricsInc(http.FileServer(http.Dir("."))))
	mux.Handle("/app/", fsHandler)
	mux.HandleFunc("GET /healthz", handlerReadiness)
	mux.HandleFunc("GET /metrics", apiConf.handlerMetrics)
	mux.HandleFunc("POST /reset", apiConf.handlerReset)

	server := &http.Server{Handler: mux, Addr: port}
	log.Printf("Serving on port: %s", port)
	log.Fatal(server.ListenAndServe())
}
