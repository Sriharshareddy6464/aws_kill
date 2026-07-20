package services

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/Sriharshareddy6464/aws-kill/models"
)

type RDSService struct {
	Client *rds.Client
}

func NewRDSService(client *rds.Client) *RDSService {
	return &RDSService{Client: client}
}

func (s *RDSService) Scan(ctx context.Context, tagFilter string) ([]models.Resource, error) {
	var resources []models.Resource
	input := &rds.DescribeDBInstancesInput{}
	result, err := s.Client.DescribeDBInstances(ctx, input)
	if err != nil {
		return nil, err
	}

	for _, db := range result.DBInstances {
		// Skip deleted/deleting DB instances
		if db.DBInstanceStatus != nil && (*db.DBInstanceStatus == "deleting" || *db.DBInstanceStatus == "deleted") {
			continue
		}

		resources = append(resources, models.Resource{
			ID:     *db.DBInstanceIdentifier,
			Name:   *db.DBInstanceIdentifier,
			Type:   "RDS",
			Region: "",
		})
	}
	return resources, nil
}

func (s *RDSService) Delete(ctx context.Context, id string) error {
	_, err := s.Client.DeleteDBInstance(ctx, &rds.DeleteDBInstanceInput{
		DBInstanceIdentifier: &id,
		SkipFinalSnapshot:    aws.Bool(true),
	})
	return err
}

