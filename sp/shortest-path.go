package sp

import "container/heap"

type Node struct {
	outgoing, incomming []*edge
}

func NewNode() *Node {
	return &Node{}
}

func (n *Node) AddUndirectedEdge(other *Node, distance int) {
	edge := edge{
		from:   n,
		to:     other,
		length: distance,
	}
	n.outgoing = append(n.outgoing, &edge)
	n.incomming = append(n.incomming, &edge)
	other.incomming = append(other.incomming, &edge)
	other.outgoing = append(other.outgoing, &edge)
}

func (n *Node) AddDirectedEdge(other *Node, distance int) {
	edge := edge{
		from:   n,
		to:     other,
		length: distance,
	}
	n.outgoing = append(n.outgoing, &edge)
	other.incomming = append(other.incomming, &edge)
}

type edge struct {
	from, to *Node
	length   int
}

type Path struct {
	edges []*edge
}

func (p *Path) lastNode() *Node {
	if len(p.edges) == 0 {
		return nil
	}

	return p.edges[len(p.edges)-1].to
}

func (p *Path) TotalDistance() int {
	distance := 0
	for _, e := range p.edges {
		distance += e.length
	}
	return distance
}

type priorityQueue []*Path

func (pq priorityQueue) Len() int {
	return len(pq)
}

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].TotalDistance() < pq[j].TotalDistance()
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *priorityQueue) Push(el any) {
	n := el.(*Path)
	*pq = append(*pq, n)
}

func (pq *priorityQueue) Pop() any {
	index := pq.Len() - 1
	el := (*pq)[index]
	*pq = (*pq)[0:index]
	return el
}

func newPaths(currentPath *Path, visited map[*Node]int) []*Path {
	paths := make([]*Path, 0)

	currentEnd := currentPath.lastNode()
	for _, e := range currentEnd.outgoing {
		distance := currentPath.TotalDistance() + e.length
		if currentShortestDistance, ok := visited[e.to]; !ok || currentShortestDistance > distance {
			paths = append(paths, &Path{
				edges: append(append(make([]*edge, 0, len(currentPath.edges)+1), currentPath.edges...), e),
			})
			visited[currentEnd] = distance
		}
	}

	return paths
}

func ShortestPath(from, to *Node) *Path {
	visited := make(map[*Node]int)

	pq := make(priorityQueue, 0)
	for _, e := range from.outgoing {
		pq = append(pq, &Path{edges: []*edge{e}})
	}
	heap.Init(&pq)

	current := heap.Pop(&pq).(*Path)
	for current.lastNode() != to {
		for _, p := range newPaths(current, visited) {
			heap.Push(&pq, p)
		}
		current = heap.Pop(&pq).(*Path)
	}

	return current
}
