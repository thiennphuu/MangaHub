package stats

import "github.com/spf13/cobra"

// StatsCmd is the main stats command
var StatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Reading statistics",
	Long:  `View personal manga reading statistics and analytics.`,
}

var overviewCmd = &cobra.Command{
	Use:   "overview",
	Short: "View reading overview",
	RunE: func(cmd *cobra.Command, args []string) error {
		println("Reading Statistics Overview")
		println("═══════════════════════════════════════")
		println("Total Manga: 47")
		println("Currently Reading: 8")
		println("Completed: 15")
		println("Plan to Read: 18")
		println("\nTotal Chapters Read: 3,547")
		println("Total Hours Reading: ~892")
		println("Average Rating: 7.8/10")
		println("Reading Streak: 45 days")
		println("Most Active Day: Saturday")
		return nil
	},
}

var detailedCmd = &cobra.Command{
	Use:   "detailed",
	Short: "Detailed statistics",
	RunE: func(cmd *cobra.Command, args []string) error {
		println("Detailed Reading Statistics")
		println("═══════════════════════════════════════")
		println("By Status:")
		println("  Reading: 1,234 chapters")
		println("  Completed: 2,313 chapters")
		println("  Total: 3,547 chapters")
		println("\nTop Genres:")
		println("  Action: 18%")
		println("  Adventure: 15%")
		println("  Drama: 12%")
		println("  Shounen: 20%")
		return nil
	},
}

func init() {
	StatsCmd.AddCommand(overviewCmd)
	StatsCmd.AddCommand(detailedCmd)
}
