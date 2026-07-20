package services

import (
	"context"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/Sriharshareddy6464/aws-kill/models"
	"github.com/Sriharshareddy6464/aws-kill/utils"
)

type ECSService struct {
	Client *ecs.Client
}

func NewECSService(client *ecs.Client) *ECSService {
	return &ECSService{Client: client}
}

func (s *ECSService) Scan(ctx context.Context, tagFilter string) ([]models.Resource, map[string]int, error) {
	var resources []models.Resource
	counts := map[string]int{
		"Clusters":         0,
		"Services":         0,
		"Task Definitions": 0,
		"Running Tasks":    0,
	}

	// 1. List Task Definitions
	tdOutput, err := s.Client.ListTaskDefinitions(ctx, &ecs.ListTaskDefinitionsInput{})
	if err == nil {
		counts["Task Definitions"] = len(tdOutput.TaskDefinitionArns)
	} else {
		utils.Logger.Warn("Skipped ECS Task Definitions list", slog.Any("error", err))
	}

	// 2. List Clusters
	listInput := &ecs.ListClustersInput{}
	listOutput, err := s.Client.ListClusters(ctx, listInput)
	if err != nil {
		return nil, nil, err
	}

	if len(listOutput.ClusterArns) == 0 {
		return resources, counts, nil
	}

	describeOutput, err := s.Client.DescribeClusters(ctx, &ecs.DescribeClustersInput{
		Clusters: listOutput.ClusterArns,
	})
	if err != nil {
		return nil, nil, err
	}

	for _, c := range describeOutput.Clusters {
		if c.Status != nil && *c.Status == "INACTIVE" {
			continue
		}

		counts["Clusters"]++
		resources = append(resources, models.Resource{
			ID:     *c.ClusterArn,
			Name:   *c.ClusterName,
			Type:   "ECS",
			Region: "",
		})

		// 3. For each active cluster, find Services
		svcOutput, err := s.Client.ListServices(ctx, &ecs.ListServicesInput{Cluster: c.ClusterArn})
		if err == nil {
			counts["Services"] += len(svcOutput.ServiceArns)
		} else {
			utils.Logger.Warn("Skipped ECS Services list for cluster "+*c.ClusterName, slog.Any("error", err))
		}

		// 4. For each active cluster, find Running Tasks
		taskOutput, err := s.Client.ListTasks(ctx, &ecs.ListTasksInput{
			Cluster:       c.ClusterArn,
			DesiredStatus: "RUNNING",
		})
		if err == nil {
			counts["Running Tasks"] += len(taskOutput.TaskArns)
		} else {
			utils.Logger.Warn("Skipped ECS Tasks list for cluster "+*c.ClusterName, slog.Any("error", err))
		}
	}

	return resources, counts, nil
}

func (s *ECSService) Delete(ctx context.Context, id string) error {
	_, err := s.Client.DeleteCluster(ctx, &ecs.DeleteClusterInput{
		Cluster: &id,
	})
	return err
}
