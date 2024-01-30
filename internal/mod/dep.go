package mod

import "slices"

type DepNode struct {
	Node int
	Deps []int
}

func Dep(node int, deps ...int) *DepNode {
	return &DepNode{
		Node: node,
		Deps: deps,
	}
}

type DepTree struct {
	nodesMap map[int]*DepNode
}

func NewDepTree(nodes ...*DepNode) *DepTree {
	dt := &DepTree{
		nodesMap: make(map[int]*DepNode),
	}
	for _, node := range nodes {
		dt.AddNode(node)
	}
	return dt
}

func (dt *DepTree) AddNode(node *DepNode) {
	if _, ok := dt.nodesMap[node.Node]; !ok {
		dt.nodesMap[node.Node] = new(DepNode)
	}
	for _, dep := range node.Deps {
		if d, ok := dt.nodesMap[dep]; !ok {
			dt.nodesMap[dep] = &DepNode{
				Node: dep,
				Deps: []int{node.Node},
			}
		} else {
			d.Deps = append(d.Deps, node.Node)
		}
	}
}

func (dt *DepTree) TopSort() []int {
	var sorted []int
	visited := make(map[int]bool)
	for node := range dt.nodesMap {
		dt.topSort(node, visited, &sorted)
	}
	slices.Reverse(sorted)
	return sorted
}

func (dt *DepTree) topSort(node int, visited map[int]bool, sorted *[]int) {
	if visited[node] {
		return
	}
	visited[node] = true
	for _, dep := range dt.nodesMap[node].Deps {
		dt.topSort(dep, visited, sorted)
	}
	*sorted = append(*sorted, node)
}
