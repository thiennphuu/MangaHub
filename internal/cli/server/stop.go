package server

import (
	"fmt"

	"mangahub/pkg/config"

	"github.com/spf13/cobra"
)

// stopCmd checks which components are running and tells the user how to stop them.
// Actual termination is done in the terminal where each server was started (Ctrl+C).
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop server components",
	Long:  `Check running MangaHub server components and show how to stop them (typically via Ctrl+C in their terminal).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		component, _ := cmd.Flags().GetString("component")

		cfg, err := config.LoadConfig("config.yaml")
		if err != nil {
			fmt.Println("⚠️  Failed to load config.yaml, using built-in defaults:", err)
			cfg = config.DefaultConfig()
		}

		fmt.Println("Checking running MangaHub server components...")
		fmt.Println()

		if component == "" || component == "http" {
			url := fmt.Sprintf("http://%s:%d/health", cfg.HTTP.Host, cfg.HTTP.Port)
			if checkHTTP(url, "HTTP API") {
				fmt.Println("   Stop by pressing Ctrl+C in the terminal running: go run ./cmd/api-server")
			} else if component == "http" {
				fmt.Println("   HTTP API server is not running.")
			}
			fmt.Println()
		}

		if component == "" || component == "tcp" {
			addr := fmt.Sprintf("%s:%d", cfg.TCP.Host, cfg.TCP.Port)
			if checkTCP(addr, "TCP Sync") {
				fmt.Println("   Stop by pressing Ctrl+C in the terminal running: go run ./cmd/tcp-server")
			} else if component == "tcp" {
				fmt.Println("   TCP sync server is not running.")
			}
			fmt.Println()
		}

		if component == "" || component == "udp" {
			addr := fmt.Sprintf("%s:%d", cfg.UDP.Host, cfg.UDP.Port)
			if checkUDP(addr, "UDP Notify") {
				fmt.Println("   Stop by pressing Ctrl+C in the terminal running: go run ./cmd/udp-server")
			} else if component == "udp" {
				fmt.Println("   UDP notification server is not running.")
			}
			fmt.Println()
		}

		if component == "" || component == "grpc" {
			addr := fmt.Sprintf("%s:%d", cfg.GRPC.Host, cfg.GRPC.Port)
			if checkTCP(addr, "gRPC") {
				fmt.Println("   Stop by pressing Ctrl+C in the terminal running: go run ./cmd/grpc-server")
			} else if component == "grpc" {
				fmt.Println("   gRPC server is not running.")
			}
			fmt.Println()
		}

		if component == "" || component == "ws" {
			url := fmt.Sprintf("http://%s:%d/health", cfg.WebSocket.Host, cfg.WebSocket.Port)
			if checkHTTP(url, "WebSocket") {
				fmt.Println("   Stop by pressing Ctrl+C in the terminal running: go run ./cmd/websocket-server")
			} else if component == "ws" {
				fmt.Println("   WebSocket server is not running.")
			}
			fmt.Println()
		}

		return nil
	},
}

func init() {
	stopCmd.Flags().StringP("component", "c", "", "Specific component to check/stop (http, tcp, udp, grpc, ws)")
}
