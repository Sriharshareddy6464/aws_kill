package services

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/Sriharshareddy6464/aws-kill/models"
)

type TargetGroupService struct {
	Client *elasticloadbalancingv2.Client
}

func NewTargetGroupService(client *elasticloadbalancingv2.Client) *TargetGroupService {
	return &TargetGroupService{Client: client}
}

func (s *TargetGroupService) Scan(ctx context.Context, tagFilter string) ([]models.Resource, error) {
	var resources []models.Resource
	input := &elasticloadbalancingv2.DescribeTargetGroupsInput{}
	result, err := s.Client.DescribeTargetGroups(ctx, input)
	if err != nil {
		return nil, err
	}

	for _, tg := range result.TargetGroups {
		resources = append(resources, models.Resource{
			ID:     *tg.TargetGroupArn,
			Name:   *tg.TargetGroupName,
			Type:   "Target Groups",
			Region: "",
		})
	}
	return resources, nil
}

func (s *TargetGroupService) Delete(ctx context.Context, id string) error {
	_, err := s.Client.DeleteTargetGroup(ctx, &elasticloadbalancingv2.DeleteTargetGroupInput{
		TargetGroupArn: &id,
	})
	return err
}
