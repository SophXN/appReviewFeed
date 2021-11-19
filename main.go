package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Feeds struct {
	Feeds Feed `json:"feed"`
}

type Feed struct {
	Entries []Entry `json:"entry"`
}

type Entry struct {
	Reviewer struct {
		NameLabel struct {
			Name string `json:"label"`
		} `json:"name"`
	} `json:"author"`
	DateLabel struct {
		Date string `json:"label"`
	} `json:"updated"`
	ContentLabel struct {
		Content string `json:"label"`
	} `json:"content"`
	ScoreLabel struct {
		Score string `json:"label"`
	} `json:"im:rating"`
}

func main() {
	doEvery(20000*time.Millisecond, getJson)
}

func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

func getJson(t time.Time) {
	response, err := http.Get("https://itunes.apple.com/us/rss/customerreviews/id=595068606/sortBy=mostRecent/page=1/json")
	if err != nil {
		fmt.Printf("No works with error %s\n", err)
	} else {
		parseJson("review.json")
		data, _ := ioutil.ReadAll(response.Body)
		err := ioutil.WriteFile("review.json", data, 0777)
		fmt.Printf("\nCreated a new file at: %v", t)
		if err != nil {
			fmt.Println(err)
		}
	}
}

type Reviews struct {
	Reviewer string
	Date     string
	Score    string
	Content  string
}

func parseJson(f string) {

	jsonFile, err := os.Open(f)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Successfully Opened %v", f)

	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var review Feeds
	json.Unmarshal(byteValue, &review)

	newReviewList := make([]Reviews, 0, 10)

	for i := 0; i < len(review.Feeds.Entries); i++ {
		newReview := Reviews{
			(review.Feeds.Entries[i].Reviewer.NameLabel.Name),
			(review.Feeds.Entries[i].DateLabel.Date),
			(review.Feeds.Entries[i].ScoreLabel.Score),
			(review.Feeds.Entries[i].ContentLabel.Content),
		}
		newReviewList = append(newReviewList, newReview)
	}

	file, _ := json.MarshalIndent(newReviewList, "", " ")

	_ = ioutil.WriteFile("output.json", file, 0644)
}
