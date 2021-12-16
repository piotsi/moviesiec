package model

import (
	"github.com/rs/xid"
)

type User struct {
	ID                 xid.ID  `json:"id"`
	Name               string  `json:"name"`
	Nickname           string  `json:"nickname"`
	IsAdmin            bool    `json:"isAdmin"`
	Email              string  `json:"email"`
	Created            string  `json:"created"`
	Password           string  `json:"password"`
	RatingsJSON        string  `json:"ratingsJSON"`
	WantToSeeJSON      string  `json:"wantToSeeJSON"`
	RatingsAverage     float64 `json:"ratingsAverage"`
	RatingsCount       int     `json:"ratingsCount"`
	ProfilePictureURL  string  `json:"profilePictureURL"`
	ProfileDescription string  `json:"profileDescription"`
}

type Rating struct {
	MovieID  string `json:"movieID"`
	UserID   string `json:"userID"`
	Rating   int    `json:"rating"`
	Favorite bool   `json:"favorite"`
}

type Authentication struct {
	Success bool `json:"success"`
}
