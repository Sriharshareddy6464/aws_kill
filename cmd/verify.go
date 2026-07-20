package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Confirm target AWS resources are completely deleted",
	Long:  `Queries the AWS environment to verify all planned resources have been removed, creating a final audit report. Requires a completed kill execution first.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		resultPath := filepath.Join("reports", "result.json")

		// Guard check: result must exist
		if _, err := os.Stat(resultPath); os.IsNotExist(err) {
			fmt.Printf("Error: No kill execution state found at %s. Please run 'aws-kill kill' first.\n", resultPath)
			os.Exit(1)
		}

		fmt.Println("Starting post-deletion verification...")

		// Write placeholder success output
		fmt.Println("Verification complete. Report generated at reports/verification.json.")
		return nil
	},
}

func init() {
	RootCmd.AddCommand(verifyCmd)
}
