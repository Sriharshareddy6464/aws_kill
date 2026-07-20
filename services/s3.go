package services

import (
	"context"
)

type S3Service struct{}

func NewS3Service() *S3Service {
	return &S3Service{}
}

func (s *S3Service) Delete(ctx context.Context, id string) error {
	return nil
}
