// Package atom implements decoding encoding of Atom feeds as
// specified in RFC 4287.
package atom

import (
	"encoding/xml"
	"time"
)

// MediaType is Atom's IANA media type.
const MediaType = "application/atom+xml"

type Feed struct {
	ID       string    `xml:"id"`
	Title    string    `xml:"title"`
	Updated  time.Time `xml:"updated"`
	Author   *Author   `xml:"author,omitempty"`
	Link     []Link    `xml:"link,omitempty"`
	Subtitle string    `xml:"subtitle,omitempty"`
	Entries  []Entry   `xml:"entry"`
}

var rootElement = xml.StartElement{
	Name: xml.Name{
		Space: "http://www.w3.org/2005/Atom",
		Local: "feed",
	},
}

type alias Feed

func (f *Feed) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(alias(*f), rootElement)
}

func (f *Feed) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return d.DecodeElement((*alias)(f), &rootElement)
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
