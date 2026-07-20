package services

import (
	"context"
)

type ECRService struct{}

func NewECRService() *ECRService {
	return &ECRService{}
}

func (s *ECRService) Delete(ctx context.Context, id string) error {
	return nil
}
