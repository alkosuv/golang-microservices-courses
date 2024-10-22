package handler

import (
	"context"
	"encoding/json"
	"github.com/alkosuv/golang-microservices-courses/docker/storage/internal/model"
	"github.com/alkosuv/golang-microservices-courses/docker/storage/internal/service"
	"net"
	"net/http"
)

type Handler struct {
	s   *service.Service
	mux *http.ServeMux
}

func New(s *service.Service) *Handler {
	h := &Handler{s: s}

	h.mux = http.NewServeMux()
	h.mux.HandleFunc("/books", h.getBooks)
	h.mux.HandleFunc("/books/add", h.addBook)

	return h
}

func (h *Handler) ListenAndServe(host string, port string) error {
	addr := net.JoinHostPort(host, port)
	return http.ListenAndServe(addr, h.mux)
}

func (h *Handler) getBooks(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	books, err := h.s.GetBooks(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func (h *Handler) addBook(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var newBook model.Book
	err := json.NewDecoder(r.Body).Decode(&newBook)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.s.AddBook(ctx, &newBook); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
