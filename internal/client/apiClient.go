package dremioClient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HostURL - Default Hashicups URL
const HostURL string = "https://api.dremio.cloud"

// Client -
type Client struct {
	HostURL             string
	HTTPClient          *http.Client
	PersonalAccessToken string
	Type                string
	ProjectId           string
}

// NewClient -
func NewClient(host, personalAccessToken *string, ptype *string, projectId *string) (*Client, error) {
	c := Client{
		HTTPClient:          &http.Client{Timeout: 30 * time.Second},
		HostURL:             HostURL,
		PersonalAccessToken: *personalAccessToken,
		Type:                *ptype,
		ProjectId:           *projectId,
	}

	if host != nil {
		c.HostURL = *host
	}

	// If username or password not provided, return empty client
	if host == nil || personalAccessToken == nil {
		return &c, nil
	}

	err := c.testPAT()
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (c *Client) RequestToDremio(method, path string, body interface{}, isGlobalEndpoint ...bool) (*http.Response, error) {
	// Default to v3 API for dremio software
	url := fmt.Sprintf("%s/api/v3%s", c.HostURL, path)

	if c.Type == "cloud" {
		// Override to v0 API if isGlobalEndpoint is true
		if len(isGlobalEndpoint) > 0 && isGlobalEndpoint[0] {
			url = fmt.Sprintf("%s/v0%s", c.HostURL, path)
		} else {
			// Default to v0 API for dremio cloud
			url = fmt.Sprintf("%s/v0/projects/%s%s", c.HostURL, c.ProjectId, path)
		}
	}

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.PersonalAccessToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	// Check for non-2xx status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return resp, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return resp, nil
}
