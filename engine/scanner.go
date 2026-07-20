package engine

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	clientRegistry "github.com/Sriharshareddy6464/aws-kill/aws"
	"github.com/Sriharshareddy6464/aws-kill/models"
	"github.com/Sriharshareddy6464/aws-kill/services"
	"github.com/Sriharshareddy6464/aws-kill/utils"
)

type Scanner struct {
	Registry *clientRegistry.ClientRegistry
	Config   aws.Config
}

func NewScanner(cfg aws.Config) *Scanner {
	return &Scanner{
		Registry: clientRegistry.NewClientRegistry(cfg),
		Config:   cfg,
	}
}

// Scan queries AWS for active resources, returning both the inventory and a status report.
func (s *Scanner) Scan(ctx context.Context, tagFilter string) (*models.Inventory, *models.StatusReport, error) {
	inventory := &models.Inventory{
		Resources: make([]models.Resource, 0),
	}

	// Initialize status maps with zero values to maintain order/presentation later
	ec2Counts := map[string]int{
		"Instances":             0,
		"Running Instances":     0,
		"Stopped Instances":     0,
		"Elastic IPs":           0,
		"Volumes":               0,
		"Snapshots":             0,
		"Key Pairs":             0,
		"Security Groups":       0,
		"Network Interfaces":    0,
		"Launch Templates":      0,
		"Placement Groups":      0,
		"Dedicated Hosts":       0,
		"Capacity Reservations": 0,
	}
	vpcCounts := map[string]int{
		"VPCs":              0,
		"Subnets":           0,
		"Route Tables":      0,
		"Internet Gateways": 0,
		"NAT Gateways":      0,
	}
	albCounts := map[string]int{
		"Load Balancers": 0,
		"Target Groups":  0,
	}
	ecsCounts := map[string]int{
		"Clusters":         0,
		"Services":         0,
		"Task Definitions": 0,
		"Running Tasks":    0,
	}
	ecrCounts := map[string]int{
		"Repositories": 0,
		"Images":       0,
	}
	rdsCounts := map[string]int{
		"DB Instances":  0,
		"DB Snapshots":  0,
		"Subnet Groups": 0,
	}
	s3Counts := map[string]int{
		"Buckets": 0,
	}
	cfCounts := map[string]int{
		"Distributions": 0,
	}

	// 1. EC2
	utils.Logger.Info("Scanning EC2 Service...")
	if ec2Res, ec2C, err := services.NewEC2Service(s.Registry.EC2).Scan(ctx, tagFilter); err == nil {
		inventory.Resources = append(inventory.Resources, ec2Res...)
		for k, v := range ec2C {
			ec2Counts[k] = v
		}
	} else {
		utils.Logger.Warn("Skipped EC2 scan", slog.Any("error", err))
	}

	// 2. VPC
	utils.Logger.Info("Scanning VPCs...")
	if vpcRes, vpcC, err := services.NewVPCService(s.Registry.EC2).Scan(ctx, tagFilter); err == nil {
		inventory.Resources = append(inventory.Resources, vpcRes...)
		for k, v := range vpcC {
			vpcCounts[k] = v
		}
	} else {
		utils.Logger.Warn("Skipped VPC scan", slog.Any("error", err))
	}

	// 3. Subnets
	utils.Logger.Info("Scanning Subnets...")
	if subRes, subC, err := services.NewSubnetService(s.Registry.EC2).Scan(ctx, tagFilter); err == nil {
		inventory.Resources = append(inventory.Resources, subRes...)
		for k, v := range subC {
			vpcCounts[k] = v
		}
	} else {
		utils.Logger.Warn("Skipped Subnets scan", slog.Any("error", err))
	}

	// 4. Security Groups
	utils.Logger.Info("Scanning Security Groups...")
	if sgRes, sgC, err := services.NewSecurityGroupService(s.Registry.EC2).Scan(ctx, tagFilter); err == nil {
		inventory.Resources = append(inventory.Resources, sgRes...)
		for k, v := range sgC {
			ec2Counts[k] = v // SGs are grouped under EC2 in summary output
		}
	} else {
		utils.Logger.Warn("Skipped Security Groups scan", slog.Any("error", err))
	}

	// 5. Internet Gateways
	utils.Logger.Info("Scanning Internet Gateways...")
	if igwRes, igwC, err := services.NewInternetGatewayService(s.Registry.EC2).Scan(ctx, tagFilter); err == nil {
		inventory.Resources = append(inventory.Resources, igwRes...)
		for k, v := range igwC {
			vpcCounts[k] = v
		}
	} else {
		utils.Logger.Warn("Skipped Internet Gateways scan", slog.Any("error", err))
	}

	// 6. NAT Gateways
	utils.Logger.Info("Scanning NAT Gateways...")
	if ngwRes, ngwC, err := services.NewNATGatewayService(s.Registry.EC2).Scan(ctx, tagFilter); err == nil {
		inventory.Resources = append(inventory.Resources, ngwRes...)
		for k, v := range ngwC {
			vpcCounts[k] = v
		}
	} else {
		utils.Logger.Warn("Skipped NAT Gateways scan", slog.Any("error", err))
	}

	// 7. Route Tables
	utils.Logger.Info("Scanning Route Tables...")
	if rtRes, rtC, err := services.NewRouteTableService(s.Registry.EC2).Scan(ctx, tagFilter); err == nil {
		inventory.Resources = append(inventory.Resources, rtRes...)
		for k, v := range rtC {
			vpcCounts[k] = v
		}
	} else {
		utils.Logger.Warn("Skipped Route Tables scan", slog.Any("error", err))
	}

	// 8. Elastic IPs
	utils.Logger.Info("Scanning Elastic IPs...")
	if eipRes, eipC, err := services.NewElasticIPService(s.Registry.EC2).Scan(ctx, tagFilter); err == nil {
		inventory.Resources = append(inventory.Resources, eipRes...)
		for k, v := range eipC {
			ec2Counts[k] = v // EIPs are grouped under EC2 in summary output
		}
	} else {
		utils.Logger.Warn("Skipped Elastic IPs scan", slog.Any("error", err))
	}

	// 9. ALB
	utils.Logger.Info("Scanning Application Load Balancers...")
	if albRes, albC, err := services.NewALBService(s.Registry.ELB).Scan(ctx, tagFilter); err == nil {
		inventory.Resources = append(inventory.Resources, albRes...)
		for k, v := range albC {
			albCounts[k] = v
		}
	} else {
		utils.Logger.Warn("Skipped Application Load Balancers scan", slog.Any("error", err))
	}

	// 10. Target Groups
	utils.Logger.Info("Scanning Target Groups...")
	if tgRes, tgC, err := services.NewTargetGroupService(s.Registry.ELB).Scan(ctx, tagFilter); err == nil {
		inventory.Resources = append(inventory.Resources, tgRes...)
		for k, v := range tgC {
			albCounts[k] = v
		}
	} else {
		utils.Logger.Warn("Skipped Target Groups scan", slog.Any("error", err))
	}

	// 11. ECS
	utils.Logger.Info("Scanning ECS Clusters...")
	if ecsRes, ecsC, err := services.NewECSService(s.Registry.ECS).Scan(ctx, tagFilter); err == nil {
		inventory.Resources = append(inventory.Resources, ecsRes...)
		for k, v := range ecsC {
			ecsCounts[k] = v
		}
	} else {
		utils.Logger.Warn("Skipped ECS Clusters scan", slog.Any("error", err))
	}

	// 12. ECR
	utils.Logger.Info("Scanning ECR Repositories...")
	if ecrRes, ecrC, err := services.NewECRService(s.Registry.ECR).Scan(ctx, tagFilter); err == nil {
		inventory.Resources = append(inventory.Resources, ecrRes...)
		for k, v := range ecrC {
			ecrCounts[k] = v
		}
	} else {
		utils.Logger.Warn("Skipped ECR Repositories scan", slog.Any("error", err))
	}

	// 13. RDS
	utils.Logger.Info("Scanning RDS DB Instances...")
	if rdsRes, rdsC, err := services.NewRDSService(s.Registry.RDS).Scan(ctx, tagFilter); err == nil {
		inventory.Resources = append(inventory.Resources, rdsRes...)
		for k, v := range rdsC {
			rdsCounts[k] = v
		}
	} else {
		utils.Logger.Warn("Skipped RDS DB Instances scan", slog.Any("error", err))
	}

	// 14. S3
	utils.Logger.Info("Scanning S3 Buckets...")
	if s3Res, s3C, err := services.NewS3Service(s.Registry.S3).Scan(ctx, tagFilter); err == nil {
		inventory.Resources = append(inventory.Resources, s3Res...)
		for k, v := range s3C {
			s3Counts[k] = v
		}
	} else {
		utils.Logger.Warn("Skipped S3 Buckets scan", slog.Any("error", err))
	}

	// 15. CloudFront
	utils.Logger.Info("Scanning CloudFront Distributions...")
	if cfRes, cfC, err := services.NewCloudFrontService(s.Registry.CloudFront).Scan(ctx, tagFilter); err == nil {
		inventory.Resources = append(inventory.Resources, cfRes...)
		for k, v := range cfC {
			cfCounts[k] = v
		}
	} else {
		utils.Logger.Warn("Skipped CloudFront Distributions scan", slog.Any("error", err))
	}

	// Fill region details
	for i := range inventory.Resources {
		inventory.Resources[i].Region = s.Config.Region
	}

	// Build status report grouping
	statusReport := &models.StatusReport{
		ScanTime: time.Now().UTC().Format("2006-01-02 15:04 MST"),
		Services: []models.ServiceStatus{},
	}

	// Helper function to sum map values
	totalCount := func(m map[string]int) int {
		sum := 0
		for _, v := range m {
			sum += v
		}
		return sum
	}

	// Add service if any resource exists
	addServiceIfDiscovered := func(name string, counts map[string]int) {
		if totalCount(counts) > 0 {
			statusReport.Services = append(statusReport.Services, models.ServiceStatus{
				ServiceName: name,
				Counts:      counts,
			})
		}
	}

	addServiceIfDiscovered("EC2", ec2Counts)
	addServiceIfDiscovered("VPC", vpcCounts)
	addServiceIfDiscovered("Application Load Balancer", albCounts)
	addServiceIfDiscovered("ECS", ecsCounts)
	addServiceIfDiscovered("ECR", ecrCounts)
	addServiceIfDiscovered("RDS", rdsCounts)
	addServiceIfDiscovered("S3", s3Counts)
	addServiceIfDiscovered("CloudFront", cfCounts)

	fmt.Printf("Scan completed. Discovered %d total active resources.\n", len(inventory.Resources))
	return inventory, statusReport, nil
}
