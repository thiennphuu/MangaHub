package chat

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"mangahub/internal/cli/progress"
	"mangahub/pkg/client"
	"mangahub/pkg/session"
)

var sendCmd = &cobra.Command{
	Use:   "send <message>",
	Short: "Send a chat message",
	Long:  `Send a message to the general chat or a specific manga discussion room.`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		mangaID, _ := cmd.Flags().GetString("manga-id")
		message := strings.Join(args, " ")

		// Load session
		sess, err := session.Load()
		if err != nil || sess.Token == "" {
			fmt.Println("⚠ Not logged in. Please login first.")
			fmt.Println("  go run ./cmd/cli auth login --username <your-username>")
			return nil
		}

		// Determine room
		roomID := "general"
		if mangaID != "" {
			roomID = mangaID
		}

		fmt.Printf("Sending message to #%s...\n", roomID)

		// Create temporary client
		wsClient := client.NewWebSocketClient("ws://10.238.53.72:9093", sess.UserID, sess.Username)

		if err := wsClient.Connect(roomID); err != nil {
			fmt.Printf("❌ Failed to connect: %v\n", err)
			fmt.Println("\nMake sure the WebSocket server is running:")
			fmt.Println("  go run ./cmd/websocket-server")
			return nil
		}

		// Send message
		if err := wsClient.SendMessage(message); err != nil {
			fmt.Printf("❌ Failed to send: %v\n", err)
			wsClient.Disconnect()
			return nil
		}

		fmt.Println("✓ Message sent successfully")
		fmt.Printf("Chat room: #%s\n", roomID)
		fmt.Printf("Message: %s\n", message)

		// Store message in SQLite
		// Import progress.RequireDatabase and models.ChatMessage at top if not already
		dbPath := "./data/mangahub.db"
		db, err := progress.RequireDatabase(dbPath)
		if err == nil {
			defer db.Close()
			msgID := fmt.Sprintf("%s-%d", sess.UserID, time.Now().UnixNano())
			timestamp := time.Now().Unix()
			_, err := db.Exec(`INSERT INTO chat_messages (id, user_id, username, room_id, message, timestamp) VALUES (?, ?, ?, ?, ?, ?)`,
				msgID, sess.UserID, sess.Username, roomID, message, timestamp)
			if err != nil {
				fmt.Printf("(SQLite) Failed to store message: %v\n", err)
			}
		} else {
			fmt.Printf("(SQLite) Failed to open DB: %v\n", err)
		}

		wsClient.Disconnect()
		return nil
	},
}

func init() {
	sendCmd.Flags().StringP("manga-id", "m", "", "Send to specific manga chat room")
	ChatCmd.AddCommand(sendCmd)
}
