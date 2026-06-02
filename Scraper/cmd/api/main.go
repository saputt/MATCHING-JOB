package main

import (
	"fmt"
	"job-matching-scraper/internal/client"
	"job-matching-scraper/internal/config"
	"job-matching-scraper/internal/database"
	"job-matching-scraper/internal/scraper"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func main() {
	cfg := config.Load()

	//melakukan cek apakah state.jso sudah ada, jika belum arahkan user untuk login dahulu
	_, err := os.Stat("state.json")
	if os.IsNotExist(err) {
		err := client.LoginManual()
		if err != nil {
			log.Fatal("login manual gagal")
		}
		return
	}

	db, err := database.NewPostgresDB(cfg.DatabaseUrl)
	if err != nil {
		log.Fatal("failed to connect to database : ", err)
	}

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:8080"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	scraperRepo := scraper.NewRepository(db)
	scraperService, err := scraper.NewService(cfg.Headless, scraperRepo)
	if err != nil {
		log.Fatal("failed to init scraper service", err)
	}
	defer scraperService.Close()

	scraperHandler := scraper.NewHandler(scraperService)

	scraperHandler.RegisterRoutes(r)

	port := ":" + cfg.AppPort

	fmt.Println("service running on port " + port)

	err = http.ListenAndServe(port, r)

	if err != nil {
		log.Fatal("failed to start server", err)
	}
}
