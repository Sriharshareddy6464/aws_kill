package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/Sriharshareddy6464/aws-kill/models"
	"github.com/Sriharshareddy6464/aws-kill/utils"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List discovered AWS services and their resource counts",
	Long:  `Displays a structured summary of the infrastructure discovered in the latest scan from reports/status.json.`,
	Run: func(cmd *cobra.Command, args []string) {
		statusPath := filepath.Join("reports", "status.json")

		var status models.StatusReport
		err := utils.ReadJSON(statusPath, &status)
		if err != nil {
			fmt.Printf("Error: No scan status report found at %s. Please run 'aws-kill scan' first.\n", statusPath)
			os.Exit(1)
		}

		fmt.Println("AWS Infrastructure Summary")
		fmt.Println()
		fmt.Println("Scan Time")
		fmt.Println(status.ScanTime)
		fmt.Println()
		fmt.Println("------------------------------------------------")

		servicesFound := len(status.Services)
		resourcesFound := 0

		// Ordered resource keys for pretty printing
		serviceKeys := map[string][]string{
			"EC2": {
				"Instances", "Running Instances", "Stopped Instances", "Elastic IPs", "Volumes",
				"Snapshots", "Key Pairs", "Security Groups", "Network Interfaces",
				"Launch Templates", "Placement Groups", "Dedicated Hosts", "Capacity Reservations",
			},
			"VPC": {
				"VPCs", "Subnets", "Route Tables", "Internet Gateways", "NAT Gateways",
			},
			"Application Load Balancer": {
				"Load Balancers", "Target Groups",
			},
			"ECS": {
				"Clusters", "Services", "Task Definitions", "Running Tasks",
			},
			"ECR": {
				"Repositories", "Images",
			},
			"RDS": {
				"DB Instances", "DB Snapshots", "Subnet Groups",
			},
			"S3": {
				"Buckets",
			},
			"CloudFront": {
				"Distributions",
			},
		}

		// Helper to check if a key is a duplicate count (which shouldn't add to Resources Found)
		isDuplicateKey := func(k string) bool {
			return k == "Running Instances" || k == "Stopped Instances" || k == "Running Tasks" || k == "Images"
		}

		for _, svc := range status.Services {
			fmt.Println()
			fmt.Println(svc.ServiceName)
			fmt.Println("------------------------------------------------")

			orderedKeys, exists := serviceKeys[svc.ServiceName]
			if !exists {
				// Fallback to alphabetical if not defined in order map
				for k, v := range svc.Counts {
					fmt.Printf("%-30s(%d)\n", k, v)
					if !isDuplicateKey(k) {
						resourcesFound += v
					}
				}
				continue
			}

			for _, k := range orderedKeys {
				val, ok := svc.Counts[k]
				if ok {
					fmt.Printf("%-30s(%d)\n", k, val)
					if !isDuplicateKey(k) {
						resourcesFound += val
					}
				}
			}
		}

		fmt.Println()
		fmt.Println("------------------------------------------------")
		fmt.Println("SUMMARY")
		fmt.Println()
		fmt.Printf("%-30s: %d\n", "Services Found", servicesFound)
		fmt.Printf("%-30s: %d\n", "Resources Found", resourcesFound)
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
