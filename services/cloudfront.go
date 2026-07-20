package services

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/Sriharshareddy6464/aws-kill/models"
)

type CloudFrontService struct {
	Client *cloudfront.Client
}

func NewCloudFrontService(client *cloudfront.Client) *CloudFrontService {
	return &CloudFrontService{Client: client}
}

func (s *CloudFrontService) Scan(ctx context.Context, tagFilter string) ([]models.Resource, map[string]int, error) {
	var resources []models.Resource
	counts := map[string]int{"Distributions": 0}
	input := &cloudfront.ListDistributionsInput{}
	result, err := s.Client.ListDistributions(ctx, input)
	if err != nil {
		return nil, nil, err
	}

	if result.DistributionList == nil {
		return resources, counts, nil
	}

	for _, dist := range result.DistributionList.Items {
		counts["Distributions"]++
		resources = append(resources, models.Resource{
			ID:     *dist.Id,
			Name:   *dist.DomainName,
			Type:   "CloudFront",
			Region: "",
		})
	}
	return resources, counts, nil
}

func (s *CloudFrontService) Delete(ctx context.Context, id string) error {
	_, err := s.Client.DeleteDistribution(ctx, &cloudfront.DeleteDistributionInput{
		Id: &id,
	})
	return err
}
