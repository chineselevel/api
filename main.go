package main

import (
	"fmt"
	"github.com/hermanschaaf/chineselevel/api"
	"github.com/jessevdk/go-flags"
	"log"
	"net/http"
)

var opts struct {
	Port int `short:"p" long:"port" description:"Port number" default:"7000"`
}

func main() {
	flags.Parse(&opts)

	// initialize operations
	o := api.NewOperations()
	fmt.Println(o.Redis)

	// text URL handlers
	http.HandleFunc("/rank", api.RankHandler)
	http.HandleFunc("/split", api.SplitHandler)

	fmt.Printf("Server running on port %d!\n", opts.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", opts.Port), nil))
}
