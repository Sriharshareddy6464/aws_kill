package engine

import (
	"context"
	"fmt"
	"log/slog"

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

// Scan queries AWS for active resources across all 15 supported services.
func (s *Scanner) Scan(ctx context.Context, tagFilter string) (*models.Inventory, error) {
	inventory := &models.Inventory{
		Resources: make([]models.Resource, 0),
	}

	// 1. EC2
	utils.Logger.Info("Scanning EC2 Instances...")
	if ec2Res, err := services.NewEC2Service(s.Registry.EC2).Scan(ctx, tagFilter); err == nil {
		inventory.Resources = append(inventory.Resources, ec2Res...)
	} else {
		utils.Logger.Warn("Skipped EC2 Instances scan", slog.Any("error", err))
	}

	// 2. VPC
	utils.Logger.Info("Scanning VPCs...")
	if vpcRes, err := services.NewVPCService(s.Registry.EC2).Scan(ctx, tagFilter); err == nil {
		inventory.Resources = append(inventory.Resources, vpcRes...)
	} else {
		utils.Logger.Warn("Skipped VPC scan", slog.Any("error", err))
	}

	// 3. Subnets
	utils.Logger.Info("Scanning Subnets...")
	if subRes, err := services.NewSubnetService(s.Registry.EC2).Scan(ctx, tagFilter); err == nil {
		inventory.Resources = append(inventory.Resources, subRes...)
	} else {
		utils.Logger.Warn("Skipped Subnets scan", slog.Any("error", err))
	}

	// 4. Security Groups
	utils.Logger.Info("Scanning Security Groups...")
	if sgRes, err := services.NewSecurityGroupService(s.Registry.EC2).Scan(ctx, tagFilter); err == nil {
		inventory.Resources = append(inventory.Resources, sgRes...)
	} else {
		utils.Logger.Warn("Skipped Security Groups scan", slog.Any("error", err))
	}

	// 5. Internet Gateways
	utils.Logger.Info("Scanning Internet Gateways...")
	if igwRes, err := services.NewInternetGatewayService(s.Registry.EC2).Scan(ctx, tagFilter); err == nil {
		inventory.Resources = append(inventory.Resources, igwRes...)
	} else {
		utils.Logger.Warn("Skipped Internet Gateways scan", slog.Any("error", err))
	}

	// 6. NAT Gateways
	utils.Logger.Info("Scanning NAT Gateways...")
	if ngwRes, err := services.NewNATGatewayService(s.Registry.EC2).Scan(ctx, tagFilter); err == nil {
		inventory.Resources = append(inventory.Resources, ngwRes...)
	} else {
		utils.Logger.Warn("Skipped NAT Gateways scan", slog.Any("error", err))
	}

	// 7. Route Tables
	utils.Logger.Info("Scanning Route Tables...")
	if rtRes, err := services.NewRouteTableService(s.Registry.EC2).Scan(ctx, tagFilter); err == nil {
		inventory.Resources = append(inventory.Resources, rtRes...)
	} else {
		utils.Logger.Warn("Skipped Route Tables scan", slog.Any("error", err))
	}

	// 8. Elastic IPs
	utils.Logger.Info("Scanning Elastic IPs...")
	if eipRes, err := services.NewElasticIPService(s.Registry.EC2).Scan(ctx, tagFilter); err == nil {
		inventory.Resources = append(inventory.Resources, eipRes...)
	} else {
		utils.Logger.Warn("Skipped Elastic IPs scan", slog.Any("error", err))
	}

	// 9. ALB
	utils.Logger.Info("Scanning Application Load Balancers...")
	if albRes, err := services.NewALBService(s.Registry.ELB).Scan(ctx, tagFilter); err == nil {
		inventory.Resources = append(inventory.Resources, albRes...)
	} else {
		utils.Logger.Warn("Skipped Application Load Balancers scan", slog.Any("error", err))
	}

	// 10. Target Groups
	utils.Logger.Info("Scanning Target Groups...")
	if tgRes, err := services.NewTargetGroupService(s.Registry.ELB).Scan(ctx, tagFilter); err == nil {
		inventory.Resources = append(inventory.Resources, tgRes...)
	} else {
		utils.Logger.Warn("Skipped Target Groups scan", slog.Any("error", err))
	}

	// 11. ECS
	utils.Logger.Info("Scanning ECS Clusters...")
	if ecsRes, err := services.NewECSService(s.Registry.ECS).Scan(ctx, tagFilter); err == nil {
		inventory.Resources = append(inventory.Resources, ecsRes...)
	} else {
		utils.Logger.Warn("Skipped ECS Clusters scan", slog.Any("error", err))
	}

	// 12. ECR
	utils.Logger.Info("Scanning ECR Repositories...")
	if ecrRes, err := services.NewECRService(s.Registry.ECR).Scan(ctx, tagFilter); err == nil {
		inventory.Resources = append(inventory.Resources, ecrRes...)
	} else {
		utils.Logger.Warn("Skipped ECR Repositories scan", slog.Any("error", err))
	}

	// 13. RDS
	utils.Logger.Info("Scanning RDS DB Instances...")
	if rdsRes, err := services.NewRDSService(s.Registry.RDS).Scan(ctx, tagFilter); err == nil {
		inventory.Resources = append(inventory.Resources, rdsRes...)
	} else {
		utils.Logger.Warn("Skipped RDS DB Instances scan", slog.Any("error", err))
	}

	// 14. S3
	utils.Logger.Info("Scanning S3 Buckets...")
	if s3Res, err := services.NewS3Service(s.Registry.S3).Scan(ctx, tagFilter); err == nil {
		inventory.Resources = append(inventory.Resources, s3Res...)
	} else {
		utils.Logger.Warn("Skipped S3 Buckets scan", slog.Any("error", err))
	}

	// 15. CloudFront
	utils.Logger.Info("Scanning CloudFront Distributions...")
	if cfRes, err := services.NewCloudFrontService(s.Registry.CloudFront).Scan(ctx, tagFilter); err == nil {
		inventory.Resources = append(inventory.Resources, cfRes...)
	} else {
		utils.Logger.Warn("Skipped CloudFront Distributions scan", slog.Any("error", err))
	}

	// Inject the scanning region into all found resources
	for i := range inventory.Resources {
		inventory.Resources[i].Region = s.Config.Region
	}

	fmt.Printf("Scan discovered %d total active resources.\n", len(inventory.Resources))
	return inventory, nil
}
