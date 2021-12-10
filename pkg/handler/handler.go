package handler

import (
	"fmt"
	"net/http"

	"github.com/piotsik/moviesiec-go/pkg/db"
)

type Handler struct {
	DB *db.DB
}

func Init() (*Handler, error) {
	db, err := db.Init()
	if err != nil {
		return nil, err
	}

	handler := Handler{
		DB: db,
	}

	return &handler, nil
}

func (h *Handler) Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello")
}
