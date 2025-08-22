package telemetry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"claude_analysis/core/config"
)

// Client handles telemetry data submission to the API
type Client struct {
	httpClient *http.Client
	config     *config.Config
}

// New creates a new telemetry client
func New(cfg *config.Config) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: cfg.API.Timeout,
		},
		config: cfg,
	}
}

// Submit sends telemetry data to the API and returns the response
func (c *Client) Submit(data interface{}) (map[string]interface{}, error) {
	// Check if data is empty
	// 支援傳入 array 或 map
	var jsonData []byte
	var err error
	if data == nil {
		return map[string]interface{}{"status": "success", "message": "no data to submit"}, nil
	}
	jsonData, err = json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Create request
	req, err := http.NewRequest("POST", c.config.API.Endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body (但不强制要求为JSON)
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// 首先检查HTTP状态码来判断成功与否
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		// 成功的情况 - 尝试解析JSON，如果失败就返回成功状态
		var responseDict map[string]interface{}
		if len(responseBody) > 0 && json.Unmarshal(responseBody, &responseDict) == nil {
			// 成功解析JSON响应
			return responseDict, nil
		} else {
			// API成功但没有JSON响应或解析失败，这是正常的
			return map[string]interface{}{
				"status":     "success",
				"statusCode": resp.StatusCode,
				"message":    "request completed successfully",
				"response":   string(responseBody),
			}, nil
		}
	} else {
		// HTTP错误状态码
		return nil, fmt.Errorf("API returned error status %d: %s", resp.StatusCode, string(responseBody))
	}
}
