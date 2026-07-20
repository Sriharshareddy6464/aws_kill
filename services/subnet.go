package services

import (
	"context"
)

type SubnetService struct{}

func NewSubnetService() *SubnetService {
	return &SubnetService{}
}

func (s *SubnetService) Delete(ctx context.Context, id string) error {
	return nil
}
