package library

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "View your library",
	Long: `Display your manga library with filtering and sorting options.

Examples:
  mangahub library list
  mangahub library list --status reading
  mangahub library list --status completed
  mangahub library list --sort-by title
  mangahub library list --sort-by last-updated --order desc`,
	RunE: func(cmd *cobra.Command, args []string) error {
		status, _ := cmd.Flags().GetString("status")
		sortBy, _ := cmd.Flags().GetString("sort-by")
		order, _ := cmd.Flags().GetString("order")

		if status != "" {
			fmt.Printf("Your Manga Library - Status: %s\n", status)
		} else {
			fmt.Println("Your Manga Library (47 entries)")
		}
		fmt.Printf("Sort: %s (%s)\n\n", sortBy, order)

		fmt.Println("Currently Reading (8):")
		fmt.Println("• one-piece (1095/∞)")
		fmt.Println("• jujutsu-kaisen (247/?)")
		fmt.Println("• attack-on-titan (89/139)")
		fmt.Println("• demon-slayer (156/205)")
		fmt.Println("\nCompleted (15):")
		fmt.Println("• death-note")
		fmt.Println("• fullmetal-alchemist")
		fmt.Println("• naruto")
		fmt.Println("\nPlan to Read (18), On Hold (4), Dropped (2)")

		return nil
	},
}

func init() {
	LibraryCmd.AddCommand(listCmd)
	listCmd.Flags().StringP("status", "s", "", "Filter by status")
	listCmd.Flags().String("sort-by", "last-updated", "Sort field")
	listCmd.Flags().String("order", "desc", "Sort order (asc/desc)")
}
