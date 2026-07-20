package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type Config struct {
	Profile string
	Region  string
}

// NewSession initializes connection settings to AWS
func NewSession(ctx context.Context, cfg Config) (aws.Config, error) {
	var opts []func(*config.LoadOptions) error
	if cfg.Profile != "" {
		opts = append(opts, config.WithSharedConfigProfile(cfg.Profile))
	}
	if cfg.Region != "" {
		opts = append(opts, config.WithRegion(cfg.Region))
	}
	return config.LoadDefaultConfig(ctx, opts...)
}
