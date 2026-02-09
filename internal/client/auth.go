package dremioClient

import (
	"fmt"
	"io"
)

func (c *Client) testPAT() error {
	resp, _ := c.RequestToDremio("GET", "/catalog", nil)
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return fmt.Errorf("login failed. status: %d, body: %s", resp.StatusCode, body)
	}
	return nil
}
