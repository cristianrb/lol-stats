package main

import (
	"flag"
	"fmt"
	"lol-stats/cristianrb/api"
	"lol-stats/cristianrb/internal"
	"os"
)

const (
	API_KEY = "RIOT_API_KEY"
)

type config struct {
	port   int
	apiKey string
}

func main() {
	var cfg config
	cfg.apiKey = os.Getenv(API_KEY)
	flag.IntVar(&cfg.port, "port", 8080, "API server port")
	flag.Parse()

	addr := fmt.Sprintf(":%d", cfg.port)
	client := api.NewHTTPClient(cfg.apiKey)
	cache := internal.NewKeyValueCache()
	server := api.NewServer(addr, client, cache)
	server.Run()
}
