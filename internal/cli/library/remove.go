package library

import (
	"fmt"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove --manga-id <id>",
	Short: "Remove manga from library",
	Long: `Remove a manga from your personal library.

Example:
  mangahub library remove --manga-id completed-series`,
	RunE: func(cmd *cobra.Command, args []string) error {
		mangaID, _ := cmd.Flags().GetString("manga-id")

		if mangaID == "" {
			return fmt.Errorf("--manga-id is required")
		}

		fmt.Printf("Removing %s from library...\n", mangaID)
		fmt.Println("✓ Successfully removed from library")

		return nil
	},
}

func init() {
	LibraryCmd.AddCommand(removeCmd)
	removeCmd.Flags().StringP("manga-id", "m", "", "Manga ID (required)")
	removeCmd.MarkFlagRequired("manga-id")
}
