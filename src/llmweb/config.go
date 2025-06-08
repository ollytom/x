package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	BaseURL      string
	Token        string
	DefaultModel string
}

func readConfig(name string) (*Config, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var conf Config
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		if strings.HasPrefix(sc.Text(), "#") {
			continue
		} else if sc.Text() == "" {
			continue
		}

		k, v, ok := strings.Cut(strings.TrimSpace(sc.Text()), " ")
		if !ok {
			return nil, fmt.Errorf("key %q: expected space after key", k)
		}
		v = strings.TrimSpace(v)
		switch k {
		case "token":
			conf.Token = v
		case "url":
			conf.BaseURL = v
		case "model":
			conf.DefaultModel = v
		default:
			return nil, fmt.Errorf("unknown configuration key %q", k)
		}
	}
	return &conf, sc.Err()
}
