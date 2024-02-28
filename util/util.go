package util

import (
	"math/rand"
	"time"
)

func ParseDate(date string) (time.Time, error) {
	return time.Parse("2006-01-02", date)
}

func IsMemberOf[T comparable](ts []T, t T) bool {
	for _, i := range ts {
		if i == t {
			return true
		}
	}
	return false
}

func GenerateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())

	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, 6)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func IsUnique(array []string) bool {
	seen := make(map[string]bool)
	for _, value := range array {
		if seen[value] {
			return false
		}
		seen[value] = true
	}
	return true
}
