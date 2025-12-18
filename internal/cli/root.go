package cli

import (
	"mangahub/internal/cli/auth"
	"mangahub/internal/cli/chat"
	"mangahub/internal/cli/export"
	"mangahub/internal/cli/grpc"
	"mangahub/internal/cli/library"
	"mangahub/internal/cli/manga"
	"mangahub/internal/cli/notify"
	"mangahub/internal/cli/profile"
	"mangahub/internal/cli/progress"
	"mangahub/internal/cli/server"
	"mangahub/internal/cli/stats"
	"mangahub/internal/cli/sync"

	"github.com/spf13/cobra"
)

var (
	token   string
	apiURL  string = "http://localhost:8080"
	verbose bool
)

var rootCmd = &cobra.Command{
	Use:   "mangahub",
	Short: "MangaHub - Manga tracking and discovery CLI",
	Long: `MangaHub is a comprehensive manga tracking system with support for:
- Searching and discovering manga
- Managing your personal library
- Tracking reading progress
- Real-time chat and notifications
- Cross-device synchronization`,
	Version: "3.0.0",
}

func init() {
	rootCmd.PersistentFlags().StringVar(&token, "token", "", "Authentication token")
	rootCmd.PersistentFlags().StringVar(&apiURL, "api", apiURL, "API server URL")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Enable verbose output")

	// Add subcommands
	rootCmd.AddCommand(manga.MangaCmd)
	rootCmd.AddCommand(library.LibraryCmd)
	rootCmd.AddCommand(progress.ProgressCmd)
	rootCmd.AddCommand(sync.SyncCmd)
	rootCmd.AddCommand(notify.NotifyCmd)
	rootCmd.AddCommand(chat.ChatCmd)
	rootCmd.AddCommand(grpc.GRPCCmd)
	rootCmd.AddCommand(stats.StatsCmd)
	rootCmd.AddCommand(export.ExportCmd)
	rootCmd.AddCommand(auth.AuthCmd)
	rootCmd.AddCommand(server.ServerCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(profile.ProfileCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
