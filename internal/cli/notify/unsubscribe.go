package notify

import (
	"fmt"

	"github.com/spf13/cobra"
)

var unsubscribeCmd = &cobra.Command{
	Use:   "unsubscribe",
	Short: "Unsubscribe from notifications",
	Long: `Unsubscribe from manga chapter release notifications.

Examples:
  mangahub notify unsubscribe
  mangahub notify unsubscribe --manga-id one-piece
  mangahub notify unsubscribe --all`,
	RunE: func(cmd *cobra.Command, args []string) error {
		mangaID, _ := cmd.Flags().GetString("manga-id")
		all, _ := cmd.Flags().GetBool("all")

		// Check if user is logged in
		session, err := loadSession()
		if err != nil {
			fmt.Println("You are not logged in.")
			fmt.Println("\nPlease login first:")
			fmt.Println("  mangahub auth login --username <username>")
			return nil
		}

		// Get notification service
		notifySvc, err := getNotificationService()
		if err != nil {
			return fmt.Errorf("database error: %w", err)
		}

		if mangaID != "" {
			// Unsubscribe from specific manga
			err = notifySvc.Unsubscribe(session.UserID, mangaID)
			if err != nil {
				return fmt.Errorf("failed to unsubscribe: %w", err)
			}
			fmt.Printf("✓ Unsubscribed from notifications for '%s'\n", mangaID)
		} else if all {
			// Unsubscribe from all
			err = notifySvc.UnsubscribeAll(session.UserID)
			if err != nil {
				return fmt.Errorf("failed to unsubscribe: %w", err)
			}
			err = notifySvc.DisableNotifications(session.UserID)
			if err != nil {
				return fmt.Errorf("failed to disable notifications: %w", err)
			}
			fmt.Println("✓ Unsubscribed from all notifications")
		} else {
			// Disable general notifications
			err = notifySvc.DisableNotifications(session.UserID)
			if err != nil {
				return fmt.Errorf("failed to disable notifications: %w", err)
			}
			fmt.Println("✓ Disabled chapter release notifications")
		}

		return nil
	},
}

func init() {
	NotifyCmd.AddCommand(unsubscribeCmd)
	unsubscribeCmd.Flags().StringP("manga-id", "m", "", "Unsubscribe from specific manga")
	unsubscribeCmd.Flags().BoolP("all", "a", false, "Unsubscribe from all notifications")
}
