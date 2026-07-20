package services

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/Sriharshareddy6464/aws-kill/models"
)

type EC2Service struct {
	Client *ec2.Client
}

func NewEC2Service(client *ec2.Client) *EC2Service {
	return &EC2Service{Client: client}
}

func (s *EC2Service) Scan(ctx context.Context, tagFilter string) ([]models.Resource, error) {
	var resources []models.Resource
	input := &ec2.DescribeInstancesInput{}
	result, err := s.Client.DescribeInstances(ctx, input)
	if err != nil {
		return nil, err
	}

	for _, reservation := range result.Reservations {
		for _, inst := range reservation.Instances {
			// Skip terminated instances
			if inst.State != nil && inst.State.Name == "terminated" {
				continue
			}

			tags := make(map[string]string)
			for _, t := range inst.Tags {
				tags[*t.Key] = *t.Value
			}

			// Add filter check here if needed
			// ...

			resources = append(resources, models.Resource{
				ID:     *inst.InstanceId,
				Name:   tags["Name"],
				Type:   "EC2 Instances",
				Region: "", // Filled in by caller
				Tags:   tags,
			})
		}
	}
	return resources, nil
}

func (s *EC2Service) Delete(ctx context.Context, id string) error {
	_, err := s.Client.TerminateInstances(ctx, &ec2.TerminateInstancesInput{
		InstanceIds: []string{id},
	})
	return err
}
