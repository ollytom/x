// Command hnatom writes to the standard output a RFC 4287 Atom feed
// of current top stories from the Hacker News front page.
// The flags are:
//
//	-n count
//		Include count items in the feed. The default is 30.
//
// See also the [Hacker News API].
//
// [Hacker News API]: https://github.com/HackerNews/API
package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"olowe.co/x/atom"
)

const apiRoot = "https://hacker-news.firebaseio.com/v0"

type Item struct {
	ID          int
	Type        string
	By          string
	Time        int
	Text        string
	Parent      int
	URL         string
	Title       string
	Score       int
	Descendants int
}

func Get(id int) (*Item, error) {
	u := fmt.Sprintf("%s/item/%s.json", apiRoot, strconv.Itoa(id))
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var item Item
	err = json.NewDecoder(resp.Body).Decode(&item)
	return &item, err
}

func Top() ([]int, error) {
	resp, err := http.Get(apiRoot + "/topstories.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	ids := make([]int, 500) // we know the API returns at most 500 items.
	err = json.NewDecoder(resp.Body).Decode(&ids)
	return ids, err
}

const maxEntries = 5

func entryContent(item *Item) []byte {
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "Score: %d\n", item.Score)
	fmt.Fprintln(buf, "<br>")
	comments := fmt.Sprintf("https://news.ycombinator.com/item?id=%d", item.ID)
	fmt.Fprintf(buf, "<a href=%q>Comments: %d", comments, item.Descendants)
	return buf.Bytes()
}

// The most number of items the top API endpoint will return.
const maxItems = 500

// 30 is the item count on the front page of Hacker News.
var numItems = flag.Uint("n", 30, "number of items to fetch")

func init() {
	flag.Parse()
	if *numItems > maxItems {
		*numItems = maxItems
		fmt.Fprintln(os.Stderr, "warning: maximum of 500 entries can be fetched")
	}
}

func main() {
	top, err := Top()
	if err != nil {
		log.Fatal(err)
	}

	feed := &atom.Feed{
		ID:       "http://home.olowe.co/hnatom/feed.atom",
		Title:    "HN Atom",
		Subtitle: "Top posts from Hacker News",
		Link: []atom.Link{
			{
				Rel:  "alternate",
				Type: "html",
				HRef: "https://news.ycombinator.com",
			},
		},
		Updated: time.Now(),
		Entries: make([]atom.Entry, *numItems),
	}

	for i := range top[:len(feed.Entries)] {
		item, err := Get(top[i])
		if err != nil {
			log.Println(err)
			continue
		}
		feed.Entries[i] = atom.Entry{
			ID:      fmt.Sprintf("%s/item/%d.json", apiRoot, top[i]),
			Title:   item.Title,
			Updated: time.Unix(int64(item.Time), 0),
			Author: &atom.Author{
				Name: item.By,
				URI:  "https://news.ycombinator.com/user?id=" + item.By,
			},
			Content: []byte(entryContent(item)),
			Links:   []atom.Link{{HRef: item.URL}},
		}
	}

	b, err := xml.MarshalIndent(feed, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	os.Stdout.Write(b)
}
