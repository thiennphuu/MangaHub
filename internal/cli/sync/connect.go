package sync

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"mangahub/pkg/client"
)

const (
	defaultTCPHost = "10.238.53.72" // Server IP
	defaultTCPPort = 9090
)

// getTCPClient returns a TCP client configured for the sync server.
// Later this can be wired to real CLI config/profile if needed.
func getTCPClient() *client.TCPClient {
	return client.NewTCPClient(defaultTCPHost, defaultTCPPort)
}

// connectCmd handles `mangahub sync connect`.
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

func init() {
	SyncCmd.AddCommand(connectCmd)
}
