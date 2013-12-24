package main

import (
	"encoding/json"
	"fmt"
	"github.com/hermanschaaf/mafan"
	"github.com/jessevdk/go-flags"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var opts struct {
	Port int `short:"p" long:"port" description:"Port number" default:"7000"`
}

type Response map[string]interface{}

func (r Response) String() (s string) {
	b, err := json.Marshal(r)
	if err != nil {
		s = ""
		return
	}
	s = string(b)
	// unescape doubly-escaped unicode characters
	s = strings.Replace(s, "\\u", `u`, -1)
	return
}

// JSONResponse returns a JSON-formatted response of a Response object, with the appropriate
// content-type header set.
func JSONResponse(rw http.ResponseWriter, response *Response) {
	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, response)
}

// RankHandler returns the rank and other information about a text
func RankHandler(rw http.ResponseWriter, r *http.Request) {
	text := r.FormValue("text")
	JSONResponse(rw, &Response{"rank": text})
}

// SplitHandler returns a tokenized version of the Chinese text
// supplied to it.
func SplitHandler(rw http.ResponseWriter, r *http.Request) {
	text := r.FormValue("text")
	s := mafan.Split(text)
	for i := range s {
		s[i] = strconv.QuoteToASCII(s[i])
		s[i] = s[i][1 : len(s[i])-1]
	}
	JSONResponse(rw, &Response{"text": s})
}

func main() {
	flags.Parse(&opts)

	// initialize operations
	o := NewOperations()

	// text URL handlers
	http.HandleFunc("/rank", RankHandler)
	http.HandleFunc("/split", SplitHandler)

	fmt.Printf("Server running on port %d!\n", opts.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", opts.Port), nil))
}
