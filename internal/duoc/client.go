package duoc

import (
	"bytes"
	"fmt"
	"github.com/imhinotori/duoc-plus/internal/config"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	mobileBaseURL string
	webBaseURL    string
	ssoBaseURL    string
	Client        *http.Client
}

func NewHost(cfg *config.Config) (*Client, error) {
	return &Client{
		mobileBaseURL: cfg.Duoc.MobileAPIUrl,
		webBaseURL:    cfg.Duoc.WebAPIUrl,
		ssoBaseURL:    cfg.Duoc.SSOURL,
		Client:        http.DefaultClient,
	}, nil
}

func (c *Client) request(url, method string, data []byte, auth Credentials) ([]byte, int, error) {
	if method != http.MethodGet && method != http.MethodPost && method != http.MethodPut && method != http.MethodDelete {
		return nil, 0, fmt.Errorf("invalid method: %s", method)
	}

	request, _ := http.NewRequest(method, url, bytes.NewBuffer(data))
	request.Header.Set("Accept", "application/json")

	if strings.Contains(url, c.ssoBaseURL) {
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		request.Header.Set("Content-Type", "application/json")
	}

	if auth != nil {
		request.Header.Set("Authorization", "Bearer "+auth.JWT())
	}

	response, err := c.Client.Do(request)
	if err != nil {
		return nil, 0, err
	}

	if response.StatusCode == http.StatusInternalServerError {
		return nil, response.StatusCode, fmt.Errorf("internal server error: %s", response.Status)
	}

	if response.StatusCode == http.StatusNotFound {
		return nil, response.StatusCode, fmt.Errorf("not found: %s", response.Status)
	}

	var body []byte
	body, err = io.ReadAll(response.Body)
	if err != nil {
		return nil, response.StatusCode, err
	}

	return body, response.StatusCode, nil
}
