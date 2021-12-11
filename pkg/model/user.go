package model

import (
	"github.com/rs/xid"
)

type User struct {
	ID                 xid.ID  `json:"id"`
	Name               string  `json:"name"`
	Nickname           string  `json:"nickname"`
	Email              string  `json:"email"`
	Created            string  `json:"created"`
	Password           string  `json:"password"`
	RatingsJSON        string  `json:"-"`
	WantToSeeJSON      string  `json:"-"`
	RatingsAverage     float64 `json:"ratingsAverage"`
	RatingsCount       int     `json:"ratingsCount"`
	ProfilePictureURL  string  `json:"profilePictureURL"`
	ProfileDescription string  `json:"profileDescription"`
}

type UserOutput struct {
	User
	Ratings   []Rating `json:"ratings"`
	WantToSee []Movie  `json:"wantToSee"`
}

type Rating struct {
	Movie    Movie `json:"movie"`
	Rating   int   `json:"rating"`
	Favorite bool  `json:"favorite"`
}

type Authentication struct {
	Success bool `json:"success"`
}
