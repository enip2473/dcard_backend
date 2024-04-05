package ads

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAdInputValidate(t *testing.T) {
	// Assuming `AdInput` and related types are defined in this package
	age1, age2 := 18, 30
	ptrAge1, ptrAge2 := &age1, &age2
	tests := []struct {
		name    string
		adInput AdInput
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid input with all fields",
			adInput: AdInput{
				Title:   "Complete Ad",
				StartAt: time.Now(),
				EndAt:   time.Now().Add(24 * time.Hour), // Valid: EndAt is after StartAt
				Conditions: Condition{
					AgeStart: ptrAge1,
					AgeEnd:   ptrAge2,
					Gender:   []GenderTarget{GenderMale, GenderFemale},
					Country:  []string{"US", "GB"},
					Platform: []string{"android", "ios"},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid gender in conditions",
			adInput: AdInput{
				Title:   "Ad With Invalid Gender",
				StartAt: time.Now(),
				EndAt:   time.Now().Add(24 * time.Hour),
				Conditions: Condition{
					Gender: []GenderTarget{"X"}, // Invalid gender
				},
			},
			wantErr: true,
			errMsg:  "gender must be either 'M' or 'F'",
		},
		{
			name: "invalid country code in conditions",
			adInput: AdInput{
				Title:   "Ad With Invalid Country Code",
				StartAt: time.Now(),
				EndAt:   time.Now().Add(24 * time.Hour),
				Conditions: Condition{
					Country: []string{"XX"}, // Assuming "XX" is not a valid country code
				},
			},
			wantErr: true,
			errMsg:  "invalid country code: XX",
		},
		{
			name: "invalid platform in conditions",
			adInput: AdInput{
				Title:   "Ad With Invalid Platform",
				StartAt: time.Now(),
				EndAt:   time.Now().Add(24 * time.Hour),
				Conditions: Condition{
					Platform: []string{"gameboy"}, // Assuming "gameboy" is not a valid platform
				},
			},
			wantErr: true,
			errMsg:  "invalid platform: gameboy",
		},
		{
			name: "age start greater than age end",
			adInput: AdInput{
				Title:   "Invalid Age Range",
				StartAt: time.Now(),
				EndAt:   time.Now().Add(24 * time.Hour),
				Conditions: Condition{
					AgeStart: ptrAge2,
					AgeEnd:   ptrAge1, // Invalid: AgeStart is greater than AgeEnd
				},
			},
			wantErr: true,
			errMsg:  "age start must not be greater than age end",
		},
		{
			name: "missing start and end dates",
			adInput: AdInput{
				Title: "Missing Dates",
				// Missing StartAt and EndAt, which is valid since validation does not require them
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.adInput.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestQueryValidate(t *testing.T) {
	tests := []struct {
		name    string
		query   Query
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid query with minimum age",
			query: Query{
				Age:   1,
				Limit: 5,
			},
			wantErr: false,
		},
		{
			name: "valid query with maximum age",
			query: Query{
				Age:   100,
				Limit: 5,
			},
			wantErr: false,
		},
		{
			name: "invalid query with age below minimum",
			query: Query{
				Age:   0,
				Limit: 5,
			},
			wantErr: true,
			errMsg:  "age must be between 1 and 100",
		},
		{
			name: "invalid query with age above maximum",
			query: Query{
				Age:   101,
				Limit: 5,
			},
			wantErr: true,
			errMsg:  "age must be between 1 and 100",
		},
		{
			name: "invalid gender",
			query: Query{
				Gender: "X",
				Limit:  5,
				Age:    -1,
			},
			wantErr: true,
			errMsg:  "gender must be either 'M' or 'F'",
		},
		{
			name: "valid gender male",
			query: Query{
				Gender: "M",
				Limit:  5,
				Age:    -1,
			},
			wantErr: false,
		},
		{
			name: "valid gender female",
			query: Query{
				Gender: "F",
				Limit:  5,
				Age:    -1,
			},
			wantErr: false,
		},
		{
			name: "valid offset",
			query: Query{
				Offset: 10,
				Limit:  5,
				Age:    -1,
			},
			wantErr: false,
		},
		{
			name: "invalid negative offset",
			query: Query{
				Offset: -1,
				Limit:  5,
				Age:    -1,
			},
			wantErr: true,
			errMsg:  "offset must be greater than or equal to 0",
		},
		{
			name: "valid limit",
			query: Query{
				Limit: 50,
				Age:   -1,
			},
			wantErr: false,
		},
		{
			name: "invalid limit above maximum",
			query: Query{
				Limit: 101,
				Age:   -1,
			},
			wantErr: true,
			errMsg:  "limit must be between 1 and 100",
		},
		{
			name: "invalid limit below minimum",
			query: Query{
				Limit: 0,
				Age:   -1,
			},
			wantErr: true,
			errMsg:  "limit must be between 1 and 100",
		},
		{
			name: "invalid country code",
			query: Query{
				Country: "INVALID",
				Limit:   5,
				Age:     -1,
			},
			wantErr: true,
			errMsg:  "invalid country code: INVALID",
		},
		{
			name: "invalid platform",
			query: Query{
				Platform: "unknown",
				Limit:    5,
				Age:      -1,
			},
			wantErr: true,
			errMsg:  "invalid platform: unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.query.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
