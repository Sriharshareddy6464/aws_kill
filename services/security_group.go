package services

import (
	"context"
)

type SecurityGroupService struct{}

func NewSecurityGroupService() *SecurityGroupService {
	return &SecurityGroupService{}
}

func (s *SecurityGroupService) Delete(ctx context.Context, id string) error {
	return nil
}
