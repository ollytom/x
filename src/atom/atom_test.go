package atom

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
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
	got, err := xml.MarshalIndent(f, "", "  ")
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

func TestDecode(t *testing.T) {
	names, err := filepath.Glob("testdata/*.xml")
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range names {
		f, err := os.Open(name)
		if err != nil {
			t.Fatal(err)
		}
		var feed Feed
		if err := xml.NewDecoder(f).Decode(&feed); err != nil {
			t.Errorf("decode %s: %v", name, err)
		}
		f.Close()
	}
}
