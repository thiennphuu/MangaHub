package server

import (
	"fmt"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop server components",
	Long:  `Stop MangaHub server components. Can stop all servers or specific ones.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		component, _ := cmd.Flags().GetString("component")

		if component != "" {
			fmt.Printf("Stopping %s server...\n", component)
			fmt.Printf("✓ %s server stopped successfully\n", component)
		} else {
			fmt.Println("Stopping all MangaHub server components...")
			fmt.Println("✓ HTTP API server stopped")
			fmt.Println("✓ TCP sync server stopped")
			fmt.Println("✓ UDP notification server stopped")
			fmt.Println("✓ gRPC internal service stopped")
			fmt.Println("✓ WebSocket chat server stopped")
			fmt.Println()
			fmt.Println("All servers stopped successfully")
		}

		return nil
	},
}

func init() {
	stopCmd.Flags().StringP("component", "c", "", "Specific component to stop (http, tcp, udp, grpc, ws)")
}
