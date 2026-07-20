package services

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/Sriharshareddy6464/aws-kill/models"
)

type ECSService struct {
	Client *ecs.Client
}

func NewECSService(client *ecs.Client) *ECSService {
	return &ECSService{Client: client}
}

func (s *ECSService) Scan(ctx context.Context, tagFilter string) ([]models.Resource, error) {
	var resources []models.Resource
	listInput := &ecs.ListClustersInput{}
	listOutput, err := s.Client.ListClusters(ctx, listInput)
	if err != nil {
		return nil, err
	}

	if len(listOutput.ClusterArns) == 0 {
		return resources, nil
	}

	describeOutput, err := s.Client.DescribeClusters(ctx, &ecs.DescribeClustersInput{
		Clusters: listOutput.ClusterArns,
	})
	if err != nil {
		return nil, err
	}

	for _, c := range describeOutput.Clusters {
		// Skip INACTIVE clusters
		if c.Status != nil && *c.Status == "INACTIVE" {
			continue
		}

		resources = append(resources, models.Resource{
			ID:     *c.ClusterArn,
			Name:   *c.ClusterName,
			Type:   "ECS",
			Region: "",
		})
	}
	return resources, nil
}

func (s *ECSService) Delete(ctx context.Context, id string) error {
	_, err := s.Client.DeleteCluster(ctx, &ecs.DeleteClusterInput{
		Cluster: &id,
	})
	return err
}
