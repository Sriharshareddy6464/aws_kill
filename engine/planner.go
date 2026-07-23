package engine

import (
	"context"
	"fmt"
	"log/slog"
	"sort"

	"github.com/Sriharshareddy6464/aws-kill/models"
	"github.com/Sriharshareddy6464/aws-kill/utils"
)

// Deletion Tiers (lower value = deleted first)
var typePriorities = map[string]int{
	"CloudFront":                10,
	"ECS":                       20,
	"EC2 Instances":             30,
	"Application Load Balancer": 40,
	"RDS":                       50,
	"NAT Gateway":               60,
	"NetworkInterface":          70,
	"Volume":                    80,
	"Snapshot":                  85,
	"Target Groups":             90,
	"Elastic IP":                100,
	"KeyPair":                   110,
	"LaunchTemplate":            115,
	"PlacementGroup":            115,
	"DedicatedHost":             115,
	"CapacityReservation":       115,
	"Security Groups":           120,
	"Subnets":                   130,
	"Route Tables":              140,
	"Internet Gateway":          150,
	"VPC":                       160,
	"S3":                        170,
}

func getTypePriority(t string) int {
	if p, ok := typePriorities[t]; ok {
		return p
	}
	return 80 // Default priority for unknown types
}

type Planner struct{}

func NewPlanner() *Planner {
	return &Planner{}
}

// Plan builds a dependency graph and calculates a safe deletion order using topological sort.
func (p *Planner) Plan(ctx context.Context, inventory *models.Inventory) (*models.Plan, error) {
	if inventory == nil || len(inventory.Resources) == 0 {
		return &models.Plan{Steps: []models.Resource{}}, nil
	}

	// 1. Initialize node lookup
	nodes := make(map[string]models.Resource)
	for _, res := range inventory.Resources {
		nodes[res.ID] = res
	}

	// 2. Build graph and calculate initial in-degrees
	adjList := make(map[string][]string)
	inDegree := make(map[string]int)

	// Initialize in-degree for all nodes
	for id := range nodes {
		inDegree[id] = 0
	}

	for _, res := range inventory.Resources {
		for _, depID := range res.Dependencies {
			// Only consider dependencies that are actually in the inventory
			if _, exists := nodes[depID]; exists {
				adjList[res.ID] = append(adjList[res.ID], depID)
				inDegree[depID]++
			}
		}
	}

	// 3. Cycle Detection & Breaking using DFS
	// If a cycle is detected, we log a warning and break the cycle by removing the edge.
	if hasCycle(nodes, adjList) {
		utils.Logger.Warn("Dependency cycle detected! Attempting to break cycles to ensure deletion continues.")
		breakCycles(nodes, adjList, inDegree)
	}

	// 4. Kahn's Algorithm with Type-Based Priority
	// We want to process nodes with in-degree 0.
	// To enforce implicit tier-based ordering, we always pop the node with the lowest type priority.
	var zeroInDegree []string
	for id, degree := range inDegree {
		if degree == 0 {
			zeroInDegree = append(zeroInDegree, id)
		}
	}

	var orderedSteps []models.Resource

	for len(zeroInDegree) > 0 {
		// Sort the zero-in-degree nodes:
		// 1. Primary: Type priority (tier) ascending (lower value = deleted first)
		// 2. Secondary: Resource ID alphabetical (for deterministic sorting)
		sort.Slice(zeroInDegree, func(i, j int) bool {
			pi := getTypePriority(nodes[zeroInDegree[i]].Type)
			pj := getTypePriority(nodes[zeroInDegree[j]].Type)
			if pi != pj {
				return pi < pj
			}
			return zeroInDegree[i] < zeroInDegree[j]
		})

		// Pop the highest priority node
		u := zeroInDegree[0]
		zeroInDegree = zeroInDegree[1:]

		orderedSteps = append(orderedSteps, nodes[u])

		// Decrement in-degree for neighbors
		for _, v := range adjList[u] {
			inDegree[v]--
			if inDegree[v] == 0 {
				zeroInDegree = append(zeroInDegree, v)
			}
		}
	}

	// If the output doesn't contain all nodes, there is an unresolved cycle
	if len(orderedSteps) < len(nodes) {
		return nil, fmt.Errorf("failed to plan deletion sequence: unresolved cycles in dependency graph (%d of %d planned)", len(orderedSteps), len(nodes))
	}

	return &models.Plan{
		Steps: orderedSteps,
	}, nil
}

// DFS cycle check helper
func hasCycle(nodes map[string]models.Resource, adjList map[string][]string) bool {
	visited := make(map[string]int) // 0 = unvisited, 1 = visiting, 2 = visited

	var dfs func(u string) bool
	dfs = func(u string) bool {
		visited[u] = 1
		for _, v := range adjList[u] {
			if visited[v] == 1 {
				return true
			}
			if visited[v] == 0 {
				if dfs(v) {
					return true
				}
			}
		}
		visited[u] = 2
		return false
	}

	for id := range nodes {
		if visited[id] == 0 {
			if dfs(id) {
				return true
			}
		}
	}
	return false
}

// breakCycles performs a DFS to detect back-edges and removes them, updating in-degrees.
func breakCycles(nodes map[string]models.Resource, adjList map[string][]string, inDegree map[string]int) {
	visited := make(map[string]int) // 0 = unvisited, 1 = visiting, 2 = visited

	var dfs func(u string)
	dfs = func(u string) {
		visited[u] = 1
		var cleanNeighbors []string
		for _, v := range adjList[u] {
			if visited[v] == 1 {
				// Back-edge detected! Break it.
				utils.Logger.Warn("Breaking cyclic dependency edge", slog.String("from", u), slog.String("to", v))
				inDegree[v]--
			} else {
				cleanNeighbors = append(cleanNeighbors, v)
				if visited[v] == 0 {
					dfs(v)
				}
			}
		}
		adjList[u] = cleanNeighbors
		visited[u] = 2
	}

	for id := range nodes {
		if visited[id] == 0 {
			dfs(id)
		}
	}
}
