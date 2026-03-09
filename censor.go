package main

import (
	"strings"
)

func censorString(s string, badwords map[string]struct{}) string {
	words := strings.Split(s, " ")
	for wIdx, word := range words {
		if _, ok := badwords[strings.ToLower(word)]; ok {
			words[wIdx] = "****"
		}
	}
	return strings.Join(words, " ")
}
