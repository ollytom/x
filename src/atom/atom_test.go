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
