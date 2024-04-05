package api

import (
	"bytes"
	"dcard_backend/pkg/ads"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

type mockRepo struct{}

func (m mockRepo) CreateAd(ad *ads.Ad) error {
	return nil
}

func (m mockRepo) ListActiveAds(query ads.Query) ([]ads.Response, error) {
	return []ads.Response{{Title: "Mock Ad"}}, nil
}

func (m mockRepo) FirstOrCreateCountry(code string) (ads.Country, error) {
	return ads.Country{Code: code}, nil
}

func (m mockRepo) FirstOrCreatePlatform(name string) (ads.Platform, error) {
	return ads.Platform{Name: name}, nil
}

func TestCreateAd(t *testing.T) {
	// Setup Fiber app
	app := fiber.New()
	app.Post("/api/v1/ad", CreateAd(mockRepo{})) // Use a mock repo implementation

	// Test case: Valid ad creation
	adInput := ads.AdInput{
		Title:   "Test Ad",
		StartAt: time.Now(),
		EndAt:   time.Now().Add(24 * time.Hour),
		Conditions: ads.Condition{
			AgeStart: func(i int) *int { return &i }(18),
			AgeEnd:   func(i int) *int { return &i }(25),
			Gender:   []ads.GenderTarget{ads.GenderMale},
			Country:  []string{"US", "CA"},
			Platform: []string{"android", "ios"},
		},
	}
	body, _ := json.Marshal(adInput)

	req := httptest.NewRequest("POST", "/api/v1/ad", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
}

func TestListAds(t *testing.T) {
	app := fiber.New()
	app.Get("/api/v1/ad", ListAds(mockRepo{})) // Use a mock repo implementation

	req := httptest.NewRequest("GET", "/api/v1/ad", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var result struct {
		Items []ads.Ad
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}
	assert.NotEmpty(t, result.Items)
}

func TestCreateAdInvalidPayload(t *testing.T) {
	app := fiber.New()
	app.Post("/api/v1/ad", CreateAd(mockRepo{})) // Assuming CreateAd accepts a mockRepo

	// Test case: Invalid JSON payload
	invalidJSON := "{invalid}"
	req := httptest.NewRequest("POST", "/api/v1/ad", strings.NewReader(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode, "Should return 400 for bad JSON")

	// Test case: Invalid ad data (e.g., end date before start date)
	adInput := ads.AdInput{
		Title:   "Test Ad",
		StartAt: time.Now().Add(24 * time.Hour), // Start date
		EndAt:   time.Now(),                     // End date before start date
	}
	body, _ := json.Marshal(adInput)

	req = httptest.NewRequest("POST", "/api/v1/ad", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ = app.Test(req)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode, "Should return 400 for invalid ad data")
}

func TestListAdsInvalidQuery(t *testing.T) {
	app := fiber.New()
	app.Get("/api/v1/ad", ListAds(mockRepo{})) // Assuming ListAds accepts a mockRepo

	// Test case: Invalid 'limit' parameter
	req := httptest.NewRequest("GET", "/api/v1/ad?limit=-1", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode, "Should return 400 for invalid 'limit' parameter")

	// Test case: Invalid 'age' parameter
	req = httptest.NewRequest("GET", "/api/v1/ad?age=invalid", nil)
	resp, _ = app.Test(req)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode, "Should return 400 for invalid 'age' parameter")

	// Test case: Unsupported 'platform' parameter
	req = httptest.NewRequest("GET", "/api/v1/ad?platform=unsupported_platform", nil)
	resp, _ = app.Test(req)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode, "Should return 400 for unsupported 'platform'")
}
