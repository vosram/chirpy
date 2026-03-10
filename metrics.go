package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	count := cfg.fileserverHits.Load()
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>
`, count)))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Reset is only allowed in dev enviroment."))
		return
	}

	cfg.fileserverHits.Store(0)
	err := cfg.db.DeleteAllUsers(r.Context())
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to reset the database: " + err.Error()))
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0 and database reset to initial state."))
}
