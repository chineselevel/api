package api

import (
	"encoding/json"
	"fmt"
	"github.com/hermanschaaf/algorithms/median"
	"github.com/hermanschaaf/go-mafan"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type Response map[string]interface{}

type rankRespRank struct {
	Total   int     `json:"total"`
	Median  float64 `json:"median"`
	Average int     `json:"average"`
}

type rankRespWords struct {
	Total   int `json:"total"`
	Unknown int `json:"unknown"`
	Known   int `json:"known"`
}

type RankResponse struct {
	Rank  rankRespRank  `json:"rank"`
	Words rankRespWords `json:"words"`
}

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

// RankHandler returns the average rank and other information about a text
func RankHandler(rw http.ResponseWriter, r *http.Request) {
	text := r.FormValue("text")
	totalRank := 0
	words := mafan.Split(text)
	sm := median.StreamingMedian{}
	numWords := 0

	// todo: should do this in goroutines
	for _, word := range words {
		r := Ops.GetRank(word)
		totalRank += r
		if r > 0 {
			numWords += 1
			sm.Add(r)
		}
	}
	avg := 0
	if numWords > 0 {
		avg = totalRank / numWords
	}
	rr := RankResponse{}
	rrRank := rankRespRank{totalRank, sm.Median, avg}
	rrWords := rankRespWords{len(words), len(words) - numWords, numWords}
	rr.Rank = rrRank
	rr.Words = rrWords
	rw.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(rr)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.Write(j)
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

// WordsHandler returns a tokenized version of the Chinese text
// supplied to it (like SplitHandler), along with information
// on these words, like individual rank
func WordsHandler(rw http.ResponseWriter, r *http.Request) {
	text := r.FormValue("text")
	s := mafan.Split(text)
	wordsInfo := []Response{}
	for i := range s {
		// convert characters to proper encoding for json
		w := strconv.QuoteToASCII(s[i])
		w = w[1 : len(w)-1]

		// create info object for this word
		info := Response{
			"word": w,
			"rank": Ops.GetRank(s[i]),
		}
		wordsInfo = append(wordsInfo, info)
	}

	JSONResponse(rw, &Response{"words": wordsInfo})
}

type Word struct {
	Value string
	Rank  int
}

// getPercentile gets the value of the rank at a percentile position.
// For example, for perc = 80, the returned value will be the smallest
// value for which 80 percent of the values in the array are smaller or equal.
// It expects the passed array to be already sorted, which means this function
// performs in O(1) constant time.
func getPercentile(values []int, perc int) int {
	l := len(values)
	if l == 0 {
		return 0
	}
	pos := (l * perc) / 100
	log.Println(l, pos, values)
	return values[pos]
}

// AnalyzeHandler takes a text and returns statistics on the
// composition: number of characters, words, rank and more.
func AnalyzeHandler(rw http.ResponseWriter, r *http.Request) {
	text := r.FormValue("text")
	words := mafan.Split(text)

	ranks := Ops.GetRanks(words)

	fmt.Println(words, ranks)

	// sort ranks
	sort.Ints(ranks)

	// get different percentiles
	p80, p90, p95, p99 := getPercentile(ranks, 80), getPercentile(ranks, 90),
		getPercentile(ranks, 95), getPercentile(ranks, 99)

	// number of words we expect the average fluent speaker to know
	maxRank := 30000.0

	// calculate the ChineseLevel score out of 100
	score := math.Min(float64(p90), maxRank) / maxRank * 100.0

	// calculate the estimated HSK score; TODO: improve
	hsk := math.Max(1.0, math.Min(float64(p99), maxRank)/maxRank*6.0)

	JSONResponse(rw, &Response{
		"score": score,
		"hsk":   hsk,
		"percentile": &Response{
			"80": p80,
			"90": p90,
			"95": p95,
			"99": p99,
		},
	})
}
