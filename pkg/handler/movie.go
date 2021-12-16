package handler

import (
	"encoding/json"
	"net/http"
	"strings"

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
	MovieUpdate(w http.ResponseWriter, r *http.Request)
	MovieRate(w http.ResponseWriter, r *http.Request)
}

// MovieGetAll queries all or one movie by it's ID,
// in query parameters paginate with 'page_size' and 'page', order by 'order_by', order asc or desc with 'order_in'
func (h *Handler) MovieGetAll(w http.ResponseWriter, r *http.Request) {
	orderBy := r.URL.Query().Get("order_by")
	orderIn := r.URL.Query().Get("order_in")
	var order string
	if orderBy != "" && orderIn != "" {
		order = strings.Join([]string{orderBy, orderIn}, " ")
	}

	movies := []model.Movie{}
	result := h.DB.Conn.Scopes(db.Paginate(r)).Order(order).Find(&movies)
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

// MovieGetByID queries one movie by it's ID
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

// MovieAdd adds new movie and accepts Movie struct
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

// MovieDelete deletes movie with given ID
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

// MovieRate changes movie's rating and ratingCount column
// accepts Rating struct
func (h *Handler) MovieRate(w http.ResponseWriter, r *http.Request) {
	var rating model.Rating
	err := json.NewDecoder(r.Body).Decode(&rating)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	movie := model.Movie{}
	movieResult := h.DB.Conn.First(&movie, "id = ?", rating.MovieID)
	if movieResult.Error == gorm.ErrRecordNotFound {
		http.Error(w, movieResult.Error.Error(), http.StatusNotFound)
		return
	}

	user := model.User{}
	userResult := h.DB.Conn.First(&user, "id = ?", rating.UserID)
	if userResult.Error == gorm.ErrRecordNotFound {
		http.Error(w, userResult.Error.Error(), http.StatusNotFound)
		return
	}

	userRatings := []model.Rating{}
	if user.RatingsJSON != "" {
		err = json.Unmarshal([]byte(user.RatingsJSON), &userRatings)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if rating.Rating < 0 || rating.Rating > 10 {
		http.Error(w, "invalid rating, must be between 0 and 10", http.StatusForbidden)
		return
	}

	ratingCount := 0
	ratingRating := rating.Rating
	for i, ur := range userRatings {
		if ur.MovieID == rating.MovieID {
			if ratingRating == 0 {
				userRatings = append(userRatings[:i], userRatings[i+1:]...)
				ratingCount = -1
				ratingRating = -ur.Rating
			} else {
				http.Error(w, "OK", http.StatusOK)
				return
			}
		}
	}
	if ratingRating > 0 {
		ratingCount = 1
		userRatings = append(userRatings, rating)
	}

	movie.Rating = ((movie.Rating * float64(movie.RatingCount)) + float64(ratingRating)) / (float64(movie.RatingCount + ratingCount))
	movie.RatingCount += ratingCount

	userRatingsString, err := json.Marshal(userRatings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.RatingsJSON = string(userRatingsString)

	err = h.DB.Conn.Transaction(func(tx *gorm.DB) error {
		movieUpdateTX := h.DB.Conn.Model(&model.Movie{}).Select("*").Where("id = ?", rating.MovieID).UpdateColumns(movie)
		if movieUpdateTX.Error != nil {
			return movieUpdateTX.Error
		}

		userUpdateTX := h.DB.Conn.Model(&model.User{}).Select("*").Where("id = ?", user.ID).UpdateColumns(user)
		if userUpdateTX.Error != nil {
			return userUpdateTX.Error
		}

		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ratingJSON, err := json.Marshal(rating)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(ratingJSON)
}

// MovieUpdate updates movie with given id, accepts Movie struct
func (h *Handler) MovieUpdate(w http.ResponseWriter, r *http.Request) {
	var movie model.Movie
	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := chi.URLParam(r, "id")

	h.DB.Conn.Model(&model.Movie{}).Select("*").Omit("id").Where("id = ?", id).UpdateColumns(movie)

	movieJSON, err := json.Marshal(movie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(movieJSON)
}
