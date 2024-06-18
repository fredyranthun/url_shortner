package main

import (
	"github.com/fredyranthun/url-shortner/config"
	"github.com/fredyranthun/url-shortner/database"
	"github.com/fredyranthun/url-shortner/server"
)

type Response struct {
	Key string `json:"key"`
	LongUrl string `json:"long_url"`
	ShortUrl string `json:"short_url"`
}

type Request struct {
	Url string `json:"url"`
}

type Config struct {
	BaseUrl string
	RdbAddr string
	RdbPassword string
}

func main() {
	conf := config.GetConfig()	
	rdb := database.NewRedisClient(conf)
	server.NewServer(conf, rdb).Start()
}