package hn

import (
	"io"
	"os"
	"testing"
)

func TestFS(t *testing.T) {
	fsys := &FS{}
	f, err := fsys.Open("39845126")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if _, err := io.Copy(os.Stderr, f); err != nil {
		t.Fatal(err)
	}
}
