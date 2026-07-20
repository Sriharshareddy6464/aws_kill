package services

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/Sriharshareddy6464/aws-kill/models"
)

type ALBService struct {
	Client *elasticloadbalancingv2.Client
}

func NewALBService(client *elasticloadbalancingv2.Client) *ALBService {
	return &ALBService{Client: client}
}

func (s *ALBService) Scan(ctx context.Context, tagFilter string) ([]models.Resource, map[string]int, error) {
	var resources []models.Resource
	counts := map[string]int{"Load Balancers": 0}
	input := &elasticloadbalancingv2.DescribeLoadBalancersInput{}
	result, err := s.Client.DescribeLoadBalancers(ctx, input)
	if err != nil {
		return nil, nil, err
	}

	for _, lb := range result.LoadBalancers {
		counts["Load Balancers"]++
		resources = append(resources, models.Resource{
			ID:     *lb.LoadBalancerArn,
			Name:   *lb.LoadBalancerName,
			Type:   "Application Load Balancer",
			Region: "",
			Tags:   nil,
		})
	}
	return resources, counts, nil
}

func (s *ALBService) Delete(ctx context.Context, id string) error {
	_, err := s.Client.DeleteLoadBalancer(ctx, &elasticloadbalancingv2.DeleteLoadBalancerInput{
		LoadBalancerArn: &id,
	})
	return err
}
