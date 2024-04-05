package ads

import (
	"errors"
	"strings"

	"github.com/biter777/countries"
)

func (a *AdInput) Validate() error {
	if a.Title == "" {
		return errors.New("the title cannot be empty")
	}

	if a.EndAt.Before(a.StartAt) {
		return errors.New("the end date must be after the start date")
	}

	// Age range validation
	if a.Conditions.AgeStart != nil && (*a.Conditions.AgeStart < 1 || *a.Conditions.AgeStart > 100) {
		return errors.New("age start must be between 1 and 100")
	}

	if a.Conditions.AgeEnd != nil && (*a.Conditions.AgeEnd < 1 || *a.Conditions.AgeEnd > 100) {
		return errors.New("age end must be between 1 and 100")
	}

	if a.Conditions.AgeStart != nil && a.Conditions.AgeEnd != nil && *a.Conditions.AgeStart > *a.Conditions.AgeEnd {
		return errors.New("age start must not be greater than age end")
	}

	// Gender validation
	for _, g := range a.Conditions.Gender {
		if g != GenderMale && g != GenderFemale {
			return errors.New("gender must be either 'M' or 'F'")
		}
	}

	// Country code validation
	for _, country := range a.Conditions.Country {
		if !isValidCountryCode(strings.ToUpper(country)) {
			return errors.New("invalid country code: " + country)
		}
	}

	// Assuming you implement isValidPlatform similar to isValidCountryCode
	for _, platform := range a.Conditions.Platform {
		if !isValidPlatform(strings.ToLower(platform)) {
			return errors.New("invalid platform: " + platform)
		}
	}

	return nil
}

func (q *Query) Validate() error {

	if q.Offset < 0 {
		return errors.New("offset must be greater than or equal to 0")
	}

	if q.Age != -1 {
		if q.Age < 1 || q.Age > 100 {
			return errors.New("age must be between 1 and 100")
		}
	}

	if q.Gender != "" && q.Gender != "M" && q.Gender != "F" {
		return errors.New("gender must be either 'M' or 'F'")
	}

	if q.Limit < 1 || q.Limit > 100 {
		return errors.New("limit must be between 1 and 100")
	}

	if q.Country != "" && !isValidCountryCode(q.Country) {
		return errors.New("invalid country code: " + q.Country)
	}

	if q.Platform != "" && !isValidPlatform(q.Platform) {
		return errors.New("invalid platform: " + q.Platform)
	}

	return nil
}

var validPlatforms = map[string]bool{
	"android": true,
	"ios":     true,
	"web":     true,
}

func isValidPlatform(platform string) bool {
	_, exists := validPlatforms[platform]
	return exists
}

func isValidCountryCode(code string) bool {
	return countries.ByName(code) != countries.Unknown
}
