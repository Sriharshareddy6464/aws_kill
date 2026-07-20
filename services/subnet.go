package services

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/Sriharshareddy6464/aws-kill/models"
)

type SubnetService struct {
	Client *ec2.Client
}

func NewSubnetService(client *ec2.Client) *SubnetService {
	return &SubnetService{Client: client}
}

func (s *SubnetService) Scan(ctx context.Context, tagFilter string) ([]models.Resource, error) {
	var resources []models.Resource
	input := &ec2.DescribeSubnetsInput{}
	result, err := s.Client.DescribeSubnets(ctx, input)
	if err != nil {
		return nil, err
	}

	for _, sub := range result.Subnets {
		// Skip default subnet's VPC parent subnets if needed, or scan all
		if sub.DefaultForAz != nil && *sub.DefaultForAz {
			continue
		}

		tags := make(map[string]string)
		for _, t := range sub.Tags {
			tags[*t.Key] = *t.Value
		}

		resources = append(resources, models.Resource{
			ID:           *sub.SubnetId,
			Name:         tags["Name"],
			Type:         "Subnets",
			Region:       "",
			Dependencies: []string{*sub.VpcId},
			Tags:         tags,
		})
	}
	return resources, nil
}

func (s *SubnetService) Delete(ctx context.Context, id string) error {
	_, err := s.Client.DeleteSubnet(ctx, &ec2.DeleteSubnetInput{
		SubnetId: &id,
	})
	return err
}
