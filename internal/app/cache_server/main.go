package main

import (
	"Cache/internal/app/cache_server/cache"
	"Cache/internal/app/cache_server/server"
)

func main() {
	c := cache.New(cache.TypInMemory)
	server.New(server.SvrCSP, c).Listen("localhost:45678")
}
