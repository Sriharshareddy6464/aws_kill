package aws

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type ClientRegistry struct {
	EC2        *ec2.Client
	S3         *s3.Client
	RDS        *rds.Client
	ECS        *ecs.Client
	ECR        *ecr.Client
	ELB        *elasticloadbalancingv2.Client
	CloudFront *cloudfront.Client
}

func NewClientRegistry(cfg aws.Config) *ClientRegistry {
	return &ClientRegistry{
		EC2:        ec2.NewFromConfig(cfg),
		S3:         s3.NewFromConfig(cfg),
		RDS:        rds.NewFromConfig(cfg),
		ECS:        ecs.NewFromConfig(cfg),
		ECR:        ecr.NewFromConfig(cfg),
		ELB:        elasticloadbalancingv2.NewFromConfig(cfg),
		CloudFront: cloudfront.NewFromConfig(cfg),
	}
}
