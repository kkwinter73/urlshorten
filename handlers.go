package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

// Store: ハンドラが必要とするデータ層の抽象
type Store interface {
	Create(url string) (string, error)
	Get(id string) (string, error)
}

type URLHandler struct {
	store   Store
	baseURL string // 短縮URLを組み立てるためのベース (例: "http://localhost:8080")
}

func NewURLHandler(store Store, baseURL string) *URLHandler {
	return &URLHandler{store: store, baseURL: baseURL}
}

// --- DTO ---

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	ID       string `json:"id"`
	ShortURL string `json:"short_url"`
}

type errorResponse struct {
	Error string `json:"error"`
}

// --- ヘルパー ---

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, errorResponse{Error: msg})
}

// --- ハンドラ ---

// POST /shorten
func (h *URLHandler) Shorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req shortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if req.URL == "" {
		writeError(w, http.StatusBadRequest, "url is required")
		return
	}

	id, err := h.store.Create(req.URL)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create short url")
		return
	}

	writeJSON(w, http.StatusCreated, shortenResponse{
		ID:       id,
		ShortURL: h.baseURL + "/" + id,
	})
}

// GET /{id}
func (h *URLHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// "/abc123" から "abc123" を取り出す
	id := strings.TrimPrefix(r.URL.Path, "/")
	if id == "" {
		writeError(w, http.StatusNotFound, "not found")
		return
	}

	url, err := h.store.Get(id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			writeError(w, http.StatusNotFound, "short url not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}
