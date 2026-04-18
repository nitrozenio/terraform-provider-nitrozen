package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const BaseURL = "https://nitrozen.io/api/v1"

type Client struct {
	Token      string
	HTTPClient *http.Client
}

func NewClient(token string) *Client {
	return &Client{
		Token:      token,
		HTTPClient: &http.Client{},
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

	// ✅ Headers
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json") // 🔥 IMPORTANT (fixes HTML issue)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 🔍 DEBUG (VERY IMPORTANT)
	fmt.Println("---- API RESPONSE ----")
	fmt.Println("Status:", resp.StatusCode)
	fmt.Println(string(respBytes))
	fmt.Println("----------------------")

	// ❌ Handle non-2xx properly
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

	err := json.Unmarshal(body, &result)
	if err != nil {
		// 🔥 SHOW RAW RESPONSE IF JSON FAILS
		return 0, fmt.Errorf("failed to parse JSON. Raw response: %s", string(body))
	}

	if result.Data.ID == 0 {
		return 0, fmt.Errorf("invalid response: missing ID. Raw: %s", string(body))
	}

	return result.Data.ID, nil
}

func (c *Client) GetProject(id int64) ([]byte, error) {
	path := fmt.Sprintf("/projects/%d", id)
	return c.DoRequest("GET", path, nil)
}

func (c *Client) DeleteProject(id int64) error {
	path := fmt.Sprintf("/projects/%d", id)
	_, err := c.DoRequest("DELETE", path, nil)
	return err
}

func (c *Client) UpdateProject(id int64, name, description string) ([]byte, error) {
	path := fmt.Sprintf("/projects/%d", id)

	body := map[string]string{
		"name":        name,
		"description": description,
	}

	return c.DoRequest("PUT", path, body)
}
