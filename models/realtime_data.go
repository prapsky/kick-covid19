package models

import (
	"time"
)

type Request struct {
	Snapshots   []Snapshots `json:"snapshots"`
	CountryName string      `json:"name"`
	CountryCode string      `json:"code"`
}

type Snapshots struct {
	Cases       int       `json:"cases"`
	TodayCases  int       `json:"todayCases"`
	Deaths      int       `json:"deaths"`
	TodayDeaths int       `json:"todayDeaths"`
	Recovered   int       `json:"recovered"`
	Active      int       `json:"active"`
	Critical    int       `json:"critical"`
	Timestamp   time.Time `json:"timestamp"`
}

type Response struct {
	Timestamp   string `json:"timestamp"`
	CountryName string `json:"country"`
	Cases       int    `json:"cases"`
	Deaths      int    `json:"deaths"`
	Recovered   int    `json:"recovered"`
}
