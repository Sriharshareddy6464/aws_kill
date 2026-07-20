package services

import (
	"context"
)

type ALBService struct{}

func NewALBService() *ALBService {
	return &ALBService{}
}

func (s *ALBService) Delete(ctx context.Context, id string) error {
	return nil
}
