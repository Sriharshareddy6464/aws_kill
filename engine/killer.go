package engine

import (
	"context"

	"github.com/Sriharshareddy6464/aws-kill/models"
)

type Killer struct{}

func NewKiller() *Killer {
	return &Killer{}
}

// Kill processes deletion steps, retrying failures and waiting for async deletion states
func (k *Killer) Kill(ctx context.Context, plan *models.Plan) (*models.Result, error) {
	result := &models.Result{
		DeletedResources: make([]models.Resource, 0),
		FailedResources:  make([]models.Resource, 0),
	}

	// Placeholder logic to run deletion on services sequentially, wait, and retry
	// for _, step := range plan.Steps {
	//    err := step.Delete(ctx)
	//    ...
	// }

	return result, nil
}
