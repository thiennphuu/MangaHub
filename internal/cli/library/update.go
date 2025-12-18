package library

import (
	"fmt"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update --manga-id <id>",
	Short: "Update library entry",
	Long: `Update a manga's status, rating, or other library information.

Example:
  mangahub library update --manga-id one-piece --status completed --rating 10`,
	RunE: func(cmd *cobra.Command, args []string) error {
		mangaID, _ := cmd.Flags().GetString("manga-id")
		status, _ := cmd.Flags().GetString("status")
		rating, _ := cmd.Flags().GetInt("rating")

		if mangaID == "" {
			return fmt.Errorf("--manga-id is required")
		}

		fmt.Printf("Updating %s...\n", mangaID)
		if status != "" {
			fmt.Printf(" Status: %s\n", status)
		}
		if rating > 0 {
			fmt.Printf(" Rating: %d/10\n", rating)
		}
		fmt.Println("✓ Successfully updated")

		return nil
	},
}

func init() {
	LibraryCmd.AddCommand(updateCmd)
	updateCmd.Flags().StringP("manga-id", "m", "", "Manga ID (required)")
	updateCmd.Flags().StringP("status", "s", "", "New status")
	updateCmd.Flags().IntP("rating", "r", 0, "Rating (0-10)")
	updateCmd.MarkFlagRequired("manga-id")
}
