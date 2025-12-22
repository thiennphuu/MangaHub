package server

import (
	"fmt"

	"mangahub/pkg/config"

	"github.com/spf13/cobra"
)

// startCmd now verifies whether components are already running and gives
// concrete commands to start missing ones.
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start server components",
	Long: `Check and help start MangaHub server components. 
This command checks if each component is reachable and prints the exact 'go run' command to start any missing ones:`,
	RunE: func(cmd *cobra.Command, args []string) error {
		httpOnly, _ := cmd.Flags().GetBool("http-only")
		tcpOnly, _ := cmd.Flags().GetBool("tcp-only")
		udpOnly, _ := cmd.Flags().GetBool("udp-only")
		grpcOnly, _ := cmd.Flags().GetBool("grpc-only")
		wsOnly, _ := cmd.Flags().GetBool("ws-only")

		cfg, err := config.LoadConfig("config.yaml")
		if err != nil {
			fmt.Println("⚠️  Failed to load config.yaml, using built-in defaults:", err)
			cfg = config.DefaultConfig()
		}

		startHTTP := (!tcpOnly && !udpOnly && !grpcOnly && !wsOnly) || httpOnly
		startTCP := (!httpOnly && !udpOnly && !grpcOnly && !wsOnly) || tcpOnly
		startUDP := (!httpOnly && !tcpOnly && !grpcOnly && !wsOnly) || udpOnly
		startGRPC := (!httpOnly && !tcpOnly && !udpOnly && !wsOnly) || grpcOnly
		startWS := (!httpOnly && !tcpOnly && !udpOnly && !grpcOnly) || wsOnly

		fmt.Println("Checking MangaHub Server Components...")
		fmt.Println()

		if startHTTP {
			url := fmt.Sprintf("http://%s:%d/health", cfg.HTTP.Host, cfg.HTTP.Port)
			if checkHTTP(url, "HTTP API") {
				fmt.Println("   Already running (no action needed).")
			} else {
				fmt.Println("   To start HTTP API server: go run ./cmd/api-server")
			}
			fmt.Println()
		}

		if startTCP {
			addr := fmt.Sprintf("%s:%d", cfg.TCP.Host, cfg.TCP.Port)
			if checkTCP(addr, "TCP Sync") {
				fmt.Println("   Already running (no action needed).")
			} else {
				fmt.Println("   To start TCP sync server: go run ./cmd/tcp-server")
			}
			fmt.Println()
		}

		if startUDP {
			addr := fmt.Sprintf("%s:%d", cfg.UDP.Host, cfg.UDP.Port)
			if checkUDP(addr, "UDP Notify") {
				fmt.Println("   Already running (no action needed).")
			} else {
				fmt.Println("   To start UDP notification server: go run ./cmd/udp-server")
			}
			fmt.Println()
		}

		if startGRPC {
			addr := fmt.Sprintf("%s:%d", cfg.GRPC.Host, cfg.GRPC.Port)
			if checkTCP(addr, "gRPC") {
				fmt.Println("   Already running (no action needed).")
			} else {
				fmt.Println("   To start gRPC server: go run ./cmd/grpc-server")
			}
			fmt.Println()
		}

		if startWS {
			url := fmt.Sprintf("http://%s:%d/health", cfg.WebSocket.Host, cfg.WebSocket.Port)
			if checkHTTP(url, "WebSocket") {
				fmt.Println("   Already running (no action needed).")
			} else {
				fmt.Println("   To start WebSocket server: go run ./cmd/websocket-server")
			}
			fmt.Println()
		}

		fmt.Println("Logs: mangahub server logs --follow")
		return nil
	},
}

func init() {
	startCmd.Flags().Bool("http-only", false, "Check/start only HTTP API server")
	startCmd.Flags().Bool("tcp-only", false, "Check/start only TCP sync server")
	startCmd.Flags().Bool("udp-only", false, "Check/start only UDP notification server")
	startCmd.Flags().Bool("grpc-only", false, "Check/start only gRPC internal service")
	startCmd.Flags().Bool("ws-only", false, "Check/start only WebSocket chat server")
}
