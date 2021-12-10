package model

import (
	"github.com/rs/xid"
)

type Movie struct {
	ID                xid.ID  `json:"id"`
	Title             string  `json:"title"`
	Director          string  `json:"director"`
	ScreenWriter      string  `json:"screenWriter"`
	Description       string  `json:"description"`
	ReleaseYear       int     `json:"releaseYear"`
	Rating            float64 `json:"rating"`
	RatingCount       int     `json:"ratingCount"`
	Ranking           int     `json:"ranking"`
	Genre             string  `json:"genre"`
	ProductionCountry string  `json:"productionCountry"`
	BoxOffice         int     `json:"boxOffice"`
	Length            int     `json:"length"`
	PosterURL         string  `json:"posterUrl"`
}
