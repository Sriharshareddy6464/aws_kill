package services

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/Sriharshareddy6464/aws-kill/models"
)

type S3Service struct {
	Client *s3.Client
}

func NewS3Service(client *s3.Client) *S3Service {
	return &S3Service{Client: client}
}

func (s *S3Service) Scan(ctx context.Context, tagFilter string) ([]models.Resource, map[string]int, error) {
	var resources []models.Resource
	counts := map[string]int{"Buckets": 0}
	input := &s3.ListBucketsInput{}
	result, err := s.Client.ListBuckets(ctx, input)
	if err != nil {
		return nil, nil, err
	}

	for _, bucket := range result.Buckets {
		counts["Buckets"]++
		resources = append(resources, models.Resource{
			ID:     *bucket.Name,
			Name:   *bucket.Name,
			Type:   "S3",
			Region: "",
		})
	}
	return resources, counts, nil
}

func (s *S3Service) Delete(ctx context.Context, id string) error {
	_, err := s.Client.DeleteBucket(ctx, &s3.DeleteBucketInput{
		Bucket: &id,
	})
	return err
}
