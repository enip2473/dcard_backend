package ads

import (
	"errors"
	"time"
)

// Ad represents an advertisement with a title, start and end dates, and conditions.
type Ad struct {
	ID      uint      `gorm:"primaryKey" json:"id"`
	Title   string    `json:"title"`
	StartAt time.Time `json:"startAt"`
	EndAt   time.Time `json:"endAt"`
	// Conditions map[string]interface{} `json:"conditions"`
}

// Validate checks if the Ad instance has valid data.
func (a *Ad) Validate() error {
	// Ensure the title is not empty
	if a.Title == "" {
		return errors.New("the title cannot be empty")
	}

	// Ensure the end date is after the start date
	if a.EndAt.Before(a.StartAt) {
		return errors.New("the end date must be after the start date")
	}

	// Optionally, add more validations for the Conditions field if needed

	return nil
}
