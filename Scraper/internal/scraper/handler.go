package scraper

import (
	"encoding/json"
	"fmt"
	"job-matching-scraper/internal/httpx"
	"job-matching-scraper/internal/model"
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
	})
}

func (h *Handler) Scrape(w http.ResponseWriter, r *http.Request) {
	var req model.ScrapeRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.UserID == "" {
		httpx.WriteError(w, http.StatusBadRequest, "userid is empty")
		return
	}

	targetPerKeyword := req.Limit
	if targetPerKeyword <= 0 {
		targetPerKeyword = 30
	}

	response, err := h.service.ScrapeAndSave(r.Context(), req.UserID, targetPerKeyword)
	fmt.Println(err)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	httpx.WriteSuccess(w, http.StatusOK, "Scraping success", response)
}
