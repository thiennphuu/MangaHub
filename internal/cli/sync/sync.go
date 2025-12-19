package sync

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/cobra"

	"mangahub/pkg/client"
	"mangahub/pkg/models"
)

const (
	defaultTCPHost = "localhost"
	defaultTCPPort = 9090
)

// SyncCmd is the main sync command (parent/root for TCP sync subcommands).
var SyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "TCP progress synchronization",
	Long:  `Connect and manage TCP synchronization with the TCP progress sync server.`,
}

// getTCPClient returns a TCP client configured for the sync server.
// Later this can be wired to real CLI config/profile if needed.
func getTCPClient() *client.TCPClient {
	return client.NewTCPClient(defaultTCPHost, defaultTCPPort)
}

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to sync server",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Connecting to TCP sync server at %s:%d...\n", defaultTCPHost, defaultTCPPort)

		c := getTCPClient()
		conn, err := c.Connect()
		if err != nil {
			return fmt.Errorf("failed to connect to sync server: %w", err)
		}
		defer conn.Close()

		now := time.Now().UTC().Format("2006-01-02 15:04:05 MST")

		fmt.Println("✓ Connected successfully!")
		fmt.Println("\nConnection Details:")
		fmt.Printf(" Server: %s:%d\n", defaultTCPHost, defaultTCPPort)
		fmt.Println(" Connection: TCP")
		fmt.Printf(" Connected at: %s\n", now)
		return nil
	},
}

var disconnectCmd = &cobra.Command{
	Use:   "disconnect",
	Short: "Disconnect from sync server",
	Long: `Disconnect from the TCP sync server.

Note: Each CLI command opens a short-lived connection, so this command
currently acts as a helper to verify that the server is reachable and
then closes the connection immediately.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := getTCPClient()

		// We simply open a connection and close it to simulate an explicit disconnect.
		conn, err := c.Connect()
		if err != nil {
			return fmt.Errorf("failed to connect to sync server for disconnect: %w", err)
		}
		_ = conn.Close()

		fmt.Println("✓ Disconnected from sync server")
		return nil
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check sync status",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := getTCPClient()

		fmt.Println("Checking TCP sync server status...")
		if err := c.CheckStatus(); err != nil {
			fmt.Println("TCP Sync Status:")
			fmt.Println(" Connection: ✗ Inactive")
			fmt.Printf(" Error: %v\n", err)
			return nil
		}

		fmt.Println("TCP Sync Status:")
		fmt.Println(" Connection: ✓ Active")
		fmt.Printf(" Server: %s:%d\n", defaultTCPHost, defaultTCPPort)
		fmt.Println(" Mode: Progress broadcast (multi-device)")
		return nil
	},
}

var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Monitor real-time sync updates",
	Long:  "Connect to the TCP sync server and stream real-time reading progress updates for all connected devices.",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := getTCPClient()

		fmt.Printf("Connecting to TCP sync server at %s:%d...\n", defaultTCPHost, defaultTCPPort)

		// Channel used to signal graceful shutdown (Ctrl+C)
		stop := make(chan struct{})

		// Handle Ctrl+C
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt)
		go func() {
			<-sigCh
			fmt.Println("\nStopping monitor...")
			close(stop)
		}()

		fmt.Println("Monitoring real-time sync updates... (Press Ctrl+C to exit)")

		err := c.MonitorUpdates(stop, func(update models.ProgressUpdate) {
			t := time.Unix(update.Timestamp, 0).UTC().Format("15:04:05")
			fmt.Printf("[%s] User %s updated %s → Chapter %d (device: %s)\n",
				t, update.UserID, update.MangaID, update.Chapter, update.DeviceID)
		})
		if err != nil {
			return fmt.Errorf("monitoring failed: %w", err)
		}

		fmt.Println("Monitor stopped.")
		return nil
	},
}

func init() {
	SyncCmd.AddCommand(connectCmd)
	SyncCmd.AddCommand(disconnectCmd)
	SyncCmd.AddCommand(statusCmd)
	SyncCmd.AddCommand(monitorCmd)
}
