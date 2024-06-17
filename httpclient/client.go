package httpclient

import (
	"context"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type Client struct {
	client *http.Client
}

func New(ctx context.Context) *Client {
	cli := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
		},
		Timeout: 5 * time.Second,
		Jar:     new(cookiejar.Jar),
	}

	context.AfterFunc(ctx, func() {
		cli.CloseIdleConnections()
	})

	return &Client{
		client: cli,
	}
}

func (c *Client) WithCookie(url *url.URL, cookie ...*http.Cookie) *Client {
	c.client.Jar.SetCookies(url, cookie)
	return c
}

func (c *Client) ResetCookie() *Client {
	c.client.Jar = new(cookiejar.Jar)
	return c
}

func (c *Client) Get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	context.AfterFunc(ctx, func() {
		resp.Body.Close()
	})

	return resp, nil
}

func (c *Client) Post(ctx context.Context, url string, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	context.AfterFunc(ctx, func() {
		resp.Body.Close()
	})

	return resp, nil
}

func (c *Client) Put(ctx context.Context, url string, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	context.AfterFunc(ctx, func() {
		resp.Body.Close()
	})

	return resp, nil
}

func (c *Client) Patch(ctx context.Context, url string, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	context.AfterFunc(ctx, func() {
		resp.Body.Close()
	})

	return resp, nil
}

func (c *Client) Delete(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	context.AfterFunc(ctx, func() {
		resp.Body.Close()
	})

	return resp, nil
}
