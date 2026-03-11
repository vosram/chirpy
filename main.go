package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/vosram/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

func main() {
	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL must be set")
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}

	dbConn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("couldn't open db: %s\n%v\n", dbUrl, err)
	}
	dbQueries := database.New(dbConn)

	const port = ":8080"
	mux := http.NewServeMux()
	apiConf := apiConfig{db: dbQueries, platform: platform}

	fsHandler := http.StripPrefix("/app/", apiConf.middlewareMetricsInc(http.FileServer(http.Dir("."))))
	mux.Handle("/app/", fsHandler)
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/users", apiConf.handlerCreateUser)
	mux.HandleFunc("POST /api/chirps", apiConf.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", apiConf.handlerGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiConf.handlerGetChirpById)
	mux.HandleFunc("GET /admin/metrics", apiConf.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiConf.handlerReset)

	server := &http.Server{Handler: mux, Addr: port}
	log.Printf("Serving on port: %s", port)
	log.Fatal(server.ListenAndServe())
}
