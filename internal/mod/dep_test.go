package mod

import (
	"testing"
)

func TestDepTree(t *testing.T) {
	deps := []*DepNode[int]{
		Dep(1, 4),
		Dep(2, 4),
		Dep(3, 1),
		Dep(3, 2),
		Dep(4, 7),
		Dep(4, 8),
		Dep(5, 1),
		Dep(5, 4),
		Dep(5, 6),
		Dep(6, 10),
		Dep(6, 11),
		Dep(7, 9),
		Dep(8, 9),
		Dep(8, 10),
		Dep(9, 12),
		Dep(10, 12),
		Dep(10, 13),
		Dep(11, 10),
		Dep(6, 10),
		Dep(6, 11),
		Dep(7, 9),
		Dep(8, 9),
		Dep(8, 10),
		Dep(9, 12),
		Dep(10, 12),
		Dep(10, 13),
		Dep(11, 10),
	}
	depMap := make(map[int]map[int]bool)
	for _, dep := range deps {
		if _, ok := depMap[dep.Node]; !ok {
			depMap[dep.Node] = make(map[int]bool)
		}
		for _, d := range dep.Deps {
			depMap[dep.Node][d] = true
		}
	}

	ts := NewDepTree(deps...).TopSort()
	exist := make(map[int]bool)
	for _, node := range ts {
		for dep := range depMap[node] {
			if !exist[dep] {
				t.Errorf("Node %d expect dep %d exist but not", node, dep)
			}
		}
		exist[node] = true
	}
}
