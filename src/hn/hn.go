// Package hn provides a filesystem interface to items on Hacker News.
package hn

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

const APIRoot = "https://hacker-news.firebaseio.com/v0"

type Item struct {
	ID     int
	Type   string
	By     string
	Time   int
	Text   string
	Parent int
	URL    string
	Title  string
}

func (it *Item) Name() string       { return strconv.Itoa(it.ID) }
func (it *Item) Size() int64        { r := toMessage(it); return r.Size() }
func (it *Item) Mode() fs.FileMode  { return 0o444 }
func (it *Item) ModTime() time.Time { return time.Unix(int64(it.Time), 0) }
func (it *Item) IsDir() bool        { return false }
func (it *Item) Sys() any           { return nil }

type FS struct {
	cache fs.FS
}

func CacheDirFS(name string) *FS {
	return &FS{cache: os.DirFS(name)}
}

func (fsys *FS) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{"open", name, fs.ErrInvalid}
	}
	name = path.Clean(name)
	switch name {
	case ".":
		return nil, fmt.Errorf("TODO")
	default:
		if _, err := strconv.Atoi(name); err != nil {
			return nil, &fs.PathError{"open", name, fs.ErrNotExist}
		}
	}
	if fsys.cache != nil {
		if f, err := fsys.cache.Open(name); err == nil {
			return f, nil
		}
	}

	u := fmt.Sprintf("%s/item/%s.json", APIRoot, name)
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, err
	}
	return &file{rc: resp.Body}, nil
}

type file struct {
	rc   io.ReadCloser
	item *Item
	msg  *bytes.Reader
}

func (f *file) Read(p []byte) (int, error) {
	var n int
	if f.item == nil {
		if err := json.NewDecoder(f.rc).Decode(&f.item); err != nil {
			return n, fmt.Errorf("decode item: %v", err)
		}
	}
	if f.msg == nil {
		f.msg = toMessage(f.item)
	}
	return f.msg.Read(p)
}

func (f *file) Stat() (fs.FileInfo, error) { return f.item, nil }

func (f *file) Close() error {
	f.msg = nil
	return f.rc.Close()
}

func toMessage(item *Item) *bytes.Reader {
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "From: %s\n", item.By)
	fmt.Fprintf(buf, "Message-ID: <%d@news.ycombinator.com>\n", item.ID)
	fmt.Fprintf(buf, "Date: %s\n", time.Unix(int64(item.Time), 0).Format(time.RFC1123Z))
	if item.Parent != 0 {
		fmt.Fprintf(buf, "References: <%d@news.ycombinator.com>\n", item.Parent)
	}
	if item.Title != "" {
		fmt.Fprintf(buf, "Subject: %s\n", item.Title)
	}
	fmt.Fprintln(buf)
	if item.URL != "" {
		fmt.Fprintln(buf, item.URL)
	}
	if item.Text != "" {
		fmt.Fprintln(buf, strings.ReplaceAll(html.UnescapeString(item.Text), "<p>", "\n\n"))
	}
	return bytes.NewReader(buf.Bytes())
}
