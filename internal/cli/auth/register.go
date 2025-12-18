package auth

import (
	"fmt"

	"github.com/spf13/cobra"

	"mangahub/pkg/models"
	"mangahub/pkg/utils"
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new MangaHub account",
	Long: `Register a new account on MangaHub.

Examples:
  mangahub auth register --username johndoe --email john@example.com`,
	RunE: func(cmd *cobra.Command, args []string) error {
		username, _ := cmd.Flags().GetString("username")
		email, _ := cmd.Flags().GetString("email")

		if username == "" {
			return fmt.Errorf("please provide --username")
		}
		if email == "" {
			return fmt.Errorf("please provide --email")
		}

		// Get services
		userSvc, err := getUserService()
		if err != nil {
			return fmt.Errorf("database error: %w", err)
		}
		authSvc := getAuthService()

		// Check if username already exists
		if _, err := userSvc.GetByUsername(username); err == nil {
			return fmt.Errorf("username '%s' is already taken", username)
		}

		// Check if email already exists
		if _, err := userSvc.GetByEmail(email); err == nil {
			return fmt.Errorf("email '%s' is already registered", email)
		}

		// Prompt for password
		prompt := utils.NewPrompt()
		password, err := prompt.Password("Password: ")
		if err != nil {
			return fmt.Errorf("failed to read password: %w", err)
		}
		if len(password) < 6 {
			return fmt.Errorf("password must be at least 6 characters")
		}

		confirmPassword, err := prompt.Password("Confirm Password: ")
		if err != nil {
			return fmt.Errorf("failed to read password: %w", err)
		}
		if password != confirmPassword {
			return fmt.Errorf("passwords do not match")
		}

		// Hash password
		hashedPassword, err := authSvc.HashPassword(password)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		// Create user
		newUser := &models.User{
			ID:           authSvc.GenerateUserID(),
			Username:     username,
			Email:        email,
			PasswordHash: hashedPassword,
		}

		if err := userSvc.Create(newUser); err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		fmt.Println("✓ Registration successful!")
		fmt.Println()
		fmt.Printf("Username: %s\n", newUser.Username)
		fmt.Printf("Email:    %s\n", newUser.Email)
		fmt.Printf("User ID:  %s\n", newUser.ID)
		fmt.Println()
		fmt.Println("You can now login with:")
		fmt.Printf("  mangahub auth login --username %s\n", username)

		return nil
	},
}

func init() {
	AuthCmd.AddCommand(registerCmd)

	registerCmd.Flags().StringP("username", "u", "", "Username (required)")
	registerCmd.Flags().StringP("email", "e", "", "Email address (required)")
}
