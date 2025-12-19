package chat

import (
	"fmt"

	"github.com/spf13/cobra"
)

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "View chat history",
	Long:  `View recent chat messages from a room.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		mangaID, _ := cmd.Flags().GetString("manga-id")
		limit, _ := cmd.Flags().GetInt("limit")

		roomID := "general"
		if mangaID != "" {
			roomID = mangaID
		}

		if limit == 0 {
			limit = 20
		}

		fmt.Printf("Chat History for #%s (last %d messages):\n", roomID, limit)
		fmt.Println("─────────────────────────────────────────────────────")

		// Note: In a real implementation, this would fetch from the API
		fmt.Println("(Chat history requires API server support)")
		fmt.Println()
		fmt.Println("To view live messages, use:")
		fmt.Printf("  go run ./cmd/cli chat join")
		if mangaID != "" {
			fmt.Printf(" --manga-id %s", mangaID)
		}
		fmt.Println()

		return nil
	},
}

func init() {
	historyCmd.Flags().StringP("manga-id", "m", "", "View history for specific manga")
	historyCmd.Flags().IntP("limit", "l", 20, "Number of messages to show")
	ChatCmd.AddCommand(historyCmd)
}
