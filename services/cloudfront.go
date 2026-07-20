package services

import (
	"context"
)

type CloudFrontService struct{}

func NewCloudFrontService() *CloudFrontService {
	return &CloudFrontService{}
}

func (s *CloudFrontService) Delete(ctx context.Context, id string) error {
	return nil
}
