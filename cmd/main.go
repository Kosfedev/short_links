package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

const (
	baseURL        = "localhost:8081"
	defaultTimeout = time.Second * 5
)

// RequestBody is...
type RequestBody struct {
	Link string `json:"link"`
}

// ResponseBody is...
type ResponseBody struct {
	Link string `json:"link"`
}

var links = map[string]string{}

func createLinkHashHandler(w http.ResponseWriter, r *http.Request) {
	requestBody := &RequestBody{}
	if err := json.NewDecoder(r.Body).Decode(requestBody); err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	shortLink := getShortLink(requestBody.Link)
	links[requestBody.Link] = shortLink

	responseBody := &ResponseBody{
		Link: shortLink,
	}
	if err := json.NewEncoder(w).Encode(responseBody); err != nil {
		http.Error(w, "Failed to encode response body", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func getShortLink(baseLink string) string {
	hash := getStringHash(baseLink)

	return fmt.Sprintf("%s/%s", baseLink, hash)
}

func getStringHash(str string) string {
	h := sha1.New()
	h.Write([]byte(str))

	return hex.EncodeToString(h.Sum(nil))[:7]
}

func main() {
	r := chi.NewRouter()
	r.Post("/", createLinkHashHandler)

	server := http.Server{
		Addr:         baseURL,
		Handler:      r,
		ReadTimeout:  defaultTimeout,
		WriteTimeout: defaultTimeout,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
