package main

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/user"
	"path"
	"strconv"

	"olowe.co/x/atom"
)

func Feed(prs []PullRequest) *atom.Feed {
	feed := atom.Feed{
		Entries: make([]atom.Entry, len(prs)),
	}
	for i, pr := range prs {
		pr := pr
		feed.Entries[i] = atom.Entry{
			Title: pr.Title,
			ID:    pr.Links.Self.HRef,
			Links: []atom.Link{
				{pr.Links.Self.HRef + "/patch", "alternate", "application/mbox"},
				{pr.Links.HTML.HRef, "canonical", "text/html"},
			},
			Updated:   pr.Updated,
			Published: &pr.Created,
			Summary:   pr.Description,
			Author:    atom.Author{Name: pr.Author.DisplayName},
			Content:   []byte(pr.Summary.HTML),
		}
		if pr.Updated.After(feed.Updated) {
			feed.Updated = pr.Updated
		}
	}
	return &feed
}

func readAuth() (username, password string, err error) {
	confDir, err := os.UserConfigDir()
	if err != nil {
		return "", "", err
	}

	b, err := os.ReadFile(path.Join(confDir, "atlassian/bitbucket.org"))
	if err != nil {
		return "", "", err
	}
	b = bytes.TrimSpace(b)
	u, p, ok := bytes.Cut(b, []byte(":"))
	if !ok {
		return "", "", fmt.Errorf("parse credentials: missing %q separator", ":")
	}
	return string(u), string(p), nil
}

const usage string = "usage: bbfeed"

func main() {
	if len(os.Args) != 1 {
		fmt.Println(usage)
		os.Exit(2)
	}

	username, password, err := readAuth()
	if err != nil {
		log.Fatalf("read credentials: %v", err)
	}
	client := &Client{username, password, http.DefaultClient}
	srv := &Server{client}
	http.HandleFunc("/", srv.handleReq)
	log.Fatal(http.ListenAndServe(":8069", nil))
	return


	workspace := os.Args[1]
	repos, err := client.Repositories(workspace)
	if err != nil {
		log.Fatalf("get repositories in %s: %v", os.Args[1], err)
	}

	u, err := user.Current()
	if err != nil {
		log.Fatalf("lookup current user: %v", err)
	}
	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		log.Fatalf("parse uid: %v", err)
	}
	gid, err := strconv.Atoi(u.Gid)
	if err != nil {
		log.Fatalf("parse gid: %v", err)
	}
	tw := tar.NewWriter(os.Stdout)
	tw.WriteHeader(&tar.Header{
		Name: path.Dir(os.Args[1]) + "/",
		Mode: int64(0o444 | fs.ModeDir),
		Uid:  uid,
		Gid:  gid,
	})
	for _, repo := range repos {
		name := path.Join(workspace, repo.Name)
		prs, err := client.PullRequests(workspace, repo.Name)
		if err != nil {
			log.Printf("get %s pull requests: %v", name, err)
			continue
		}
		feed := Feed(prs)
		feed.Title = name + " pull requests"
		feed.ID = fmt.Sprintf("https://bitbucket.org/" + name)
		b, err := atom.Marshal(feed)
		if err != nil {
			log.Printf("marshal %s feed: %v", name, err)
			continue
		}
		hdr := tar.Header{
			Name: name,
			Size: int64(len(b)),
			Mode: 0o644,
			Uid:  uid,
			Gid:  gid,
		}
		if err := tw.WriteHeader(&hdr); err != nil {
			log.Fatal(err)
		}
		tw.Write(b)
	}
	tw.Flush()
	tw.Close()
}
