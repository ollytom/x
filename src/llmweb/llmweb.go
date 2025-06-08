package main

import (
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"path"

	"olowe.co/x/openai"
)

type Chat struct {
	client   *openai.Client
	template *template.Template
}

func (c *Chat) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Println(req.RemoteAddr, req.Method, req.URL)
	chat := openai.Chat{
		Model: "mistral-small-latest",
		Messages: []openai.Message{
			{openai.RoleSystem, ""},
		},
	}

	if req.Method == http.MethodGet {
		if err := c.template.Execute(w, &chat); err != nil {
			log.Println(err)
		}
		return
	} else if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if err := req.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("%+v\n", req.PostForm)

	if sys, ok := req.PostForm[openai.RoleSystem]; ok {
		chat.Messages[0].Content = sys[0]
	}

	nuser := len(req.PostForm[openai.RoleUser])
	nassistant := len(req.PostForm[openai.RoleAssistant])
	if nuser != nassistant+1 {
		e := fmt.Sprintf("expected %d user messages for %d assistant messages, got %d", nassistant+1, nassistant, nuser)
		http.Error(w, e, http.StatusBadRequest)
		return
	}

	for i := 0; i < nassistant; i++ {
		user := openai.Message{openai.RoleUser, req.PostForm[openai.RoleUser][i]}
		chat.Messages = append(chat.Messages, user)
		reply := openai.Message{openai.RoleAssistant, req.PostForm[openai.RoleAssistant][i]}
		chat.Messages = append(chat.Messages, reply)
	}
	latest := openai.Message{openai.RoleUser, req.PostForm[openai.RoleUser][nuser-1]}
	chat.Messages = append(chat.Messages, latest)

	reply, err := c.client.Complete(&chat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	chat.Messages = append(chat.Messages, *reply)
	c.template.Execute(w, &chat)
}

func servePWA(w http.ResponseWriter, req *http.Request) {
	log.Println(req.RemoteAddr, req.Method, req.URL)
	http.ServeFile(w, req, "manifest.json")
}

func main() {
	confDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}
	config, err := readConfig(path.Join(confDir, "openai"))
	if err != nil {
		log.Fatalf("read configuration: %v", err)
	}
	client := &openai.Client{http.DefaultClient, config.Token, config.BaseURL}

	tmpl, err := template.ParseGlob("*.html")
	if err != nil {
		log.Fatal(err)
	}

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/", &Chat{client, tmpl})
	http.HandleFunc("/manifest.json", servePWA)
	log.Fatal(http.Serve(ln, nil))
}
