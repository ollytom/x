package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"olowe.co/x/openai"
)

var model = flag.String("m", "ministral-8b-latest", "model")
var baseURL = flag.String("u", "http://127.0.0.1:8080", "openai API base URL")
var sysPrompt = flag.String("s", "", "system prompt")
var converse = flag.Bool("c", false, "start a back-and-forth chat")

func readToken() (string, error) {
	confDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	b, err := os.ReadFile(path.Join(confDir, "openai/token"))
	return string(bytes.TrimSpace(b)), err
}

func copyAll(w io.Writer, paths []string) (n int64, err error) {
	if len(paths) == 0 {
		return io.Copy(w, os.Stdin)
	}
	var errs []error
	for _, name := range paths {
		f, err := os.Open(name)
		if err != nil {
			return n, err
		}
		nn, err := io.Copy(w, f)
		if err != nil {
			errs = append(errs, fmt.Errorf("copy %s: %w", name, err))
		}
		n += nn
	}
	return n, errors.Join(errs...)
}

func init() {
	log.SetFlags(0)
	log.SetPrefix("llm: ")
	flag.Parse()
}

func main() {
	token, err := readToken()
	if err != nil {
		log.Fatalf("read auth token: %v", err)
	}
	client := &openai.Client{http.DefaultClient, token, *baseURL}

	chat := openai.Chat{Model: *model}
	if *sysPrompt != "" {
		chat.Messages =  []openai.Message{
			{openai.RoleSystem, *sysPrompt},
		}
	}
	buf := &bytes.Buffer{}
	if !*converse {
		_, err := copyAll(buf, flag.Args())
		if err != nil {
			log.Fatalln("construct prompt:", err)
		}
		msg := openai.Message{openai.RoleUser, buf.String()}
		chat.Messages = append(chat.Messages, msg)
		reply, err := client.Complete(&chat)
		if err != nil {
			log.Fatalln("llm complete:", err)
		}
		fmt.Println(reply.Content)
		return
	}

	sc := bufio.NewScanner(os.Stdin)
	if len(flag.Args()) > 0 {
		log.Println("conversation mode, ignoring arguments")
	}
	for sc.Scan() {
		if sc.Text() == "." {
			msg := openai.Message{openai.RoleUser, buf.String()}
			chat.Messages = append(chat.Messages, msg)
			reply, err := client.Complete(&chat)
			if err != nil {
				fmt.Fprintln(os.Stderr, "chat not completed:", err)
				continue // try again, allowing a retry with another "." line
			}
			buf.Reset()
			fmt.Println(reply.Content)
			chat.Messages = append(chat.Messages, *reply)
			continue
		}
		fmt.Fprintln(buf, sc.Text())
	}
}
