package services

import (
	"context"
)

type InternetGatewayService struct{}

func NewInternetGatewayService() *InternetGatewayService {
	return &InternetGatewayService{}
}

func (s *InternetGatewayService) Delete(ctx context.Context, id string) error {
	return nil
}
