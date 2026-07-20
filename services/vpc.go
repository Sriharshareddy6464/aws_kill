package services

import (
	"context"
)

type VPCService struct{}

func NewVPCService() *VPCService {
	return &VPCService{}
}

func (s *VPCService) Delete(ctx context.Context, id string) error {
	return nil
}
