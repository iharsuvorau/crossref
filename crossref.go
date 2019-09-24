// Package crossref provides a wrapper on top of CrossRef REST API,
// read more at https://github.com/CrossRef/rest-api-doc (old, untrue)
// and https://gitlab.com/crossref/rest_api/tree/master/src/cayenne
// (new). The CrossRef API is unstable, so use on your own risk.
package crossref

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Library specific

// Client is a crossref client which handles all further requests.
type Client struct {
	apiBase   *url.URL
	worksPath *url.URL
}

// New returns a new client with generated internal API URLs.
func New(apiBase string) (*Client, error) {
	apiBase = strings.TrimRight(apiBase, "/") + "/"

	u, err := url.Parse(apiBase)
	if err != nil {
		return nil, err
	}

	w, err := url.Parse("works")
	if err != nil {
		return nil, err
	}
	worksPath := u.ResolveReference(w)

	return &Client{
		apiBase:   u,
		worksPath: worksPath,
	}, nil
}

// APIBase returns the base URL.
func (c *Client) APIBase() *url.URL {
	return c.apiBase
}

// WorksPath returns the works URL.
func (c *Client) WorksPath() *url.URL {
	return c.worksPath
}

// CrossRef specific

// DOI is a unique identifier of a publication.
type DOI string

// Response is a CrossRef REST API response type.
type Response struct {
	Status         string
	MessageType    string `json:"message-type"`
	MessageVersion string `json:"message-version"`
	Message        interface{}
}

// Work is a CrossRef work type.
type Work struct {
	Title           string
	ReferencesCount int
	Authors         []string
}

// GetWork returns a work by DOI.
func GetWork(c *Client, id DOI) (*Work, error) {
	path := fmt.Sprintf("%s/%s", c.WorksPath(), id)
	resp, err := http.Get(path)
	if err != nil || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("failed to get %s: %v, %v", path, resp.StatusCode, err)
	}
	defer resp.Body.Close()
	return decodeWork(resp.Body)
}

func decodeWork(r io.Reader) (*Work, error) {
	resp := Response{}
	err := json.NewDecoder(r).Decode(&resp)
	if err != nil {
		return nil, err
	}
	if resp.Status != "ok" {
		return nil, fmt.Errorf("bad response status: %+v", resp)
	}
	if resp.MessageType != "work" {
		return nil, fmt.Errorf("bad response type: %+v", resp)
	}
	workInt, ok := resp.Message.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed work type conversion: %+v", resp.Message)
	}

	titleParts := workInt["title"].([]interface{})
	var title string
	for _, v := range titleParts {
		title += v.(string)
	}

	refcount := int(workInt["reference-count"].(float64))

	authorsInt := workInt["author"].([]interface{})
	authors := make([]string, len(authorsInt))
	for i, v := range authorsInt {
		author := v.(map[string]interface{})
		authors[i] = author["name"].(string)
	}

	work := Work{
		Title:           title,
		ReferencesCount: refcount,
		Authors:         authors,
	}
	return &work, nil
}
