package server

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"mangahub/pkg/config"
	"mangahub/pkg/database"

	"github.com/spf13/cobra"
)

// healthCmd performs real connectivity checks against running server components.
var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Detailed health check",
	Long:  `Perform a detailed health check on all server components using HTTP and TCP reachability checks.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load config (same default path as servers)
		cfg, err := config.LoadConfig("config.yaml")
		if err != nil {
			fmt.Println("⚠️  Failed to load config.yaml, using built-in defaults:", err)
			cfg = config.DefaultConfig()
		}

		fmt.Println("MangaHub Server Health Check")
		fmt.Println("════════════════════════════")
		fmt.Println()

		httpOK := checkHTTP(fmt.Sprintf("http://%s:%d/health", cfg.HTTP.Host, cfg.HTTP.Port), "HTTP API")
		wsOK := checkHTTP(fmt.Sprintf("http://%s:%d/health", cfg.WebSocket.Host, cfg.WebSocket.Port), "WebSocket")
		tcpOK := checkTCP(fmt.Sprintf("%s:%d", cfg.TCP.Host, cfg.TCP.Port), "TCP Sync")
		udpOK := checkUDP(fmt.Sprintf("%s:%d", cfg.UDP.Host, cfg.UDP.Port), "UDP Notify")
		grpcOK := checkTCP(fmt.Sprintf("%s:%d", cfg.GRPC.Host, cfg.GRPC.Port), "gRPC")

		fmt.Println()
		dbOK := checkDatabase(cfg.Database.Path)

		fmt.Println()
		overall := httpOK && wsOK && tcpOK && udpOK && grpcOK && dbOK
		if overall {
			fmt.Println("Overall Health: ✓ Healthy")
		} else {
			fmt.Println("Overall Health: ✗ Degraded (see failed checks above)")
		}

		return nil
	},
}

func checkHTTP(url, name string) bool {
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf(" ✗ %s: unreachable (%v)\n", name, err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf(" ✗ %s: HTTP %d\n", name, resp.StatusCode)
		return false
	}

	fmt.Printf(" ✓ %s: healthy (%s)\n", name, url)
	return true
}

func checkTCP(addr, name string) bool {
	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		fmt.Printf(" ✗ %s: unreachable (%v)\n", name, err)
		return false
	}
	_ = conn.Close()
	fmt.Printf(" ✓ %s: listening (%s)\n", name, addr)
	return true
}

func checkUDP(addr, name string) bool {
	conn, err := net.DialTimeout("udp", addr, 2*time.Second)
	if err != nil {
		fmt.Printf(" ✗ %s: unreachable (%v)\n", name, err)
		return false
	}
	_ = conn.Close()
	fmt.Printf(" ✓ %s: reachable (%s)\n", name, addr)
	return true
}

func checkDatabase(path string) bool {
	db, err := database.New(path)
	if err != nil {
		fmt.Printf(" ✗ Database: failed to open (%v)\n", err)
		return false
	}
	defer db.Close()

	// Simple readiness: try a quick Init in case schema is missing
	if err := db.Init(); err != nil {
		fmt.Printf(" ✗ Database: schema/init error (%v)\n", err)
		return false
	}

	fmt.Printf(" ✓ Database: OK (%s)\n", path)
	return true
}
