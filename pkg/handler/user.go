package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/piotsik/moviesiec-go/pkg/db"
	"github.com/piotsik/moviesiec-go/pkg/hash"
	"github.com/piotsik/moviesiec-go/pkg/model"
	"github.com/rs/xid"
	"gorm.io/gorm"
)

type User interface {
	UserAdd(w http.ResponseWriter, r *http.Request)
	UserGetAll(w http.ResponseWriter, r *http.Request)
	UserDelete(w http.ResponseWriter, r *http.Request)
	UserGetByID(w http.ResponseWriter, r *http.Request)
	UserAuthenticate(w http.ResponseWriter, r *http.Request)
}

func (h *Handler) UserGetAll(w http.ResponseWriter, r *http.Request) {
	orderBy := r.URL.Query().Get("order_by")
	orderIn := r.URL.Query().Get("order_in")
	var order string
	if orderBy != "" && orderIn != "" {
		order = strings.Join([]string{orderBy, orderIn}, " ")
	}

	users := []model.User{}
	result := h.DB.Conn.Scopes(db.Paginate(r)).Order(order).Find(&users)
	if result.Error == gorm.ErrRecordNotFound {
		http.Error(w, result.Error.Error(), http.StatusNotFound)
		return
	}

	usersJSON, err := json.Marshal(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(usersJSON)
}

func (h *Handler) UserGetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	user := model.User{}
	result := h.DB.Conn.First(&user, "id = ?", id)
	if result.Error == gorm.ErrRecordNotFound {
		http.Error(w, result.Error.Error(), http.StatusNotFound)
		return
	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(userJSON)
}

func (h *Handler) UserAdd(w http.ResponseWriter, r *http.Request) {
	var user model.User
	user.ID = xid.New()
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.Password, err = hash.Password(user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := h.DB.Conn.First(&model.User{}, "nickname = ?", user.Nickname)
	if result.Error != gorm.ErrRecordNotFound {
		http.Error(w, "user with this nickname alreade exists", http.StatusConflict)
		return
	}

	h.DB.Conn.Create(user)

	userJSON, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(userJSON)
}

func (h *Handler) UserDelete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	user := model.User{}
	result := h.DB.Conn.First(&user, "id = ?", id)
	if result.Error == gorm.ErrRecordNotFound {
		http.Error(w, result.Error.Error(), http.StatusNotFound)
		return
	}

	h.DB.Conn.Delete(&user)
}

func (h *Handler) UserUpdate(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.Password, err = hash.Password(user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id := chi.URLParam(r, "id")

	h.DB.Conn.Model(user).Select("*").Omit("id").Where("id = ?", id).UpdateColumns(user)

	userJSON, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(userJSON)
}

func (h *Handler) UserAuthenticate(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userCompare := model.User{}
	result := h.DB.Conn.First(&userCompare, "nickname = ?", user.Nickname)
	if result.Error == gorm.ErrRecordNotFound {
		http.Error(w, result.Error.Error(), http.StatusNotFound)
		return
	}

	authenticated := model.Authentication{Success: hash.CheckPassword(user.Password, userCompare.Password)}
	authenticatedJSON, err := json.Marshal(authenticated)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(authenticatedJSON)
}
