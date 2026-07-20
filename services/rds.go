package services

import (
	"context"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/Sriharshareddy6464/aws-kill/models"
	"github.com/Sriharshareddy6464/aws-kill/utils"
)

type RDSService struct {
	Client *rds.Client
}

func NewRDSService(client *rds.Client) *RDSService {
	return &RDSService{Client: client}
}

func (s *RDSService) Scan(ctx context.Context, tagFilter string) ([]models.Resource, map[string]int, error) {
	var resources []models.Resource
	counts := map[string]int{
		"DB Instances":  0,
		"DB Snapshots":  0,
		"Subnet Groups": 0,
	}

	// 1. DB Instances
	result, err := s.Client.DescribeDBInstances(ctx, &rds.DescribeDBInstancesInput{})
	if err == nil {
		for _, db := range result.DBInstances {
			if db.DBInstanceStatus != nil && (*db.DBInstanceStatus == "deleting" || *db.DBInstanceStatus == "deleted") {
				continue
			}

			counts["DB Instances"]++
			resources = append(resources, models.Resource{
				ID:     *db.DBInstanceIdentifier,
				Name:   *db.DBInstanceIdentifier,
				Type:   "RDS",
				Region: "",
			})
		}
	} else {
		return nil, nil, err
	}

	// 2. DB Snapshots
	snapOutput, err := s.Client.DescribeDBSnapshots(ctx, &rds.DescribeDBSnapshotsInput{})
	if err == nil {
		counts["DB Snapshots"] = len(snapOutput.DBSnapshots)
	} else {
		utils.Logger.Warn("Skipped RDS DB Snapshots describe", slog.Any("error", err))
	}

	// 3. Subnet Groups
	sngOutput, err := s.Client.DescribeDBSubnetGroups(ctx, &rds.DescribeDBSubnetGroupsInput{})
	if err == nil {
		counts["Subnet Groups"] = len(sngOutput.DBSubnetGroups)
	} else {
		utils.Logger.Warn("Skipped RDS DB Subnet Groups describe", slog.Any("error", err))
	}

	return resources, counts, nil
}

func (s *RDSService) Delete(ctx context.Context, id string) error {
	_, err := s.Client.DeleteDBInstance(ctx, &rds.DeleteDBInstanceInput{
		DBInstanceIdentifier: &id,
		SkipFinalSnapshot:    aws.Bool(true),
	})
	return err
}
