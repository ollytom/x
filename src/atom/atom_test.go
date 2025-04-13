package atom

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestMarshal(t *testing.T) {
	f := &Feed{
		Title: "Example Feed",
		Link: []Link{
			{HRef: "http://example.org/"},
		},
		Updated: time.Date(2003, time.Month(12), 13, 18, 30, 2, 0, time.UTC),
		Author:  &Author{Name: "John Doe"},
		ID:      "urn:uuid:60a76c80-d399-11d9-b93C-0003939e0af6",
		Entries: []Entry{
			{
				Title: "Atom-Powered Robots Run Amok",
				ID:    "urn:uuid:1225c695-cfb8-4ebb-aaaa-80da344efa6a",
				Links: []Link{
					{HRef: "http://example.org/2003/12/13/atom03"},
				},
				Updated: time.Date(2003, time.Month(12), 13, 18, 30, 2, 0, time.UTC),
				Summary: "Some text.",
			},
		},
	}
	feed1 := &feed{
		Namespace: xmlns,
		Feed:      f,
	}
	got, err := xml.MarshalIndent(feed1, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	got = bytes.ReplaceAll(got, []byte("></link>"), []byte("/>"))

	want, err := os.ReadFile("testdata/1.xml")
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(got, want) {
		t.Errorf("oops")
		fmt.Fprintln(os.Stderr, string(got))
	}
}

func TestEncoder(t *testing.T) {
	f, err := os.Open("testdata/1.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	var feed Feed
	err = NewDecoder(f).Decode(&feed)
	if err != nil {
		t.Fatal(err)
	}

	if feed.ID != "urn:uuid:60a76c80-d399-11d9-b93C-0003939e0af6" {
		t.Errorf("Incorrect ID %s", feed.ID)
	}
	if feed.Title != "Example Feed" {
		t.Errorf("Incorrect Title %s", feed.Title)
	}
	if feed.Updated.String() != "2003-12-13 18:30:02 +0000 UTC" {
		t.Errorf("Incorrect Updated %s", feed.Updated.String())
	}
	if feed.Author.Name != "John Doe" {
		t.Errorf("Incorrect Author Name %s", feed.Author.Name)
	}
	if feed.Link[0].HRef != "http://example.org/" {
		t.Errorf("Incorrect Link HRef %s", feed.Link[0].HRef)
	}
	if feed.Entries[0].ID != "urn:uuid:1225c695-cfb8-4ebb-aaaa-80da344efa6a" {
		t.Errorf("Incorrect Entry ID %s", feed.Entries[0].ID)
	}
	if feed.Entries[0].Title != "Atom-Powered Robots Run Amok" {
		t.Errorf("Incorrect Entry Title %s", feed.Entries[0].Title)
	}
	if feed.Entries[0].Updated.String() != "2003-12-13 18:30:02 +0000 UTC" {
		t.Errorf("Incorrect Entry Updated %s", feed.Entries[0].Updated.String())
	}
	if feed.Entries[0].Links[0].HRef != "http://example.org/2003/12/13/atom03" {
		t.Errorf("Incorrect Entry Link %s", feed.Entries[0].Links[0].HRef)
	}
	if feed.Entries[0].Summary != "Some text." {
		t.Errorf("Incorrect Entry Summary %s", feed.Entries[0].Summary)
	}
}
