package services

import (
	"context"
)

type NATGatewayService struct{}

func NewNATGatewayService() *NATGatewayService {
	return &NATGatewayService{}
}

func (s *NATGatewayService) Delete(ctx context.Context, id string) error {
	return nil
}
