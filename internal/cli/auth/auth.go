package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"mangahub/internal/auth"
	"mangahub/internal/user"
	"mangahub/pkg/database"
	"mangahub/pkg/models"
	"mangahub/pkg/utils"
)

// AuthCmd is the main auth command
var AuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication commands",
	Long:  `Login, logout, and manage authentication tokens.`,
}

// Session stores the current user session
type Session struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
}

// getDBPath returns the path to the database file
func getDBPath() string {
	paths := []string{
		"data/mangahub.db",
		"./data/mangahub.db",
		filepath.Join(os.Getenv("HOME"), ".mangahub", "mangahub.db"),
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return "data/mangahub.db"
}

// getSessionPath returns the path to the session file
func getSessionPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ".mangahub_session"
	}
	return filepath.Join(homeDir, ".mangahub", "session.json")
}

// saveSession saves the current session to file
func saveSession(session *Session) error {
	sessionPath := getSessionPath()
	if err := os.MkdirAll(filepath.Dir(sessionPath), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(sessionPath, data, 0600)
}

// loadSession loads the current session from file
func loadSession() (*Session, error) {
	sessionPath := getSessionPath()
	data, err := os.ReadFile(sessionPath)
	if err != nil {
		return nil, err
	}
	var session Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}
	return &session, nil
}

// clearSession removes the session file
func clearSession() error {
	sessionPath := getSessionPath()
	return os.Remove(sessionPath)
}

// getUserService creates and returns a user service connected to the database
func getUserService() (*user.Service, error) {
	dbPath := getDBPath()
	db, err := database.New(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return user.NewService(db), nil
}

// getAuthService creates and returns an auth service
func getAuthService() *auth.AuthService {
	return auth.NewAuthService("")
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to MangaHub",
	Long: `Login to MangaHub with your username or email.

Examples:
  mangahub auth login --username johndoe
  mangahub auth login --email john@example.com`,
	RunE: func(cmd *cobra.Command, args []string) error {
		username, _ := cmd.Flags().GetString("username")
		email, _ := cmd.Flags().GetString("email")

		if username == "" && email == "" {
			return fmt.Errorf("please provide --username or --email")
		}

		// Get user service
		userSvc, err := getUserService()
		if err != nil {
			return fmt.Errorf("database error: %w", err)
		}

		// Prompt for password
		prompt := utils.NewPrompt()
		password, err := prompt.Password("Password: ")
		if err != nil {
			return fmt.Errorf("failed to read password: %w", err)
		}

		// Find user
		var foundUser *models.User
		if username != "" {
			fmt.Printf("Logging in as %s...\n", username)
			foundUser, err = userSvc.GetByUsername(username)
			if err != nil {
				return fmt.Errorf("invalid username or password")
			}
		} else {
			fmt.Printf("Logging in with email %s...\n", email)
			foundUser, err = userSvc.GetByEmail(email)
			if err != nil {
				return fmt.Errorf("invalid email or password")
			}
		}

		// Verify password
		authSvc := getAuthService()
		if err := authSvc.VerifyPassword(password, foundUser.PasswordHash); err != nil {
			return fmt.Errorf("invalid username or password")
		}

		// Generate token
		token, expiresAt, err := authSvc.GenerateToken(foundUser.ID, foundUser.Username, foundUser.Email)
		if err != nil {
			return fmt.Errorf("failed to generate token: %w", err)
		}

		// Save session
		session := &Session{
			UserID:    foundUser.ID,
			Username:  foundUser.Username,
			Email:     foundUser.Email,
			Token:     token,
			ExpiresAt: expiresAt.Format("2006-01-02 15:04:05 MST"),
		}
		if err := saveSession(session); err != nil {
			fmt.Printf("Warning: could not save session: %v\n", err)
		}

		fmt.Println("✓ Login successful")
		fmt.Println()
		fmt.Printf("User: %s\n", foundUser.Username)
		fmt.Printf("Email: %s\n", foundUser.Email)
		fmt.Printf("Token expires: %s\n", expiresAt.Format("2006-01-02 15:04:05 MST"))

		return nil
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from MangaHub",
	RunE: func(cmd *cobra.Command, args []string) error {
		session, err := loadSession()
		if err != nil {
			fmt.Println("✓ Already logged out")
			return nil
		}

		fmt.Printf("Logging out user %s...\n", session.Username)
		if err := clearSession(); err != nil {
			return fmt.Errorf("failed to clear session: %w", err)
		}

		fmt.Println("✓ Logged out successfully")
		return nil
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check authentication status",
	RunE: func(cmd *cobra.Command, args []string) error {
		session, err := loadSession()
		if err != nil {
			fmt.Println("Authentication Status: Not logged in")
			fmt.Println("\nUse 'mangahub auth login --username <username>' to login")
			return nil
		}

		// Verify token is still valid
		authSvc := getAuthService()
		claims, err := authSvc.VerifyToken(session.Token)
		if err != nil {
			fmt.Println("Authentication Status: Session expired")
			fmt.Println("\nPlease login again with 'mangahub auth login'")
			return nil
		}

		fmt.Println("Authentication Status: ✓ Logged in")
		fmt.Println()
		fmt.Printf("User ID:  %s\n", claims.UserID)
		fmt.Printf("Username: %s\n", claims.Username)
		fmt.Printf("Email:    %s\n", claims.Email)
		fmt.Printf("Expires:  %s\n", session.ExpiresAt)

		return nil
	},
}

func init() {
	AuthCmd.AddCommand(loginCmd)
	AuthCmd.AddCommand(logoutCmd)
	AuthCmd.AddCommand(statusCmd)

	loginCmd.Flags().StringP("username", "u", "", "Username")
	loginCmd.Flags().StringP("email", "e", "", "Email address")
}
