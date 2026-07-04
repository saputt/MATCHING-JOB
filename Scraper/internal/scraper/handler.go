package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"job-matching-scraper/internal/httpx"
	"job-matching-scraper/internal/model"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/scrape", func(r chi.Router) {
		r.Post("/", h.Scrape)
		r.Patch("/missing", h.ScrapeMissingData)
	})
}

func (h *Handler) Scrape(w http.ResponseWriter, r *http.Request) {
	var req model.ScrapeRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	targetPerKeyword := req.Limit
	if targetPerKeyword <= 0 {
		targetPerKeyword = 30
	}

	response, err := h.service.ScrapeAndSave(r.Context(), targetPerKeyword)
	fmt.Println(err)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	httpx.WriteSuccess(w, http.StatusOK, "Scraping success", response)
}

func (h *Handler) ScrapeMissingData(w http.ResponseWriter, r *http.Request) {
	go func() {
		ctx := context.Background()
		_, err := h.service.ScrapeEmpty(ctx)
		if err != nil {
			log.Printf("[BACKGROUND ERROR] Failed to patch data: %v", err)
		}
	}()

	httpx.WriteSuccess(w, http.StatusOK, "Patching data with worker pool completed successfully", nil)
}
