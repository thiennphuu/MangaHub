package db

import (
	"fmt"

	"mangahub/pkg/database"

	"github.com/spf13/cobra"
)

// optimizeCmd handles `mangahub db optimize`.
// It runs lightweight SQLite optimizations: ANALYZE, REINDEX, and VACUUM.
var optimizeCmd = &cobra.Command{
	Use:   "optimize",
	Short: "Optimize database performance",
	Long: `Run SQLite maintenance commands (ANALYZE, REINDEX, VACUUM) on the MangaHub database.

Use this periodically after large imports or many updates to keep the database compact and query plans up to date.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Optimizing database at %s...\n\n", defaultDBPath)

		db, err := database.New(defaultDBPath)
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
		defer db.Close()

		fmt.Println("Running ANALYZE...")
		if _, err := db.Exec("ANALYZE;"); err != nil {
			return fmt.Errorf("ANALYZE failed: %w", err)
		}
		fmt.Println("✓ ANALYZE completed")

		fmt.Println("Running REINDEX...")
		if _, err := db.Exec("REINDEX;"); err != nil {
			return fmt.Errorf("REINDEX failed: %w", err)
		}
		fmt.Println("✓ REINDEX completed")

		fmt.Println("Running VACUUM...")
		if _, err := db.Exec("VACUUM;"); err != nil {
			return fmt.Errorf("VACUUM failed: %w", err)
		}
		fmt.Println("✓ VACUUM completed")

		fmt.Println("\n✓ Database optimization completed successfully")
		return nil
	},
}

func init() {
	DBCmd.AddCommand(optimizeCmd)
}
