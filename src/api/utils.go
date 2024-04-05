package api

import (
	"dcard_backend/pkg/ads"
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func ParseNewAdInput(input ads.AdInput) ([]string, []string, ads.Ad, error) {
	if input.Conditions.AgeStart == nil {
		defaultStartAge := 1
		input.Conditions.AgeStart = &defaultStartAge
	}

	if input.Conditions.AgeEnd == nil {
		defaultEndAge := 100
		input.Conditions.AgeEnd = &defaultEndAge
	}

	for i, country := range input.Conditions.Country {
		input.Conditions.Country[i] = strings.ToUpper(country)
	}

	// Validate request
	if err := input.Validate(); err != nil {
		return nil, nil, ads.Ad{}, fmt.Errorf("invalid request payload: %v", err)
	}

	containsM := false
	containsF := false
	for _, gender := range input.Conditions.Gender {
		if gender == ads.GenderMale {
			containsM = true
		} else if gender == ads.GenderFemale {
			containsF = true
		}
	}

	ad := ads.Ad{
		Title:    input.Title,
		StartAt:  input.StartAt,
		EndAt:    input.EndAt,
		AgeStart: *input.Conditions.AgeStart,
		AgeEnd:   *input.Conditions.AgeEnd,
	}

	if containsM && !containsF {
		ad.GenderTarget = ads.GenderMale
	} else if !containsM && containsF {
		ad.GenderTarget = ads.GenderFemale
	}

	countries := input.Conditions.Country
	platforms := input.Conditions.Platform

	for i, countryCode := range countries {
		countries[i] = strings.ToUpper(countryCode)
	}

	for i, platformName := range platforms {
		platforms[i] = strings.ToLower(platformName)
	}

	if len(countries) == 0 {
		countries = append(countries, "any")
	}

	if len(platforms) == 0 {
		platforms = append(platforms, "any")
	}

	return countries, platforms, ad, nil

}

func ParseQuery(c *fiber.Ctx) (ads.Query, error) {
	query := ads.Query{}

	if offset, err := strconv.Atoi(c.Query("offset", "0")); err == nil {
		query.Offset = offset
	} else {
		return ads.Query{}, fmt.Errorf("invalid value for query parameter 'offset': %v", err)
	}

	if limit, err := strconv.Atoi(c.Query("limit", "5")); err == nil { // Assuming default limit is 10
		query.Limit = limit
	} else {
		return ads.Query{}, fmt.Errorf("invalid value for query parameter 'limit': %v", err)
	}

	if age, err := strconv.Atoi(c.Query("age", "-1")); err == nil {
		query.Age = age
	} else {
		return ads.Query{}, fmt.Errorf("invalid value for query parameter 'age': %v", err)
	}

	query.Gender = c.Query("gender", "")
	query.Country = c.Query("country", "")
	query.Platform = c.Query("platform", "")

	if query.Gender != "M" && query.Gender != "F" && query.Gender != "" {
		return ads.Query{}, fmt.Errorf("invalid value for query parameter gender")
	}

	if query.Platform != "" && query.Platform != "android" && query.Platform != "ios" && query.Platform != "web" {
		return ads.Query{}, fmt.Errorf("invalid value for query parameter platform")
	}

	err := query.Validate()
	if err != nil {
		return ads.Query{}, fmt.Errorf("invalid query parameters: %v", err)
	}

	return query, nil
}
