package notify

import (
	"fmt"

	"github.com/spf13/cobra"
)

var subscribeCmd = &cobra.Command{
	Use:   "subscribe",
	Short: "Subscribe to notifications",
	Long: `Subscribe to manga chapter release notifications.

Examples:
  mangahub notify subscribe
  mangahub notify subscribe --manga-id one-piece
  mangahub notify subscribe --all`,
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
			// Subscribe to specific manga
			err = notifySvc.Subscribe(session.UserID, mangaID)
			if err != nil {
				return fmt.Errorf("failed to subscribe: %w", err)
			}
			fmt.Printf("✓ Subscribed to notifications for '%s'\n", mangaID)
		} else if all {
			// Subscribe to all manga in library
			libSvc, err := getLibraryService()
			if err != nil {
				return fmt.Errorf("database error: %w", err)
			}

			library, err := libSvc.GetLibrary(session.UserID, 100, 0)
			if err != nil {
				return fmt.Errorf("failed to get library: %w", err)
			}

			count := 0
			for _, entry := range library {
				if entry.Status == "reading" || entry.Status == "plan-to-read" {
					err = notifySvc.Subscribe(session.UserID, entry.MangaID)
					if err == nil {
						count++
					}
				}
			}
			fmt.Printf("✓ Subscribed to notifications for %d manga\n", count)
		} else {
			// Enable general notifications
			err = notifySvc.EnableNotifications(session.UserID)
			if err != nil {
				return fmt.Errorf("failed to enable notifications: %w", err)
			}
			fmt.Println("✓ Subscribed to chapter release notifications")
		}

		fmt.Println("\nNotification Preferences:")
		prefs, _ := notifySvc.GetPreferences(session.UserID)
		if prefs != nil {
			fmt.Printf("  Chapter releases: %s\n", boolToEnabled(prefs.ChapterReleases))
			fmt.Printf("  Email notifications: %s\n", boolToEnabled(prefs.EmailNotifications))
			fmt.Printf("  Sound enabled: %s\n", boolToEnabled(prefs.SoundEnabled))
		} else {
			fmt.Println("  Chapter releases: enabled")
			fmt.Println("  Email notifications: enabled")
		}

		return nil
	},
}

func boolToEnabled(b bool) string {
	if b {
		return "enabled"
	}
	return "disabled"
}

func init() {
	NotifyCmd.AddCommand(subscribeCmd)
	subscribeCmd.Flags().StringP("manga-id", "m", "", "Subscribe to specific manga")
	subscribeCmd.Flags().BoolP("all", "a", false, "Subscribe to all manga in library")
}
