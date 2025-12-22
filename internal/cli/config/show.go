package config

import (
	"fmt"

	"github.com/spf13/cobra"
)

// showCmd handles `mangahub config show`.
var showCmd = &cobra.Command{
	Use:   "show [section]",
	Short: "Show configuration (optionally by section)",
	Args:  cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			printFullConfig()
			return nil
		}

		section := args[0]
		switch section {
		case "server", "servers":
			fmt.Println("Server Configuration")
			fmt.Println(" API Server: http://localhost:8080")
			fmt.Println(" TCP Sync: localhost:9090")
			fmt.Println(" UDP Notify: localhost:9091")
			fmt.Println(" WebSocket Chat: localhost:9093")
			fmt.Println(" gRPC Service: localhost:9092")
		case "database", "db":
			fmt.Println("Database Configuration")
			fmt.Println(" Type: SQLite")
			fmt.Println(" Path: ~/.mangahub/data.db")
		case "profile":
			fmt.Println("Profile Configuration")
			fmt.Println(" Active profile: default")
		default:
			fmt.Printf("Unknown section: %s\n", section)
			fmt.Println("Available sections: server, database, profile")
		}
		return nil
	},
}

func init() {
	ConfigCmd.AddCommand(showCmd)
}

func printFullConfig() {
	fmt.Println("MangaHub Configuration")
	fmt.Println("═══════════════════════════════════════")
	fmt.Println("API Server: http://localhost:8080")
	fmt.Println("TCP Sync: localhost:9090")
	fmt.Println("UDP Notify: localhost:9091")
	fmt.Println("WebSocket Chat: localhost:9093")
	fmt.Println("gRPC Service: localhost:9092")
	fmt.Println("\nDatabase:")
	fmt.Println(" Type: SQLite")
	fmt.Println(" Path: ~/.mangahub/data.db")
	fmt.Println("\nProfile: default")
}
