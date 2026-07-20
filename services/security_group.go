package services

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/Sriharshareddy6464/aws-kill/models"
)

type SecurityGroupService struct {
	Client *ec2.Client
}

func NewSecurityGroupService(client *ec2.Client) *SecurityGroupService {
	return &SecurityGroupService{Client: client}
}

func (s *SecurityGroupService) Scan(ctx context.Context, tagFilter string) ([]models.Resource, map[string]int, error) {
	var resources []models.Resource
	counts := map[string]int{"Security Groups": 0}
	input := &ec2.DescribeSecurityGroupsInput{}
	result, err := s.Client.DescribeSecurityGroups(ctx, input)
	if err != nil {
		return nil, nil, err
	}

	for _, sg := range result.SecurityGroups {
		if sg.GroupName != nil && *sg.GroupName == "default" {
			continue
		}

		counts["Security Groups"]++
		tags := make(map[string]string)
		for _, t := range sg.Tags {
			tags[*t.Key] = *t.Value
		}

		resources = append(resources, models.Resource{
			ID:           *sg.GroupId,
			Name:         *sg.GroupName,
			Type:         "Security Groups",
			Region:       "",
			Dependencies: []string{*sg.VpcId},
			Tags:         tags,
		})
	}
	return resources, counts, nil
}

func (s *SecurityGroupService) Delete(ctx context.Context, id string) error {
	_, err := s.Client.DeleteSecurityGroup(ctx, &ec2.DeleteSecurityGroupInput{
		GroupId: &id,
	})
	return err
}
