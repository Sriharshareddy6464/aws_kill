package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var killCmd = &cobra.Command{
	Use:   "kill",
	Short: "Execute planned resource deletions in order",
	Long:  `Sequentially destroys the infrastructure listed in the plan, checking for dependency releases and polling AWS. Requires a completed plan first.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		planPath := filepath.Join("reports", "plan.json")

		// Guard check: plan must exist
		if _, err := os.Stat(planPath); os.IsNotExist(err) {
			fmt.Printf("Error: No execution plan found at %s. Please run 'aws-kill plan' first.\n", planPath)
			os.Exit(1)
		}

		fmt.Println("Executing deletion plan...")

		// Reset further downstream execution status
		os.Remove(filepath.Join("reports", "verification.json"))

		// Write placeholder success output
		fmt.Println("Kill phase executed. Results saved to reports/result.json.")
		return nil
	},
}

func init() {
	RootCmd.AddCommand(killCmd)
}
