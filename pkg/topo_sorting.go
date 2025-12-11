package pkg

import (
	"fmt"

	"github.com/Thanus-Kumaar/anokha-2025-backend/models"
	"github.com/google/uuid"
)

func TopoSort(dependencies []models.EventDependencyRequest) ([]uuid.UUID, error) {
	// Build adjacency list
	adj := make(map[uuid.UUID][]uuid.UUID)
	inDegree := make(map[uuid.UUID]int)

	nodes := make(map[uuid.UUID]bool)

	for _, dep := range dependencies {
		adj[dep.StartEvent] = append(adj[dep.StartEvent], dep.EndEvent)
		inDegree[dep.EndEvent]++
		nodes[dep.StartEvent] = true
		nodes[dep.EndEvent] = true
	}

	// Queue for nodes with 0 in-degree
	queue := make([]uuid.UUID, 0)
	for node := range nodes {
		if inDegree[node] == 0 {
			queue = append(queue, node)
		}
	}

	var sorted []uuid.UUID

	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		sorted = append(sorted, n)

		for _, m := range adj[n] {
			inDegree[m]--
			if inDegree[m] == 0 {
				queue = append(queue, m)
			}
		}
	}

	if len(sorted) != len(nodes) {
		return nil, fmt.Errorf("cycle detected in event dependencies")
	}

	return sorted, nil
}
