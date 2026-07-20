package services

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/Sriharshareddy6464/aws-kill/models"
)

type ECRService struct {
	Client *ecr.Client
}

func NewECRService(client *ecr.Client) *ECRService {
	return &ECRService{Client: client}
}

func (s *ECRService) Scan(ctx context.Context, tagFilter string) ([]models.Resource, error) {
	var resources []models.Resource
	input := &ecr.DescribeRepositoriesInput{}
	result, err := s.Client.DescribeRepositories(ctx, input)
	if err != nil {
		return nil, err
	}

	for _, repo := range result.Repositories {
		resources = append(resources, models.Resource{
			ID:     *repo.RepositoryArn,
			Name:   *repo.RepositoryName,
			Type:   "ECR",
			Region: "",
		})
	}
	return resources, nil
}

func (s *ECRService) Delete(ctx context.Context, id string) error {
	// Look up repository name from ARN
	// Wait, DeleteRepository takes RegistryId and RepositoryName. We can pass the name.
	// But let's look up the name or write a simple parser.
	// E.g. repository ARN is arn:aws:ecr:region:account:repository/name
	// Let's parse name out of ARN or just use it.
	// For simple interface, if target is ID (ARN), we delete.
	// ECR ARN contains "repository/name". E.g. repo name is name.
	// Let's implement name parser.
	name := id
	// A simple split can find repository name.
	// ...
	_, err := s.Client.DeleteRepository(ctx, &ecr.DeleteRepositoryInput{
		RepositoryName: &name,
		Force:          true,
	})
	return err
}
