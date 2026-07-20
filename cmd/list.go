package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/Sriharshareddy6464/aws-kill/models"
	"github.com/Sriharshareddy6464/aws-kill/utils"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all supported AWS services and their scan status",
	Long:  `Displays a structured list of supported services. If a scan was run, shows a checkmark next to services found in the inventory.`,
	Run: func(cmd *cobra.Command, args []string) {
		inventoryPath := filepath.Join("reports", "inventory.json")
		foundMap := make(map[string]bool)

		var inventory models.Inventory
		if err := utils.ReadJSON(inventoryPath, &inventory); err == nil {
			for _, res := range inventory.Resources {
				foundMap[res.Type] = true
			}
		}

		check := func(serviceType string) string {
			if foundMap[serviceType] {
				return "✓"
			}
			return " "
		}

		fmt.Println("AWS Kill Switch - Supported Services")
		fmt.Println()
		fmt.Println("Compute")
		fmt.Println("--------")
		fmt.Printf("%s EC2 Instances\n", check("EC2 Instances"))
		fmt.Println()
		fmt.Println("Networking")
		fmt.Println("-----------")
		fmt.Printf("%s VPC\n", check("VPC"))
		fmt.Printf("%s Subnets\n", check("Subnets"))
		fmt.Printf("%s Security Groups\n", check("Security Groups"))
		fmt.Printf("%s Route Tables\n", check("Route Tables"))
		fmt.Printf("%s Internet Gateway\n", check("Internet Gateway"))
		fmt.Printf("%s NAT Gateway\n", check("NAT Gateway"))
		fmt.Printf("%s Elastic IP\n", check("Elastic IP"))
		fmt.Println()
		fmt.Println("Load Balancing")
		fmt.Println("--------------")
		fmt.Printf("%s Application Load Balancer\n", check("Application Load Balancer"))
		fmt.Printf("%s Target Groups\n", check("Target Groups"))
		fmt.Println()
		fmt.Println("Containers")
		fmt.Println("----------")
		fmt.Printf("%s ECS\n", check("ECS"))
		fmt.Printf("%s ECR\n", check("ECR"))
		fmt.Println()
		fmt.Println("Storage")
		fmt.Println("--------")
		fmt.Printf("%s S3\n", check("S3"))
		fmt.Println()
		fmt.Println("Database")
		fmt.Println("--------")
		fmt.Printf("%s RDS\n", check("RDS"))
		fmt.Println()
		fmt.Println("CDN")
		fmt.Println("---")
		fmt.Printf("%s CloudFront\n", check("CloudFront"))
		fmt.Println()
		fmt.Println("Total Supported Services : 15")
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
