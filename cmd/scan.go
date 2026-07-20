package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/Sriharshareddy6464/aws-kill/aws"
	"github.com/Sriharshareddy6464/aws-kill/engine"
	"github.com/Sriharshareddy6464/aws-kill/utils"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan target AWS environment for active resources",
	Long:  `Discovers supported AWS resources in the specified account and region, saving them to reports/inventory.json and reports/status.json.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Starting AWS infrastructure scan...")

		// Reset downstream files to maintain sequence integrity
		cleanupDownstream()

		// Get flags from Viper
		prof := viper.GetString("profile")
		reg := viper.GetString("region")
		tagFilter := viper.GetString("tag")

		// Initialize AWS Session config
		awsCfg, err := aws.NewSession(cmd.Context(), aws.Config{
			Profile: prof,
			Region:  reg,
		})
		if err != nil {
			return fmt.Errorf("failed to initialize AWS config: %w", err)
		}

		// Initialize and run Scanner
		scanner := engine.NewScanner(awsCfg)
		inventory, statusReport, err := scanner.Scan(cmd.Context(), tagFilter)
		if err != nil {
			return fmt.Errorf("failed during AWS scan: %w", err)
		}

		// Ensure reports directory exists
		if err := os.MkdirAll("reports", 0755); err != nil {
			return fmt.Errorf("failed to create reports directory: %w", err)
		}

		// Write inventory report
		inventoryPath := filepath.Join("reports", "inventory.json")
		if err := utils.WriteJSON(inventoryPath, inventory); err != nil {
			return fmt.Errorf("failed to write inventory JSON: %w", err)
		}

		// Write status report
		statusPath := filepath.Join("reports", "status.json")
		if err := utils.WriteJSON(statusPath, statusReport); err != nil {
			return fmt.Errorf("failed to write status JSON: %w", err)
		}

		fmt.Printf("Scan completed. Inventory saved to %s, status saved to %s.\n", inventoryPath, statusPath)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(scanCmd)
}

func cleanupDownstream() {
	files := []string{
		filepath.Join("reports", "status.json"),
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
