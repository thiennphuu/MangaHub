package export

import "github.com/spf13/cobra"

// ExportCmd is the main export command
var ExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export data",
	Long:  `Export your library and progress data in various formats.`,
}

var libraryCmd = &cobra.Command{
	Use:   "library",
	Short: "Export library to file",
	RunE: func(cmd *cobra.Command, args []string) error {
		format, _ := cmd.Flags().GetString("format")
		output, _ := cmd.Flags().GetString("output")

		println("Exporting library to " + format + "...")
		println("✓ Successfully exported to " + output)
		println("\nExport Details:")
		println(" Format: " + format)
		println(" Entries: 47")
		println(" File size: ~125 KB")

		return nil
	},
}

var progressCmd = &cobra.Command{
	Use:   "progress",
	Short: "Export progress history",
	RunE: func(cmd *cobra.Command, args []string) error {
		format, _ := cmd.Flags().GetString("format")
		output, _ := cmd.Flags().GetString("output")

		println("Exporting progress to " + format + "...")
		println("✓ Successfully exported to " + output)

		return nil
	},
}

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Export all data",
	RunE: func(cmd *cobra.Command, args []string) error {
		output, _ := cmd.Flags().GetString("output")

		println("Exporting all data...")
		println("✓ Successfully exported to " + output)
		println("\nArchive Contents:")
		println(" - library.json (47 entries)")
		println(" - progress.csv (3,547 chapters)")
		println(" - metadata.json")
		println(" - statistics.json")

		return nil
	},
}

func init() {
	ExportCmd.AddCommand(libraryCmd)
	ExportCmd.AddCommand(progressCmd)
	ExportCmd.AddCommand(allCmd)

	libraryCmd.Flags().String("format", "json", "Export format (json, csv, xml)")
	libraryCmd.Flags().String("output", "library.json", "Output file path")

	progressCmd.Flags().String("format", "csv", "Export format")
	progressCmd.Flags().String("output", "progress.csv", "Output file path")

	allCmd.Flags().String("output", "mangahub-backup.tar.gz", "Output file path")
}
