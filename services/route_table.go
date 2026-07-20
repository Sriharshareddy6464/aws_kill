package services

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/Sriharshareddy6464/aws-kill/models"
)

type RouteTableService struct {
	Client *ec2.Client
}

func NewRouteTableService(client *ec2.Client) *RouteTableService {
	return &RouteTableService{Client: client}
}

func (s *RouteTableService) Scan(ctx context.Context, tagFilter string) ([]models.Resource, error) {
	var resources []models.Resource
	input := &ec2.DescribeRouteTablesInput{}
	result, err := s.Client.DescribeRouteTables(ctx, input)
	if err != nil {
		return nil, err
	}

	for _, rt := range result.RouteTables {
		// Skip main route tables since they are managed by VPC lifecycle
		isMain := false
		for _, assoc := range rt.Associations {
			if assoc.Main != nil && *assoc.Main {
				isMain = true
				break
			}
		}
		if isMain {
			continue
		}

		tags := make(map[string]string)
		for _, t := range rt.Tags {
			tags[*t.Key] = *t.Value
		}

		resources = append(resources, models.Resource{
			ID:           *rt.RouteTableId,
			Name:         tags["Name"],
			Type:         "Route Tables",
			Region:       "",
			Dependencies: []string{*rt.VpcId},
			Tags:         tags,
		})
	}
	return resources, nil
}

func (s *RouteTableService) Delete(ctx context.Context, id string) error {
	_, err := s.Client.DeleteRouteTable(ctx, &ec2.DeleteRouteTableInput{
		RouteTableId: &id,
	})
	return err
}
