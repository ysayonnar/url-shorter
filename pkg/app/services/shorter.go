package services

import (
	"fmt"
	"math/rand"
	"strings"
	"url-shorter/pkg/config"
	customErrors "url-shorter/pkg/errors"
)

func UrlShorter(initialUrl string, length int) (string, customErrors.DefaultError) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	tokenSymbols := []string{}
	for i := 0; i < length; i++ {
		randIndex := rand.Intn(len(charset))
		tokenSymbols = append(tokenSymbols, string(charset[randIndex]))
	}

	url := fmt.Sprintf("https://%v/%s", config.Domen, strings.Join(tokenSymbols, "")) //generating url
	return url, customErrors.DefaultError{}
}