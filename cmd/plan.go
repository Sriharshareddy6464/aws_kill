package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/Sriharshareddy6464/aws-kill/engine"
	"github.com/Sriharshareddy6464/aws-kill/models"
	"github.com/Sriharshareddy6464/aws-kill/utils"
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

		fmt.Println("Reading scan inventory...")
		var inventory models.Inventory
		if err := utils.ReadJSON(inventoryPath, &inventory); err != nil {
			return fmt.Errorf("failed to read inventory file: %w", err)
		}

		fmt.Printf("Loaded %d resources from inventory. Analyzing dependencies...\n", len(inventory.Resources))

		// Clean up further downstream files
		os.Remove(filepath.Join("reports", "result.json"))
		os.Remove(filepath.Join("reports", "verification.json"))

		// Run Planner
		planner := engine.NewPlanner()
		plan, err := planner.Plan(cmd.Context(), &inventory)
		if err != nil {
			return fmt.Errorf("planning failed: %w", err)
		}

		// Ensure reports folder exists
		if err := os.MkdirAll("reports", 0755); err != nil {
			return fmt.Errorf("failed to create reports directory: %w", err)
		}

		// Write plan to file
		planPath := filepath.Join("reports", "plan.json")
		if err := utils.WriteJSON(planPath, plan); err != nil {
			return fmt.Errorf("failed to write plan JSON: %w", err)
		}

		// Output plan details
		fmt.Println("\nGenerated AWS Deletion Order:")
		fmt.Println("------------------------------------------------")
		for i, step := range plan.Steps {
			name := step.Name
			if name == "" {
				name = "<no-name>"
			}
			fmt.Printf("%3d. %-30s [%-25s] ID: %s\n", i+1, name, step.Type, step.ID)
		}
		fmt.Println("------------------------------------------------")
		fmt.Printf("Plan created successfully. %d deletion steps generated and saved to %s.\n", len(plan.Steps), planPath)

		return nil
	},
}

func init() {
	RootCmd.AddCommand(planCmd)
}
