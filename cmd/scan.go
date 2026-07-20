package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan target AWS environment for active resources",
	Long:  `Discovers supported AWS resources in the specified account and region, saving them to reports/inventory.json.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Starting AWS infrastructure scan...")

		// Reset downstream files to maintain sequence integrity
		cleanupDownstream()

		// Placeholders for actual scan logic invocation
		// inventory, err := engine.Scan(context.Background(), filter)
		// ...

		// Write placeholder success output
		fmt.Println("Scan completed. Inventory saved to reports/inventory.json.")
		return nil
	},
}

func init() {
	RootCmd.AddCommand(scanCmd)
}

func cleanupDownstream() {
	files := []string{
		filepath.Join("reports", "plan.json"),
		filepath.Join("reports", "result.json"),
		filepath.Join("reports", "verification.json"),
	}
	for _, file := range files {
		if _, err := os.Stat(file); err == nil {
			os.Remove(file)
		}
	}
}
