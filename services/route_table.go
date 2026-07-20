package services

import (
	"context"
)

type RouteTableService struct{}

func NewRouteTableService() *RouteTableService {
	return &RouteTableService{}
}

func (s *RouteTableService) Delete(ctx context.Context, id string) error {
	return nil
}
