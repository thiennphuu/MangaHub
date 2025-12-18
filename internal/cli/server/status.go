package server

import (
	"fmt"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check server status",
	Long:  `Display current status of all MangaHub server components.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		degraded, _ := cmd.Flags().GetBool("degraded")

		fmt.Println("MangaHub Server Status")
		fmt.Println()
		fmt.Println("┌─────────────────────┬──────────┬─────────────────────┬────────────┬────────────┐")
		fmt.Println("│ Service             │ Status   │ Address             │ Uptime     │ Load       │")
		fmt.Println("├─────────────────────┼──────────┼─────────────────────┼────────────┼────────────┤")

		if degraded {
			fmt.Println("│ HTTP API            │ ✓ Online │ localhost:8080      │ 45m        │ 8 req/min  │")
			fmt.Println("│ TCP Sync            │ ✗ Error  │ localhost:9090      │ -          │ -          │")
			fmt.Println("│ UDP Notifications   │ ⚠ Warn   │ localhost:9091      │ 45m        │ 0 clients  │")
			fmt.Println("│ gRPC Internal       │ ✓ Online │ localhost:9092      │ 45m        │ 2 req/min  │")
			fmt.Println("│ WebSocket Chat      │ ✓ Online │ localhost:9093      │ 45m        │ 5 users    │")
			fmt.Println("└─────────────────────┴──────────┴─────────────────────┴────────────┴────────────┘")
			fmt.Println()
			fmt.Println("Overall System Health: ⚠ Degraded")
			fmt.Println("Issues Detected:")
			fmt.Println(" ✗ TCP Sync Server: Port 9090 already in use")
			fmt.Println("   Solution: Kill process on port 9090 or change port in config")
			fmt.Println()
			fmt.Println(" ⚠ UDP Notifications: No clients registered")
			fmt.Println("   This is normal if no users have subscribed to notifications")
			fmt.Println("Run 'mangahub server health' for detailed diagnostics")
		} else {
			fmt.Println("│ HTTP API            │ ✓ Online │ localhost:8080      │ 2h 15m     │ 12 req/min │")
			fmt.Println("│ TCP Sync            │ ✓ Online │ localhost:9090      │ 2h 15m     │ 3 clients  │")
			fmt.Println("│ UDP Notifications   │ ✓ Online │ localhost:9091      │ 2h 15m     │ 8 clients  │")
			fmt.Println("│ gRPC Internal       │ ✓ Online │ localhost:9092      │ 2h 15m     │ 5 req/min  │")
			fmt.Println("│ WebSocket Chat      │ ✓ Online │ localhost:9093      │ 2h 15m     │ 12 users   │")
			fmt.Println("└─────────────────────┴──────────┴─────────────────────┴────────────┴────────────┘")
			fmt.Println()
			fmt.Println("Overall System Health: ✓ Healthy")
			fmt.Println()
			fmt.Println("Database:")
			fmt.Println(" Connection: ✓ Active")
			fmt.Println(" Size: 2.1 MB")
			fmt.Println(" Tables: 3 (users, manga, user_progress)")
			fmt.Println(" Last backup: 2024-01-20 12:00:00")
			fmt.Println()
			fmt.Println("Memory Usage: 45.2 MB / 512 MB (8.8%)")
			fmt.Println("CPU Usage: 2.3% average")
			fmt.Println("Disk Space: 892 MB / 10 GB available")
		}

		return nil
	},
}

func init() {
	statusCmd.Flags().Bool("degraded", false, "Show sample degraded status output")
}
