package services

import (
	"context"
)

type TargetGroupService struct{}

func NewTargetGroupService() *TargetGroupService {
	return &TargetGroupService{}
}

func (s *TargetGroupService) Delete(ctx context.Context, id string) error {
	return nil
}
