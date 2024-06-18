package main

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"log"
	"net/http"

	"github.com/redis/go-redis/v9"
)

type Response struct {
	Key string `json:"key"`
	LongUrl string `json:"long_url"`
	ShortUrl string `json:"short_url"`
}

type Request struct {
	Url string `json:"url"`
}

var ctx = context.Background()
const baseUrl = "http://localhost:8080/"

func main() {
	crc32q := crc32.MakeTable(0xD5828281)
	
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})	

	http.HandleFunc("GET /{hash}", func (w http.ResponseWriter, r *http.Request) {
		hash := r.PathValue("hash")
		
		val, err := rdb.Get(ctx, hash).Result()
		
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		
		w.Header().Set("Location", val)
		w.WriteHeader(http.StatusPermanentRedirect)
	})

	http.HandleFunc("DELETE /api/{hash}", func(w http.ResponseWriter, r *http.Request) {
		hash := r.PathValue("hash")
		_, err := rdb.Get(ctx, hash).Result()

		if err != nil {
			w.WriteHeader(http.StatusOK)
			return
		}

		err = rdb.Del(ctx, hash).Err()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("POST /api/short", func(w http.ResponseWriter, r *http.Request) {
		var request Request
		
		err := json.NewDecoder(r.Body).Decode(&request)
		
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		url := request.Url

		hash := crc32.Checksum([]byte(url), crc32q)
		hashS := fmt.Sprintf("%x", hash)

		val, _ := rdb.Get(ctx, hashS).Result()

		if val != "" {
			response := Response{
				LongUrl: url,
				ShortUrl: baseUrl + hashS,
				Key: hashS,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)

			return
		}

		err = rdb.Set(ctx, hashS, url, 0).Err()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		
		response := Response{
			LongUrl: url,
			ShortUrl: baseUrl + hashS,
			Key: hashS,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	http.Handle("GET /", http.FileServer(http.Dir("static")));

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal("Server not working.")
	}
}