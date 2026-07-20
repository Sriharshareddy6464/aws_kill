package aws

import (
	"context"
)

type Config struct {
	Profile string
	Region  string
}

// NewSession initializes connection settings to AWS (placeholder for config.LoadDefaultConfig)
func NewSession(ctx context.Context, cfg Config) (interface{}, error) {
	// Placeholder return for aws.Config
	return nil, nil
}
