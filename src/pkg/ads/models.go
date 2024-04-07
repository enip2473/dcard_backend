package ads

import (
	"time"

	"gorm.io/gorm"
)

type GenderTarget string

const (
	GenderMale   GenderTarget = "M"
	GenderFemale GenderTarget = "F"
)

type Condition struct {
	AgeStart *int           `json:"ageStart"`
	AgeEnd   *int           `json:"ageEnd"`
	Gender   []GenderTarget `json:"gender"`
	Country  []string       `json:"country"`
	Platform []string       `json:"platform"`
}

type Query struct {
	Offset   int `json:"offset"`
	Limit    int `json:"limit"`
	Age      int `json:"age"`
	Gender   string
	Country  string
	Platform string
}

type Response struct {
	Title string    `json:"title"`
	EndAt time.Time `json:"endAt"`
}
type AdInput struct {
	Title      string    `json:"title"`
	StartAt    time.Time `json:"startAt"`
	EndAt      time.Time `json:"endAt"`
	Conditions Condition `json:"conditions"`
}

type Ad struct {
	gorm.Model
	Title        string       `json:"title"`
	StartAt      time.Time    `json:"startAt"`
	EndAt        time.Time    `json:"endAt" gorm:"index"`
	AgeStart     int          `json:"startAge"`
	AgeEnd       int          `json:"endAge"`
	GenderTarget GenderTarget `json:"genderTarget"`
	Countries    []Country    `gorm:"many2many:ad_countries;"`
	Platforms    []Platform   `gorm:"many2many:ad_platforms;"`
}

type Country struct {
	Code string `gorm:"primaryKey" json:"code"`
}

type Platform struct {
	Name string `gorm:"primaryKey" json:"name"`
}
