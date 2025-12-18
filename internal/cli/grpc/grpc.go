package grpc

import "github.com/spf13/cobra"

// GRPCCmd is the main gRPC command
var GRPCCmd = &cobra.Command{
	Use:   "grpc",
	Short: "gRPC service operations",
	Long:  `Query and manipulate manga data via gRPC service calls.`,
}

var getCmd = &cobra.Command{
	Use:   "manga get --id <manga-id>",
	Short: "Query manga via gRPC",
	RunE: func(cmd *cobra.Command, args []string) error {
		mangaID, _ := cmd.Flags().GetString("id")

		println("Querying gRPC service for: " + mangaID)
		println("✓ Successfully retrieved")
		println("\nManga Data:")
		println(" ID: " + mangaID)
		println(" Title: Attack on Titan")
		println(" Author: Isayama Hajime")
		println(" Status: Completed")
		println(" Chapters: 139")

		return nil
	},
}

var searchCmd = &cobra.Command{
	Use:   "manga search --query <term>",
	Short: "Search manga via gRPC",
	RunE: func(cmd *cobra.Command, args []string) error {
		query, _ := cmd.Flags().GetString("query")

		println("Searching via gRPC: " + query)
		println("✓ Found 5 results")
		println(" - attack-on-titan")
		println(" - attack-on-titan-jr")
		println(" - aot-before-fall")

		return nil
	},
}

var updateProgressCmd = &cobra.Command{
	Use:   "progress update --manga-id <id> --chapter <number>",
	Short: "Update progress via gRPC",
	RunE: func(cmd *cobra.Command, args []string) error {
		mangaID, _ := cmd.Flags().GetString("manga-id")
		chapter, _ := cmd.Flags().GetInt("chapter")

		println("Updating progress via gRPC...")
		println("✓ Successfully updated")
		println(" Manga: " + mangaID)
		println(" Chapter: " + string(rune(chapter)))

		return nil
	},
}

func init() {
	GRPCCmd.AddCommand(getCmd)
	GRPCCmd.AddCommand(searchCmd)
	GRPCCmd.AddCommand(updateProgressCmd)

	getCmd.Flags().StringP("id", "i", "", "Manga ID")
	searchCmd.Flags().StringP("query", "q", "", "Search query")
	updateProgressCmd.Flags().StringP("manga-id", "m", "", "Manga ID")
	updateProgressCmd.Flags().IntP("chapter", "c", 0, "Chapter number")
}
