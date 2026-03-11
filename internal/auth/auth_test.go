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
