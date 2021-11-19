package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Reviews struct {
	Reviewer string
	Date     string
	Score    string
	Content  string
}

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
		Date time.Time `json:"label"`
	} `json:"updated"`
	ContentLabel struct {
		Content string `json:"label"`
	} `json:"content"`
	ScoreLabel struct {
		Score string `json:"label"`
	} `json:"im:rating"`
}

func main() {
	fmt.Printf("Ello Mate! You must be looking for new app store reviews!\nLet me go check on that for you....\n\n")
	art()
	doEvery(5000*time.Millisecond, getJson)
}

func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

func getJson(t time.Time) {
	response, err := http.Get("https://itunes.apple.com/us/rss/customerreviews/id=284882215/sortBy=mostRecent/page=1/json")
	if err != nil {
		fmt.Println(err)

	} else {
		data, _ := ioutil.ReadAll(response.Body)
		err := ioutil.WriteFile("review.json", data, 0777)
		parseJson("review.json")
		e := os.Remove("review.json")
		if e != nil {
			fmt.Println(e)
		}
		if err != nil {
			fmt.Println(err)
		}
	}
}

func parseJson(f string) {
	jsonFile, err := os.Open(f)
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var review Feeds
	json.Unmarshal(byteValue, &review)

	newReviewList := make([]Reviews, 0, 10)

	for i := 0; i < len(review.Feeds.Entries); i++ {
		checkReview := CheckReview(review.Feeds.Entries[i].DateLabel.Date)

		if checkReview {
			convTime := ConvertTime(review.Feeds.Entries[i].DateLabel.Date)
			newReview := Reviews{
				review.Feeds.Entries[i].Reviewer.NameLabel.Name,
				convTime.Format("2006-01-02 15:04:00"),
				review.Feeds.Entries[i].ScoreLabel.Score,
				review.Feeds.Entries[i].ContentLabel.Content,
			}
			newReviewList = append(newReviewList, newReview)
		} else {
			continue
		}
	}

	name := GetFilenameDate()
	file, _ := json.MarshalIndent(newReviewList, "", " ")
	_ = ioutil.WriteFile(name, file, 0644)

	currentTime := time.Now()
	yest := currentTime.Add(-time.Hour * 24)
	count := len(newReviewList)
	fmt.Printf("\nWOW There are %d New reviews since: %s\n", count, yest.Format("2006-01-02 15:04:00"))
	fmt.Printf("\nThey are saved in %s for you.", name)
}

func ConvertTime(t time.Time) time.Time {
	loc, err := time.LoadLocation("America/Los_Angeles")
	localStartTime := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), loc)
	if err != nil {
		fmt.Println(err)
	}

	return localStartTime
}

func GetFilenameDate() string {
	const layout = "01-02-2006"
	t := time.Now()

	return t.Format(layout) + ".json"
}

func CheckReview(t time.Time) bool {
	currentTime := time.Now()
	yest := currentTime.Add(-time.Hour * 48)

	if t.After(yest) {
		return true
	} else {
		return false
	}
}

func art() {
	hehe := `
	██████╗ ██╗   ██╗███╗   ██╗██╗    ██╗ █████╗ ██╗   ██╗
	██╔══██╗██║   ██║████╗  ██║██║    ██║██╔══██╗╚██╗ ██╔╝
	██████╔╝██║   ██║██╔██╗ ██║██║ █╗ ██║███████║ ╚████╔╝
	██╔══██╗██║   ██║██║╚██╗██║██║███╗██║██╔══██║  ╚██╔╝
	██║  ██║╚██████╔╝██║ ╚████║╚███╔███╔╝██║  ██║   ██║
	╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═══╝ ╚══╝╚══╝ ╚═╝  ╚═╝   ╚═╝
	`
	fmt.Println(hehe)
}
