package auth

import (
	"testing"
)

func TestHashingCorrectPassword(t *testing.T) {
	password := "mycoolpassword"
	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("This password hash should have passed: %s", err.Error())
	}

	match, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Errorf("This password match should have successed with a match")
	}
	if match == false {
		t.Errorf("The password should have matched")
	}
}

func TestHashingIncorrectPassword(t *testing.T) {
	password := "mycoolpassword"
	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("This password hash should have passed: %s", err.Error())
	}

	match, err := CheckPasswordHash("myawesomepassword", hash)
	if err != nil {
		t.Errorf("This password match should have successed with a failed match")
	}
	if match == true {
		t.Errorf("The password should have not matched")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name          string
		password      string
		hash          string
		wantErr       bool
		matchPassword bool
	}{
		{
			name:          "Correct password",
			password:      password1,
			hash:          hash1,
			wantErr:       false,
			matchPassword: true,
		},
		{
			name:          "Incorrect password",
			password:      "wrongPassword",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Password doesn't match different hash",
			password:      password1,
			hash:          hash2,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Empty password",
			password:      "",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Invalid hash",
			password:      password1,
			hash:          "invalidhash",
			wantErr:       true,
			matchPassword: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && match != tt.matchPassword {
				t.Errorf("CheckPasswordHash() expects %v, got %v", tt.matchPassword, match)
			}
		})
	}
}
