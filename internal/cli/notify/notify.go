package notify

import "github.com/spf13/cobra"

// NotifyCmd is the main notify command
var NotifyCmd = &cobra.Command{
	Use:   "notify",
	Short: "UDP notification management",
	Long:  `Subscribe to and manage manga release notifications via UDP.`,
}
