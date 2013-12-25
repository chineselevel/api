package main

import (
	"fmt"
	"github.com/chineselevel/api/api"
	"github.com/jessevdk/go-flags"
	"log"
	"net/http"
)

var opts struct {
	Port int `short:"p" long:"port" description:"Port number" default:"7000"`
}

func main() {
	defer api.Ops.Redis.Close()

	flags.Parse(&opts)

	// text URL handlers
	http.HandleFunc("/rank", api.RankHandler)
	http.HandleFunc("/split", api.SplitHandler)

	fmt.Printf("Server running on port %d!\n", opts.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", opts.Port), nil))
}
