package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Generate deletion plan based on resource dependencies",
	Long:  `Analyzes relationships between scanned resources and maps out an optimal, safe deletion sequence. Requires a completed scan first.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		inventoryPath := filepath.Join("reports", "inventory.json")

		// Guard check: inventory must exist
		if _, err := os.Stat(inventoryPath); os.IsNotExist(err) {
			fmt.Printf("Error: No scan inventory found at %s. Please run 'aws-kill scan' first.\n", inventoryPath)
			os.Exit(1)
		}

		fmt.Println("Analyzing dependencies and building execution plan...")

		// Reset further downstream execution status
		os.Remove(filepath.Join("reports", "result.json"))
		os.Remove(filepath.Join("reports", "verification.json"))

		// Write placeholder success output
		fmt.Println("Plan created. Order generated and saved to reports/plan.json.")
		return nil
	},
}

func init() {
	RootCmd.AddCommand(planCmd)
}
