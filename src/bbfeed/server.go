package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"

	"olowe.co/x/atom"
)

type Server struct {
	client *Client
}

func (srv *Server) handleReq(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL)
	workspace := strings.TrimPrefix(path.Dir(r.URL.String()), "/")
	repo := path.Base(r.URL.String())
	prs, err := srv.client.PullRequests(workspace, repo)
	if errors.Is(err, errNotExist) {
		http.NotFound(w, r)
		return
	} else if err != nil {
		msg := fmt.Sprintf("get pull requests for %s/%s: %v", workspace, repo, err)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	feed := Feed(prs)
	feed.Title = fmt.Sprintf("%s/%s pull requests", workspace, repo)
	feed.ID = fmt.Sprintf("https://bitbucket.org/%s/%s", workspace, repo)
	feed.Link = []atom.Link{{feed.ID, "alternate", "text/html"}}
	b, err := atom.Marshal(feed)
	if err != nil {
		msg := fmt.Sprintf("marshal atom feed: %v", err)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/atom+xml")
	w.Write(b)
}
