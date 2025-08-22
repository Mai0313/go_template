package auth

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"

	"claude_analysis/cmd/installer/internal/env"
)

// GetGAISFToken performs login to get GAISF token from the MLOP gateway
// Returns the GAISF token string or error if login fails
func GetGAISFToken(username, password string) (string, error) {
	// Use connectivity-selected base URL and domain via unified environment selection
	environment := env.SelectAvailableURL()
	baseURL := environment.MLOPBaseURL
	loginURL := strings.TrimRight(baseURL, "/") + "/auth/login"

	// Cookie-aware HTTP client with redirect support (default)
	jar, err := cookiejar.New(nil)
	if err != nil {
		return "", fmt.Errorf("failed to create cookie jar: %w", err)
	}
	client := &http.Client{Jar: jar, Timeout: 30 * time.Second}

	// Step 1: GET login page and parse CSRF from input[name="_csrf"]
	resp, err := client.Get(loginURL)
	if err != nil {
		return "", fmt.Errorf("failed to get login page: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return "", fmt.Errorf("login page request failed, status: %d", resp.StatusCode)
	}

	csrfToken, err := extractCSRFToken(resp.Body)
	if err != nil || csrfToken == "" {
		return "", fmt.Errorf("unable to find CSRF token on login page: %w", err)
	}

	// Step 2: POST credentials to /auth/login
	form := url.Values{
		"_csrf":            {csrfToken},
		"username":         {username},
		"password":         {password},
		"expiration_hours": {"720"}, // 30 * 24
		"domain":           {environment.Config.Domain},
	}

	req, err := http.NewRequest(http.MethodPost, loginURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create login request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", loginURL)

	resp2, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("login request failed: %w", err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode < 200 || resp2.StatusCode >= 400 {
		// Read a small portion for context
		body, _ := io.ReadAll(io.LimitReader(resp2.Body, 1024))
		return "", fmt.Errorf("login failed, status %d: %s", resp2.StatusCode, string(body))
	}

	// Step 3: Parse token from first <textarea>
	token, err := extractFirstTextarea(resp2.Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse token from response: %w", err)
	}
	token = strings.TrimSpace(token)
	if token == "" {
		return "", errors.New("could not find token in login response")
	}
	return token, nil
}

// extractCSRFToken parses the HTML document and returns the value of input[name="_csrf"].
func extractCSRFToken(r io.Reader) (string, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return "", err
	}
	var csrf string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && strings.EqualFold(n.Data, "input") {
			var nameVal, val string
			for _, a := range n.Attr {
				if strings.EqualFold(a.Key, "name") {
					nameVal = a.Val
				}
				if strings.EqualFold(a.Key, "value") {
					val = a.Val
				}
			}
			if nameVal == "_csrf" {
				csrf = val
				return
			}
		}
		for c := n.FirstChild; c != nil && csrf == ""; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	if csrf == "" {
		return "", errors.New("_csrf not found")
	}
	return csrf, nil
}

// extractFirstTextarea returns the text content of the first <textarea> element.
func extractFirstTextarea(r io.Reader) (string, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return "", err
	}
	var result string
	var found bool
	textContent := func(n *html.Node) string {
		var b strings.Builder
		var walk func(*html.Node)
		walk = func(nn *html.Node) {
			if nn.Type == html.TextNode {
				b.WriteString(nn.Data)
			}
			for c := nn.FirstChild; c != nil; c = c.NextSibling {
				walk(c)
			}
		}
		walk(n)
		return b.String()
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && strings.EqualFold(n.Data, "textarea") && !found {
			result = textContent(n)
			found = true
			return
		}
		for c := n.FirstChild; c != nil && !found; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	if !found {
		return "", errors.New("no textarea found")
	}
	return result, nil
}
