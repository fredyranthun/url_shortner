package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/fredyranthun/url-shortner/config"
	"github.com/fredyranthun/url-shortner/hash"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	db *redis.Client
	conf *config.Config
	ctx context.Context
}

type Response struct {
	Key string `json:"key"`
	LongUrl string `json:"long_url"`
	ShortUrl string `json:"short_url"`
}

type Request struct {
	Url string `json:"url"`
}

func NewServer(conf *config.Config, rdb *redis.Client, ctx context.Context) *Server {
	return &Server{
		conf: conf,
		db: rdb,
		ctx: ctx,
	}
}

func (server *Server) Start() {

	http.HandleFunc("GET /{hash}", func (w http.ResponseWriter, r *http.Request) {
		hash := r.PathValue("hash")
		val, err := server.db.Get(server.ctx, hash).Result()
		
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		
		w.Header().Set("Location", val)
		w.WriteHeader(http.StatusPermanentRedirect)
	})

	http.HandleFunc("DELETE /api/{hash}", func(w http.ResponseWriter, r *http.Request) {
		hash := r.PathValue("hash")
		_, err := server.db.Get(server.ctx, hash).Result()

		if err != nil {
			w.WriteHeader(http.StatusOK)
			return
		}

		err = server.db.Del(server.ctx, hash).Err()

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

		hash := hash.Crc32Hash(url)

		val, _ := server.db.Get(server.ctx, hash).Result()

		if val != "" {
			response := Response{
				LongUrl: url,
				ShortUrl: server.conf.BaseUrl + hash,
				Key: hash,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)

			return
		}

		err = server.db.Set(server.ctx, hash, url, 0).Err()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		
		response := Response{
			LongUrl: url,
			ShortUrl: server.conf.BaseUrl + hash,
			Key: hash,
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