package db

import (
	"database/sql"
	"fmt"

	"mangahub/pkg/database"

	"github.com/spf13/cobra"
)

// repairCmd handles `mangahub db repair`.
var repairCmd = &cobra.Command{
	Use:   "repair",
	Short: "Attempt to repair the database",
	Long: `Run lightweight SQLite maintenance commands (VACUUM, integrity check)
and re-run schema initialization to ensure tables and indexes exist.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Repairing database at %s...\n\n", defaultDBPath)

		db, err := database.New(defaultDBPath)
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
		defer db.Close()

		if err := runQuickCheck(db.DB); err != nil {
			fmt.Println("⚠️  quick_check reported issues:", err)
		} else {
			fmt.Println("✓ PRAGMA quick_check: ok")
		}

		fmt.Println("Running VACUUM...")
		if _, err := db.Exec("VACUUM;"); err != nil {
			return fmt.Errorf("VACUUM failed: %w", err)
		}
		fmt.Println("✓ VACUUM completed")

		fmt.Println("Re-initializing schema (idempotent)...")
		if err := db.Init(); err != nil {
			return fmt.Errorf("schema initialization failed: %w", err)
		}

		fmt.Println("\n✓ Database repair/maintenance completed")
		return nil
	},
}

func init() {
	DBCmd.AddCommand(repairCmd)
}

func runQuickCheck(sqlDB *sql.DB) error {
	rows, err := sqlDB.Query(`PRAGMA quick_check;`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var res string
		if err := rows.Scan(&res); err != nil {
			return err
		}
		if res != "ok" {
			return fmt.Errorf("%s", res)
		}
	}
	return rows.Err()
}
