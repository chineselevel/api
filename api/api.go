package api

import (
	"encoding/json"
	"github.com/hermanschaaf/algorithms/median"
	"github.com/hermanschaaf/go-mafan"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"
)

type rank struct {
	Total   int     `json:"total"`
	Median  float64 `json:"median"`
	Average int     `json:"average"`
}

type words struct {
	Total   int `json:"total"`
	Unknown int `json:"unknown"`
	Known   int `json:"known"`
}

type RankResponse struct {
	Rank  rank  `json:"rank"`
	Words words `json:"words"`
}

// JSONResponse sets the Content-Type header to application/json
// and returns the response.
func JSONResponse(rw http.ResponseWriter, json []byte) {
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(json)
}

// RankHandler returns the average rank and other information about a text
func RankHandler(rw http.ResponseWriter, r *http.Request) {
	text := r.FormValue("text")
	totalRank := 0
	w := mafan.Split(text)
	sm := median.StreamingMedian{}
	numWords := 0

	// todo: should do this in goroutines
	for _, word := range w {
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
	rrRank := rank{totalRank, sm.Median, avg}
	rrWords := words{len(w), len(w) - numWords, numWords}
	rr.Rank = rrRank
	rr.Words = rrWords
	b, err := json.Marshal(rr)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	JSONResponse(rw, b)
}

type SplitResponse struct {
	Text []string `json:"text"`
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
	sr := SplitResponse{Text: s}
	b, err := json.Marshal(sr)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	JSONResponse(rw, b)
}

type WordsResponse struct {
	Word string `json:"word"`
	Rank int    `json:"rank"`
}

// WordsHandler returns a tokenized version of the Chinese text
// supplied to it (like SplitHandler), along with information
// on these words, like individual rank
func WordsHandler(rw http.ResponseWriter, r *http.Request) {
	text := r.FormValue("text")
	s := mafan.Split(text)
	wordsInfo := []WordsResponse{}
	for i := range s {
		// convert characters to proper encoding for json
		w := strconv.QuoteToASCII(s[i])
		w = w[1 : len(w)-1]

		// create info object for this word
		info := WordsResponse{
			Word: w,
			Rank: Ops.GetRank(s[i]),
		}
		wordsInfo = append(wordsInfo, info)
	}
	b, err := json.Marshal(wordsInfo)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	JSONResponse(rw, b)
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

type percentile struct {
	Eighty     int `json:"80"`
	Ninety     int `json:"90"`
	NinetyFive int `json:"95"`
	NinetyNine int `json:"99"`
}

type AnalyzeResponse struct {
	Score      float64    `json:"score"`
	HSK        float64    `json:"hsk"`
	Percentile percentile `json:"percentile"`
}

// AnalyzeHandler takes a text and returns statistics on the
// composition: number of characters, words, rank and more.
func AnalyzeHandler(rw http.ResponseWriter, r *http.Request) {
	text := r.FormValue("text")
	words := mafan.Split(text)

	ranks := Ops.GetRanks(words)

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

	p := percentile{p80, p90, p95, p99}
	resp := AnalyzeResponse{score, hsk, p}
	b, err := json.Marshal(resp)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	JSONResponse(rw, b)
}
