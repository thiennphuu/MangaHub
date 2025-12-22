package db

import (
	"fmt"
	"os"

	"mangahub/pkg/database"

	"github.com/spf13/cobra"
)

// statsCmd handles `mangahub db stats`.
var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show database statistics",
	Long:  `Display basic statistics for the MangaHub SQLite database such as size and row counts.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Database statistics for %s\n\n", defaultDBPath)

		info, err := os.Stat(defaultDBPath)
		if err != nil {
			return fmt.Errorf("failed to stat database file: %w", err)
		}

		fmt.Printf("File size: %.2f MB\n", float64(info.Size())/1024.0/1024.0)

		db, err := database.New(defaultDBPath)
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
		defer db.Close()

		printTableCount := func(name string) {
			var count int
			if err := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", name)).Scan(&count); err != nil {
				fmt.Printf("  %s: error (%v)\n", name, err)
				return
			}
			fmt.Printf("  %s: %d rows\n", name, count)
		}

		fmt.Println("\nRow counts:")
		printTableCount("users")
		printTableCount("manga")
		printTableCount("user_progress")
		printTableCount("chat_messages")
		printTableCount("notifications")

		return nil
	},
}

func init() {
	DBCmd.AddCommand(statsCmd)
}
