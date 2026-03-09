package main

import (
	"strings"
)

func censorString(s string, badwords []string) string {
	words := strings.Split(s, " ")
	for wIdx, word := range words {
		for _, badword := range badwords {
			if strings.ToLower(word) == badword {
				words[wIdx] = "****"
				break
			}
		}
	}
	return strings.Join(words, " ")
}
