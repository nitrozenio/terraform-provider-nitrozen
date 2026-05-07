package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const BaseURL = "https://nitrozen.io/api/v1"

type Client struct {
	Token      string
	HTTPClient *http.Client
}

func NewClient(token string) *Client {
	return &Client{
		Token: token,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) DoRequest(method, path string, body interface{}) ([]byte, error) {
	url := BaseURL + path

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(respBytes))
	}

	return respBytes, nil
}

func ExtractID(body []byte) (int64, error) {
	var result struct {
		Data struct {
			ID int64 `json:"id"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return 0, fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Data.ID == 0 {
		return 0, fmt.Errorf("invalid response: missing ID")
	}

	return result.Data.ID, nil
}

// Project methods

func (c *Client) GetProject(id int64) ([]byte, error) {
	return c.DoRequest("GET", fmt.Sprintf("/projects/%d", id), nil)
}

func (c *Client) UpdateProject(id int64, name, description string) ([]byte, error) {
	return c.DoRequest("PUT", fmt.Sprintf("/projects/%d", id), map[string]string{
		"name":        name,
		"description": description,
	})
}

func (c *Client) DeleteProject(id int64) error {
	_, err := c.DoRequest("DELETE", fmt.Sprintf("/projects/%d", id), nil)
	return err
}

// Entry methods

func (c *Client) CreateEntry(projectID int64, title, content, category string, isPublished bool) ([]byte, error) {
	return c.DoRequest("POST", fmt.Sprintf("/projects/%d/entries", projectID), map[string]interface{}{
		"title":        title,
		"content":      content,
		"category":     category,
		"is_published": isPublished,
	})
}

func (c *Client) GetEntry(projectID, entryID int64) ([]byte, error) {
	return c.DoRequest("GET", fmt.Sprintf("/projects/%d/entries/%d", projectID, entryID), nil)
}

func (c *Client) UpdateEntry(projectID, entryID int64, title, content, category string, isPublished bool) ([]byte, error) {
	return c.DoRequest("PUT", fmt.Sprintf("/projects/%d/entries/%d", projectID, entryID), map[string]interface{}{
		"title":        title,
		"content":      content,
		"category":     category,
		"is_published": isPublished,
	})
}

func (c *Client) DeleteEntry(projectID, entryID int64) error {
	_, err := c.DoRequest("DELETE", fmt.Sprintf("/projects/%d/entries/%d", projectID, entryID), nil)
	return err
}
