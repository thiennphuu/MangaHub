package db

import (
	"database/sql"
	"fmt"

	"mangahub/pkg/database"

	"github.com/spf13/cobra"
)

const defaultDBPath = "./data/mangahub.db"

// checkCmd handles `mangahub db check`.
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check database integrity",
	Long:  `Run SQLite integrity checks on the MangaHub database and verify core tables exist.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Checking database at %s...\n\n", defaultDBPath)

		db, err := database.New(defaultDBPath)
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
		defer db.Close()

		if err := runIntegrityCheck(db.DB); err != nil {
			return err
		}

		if err := verifyCoreTables(db.DB); err != nil {
			return err
		}

		fmt.Println("\n✓ Database check completed successfully")
		return nil
	},
}

func init() {
	DBCmd.AddCommand(checkCmd)
}

func runIntegrityCheck(sqlDB *sql.DB) error {
	rows, err := sqlDB.Query(`PRAGMA integrity_check;`)
	if err != nil {
		return fmt.Errorf("integrity_check failed: %w", err)
	}
	defer rows.Close()

	ok := true
	for rows.Next() {
		var res string
		if err := rows.Scan(&res); err != nil {
			return err
		}
		if res != "ok" {
			ok = false
			fmt.Println("✗ Integrity issue:", res)
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}

	if ok {
		fmt.Println("✓ PRAGMA integrity_check: ok")
	} else {
		return fmt.Errorf("database integrity check reported issues")
	}
	return nil
}

func verifyCoreTables(sqlDB *sql.DB) error {
	fmt.Println("Verifying core tables exist...")

	required := []string{
		"users",
		"manga",
		"user_progress",
	}

	for _, tbl := range required {
		var name string
		err := sqlDB.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name=?`, tbl).Scan(&name)
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("missing required table: %s", tbl)
			}
			return fmt.Errorf("failed to verify table %s: %w", tbl, err)
		}
		fmt.Printf("✓ Table %s exists\n", tbl)
	}

	return nil
}
