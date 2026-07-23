package engine

import (
	"context"
	"testing"

	"github.com/Sriharshareddy6464/aws-kill/models"
)

func TestPlan_LinearDependencies(t *testing.T) {
	inventory := &models.Inventory{
		Resources: []models.Resource{
			{ID: "vpc-1", Type: "VPC"},
			{ID: "subnet-1", Type: "Subnets", Dependencies: []string{"vpc-1"}},
			{ID: "ec2-1", Type: "EC2 Instances", Dependencies: []string{"subnet-1"}},
		},
	}

	planner := NewPlanner()
	plan, err := planner.Plan(context.Background(), inventory)
	if err != nil {
		t.Fatalf("Plan failed: %v", err)
	}

	if len(plan.Steps) != 3 {
		t.Fatalf("Expected 3 steps, got %d", len(plan.Steps))
	}

	// Deletion order must be: ec2-1 -> subnet-1 -> vpc-1
	expected := []string{"ec2-1", "subnet-1", "vpc-1"}
	for i, id := range expected {
		if plan.Steps[i].ID != id {
			t.Errorf("At step %d: expected %s, got %s", i, id, plan.Steps[i].ID)
		}
	}
}

func TestPlan_TierPriorityFallback(t *testing.T) {
	inventory := &models.Inventory{
		Resources: []models.Resource{
			{ID: "s3-bucket-1", Type: "S3"},
			{ID: "cf-dist-1", Type: "CloudFront"},
			{ID: "vpc-1", Type: "VPC"},
		},
	}

	planner := NewPlanner()
	plan, err := planner.Plan(context.Background(), inventory)
	if err != nil {
		t.Fatalf("Plan failed: %v", err)
	}

	// Both S3, CloudFront, and VPC have in-degree 0.
	// CloudFront (tier 10) must be deleted before VPC (tier 160), which must be deleted before S3 (tier 170).
	expected := []string{"cf-dist-1", "vpc-1", "s3-bucket-1"}
	for i, id := range expected {
		if plan.Steps[i].ID != id {
			t.Errorf("At step %d: expected %s, got %s", i, id, plan.Steps[i].ID)
		}
	}
}

func TestPlan_CycleDetectionAndBreaking(t *testing.T) {
	inventory := &models.Inventory{
		Resources: []models.Resource{
			{ID: "node-a", Type: "EC2 Instances", Dependencies: []string{"node-b"}},
			{ID: "node-b", Type: "Subnets", Dependencies: []string{"node-a"}},
		},
	}

	planner := NewPlanner()
	plan, err := planner.Plan(context.Background(), inventory)
	if err != nil {
		t.Fatalf("Plan failed: %v", err)
	}

	// The cycle must be broken and all nodes should be planned.
	if len(plan.Steps) != 2 {
		t.Fatalf("Expected 2 steps, got %d", len(plan.Steps))
	}
}
