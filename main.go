package main

import (
	"github.com/aqin97/cache/cache"
	"github.com/aqin97/cache/server"
)

func main() {
	c := cache.New("inmemory")
	server.New(c).Listen()
}
