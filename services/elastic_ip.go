package services

import (
	"context"
)

type ElasticIPService struct{}

func NewElasticIPService() *ElasticIPService {
	return &ElasticIPService{}
}

func (s *ElasticIPService) Delete(ctx context.Context, id string) error {
	return nil
}
