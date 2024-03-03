package storage

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"log"
	"time"
)

const (
	MaxKeyLength = 6
	oneDay       = 24 * time.Hour
)

type Memory map[string]*Link

var memory = make(Memory)

type Link struct {
	Key            string
	OriginalURL    string
	CreationDate   time.Time
	ExpirationDate time.Time
}

func createNewLink(key string, url string) *Link {
	now := time.Now()
	exp := now.Add(oneDay)
	return &Link{
		Key:            key,
		OriginalURL:    url,
		CreationDate:   now,
		ExpirationDate: exp,
	}
}

func createLinkKey(url string) (string, error) {
	sum := CreateHash(url)
	return base64.StdEncoding.EncodeToString(sum[:MaxKeyLength]), nil
}

func CreateHash(input string) []byte {
	hash := sha1.New()
	hash.Write([]byte(input))
	return hash.Sum(nil)
}

func Create(url string) (string, time.Time) {
	key, err := createLinkKey(url)
	if err != nil {
		log.Fatal(err)
	}
	var link *Link
	if memory[key] == nil {
		link = createNewLink(key, url)
		memory[key] = link
	}
	return key, link.ExpirationDate
}

func GetRedirectUrl(key string) (string, error) {
	var link *Link
	if link = memory[key]; link != nil {
		if link.ExpirationDate.Before(time.Now()) {
			delete(memory, key)
			return "", errors.New("expired link")
		}
		return link.OriginalURL, nil
	}
	return "", errors.New("link not found")
}
