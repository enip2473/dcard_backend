package api

import (
	"dcard_backend/pkg/ads"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestParseNewAdInput(t *testing.T) {
	now := time.Now()
	oneDayLater := now.Add(24 * time.Hour)
	tests := []struct {
		name              string
		input             ads.AdInput
		expectErr         bool
		expectedAd        ads.Ad
		expectedErr       string
		expectedCountries []string
		expectedPlatforms []string
	}{
		{
			name: "valid input",
			input: ads.AdInput{
				Title:   "Test Ad",
				StartAt: now,
				EndAt:   oneDayLater,
				Conditions: ads.Condition{
					AgeStart: func(i int) *int { return &i }(18),
					AgeEnd:   func(i int) *int { return &i }(25),
					Gender:   []ads.GenderTarget{ads.GenderMale},
					Country:  []string{"US", "CA"},
					Platform: []string{"android", "ios"},
				},
			},
			expectErr: false,
			expectedAd: ads.Ad{
				Title:        "Test Ad",
				StartAt:      now,
				EndAt:        oneDayLater,
				AgeStart:     18,
				AgeEnd:       25,
				GenderTarget: ads.GenderMale,
			},
			expectedCountries: []string{"US", "CA"},
			expectedPlatforms: []string{"android", "ios"},
		},
		{
			name: "invalid age range",
			input: ads.AdInput{
				Title:   "Invalid Age Range",
				StartAt: time.Now(),
				EndAt:   time.Now().Add(24 * time.Hour),
				Conditions: ads.Condition{
					AgeStart: func(i int) *int { return &i }(25),
					AgeEnd:   func(i int) *int { return &i }(18), // Invalid: AgeStart is greater than AgeEnd
				},
			},
			expectErr:   true,
			expectedErr: "age start must not be greater than age end",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			countries, platforms, ad, err := ParseNewAdInput(tc.input)

			if tc.expectErr {
				assert.Error(t, err)
				if tc.expectedErr != "" {
					assert.Contains(t, err.Error(), tc.expectedErr)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedAd, ad)
				assert.ElementsMatch(t, tc.expectedCountries, countries)
				assert.ElementsMatch(t, tc.expectedPlatforms, platforms)
			}
		})
	}
}

// Mock handler to use ParseQuery and return the parsed query as JSON for easy inspection in tests
func mockParseQueryHandler(c *fiber.Ctx) error {
	query, err := ParseQuery(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(query)
}

func TestParseQuery(t *testing.T) {
	// Setup Fiber app
	app := fiber.New()
	app.Get("/test-parse-query", mockParseQueryHandler)

	tests := []struct {
		name           string
		queryString    string
		expectStatus   int
		expectedResult ads.Query
		expectedErr    string
	}{
		{
			name:         "valid query parameters",
			queryString:  "offset=1&limit=10&age=25&gender=M&country=US&platform=android",
			expectStatus: fiber.StatusOK,
			expectedResult: ads.Query{
				Offset:   1,
				Limit:    10,
				Age:      25,
				Gender:   "M",
				Country:  "US",
				Platform: "android",
			},
		},
		{
			name:         "invalid gender",
			queryString:  "gender=X",
			expectStatus: fiber.StatusBadRequest,
			expectedErr:  "invalid value for query parameter gender",
		},
		{
			name:         "invalid platform",
			queryString:  "platform=console",
			expectStatus: fiber.StatusBadRequest,
			expectedErr:  "invalid value for query parameter platform",
		},
		{
			name:         "valid query with default values",
			queryString:  "",
			expectStatus: fiber.StatusOK,
			expectedResult: ads.Query{
				Offset:   0,
				Limit:    5,
				Age:      -1,
				Gender:   "",
				Country:  "",
				Platform: "",
			},
		},
		{
			name:         "valid age parameter",
			queryString:  "age=30",
			expectStatus: fiber.StatusOK,
			expectedResult: ads.Query{
				Offset:   0,
				Limit:    5,
				Age:      30,
				Gender:   "",
				Country:  "",
				Platform: "",
			},
		},
		{
			name:         "invalid age parameter",
			queryString:  "age=abc",
			expectStatus: fiber.StatusBadRequest,
			expectedErr:  "invalid value for query parameter 'age'",
		},
		{
			name:         "invalid offset",
			queryString:  "offset=-1",
			expectStatus: fiber.StatusBadRequest,
			expectedErr:  "invalid query parameters: offset must be greater than or equal to 0",
		},
		{
			name:         "invalid limit",
			queryString:  "limit=101",
			expectStatus: fiber.StatusBadRequest,
			expectedErr:  "invalid query parameters: limit must be between 1 and 100",
		},
		{
			name:         "valid query with max limit",
			queryString:  "limit=100",
			expectStatus: fiber.StatusOK,
			expectedResult: ads.Query{
				Offset:   0,
				Limit:    100,
				Age:      -1,
				Gender:   "",
				Country:  "",
				Platform: "",
			},
		},
		{
			name:         "valid country parameter",
			queryString:  "country=JP",
			expectStatus: fiber.StatusOK,
			expectedResult: ads.Query{
				Offset:   0,
				Limit:    5,
				Age:      -1,
				Gender:   "",
				Country:  "JP",
				Platform: "",
			},
		},
		{
			name:         "valid platform parameter",
			queryString:  "platform=ios",
			expectStatus: fiber.StatusOK,
			expectedResult: ads.Query{
				Offset:   0,
				Limit:    5,
				Age:      -1,
				Gender:   "",
				Country:  "",
				Platform: "ios",
			},
		},
		{
			name:         "multiple valid parameters",
			queryString:  "age=45&gender=F&country=US&platform=web",
			expectStatus: fiber.StatusOK,
			expectedResult: ads.Query{
				Offset:   0,
				Limit:    5,
				Age:      45,
				Gender:   "F",
				Country:  "US",
				Platform: "web",
			},
		},
		{
			name:         "mix of valid and invalid parameters",
			queryString:  "age=abc&gender=X&country=ZZ&platform=gamecube",
			expectStatus: fiber.StatusBadRequest,
			expectedErr:  "invalid value for query parameter 'age'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", fmt.Sprintf("/test-parse-query?%s", tt.queryString), nil)
			resp, _ := app.Test(req, -1)

			assert.Equal(t, tt.expectStatus, resp.StatusCode)

			if tt.expectStatus == fiber.StatusOK {
				var result ads.Query
				json.NewDecoder(resp.Body).Decode(&result)
				assert.Equal(t, tt.expectedResult, result)
			} else {
				var errorResult map[string]string
				json.NewDecoder(resp.Body).Decode(&errorResult)
				assert.Contains(t, errorResult["error"], tt.expectedErr)
			}
		})
	}
}
