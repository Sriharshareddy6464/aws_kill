package services

import (
	"context"
)

type RDSService struct{}

func NewRDSService() *RDSService {
	return &RDSService{}
}

func (s *RDSService) Delete(ctx context.Context, id string) error {
	return nil
}
