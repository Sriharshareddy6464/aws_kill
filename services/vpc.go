package services

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/Sriharshareddy6464/aws-kill/models"
)

type VPCService struct {
	Client *ec2.Client
}

func NewVPCService(client *ec2.Client) *VPCService {
	return &VPCService{Client: client}
}

func (s *VPCService) Scan(ctx context.Context, tagFilter string) ([]models.Resource, map[string]int, error) {
	var resources []models.Resource
	counts := map[string]int{"VPCs": 0}
	input := &ec2.DescribeVpcsInput{}
	result, err := s.Client.DescribeVpcs(ctx, input)
	if err != nil {
		return nil, nil, err
	}

	for _, vpc := range result.Vpcs {
		if vpc.IsDefault != nil && *vpc.IsDefault {
			continue
		}

		counts["VPCs"]++
		tags := make(map[string]string)
		for _, t := range vpc.Tags {
			tags[*t.Key] = *t.Value
		}

		resources = append(resources, models.Resource{
			ID:     *vpc.VpcId,
			Name:   tags["Name"],
			Type:   "VPC",
			Region: "",
			Tags:   tags,
		})
	}
	return resources, counts, nil
}

func (s *VPCService) Delete(ctx context.Context, id string) error {
	_, err := s.Client.DeleteVpc(ctx, &ec2.DeleteVpcInput{
		VpcId: &id,
	})
	return err
}
