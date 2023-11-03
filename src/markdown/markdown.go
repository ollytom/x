package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/yuin/goldmark"
)

func main() {
	sc := bufio.NewScanner(os.Stdin)
	var gotErr bool
	for sc.Scan() {
		if err := goldmark.Convert(sc.Bytes(), os.Stdout); err != nil {
			fmt.Fprintln(os.Stderr, err)
			gotErr = true
		}
	}
	if sc.Err() != nil {
		fmt.Fprintln(os.Stderr, sc.Err())
		gotErr = true
	}
	if gotErr {
		os.Exit(1)
	}
}
