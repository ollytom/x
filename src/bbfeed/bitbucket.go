package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"
)

// https://developer.atlassian.com/cloud/bitbucket/rest/

const apiRoot string = "https://api.bitbucket.org/2.0/"

var errNotExist = errors.New("no such repository")

type Repository struct {
	Name        string
	Description string
}

type PullRequest struct {
	ID          int
	Title       string
	Description string
	Summary     struct {
		Raw  string
		HTML string
	}
	Created time.Time `json:"created_on"`
	Updated time.Time `json:"updated_on"`
	Author  struct {
		Type        string
		DisplayName string `json:"display_name"`
	}
	Links struct {
		Self struct {
			HRef string
		}
		HTML struct {
			HRef string
		}
	}
	State string
}

type Client struct {
	Username, Password string
	*http.Client
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(c.Username, c.Password)
	return c.Client.Do(req)
}

func (c *Client) PullRequests(workspace, repo string) ([]PullRequest, error) {
	u, err := url.Parse(apiRoot + path.Join("repositories", workspace, repo, "pullrequests"))
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Add("state", "OPEN")
	q.Add("state", "MERGED")
	q.Add("state", "DECLINED")
	u.RawQuery = q.Encode()
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, errNotExist
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-ok response status: %s", resp.Status)
	}
	v := struct {
		Values []PullRequest
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}
	return v.Values, nil
}

func (c *Client) Repositories(workspace string) ([]Repository, error) {
	var repos []Repository
	type results struct {
		Next   string
		Values []Repository
	}
	next := apiRoot + path.Join("repositories", workspace)
	for next != "" {
		req, err := http.NewRequest(http.MethodGet, next, nil)
		if err != nil {
			return repos, err
		}
		resp, err := c.Do(req)
		if err != nil {
			return repos, err
		}
		if resp.StatusCode == http.StatusNotFound {
			return repos, fmt.Errorf("no such workspace")
		} else if resp.StatusCode > 399 {
			return repos, fmt.Errorf("non-ok status %s", resp.Status)
		}
		var v results
		if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
			return repos, fmt.Errorf("decode repositories: %w", err)
		}
		resp.Body.Close()
		next = v.Next
		repos = append(repos, v.Values...)
	}
	return repos, nil
}
