package services

import (
	"context"
)

type ECSService struct{}

func NewECSService() *ECSService {
	return &ECSService{}
}

func (s *ECSService) Delete(ctx context.Context, id string) error {
	return nil
}
