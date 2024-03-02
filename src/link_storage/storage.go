package storage

import (
	"bytes"
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
	byteSum := sha1.Sum([]byte(url))
	buf := new(bytes.Buffer)
	for i := 0; i < len(byteSum); i++ {
		err := buf.WriteByte(byteSum[i])
		if err != nil {
			return "", err
		}
	}
	encodedBytes := buf.Bytes()
	return base64.StdEncoding.EncodeToString(encodedBytes[:MaxKeyLength]), nil
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
