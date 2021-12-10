package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/piotsik/moviesiec-go/pkg/db"
	"github.com/piotsik/moviesiec-go/pkg/model"
	"github.com/rs/xid"
	"gorm.io/gorm"
)

type Movie interface {
	MovieGet(w http.ResponseWriter, r *http.Request)
	MovieGetByUUID(w http.ResponseWriter, r *http.Request)
	MovieAdd(w http.ResponseWriter, r *http.Request)
	MovieDelete(w http.ResponseWriter, r *http.Request)
}

func (h *Handler) MovieGetAll(w http.ResponseWriter, r *http.Request) {
	movies := []model.Movie{}
	result := h.DB.Conn.Scopes(db.Paginate(r)).Find(&movies)
	if result.Error == gorm.ErrRecordNotFound {
		http.Error(w, result.Error.Error(), http.StatusNotFound)
		return
	}

	moviesJSON, err := json.Marshal(movies)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(moviesJSON)
}

func (h *Handler) MovieGetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	movie := model.Movie{}
	result := h.DB.Conn.First(&movie, "id = ?", id)
	if result.Error == gorm.ErrRecordNotFound {
		http.Error(w, result.Error.Error(), http.StatusNotFound)
		return
	}

	movieJSON, err := json.Marshal(movie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(movieJSON)
}

func (h *Handler) MovieAdd(w http.ResponseWriter, r *http.Request) {
	var movie model.Movie
	movie.ID = xid.New()
	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.DB.Conn.Create(movie)

	movieJSON, err := json.Marshal(movie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(movieJSON)
}

func (h *Handler) MovieDelete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	movie := model.Movie{}
	result := h.DB.Conn.First(&movie, "id = ?", id)
	if result.Error == gorm.ErrRecordNotFound {
		http.Error(w, result.Error.Error(), http.StatusNotFound)
		return
	}

	h.DB.Conn.Delete(&movie)
}
