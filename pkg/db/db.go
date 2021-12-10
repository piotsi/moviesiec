package db

import (
	"net/http"
	"strconv"

	"github.com/piotsik/moviesiec-go/pkg/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DB struct {
	Conn *gorm.DB
}

func Init() (*DB, error) {
	db := &DB{}
	conn, err := gorm.Open(sqlite.Open("moviesiec.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	conn.AutoMigrate(&model.User{})
	conn.AutoMigrate(&model.Movie{})

	db.Conn = conn

	return db, nil
}

func Paginate(r *http.Request) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page == 0 {
			page = 1
		}

		pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
