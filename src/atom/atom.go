package atom

import (
	"encoding/xml"
	"time"
)

const xmlns = "http://www.w3.org/2005/Atom"

type Feed struct {
	ID       string    `xml:"id"`
	Title    string    `xml:"title"`
	Updated  time.Time `xml:"updated"`
	Author   *Author   `xml:"author,omitempty"`
	Link     []Link    `xml:"link,omitempty"`
	Subtitle string    `xml:"subtitle,omitempty"`
	Entries  []Entry   `xml:"entry"`
}

type feed struct {
	XMLName   struct{} `xml:"feed"`
	Namespace string   `xml:"xmlns,attr"`
	*Feed
}

type Author struct {
	Name  string `xml:"name"`
	URI   string `xml:"uri,omitempty"`
	Email string `xml:"email,omitempty"`
}

type Entry struct {
	ID        string     `xml:"id"`
	Title     string     `xml:"title"`
	Updated   time.Time  `xml:"updated,omitempty"`
	Author    *Author    `xml:"author,omitempty"`
	Links     []Link     `xml:"link"`
	Summary   string     `xml:"summary,omitempty"`
	Content   []byte     `xml:"content,omitempty"`
	Published *time.Time `xml:"published,omitempty"`
}

type Link struct {
	HRef string `xml:"href,attr,omitempty"`
	Rel  string `xml:"rel,attr,omitempty"`
	Type string `xml:"type,attr,omitempty"`
}

func Marshal(f *Feed) ([]byte, error) {
	f1 := &feed{
		Namespace: xmlns,
		Feed:      f,
	}
	return xml.MarshalIndent(f1, "", "\t")
	// b = bytes.ReplaceAll(b, []byte("></link>"), []byte("/>"))
	// return b, err
}
