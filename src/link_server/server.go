package main

import (
	"encoding/json"
	"fmt"
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

func startServer() {

	http.HandleFunc("/{key}", func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			if key := request.PathValue("key"); key != "" {
				url, err := storage.GetRedirectUrl(key)
				if err != nil {
					http.Error(writer, err.Error(), http.StatusBadRequest)
					break
				}
				http.Redirect(writer, request, url, http.StatusMovedPermanently)
			}
		default:
			http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodPost:
			createRequest := &ShortRequest{}
			if err := json.NewDecoder(request.Body).Decode(createRequest); err != nil {
				http.Error(writer, err.Error(), http.StatusBadRequest)
			}
			if token, exp := storage.Create(createRequest.URL); token != "" {
				sr := &ShortResponse{
					URL:            fmt.Sprintf("%s/%s", servingAddress, token),
					ExpirationDate: exp,
				}
				err := json.NewEncoder(writer).Encode(sr)
				if err != nil {
					http.Error(writer, err.Error(), http.StatusBadRequest)
				}
			}
		default:
			http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	log.Fatal(http.ListenAndServe(servingAddress, nil))
}

func main() {
	fmt.Println("shorting url server")
	startServer()
}
