package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"mangahub/pkg/models"
)

// HTTPClient represents an HTTP client for API calls
type HTTPClient struct {
	BaseURL string
	Token   string
	Client  *http.Client
}

// NewHTTPClient creates a new HTTP client
func NewHTTPClient(baseURL, token string) *HTTPClient {
	return &HTTPClient{
		BaseURL: baseURL,
		Token:   token,
		Client:  &http.Client{},
	}
}

// SetToken sets the authentication token
func (c *HTTPClient) SetToken(token string) {
	c.Token = token
}

// Register registers a new user
func (c *HTTPClient) Register(username, email, password string) (*models.User, error) {
	req := models.RegisterRequest{
		Username: username,
		Email:    email,
		Password: password,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.post("/auth/register", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("registration failed: %s", string(body))
	}

	var user models.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

// Login logs in a user
func (c *HTTPClient) Login(username, password string) (*models.LoginResponse, error) {
	req := models.LoginRequest{
		Username: username,
		Password: password,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.post("/auth/login", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("login failed: %s", string(body))
	}

	var loginResp models.LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return nil, err
	}

	c.Token = loginResp.Token
	return &loginResp, nil
}

// SearchManga searches for manga
func (c *HTTPClient) SearchManga(query string) (*models.SearchResult, error) {
	resp, err := c.get("/manga?q=" + query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search failed with status %d", resp.StatusCode)
	}

	var result models.SearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetManga gets manga details
func (c *HTTPClient) GetManga(id string) (*models.Manga, error) {
	resp, err := c.get("/manga/" + id)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get manga failed with status %d", resp.StatusCode)
	}

	var manga models.Manga
	if err := json.NewDecoder(resp.Body).Decode(&manga); err != nil {
		return nil, err
	}

	return &manga, nil
}

// Helper methods

func (c *HTTPClient) post(endpoint string, data []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", c.BaseURL+endpoint, io.NopCloser(bytes.NewBuffer(data)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	return c.Client.Do(req)
}

func (c *HTTPClient) get(endpoint string) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.BaseURL+endpoint, nil)
	if err != nil {
		return nil, err
	}

	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	return c.Client.Do(req)
}

func (c *HTTPClient) put(endpoint string, data []byte) (*http.Response, error) {
	req, err := http.NewRequest("PUT", c.BaseURL+endpoint, io.NopCloser(bytes.NewBuffer(data)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	return c.Client.Do(req)
}

func (c *HTTPClient) delete(endpoint string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", c.BaseURL+endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	return c.Client.Do(req)
}

// GetLibrary retrieves user's library
func (c *HTTPClient) GetLibrary(status string) ([]models.Progress, error) {
	endpoint := "/users/library"
	if status != "" {
		endpoint += "?status=" + status
	}

	resp, err := c.get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get library: status %d", resp.StatusCode)
	}

	var progressList []models.Progress
	if err := json.NewDecoder(resp.Body).Decode(&progressList); err != nil {
		return nil, err
	}

	return progressList, nil
}

// AddToLibrary adds a manga to the user's library
func (c *HTTPClient) AddToLibrary(mangaID, status string) error {
	payload := map[string]string{
		"manga_id": mangaID,
		"status":   status,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := c.post("/users/library", data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to add to library: status %d", resp.StatusCode)
	}

	return nil
}

// RemoveFromLibrary removes a manga from the user's library
func (c *HTTPClient) RemoveFromLibrary(mangaID string) error {
	resp, err := c.delete("/users/library/" + mangaID)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to remove from library: status %d", resp.StatusCode)
	}

	return nil
}

// UpdateProgress updates reading progress for a manga
func (c *HTTPClient) UpdateProgress(mangaID string, chapter int, status string, rating int) error {
	payload := models.Progress{
		MangaID:        mangaID,
		CurrentChapter: chapter,
		Status:         status,
		Rating:         rating,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := c.put("/users/library/"+mangaID+"/progress", data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update progress: status %d", resp.StatusCode)
	}

	return nil
}
