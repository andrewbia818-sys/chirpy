package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/andrewbia818-sys/chirpy/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	*database.Queries
	Platform string
	Secret   string
	PolkaKey string
}
type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func main() {
	godotenv.Load()
	platform := os.Getenv("PLATFORM")
	secret := os.Getenv("SECRET")
	if platform == "" {
		log.Fatal("PLATFORM environment variable is not set")
	}
	if secret == "" {
		log.Fatal("SECRET environment variable is not set")
	}

	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}

	queries := database.New(db)

	cfg := &apiConfig{
		//		Queries *database.Queries,
		Queries:  queries,
		Platform: platform,
		Secret:   secret,
		PolkaKey: os.Getenv("POLKA_KEY"),
	}

	//}

	const port = "8080"
	const filepathRoot = "./app/"

	//cfg := &apiConfig{}

	mux := http.NewServeMux()

	mux.Handle("/app/", cfg.middlewareMetricsInc(
		http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	// /metrics endpoint
	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)

	// /reset endpoint
	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)

	// readiness endpoint
	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	// validateChirp endpoint.
	//mux.HandleFunc("POST /api/validate_chirp/{rest...}", handlerValidateChirp)

	// CreateUser endpoint
	mux.HandleFunc("POST /api/users", cfg.handlerCreateUser)

	// Chirps endpoint
	mux.HandleFunc("POST /api/chirps", cfg.handlerCreateChirp)

	// GetChirps endpoint
	mux.HandleFunc("GET /api/chirps", cfg.handlerGetChirps)

	// GetChirp endpoint
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.handlerGetChirp)

	// User login endpoint
	mux.HandleFunc("POST /api/login", cfg.handlerLogin)

	// Refresh token endpoint
	mux.HandleFunc("POST /api/refresh", cfg.handlerRefresh)

	// Revoke refresh token endpoint
	mux.HandleFunc("POST /api/revoke", cfg.handlerRevoke)

	// UpdateUser endpoint
	mux.HandleFunc("PUT /api/users", cfg.handlerUpdateUser)

	// UpdateChirp endpoint
	//mux.HandleFunc("PUT /api/chirps/{id}", cfg.handlerUpdateChirp)

	// DeleteChirp endpoint DELETE /api/chirps/{chirpID}
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.handlerDeleteChirp)

	// UpgradeUser endpoint POST /api/polka/webhooks
	mux.HandleFunc("POST /api/polka/webhooks", cfg.handlerUpgradeUser)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port %s\n", port)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
