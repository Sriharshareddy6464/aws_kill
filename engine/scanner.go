package engine

import (
	"context"

	"github.com/Sriharshareddy6464/aws-kill/models"
)

type Scanner struct {
	// AWS clients/sessions placeholder
}

func NewScanner() *Scanner {
	return &Scanner{}
}

// Scan queries AWS for active resources matching filter criteria (e.g. tag)
func (s *Scanner) Scan(ctx context.Context, tagFilter string) (*models.Inventory, error) {
	inventory := &models.Inventory{
		Resources: make([]models.Resource, 0),
	}

	// Placeholder scanning logic
	// e.g. iterate over registered services, query AWS, assemble resource definitions

	return inventory, nil
}
