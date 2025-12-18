package session

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Session stores the current user session
type Session struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
}

// GetPath returns the path to the session file
func GetPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ".mangahub_session"
	}
	return filepath.Join(homeDir, ".mangahub", "session.json")
}

// Save saves the current session to file
func Save(session *Session) error {
	sessionPath := GetPath()
	if err := os.MkdirAll(filepath.Dir(sessionPath), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(sessionPath, data, 0600)
}

// Load loads the current session from file
func Load() (*Session, error) {
	sessionPath := GetPath()
	data, err := os.ReadFile(sessionPath)
	if err != nil {
		return nil, err
	}
	var sess Session
	if err := json.Unmarshal(data, &sess); err != nil {
		return nil, err
	}
	return &sess, nil
}

// Clear removes the session file
func Clear() error {
	sessionPath := GetPath()
	return os.Remove(sessionPath)
}
