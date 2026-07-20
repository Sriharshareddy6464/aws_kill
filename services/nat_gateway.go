package services

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/Sriharshareddy6464/aws-kill/models"
)

type NATGatewayService struct {
	Client *ec2.Client
}

func NewNATGatewayService(client *ec2.Client) *NATGatewayService {
	return &NATGatewayService{Client: client}
}

func (s *NATGatewayService) Scan(ctx context.Context, tagFilter string) ([]models.Resource, map[string]int, error) {
	var resources []models.Resource
	counts := map[string]int{"NAT Gateways": 0}
	input := &ec2.DescribeNatGatewaysInput{}
	result, err := s.Client.DescribeNatGateways(ctx, input)
	if err != nil {
		return nil, nil, err
	}

	for _, ngw := range result.NatGateways {
		if ngw.State == "deleted" || ngw.State == "deleting" {
			continue
		}

		counts["NAT Gateways"]++
		tags := make(map[string]string)
		for _, t := range ngw.Tags {
			tags[*t.Key] = *t.Value
		}

		var deps []string
		if ngw.SubnetId != nil {
			deps = append(deps, *ngw.SubnetId)
		}

		resources = append(resources, models.Resource{
			ID:           *ngw.NatGatewayId,
			Name:         tags["Name"],
			Type:         "NAT Gateway",
			Region:       "",
			Dependencies: deps,
			Tags:         tags,
		})
	}
	return resources, counts, nil
}

func (s *NATGatewayService) Delete(ctx context.Context, id string) error {
	_, err := s.Client.DeleteNatGateway(ctx, &ec2.DeleteNatGatewayInput{
		NatGatewayId: &id,
	})
	return err
}
