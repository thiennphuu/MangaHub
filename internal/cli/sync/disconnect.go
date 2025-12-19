package sync

import (
	"fmt"

	"github.com/spf13/cobra"
)

// disconnectCmd handles `mangahub sync disconnect`.
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

func init() {
	SyncCmd.AddCommand(disconnectCmd)
}
