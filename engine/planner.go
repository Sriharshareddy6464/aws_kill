package engine

import (
	"context"

	"github.com/Sriharshareddy6464/aws-kill/models"
)

type Planner struct{}

func NewPlanner() *Planner {
	return &Planner{}
}

// Plan builds a dependency graph and calculates deletion ordering
func (p *Planner) Plan(ctx context.Context, inventory *models.Inventory) (*models.Plan, error) {
	plan := &models.Plan{
		Steps: make([]models.Resource, 0),
	}

	// Placeholder graph processing and sorting logic (topological sort)
	// plan.Steps = reverseTopologicalSort(inventory.Resources)

	return plan, nil
}
