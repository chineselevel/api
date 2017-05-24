package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/chineselevel/api/api"
	"github.com/jessevdk/go-flags"
)

var opts struct {
	Port int `short:"p" long:"port" description:"Port number" default:"7000"`
}

func main() {
	defer api.Ops.RedisPool.Close()

	flags.Parse(&opts)

	// text URL handlers
	http.HandleFunc("/rank", api.RankHandler)
	http.HandleFunc("/split", api.SplitHandler)
	http.HandleFunc("/words", api.WordsHandler)
	http.HandleFunc("/analyze", api.AnalyzeHandler)

	fmt.Printf("Server running on port %d!\n", opts.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", opts.Port), nil))

	fmt.Printf("Closing connections...")
}
