package main

import (
	"encoding/json"
	"fmt"
	cache "github.com/bulaiocht/link_cache"
	"github.com/bulaiocht/link_storage"
	"log"
	"net/http"
	"time"
)

const servingAddress = "localhost:9999"

type ShortRequest struct {
	URL string `json:"url"`
}

type ShortResponse struct {
	URL            string    `json:"url"`
	ExpirationDate time.Time `json:"expiration_date"`
}

func redirect() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			if key := request.PathValue("key"); key != "" {
				up := cache.LookUp(key)
				if up != nil {
					http.Redirect(writer, request, fmt.Sprint(up), http.StatusMovedPermanently)
					break
				}
				url, err := storage.GetRedirectUrl(key)
				if err != nil {
					http.Error(writer, err.Error(), http.StatusBadRequest)
					break
				}
				cache.Put(key, url)
				http.Redirect(writer, request, url, http.StatusMovedPermanently)
			}
		default:
			http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}

func createShortURL() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodPost:
			createRequest := &ShortRequest{}
			if err := json.NewDecoder(request.Body).Decode(createRequest); err != nil {
				http.Error(writer, err.Error(), http.StatusBadRequest)
			}
			hash := storage.CreateHash(createRequest.URL)
			value := cache.LookUp(fmt.Sprint(hash))
			if value != nil {
				err := json.NewEncoder(writer).Encode(value)
				if err != nil {
					http.Error(writer, err.Error(), http.StatusBadRequest)
				}
				break
			}
			if token, exp := storage.Create(createRequest.URL); token != "" {
				sr := &ShortResponse{
					URL:            fmt.Sprintf("%s/%s", servingAddress, token),
					ExpirationDate: exp,
				}
				cache.Put(fmt.Sprint(hash), sr)
				err := json.NewEncoder(writer).Encode(sr)
				if err != nil {
					http.Error(writer, err.Error(), http.StatusBadRequest)
				}
			}
		default:
			http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}

func startServer() {
	http.HandleFunc("/{key}", redirect())
	http.HandleFunc("/", createShortURL())
	log.Fatal(http.ListenAndServe(servingAddress, nil))
}

func main() {
	fmt.Println("shorting url server")
	startServer()
}
