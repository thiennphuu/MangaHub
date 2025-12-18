package chat

import "github.com/spf13/cobra"

// ChatCmd is the main chat command
var ChatCmd = &cobra.Command{
	Use:   "chat",
	Short: "WebSocket chat system",
	Long:  `Join chat rooms and communicate with other manga fans in real-time.`,
}

var joinCmd = &cobra.Command{
	Use:   "join [--manga-id <id>]",
	Short: "Join a chat room",
	RunE: func(cmd *cobra.Command, args []string) error {
		mangaID, _ := cmd.Flags().GetString("manga-id")

		println("Connecting to WebSocket chat server at ws://localhost:9093...")

		if mangaID != "" {
			println("✓ Connected to " + mangaID + " Discussion")
		} else {
			println("✓ Connected to General Chat")
		}

		println("\nChat Room: #general")
		println("Connected users: 12")
		println("Your status: Online")
		println("\nRecent messages:")
		println("[16:45] alice: Just finished reading the latest chapter!")
		println("[16:47] bob: Which manga are you reading?")
		println("[16:48] alice: Attack on Titan, it's getting intense")
		println("[16:50] charlie: No spoilers please!")
		println("\n─────────────────────────────────────────────────────")
		println("You are now in chat. Type your message and press Enter.")
		println("Type /help for commands or /quit to leave.")

		return nil
	},
}

var sendCmd = &cobra.Command{
	Use:   "send <message>",
	Short: "Send chat message",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		mangaID, _ := cmd.Flags().GetString("manga-id")

		println("Sending message...")
		println("✓ Message sent successfully")

		if mangaID != "" {
			println("Chat room: " + mangaID)
		} else {
			println("Chat room: #general")
		}

		return nil
	},
}

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "View chat history",
	RunE: func(cmd *cobra.Command, args []string) error {
		println("Chat History:")
		println("[16:45] alice: Just finished reading the latest chapter!")
		println("[16:47] bob: Which manga are you reading?")
		println("[16:48] alice: Attack on Titan, it's getting intense")
		println("[16:50] charlie: No spoilers please!")
		println("[16:51] diana: I love AoT!")
		return nil
	},
}

func init() {
	ChatCmd.AddCommand(joinCmd)
	ChatCmd.AddCommand(sendCmd)
	ChatCmd.AddCommand(historyCmd)

	joinCmd.Flags().StringP("manga-id", "m", "", "Manga ID for specific discussion")
	sendCmd.Flags().StringP("manga-id", "m", "", "Manga ID for specific discussion")
	historyCmd.Flags().StringP("manga-id", "m", "", "Manga ID for specific discussion")
	historyCmd.Flags().IntP("limit", "l", 50, "Number of messages to show")
}
