package astar

import (
	"container/heap"
	"math"
)

type Node struct {
	ID  int64
	Lat float64
	Lon float64
	Adj []int64
}

type Graph struct {
	Nodes map[int64]*Node
}

type Item struct {
	NodeID   int64
	Priority float64
	Path     []int64
	Index    int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Priority < pq[j].Priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.Index = -1
	*pq = old[0 : n-1]
	return item
}

func heuristic(a, b *Node) float64 {
	dx := a.Lat - b.Lat
	dy := a.Lon - b.Lon
	return math.Sqrt(dx*dx + dy*dy)
}

func Algorithm(graph *Graph, startID, goalID int64) []int64 {
	open := make(PriorityQueue, 0)
	heap.Init(&open)
	heap.Push(&open, &Item{NodeID: startID, Priority: 0, Path: []int64{startID}})
	visited := make(map[int64]float64)

	for open.Len() > 0 {
		item := heap.Pop(&open).(*Item)
		current := graph.Nodes[item.NodeID]

		if item.NodeID == goalID {
			return item.Path
		}

		for _, neighborID := range current.Adj {
			neighbor := graph.Nodes[neighborID]
			newCost := float64(len(item.Path))
			if cost, seen := visited[neighborID]; !seen || newCost < cost {
				visited[neighborID] = newCost
				priority := newCost + heuristic(neighbor, graph.Nodes[goalID])
				newPath := append([]int64{}, item.Path...)
				newPath = append(newPath, neighborID)
				heap.Push(&open, &Item{NodeID: neighborID, Priority: priority, Path: newPath})
			}
		}
	}
	return nil
}
