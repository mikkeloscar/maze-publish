package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	pathUploadStart = "%s/api/repos/%s/upload/start"
	pathUploadFile  = "%s/api/repos/%s/upload/file/%s/%s"
	pathUploadDone  = "%s/api/repos/%s/upload/done/%s"
)

type UploadStart struct {
	SessionID string `json:"session_id"`
}

type client struct {
	client *http.Client
	base   string
}

func NewClient(uri string) *client {
	return &client{http.DefaultClient, uri}
}

func (c *client) UploadStart(repo string) (*UploadStart, error) {
	out := new(UploadStart)
	uri := fmt.Sprintf(pathUploadStart, c.base, repo)
	err := c.post(uri, nil, out)
	return out, err
}

func (c *client) UploadFile(repo, filename, sessionID string, in io.Reader) error {
	uri := fmt.Sprintf(pathUploadFile, c.base, repo, filename, sessionID)
	err := c.postRaw(uri, in, nil)
	return err
}

func (c *client) UploadDone(repo, sessionID string) error {
	uri := fmt.Sprintf(pathUploadDone, c.base, repo, sessionID)
	err := c.postRaw(uri, nil, nil)
	return err
}

func (c *client) post(rawurl string, in, out interface{}) error {
	return c.do(rawurl, "POST", in, out)
}

func (c *client) postRaw(rawurl string, in, out interface{}) error {
	return c.doRaw(rawurl, "POST", in, out)
}

func (c *client) do(rawurl, method string, in, out interface{}) error {
	body, err := c.request(rawurl, method, in, out)
	if err != nil {
		return err
	}
	defer body.Close()

	if out != nil {
		return json.NewDecoder(body).Decode(out)
	}
	return nil
}

func (c *client) doRaw(rawurl, method string, in, out interface{}) error {
	uri, err := url.Parse(rawurl)
	if err != nil {
		return err
	}

	var buf io.Reader
	if in != nil {
		buf = in.(io.Reader)
	}

	req, err := http.NewRequest(method, uri.String(), buf)
	if err != nil {
		return err
	}

	if in != nil {
		req.Header.Set("Content-Type", "application/octet-stream")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode > http.StatusPartialContent {
		fmt.Println(resp.StatusCode)
		out, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf(string(out))
	}

	if out != nil {
		return json.NewDecoder(resp.Body).Decode(out)
	}
	return nil
}

func (c *client) request(rawurl, method string, in, out interface{}) (io.ReadCloser, error) {
	uri, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if in != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(in)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, uri.String(), buf)
	if err != nil {
		return nil, err
	}
	if in != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode > http.StatusPartialContent {
		defer resp.Body.Close()
		out, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf(string(out))
	}
	return resp.Body, nil
}