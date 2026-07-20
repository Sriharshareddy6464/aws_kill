package services

import (
	"context"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/Sriharshareddy6464/aws-kill/models"
	"github.com/Sriharshareddy6464/aws-kill/utils"
)

type ECRService struct {
	Client *ecr.Client
}

func NewECRService(client *ecr.Client) *ECRService {
	return &ECRService{Client: client}
}

func (s *ECRService) Scan(ctx context.Context, tagFilter string) ([]models.Resource, map[string]int, error) {
	var resources []models.Resource
	counts := map[string]int{
		"Repositories": 0,
		"Images":       0,
	}

	input := &ecr.DescribeRepositoriesInput{}
	result, err := s.Client.DescribeRepositories(ctx, input)
	if err != nil {
		return nil, nil, err
	}

	for _, repo := range result.Repositories {
		counts["Repositories"]++
		resources = append(resources, models.Resource{
			ID:     *repo.RepositoryArn,
			Name:   *repo.RepositoryName,
			Type:   "ECR",
			Region: "",
		})

		// List images for this repository
		imgOutput, err := s.Client.ListImages(ctx, &ecr.ListImagesInput{
			RepositoryName: repo.RepositoryName,
		})
		if err == nil {
			counts["Images"] += len(imgOutput.ImageIds)
		} else {
			utils.Logger.Warn("Skipped ECR Images describe for repository "+*repo.RepositoryName, slog.Any("error", err))
		}
	}
	return resources, counts, nil
}

func (s *ECRService) Delete(ctx context.Context, id string) error {
	name := id
	_, err := s.Client.DeleteRepository(ctx, &ecr.DeleteRepositoryInput{
		RepositoryName: &name,
		Force:          true,
	})
	return err
}
