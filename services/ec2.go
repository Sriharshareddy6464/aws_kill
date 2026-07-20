package services

import (
	"context"
)

type EC2Service struct{}

func NewEC2Service() *EC2Service {
	return &EC2Service{}
}

func (s *EC2Service) Delete(ctx context.Context, id string) error {
	return nil
}
