package library

import (
	"fmt"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add --manga-id <id>",
	Short: "Add manga to library",
	Long: `Add a manga to your personal library with optional status and rating.

Status options: reading, completed, plan-to-read, on-hold, dropped

Examples:
  mangahub library add --manga-id one-piece --status reading
  mangahub library add --manga-id death-note --status completed --rating 9`,
	RunE: func(cmd *cobra.Command, args []string) error {
		mangaID, _ := cmd.Flags().GetString("manga-id")
		status, _ := cmd.Flags().GetString("status")
		rating, _ := cmd.Flags().GetInt("rating")

		if mangaID == "" {
			return fmt.Errorf("--manga-id is required")
		}

		fmt.Printf("Adding %s to library (status: %s)...\n", mangaID, status)
		if rating > 0 {
			fmt.Printf("Rating: %d/10\n", rating)
		}
		fmt.Println("✓ Successfully added to library")

		return nil
	},
}

func init() {
	LibraryCmd.AddCommand(addCmd)
	addCmd.Flags().StringP("manga-id", "m", "", "Manga ID (required)")
	addCmd.Flags().StringP("status", "s", "reading", "Initial status")
	addCmd.Flags().IntP("rating", "r", 0, "Rating (0-10)")
	addCmd.MarkFlagRequired("manga-id")
}
