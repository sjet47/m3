package mod

import (
	"cmp"
	"slices"
)

type DepNode[N cmp.Ordered] struct {
	Node N
	Deps []N
}

func Dep[N cmp.Ordered](node N, deps ...N) *DepNode[N] {
	return &DepNode[N]{
		Node: node,
		Deps: deps,
	}
}

type DepTree[N cmp.Ordered] struct {
	nodesMap map[N]*DepNode[N]
}

func NewDepTree[N cmp.Ordered](nodes ...*DepNode[N]) *DepTree[N] {
	dt := &DepTree[N]{
		nodesMap: make(map[N]*DepNode[N]),
	}
	for _, node := range nodes {
		dt.AddNode(node)
	}
	return dt
}

func (dt *DepTree[N]) AddNode(node *DepNode[N]) {
	if _, ok := dt.nodesMap[node.Node]; !ok {
		dt.nodesMap[node.Node] = new(DepNode[N])
	}
	for _, dep := range node.Deps {
		if d, ok := dt.nodesMap[dep]; !ok {
			dt.nodesMap[dep] = &DepNode[N]{
				Node: dep,
				Deps: []N{node.Node},
			}
		} else {
			d.Deps = append(d.Deps, node.Node)
		}
	}
}

func (dt *DepTree[N]) TopSort() []N {
	var sorted []N
	visited := make(map[N]bool)
	for node := range dt.nodesMap {
		dt.topSort(node, visited, &sorted)
	}
	slices.Reverse(sorted)
	return sorted
}

func (dt *DepTree[N]) topSort(node N, visited map[N]bool, sorted *[]N) {
	if visited[node] {
		return
	}
	visited[node] = true
	for _, dep := range dt.nodesMap[node].Deps {
		dt.topSort(dep, visited, sorted)
	}
	*sorted = append(*sorted, node)
}
