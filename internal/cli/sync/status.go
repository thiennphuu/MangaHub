package sync

import (
	"fmt"

	"github.com/spf13/cobra"
)

// statusCmd handles `mangahub sync status`.
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

func init() {
	SyncCmd.AddCommand(statusCmd)
}
