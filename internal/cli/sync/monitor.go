package sync

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/cobra"

	"mangahub/pkg/models"
)

// monitorCmd handles `mangahub sync monitor`.
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
	SyncCmd.AddCommand(monitorCmd)
}
