package manga

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"mangahub/internal/manga"
	"mangahub/pkg/database"
)

// MangaCmd is the main manga command
var MangaCmd = &cobra.Command{
	Use:   "manga",
	Short: "Manga search and discovery commands",
	Long:  `Search, discover, and view information about manga titles.`,
}

// getDBPath returns the path to the database file
func getDBPath() string {
	// Try multiple locations for the database
	paths := []string{
		"data/mangahub.db",
		"./data/mangahub.db",
		filepath.Join(os.Getenv("HOME"), ".mangahub", "mangahub.db"),
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return "data/mangahub.db"
}

// getMangaService creates and returns a manga service connected to the database
func getMangaService() (*manga.Service, error) {
	dbPath := getDBPath()
	db, err := database.New(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return manga.NewService(db), nil
}
