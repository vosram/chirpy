package main

import "testing"

func TestCensorChirp(t *testing.T) {
	str1 := "Man, kerfuffle this thang"
	str1Answer := "Man, **** this thang"
	str2 := "Man, kerfuffle this sharbert fornax . It's unbeliable!"
	str2Answer := "Man, **** this **** **** . It's unbeliable!"

	cleanedRes1 := censorString(str1, []string{"kerfuffle"})
	cleanedRes2 := censorString(str2, []string{"kerfuffle", "sharbert", "fornax"})
	if str1Answer != cleanedRes1 {
		t.Errorf("%s IS NOT EQUAL TO %s", cleanedRes1, str1Answer)
	}
	if str2Answer != cleanedRes2 {
		t.Errorf("%s IS NOT EQUAL TO %s", cleanedRes2, str2Answer)
	}
}
