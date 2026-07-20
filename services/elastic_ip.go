package services

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/Sriharshareddy6464/aws-kill/models"
)

type ElasticIPService struct {
	Client *ec2.Client
}

func NewElasticIPService(client *ec2.Client) *ElasticIPService {
	return &ElasticIPService{Client: client}
}

func (s *ElasticIPService) Scan(ctx context.Context, tagFilter string) ([]models.Resource, map[string]int, error) {
	var resources []models.Resource
	counts := map[string]int{"Elastic IPs": 0}
	input := &ec2.DescribeAddressesInput{}
	result, err := s.Client.DescribeAddresses(ctx, input)
	if err != nil {
		return nil, nil, err
	}

	for _, addr := range result.Addresses {
		counts["Elastic IPs"]++
		tags := make(map[string]string)
		for _, t := range addr.Tags {
			tags[*t.Key] = *t.Value
		}

		resources = append(resources, models.Resource{
			ID:     *addr.AllocationId,
			Name:   *addr.PublicIp,
			Type:   "Elastic IP",
			Region: "",
			Tags:   tags,
		})
	}
	return resources, counts, nil
}

func (s *ElasticIPService) Delete(ctx context.Context, id string) error {
	_, err := s.Client.ReleaseAddress(ctx, &ec2.ReleaseAddressInput{
		AllocationId: &id,
	})
	return err
}
