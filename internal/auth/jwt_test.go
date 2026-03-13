package auth

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWTs(t *testing.T) {
	tokenSecret1 := "CoolPassword"
	tokenSecret2 := "AwesomePassword"

	tests := []struct {
		name         string
		tokenSecret  string
		userID       uuid.UUID
		expiresIn    time.Duration
		wantErr      bool
		testWaitTime time.Duration
	}{
		{
			name:         "valid jwt",
			tokenSecret:  tokenSecret1,
			userID:       uuid.New(),
			expiresIn:    time.Minute * 5,
			wantErr:      false,
			testWaitTime: time.Millisecond * 50,
		},
		{
			name:         "second valid jwt",
			tokenSecret:  tokenSecret2,
			userID:       uuid.New(),
			expiresIn:    time.Minute * 5,
			wantErr:      false,
			testWaitTime: time.Millisecond * 50,
		},
		{
			name:         "expired jwt",
			tokenSecret:  tokenSecret1,
			userID:       uuid.New(),
			expiresIn:    time.Millisecond * 50,
			wantErr:      true,
			testWaitTime: time.Second * 1,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokenStr, err := MakeJWT(test.userID, test.tokenSecret, test.expiresIn)
			if err != nil {
				t.Errorf("%s failed creating token - error: %s", test.name, err.Error())
				return
			}
			time.Sleep(test.testWaitTime)
			userId, err := ValidateJWT(tokenStr, test.tokenSecret)
			if test.wantErr == true {
				if err == nil {
					t.Errorf("Expected err in creating token and none found")
					return
				}
			} else {
				if err != nil {
					fmt.Println("Error: ", err)
					t.Errorf("Expected no error and got an error")
					return
				}
			}
			if !test.wantErr && userId != test.userID {
				t.Errorf("UUID returned from validate is incorrect")
				return
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	validToken, _ := MakeJWT(userID, "secret", time.Hour)

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Valid token",
			tokenString: validToken,
			tokenSecret: "secret",
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name:        "Invalid token",
			tokenString: "invalid.token.string",
			tokenSecret: "secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Wrong secret",
			tokenString: validToken,
			tokenSecret: "wrong_secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name               string
		containsAuthHeader bool
		initialBearerToken string
		expectedToken      string
		wantErr            bool
	}{
		{
			name:               "Valid token",
			containsAuthHeader: true,
			initialBearerToken: "Bearer imagineatokenstringhere",
			expectedToken:      "imagineatokenstringhere",
			wantErr:            false,
		},
		{
			name:               "No auth header",
			containsAuthHeader: false,
			initialBearerToken: "",
			expectedToken:      "",
			wantErr:            true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			header := http.Header{}
			if test.containsAuthHeader {
				header.Add("Authorization", test.initialBearerToken)
			}

			tokenStr, err := GetBearerToken(header)
			if err != nil && !test.wantErr {
				t.Errorf("unwanted error: %s", err.Error())
			}
			if tokenStr != test.expectedToken {
				t.Errorf("Token %s is not the expected: %s", tokenStr, test.expectedToken)
			}
			if err == nil && test.wantErr {
				t.Errorf("expected error and got none")
			}

		})
	}
}
