package graphs

// Problem: Clone Graph
// Difficulty: Medium
//
// You are handed one node of a connected undirected graph. Each node has an
// integer value and a list of neighbours. Return a deep copy: brand-new nodes
// with the same values and the same connections, sharing no pointer with the
// original.
//
// Example:
//
//	1 -- 2
//	|    |      clone(node 1) -> a new 4-cycle 1-2-3-4 made of fresh nodes
//	4 -- 3
//
// Reference: https://leetcode.com/problems/clone-graph/
//
// Pattern: graph traversal + a map from original node to its copy. The map
// doubles as the visited set, which is what stops the traversal looping
// forever on cycles. Create the copy BEFORE recursing into neighbours, so a
// cycle that leads back finds the (still incomplete) copy in the map.
//
// Time: O(V + E). Space: O(V) for the map (plus recursion or queue).

// GraphNode is a vertex in an undirected graph given by pointers.
type GraphNode struct {
	Val       int
	Neighbors []*GraphNode
}

// CloneGraph deep-copies with recursive DFS.
func CloneGraph(node *GraphNode) *GraphNode {
	copies := map[*GraphNode]*GraphNode{}
	var clone func(n *GraphNode) *GraphNode
	clone = func(n *GraphNode) *GraphNode {
		if c, ok := copies[n]; ok {
			return c
		}
		c := &GraphNode{Val: n.Val}
		copies[n] = c // register first: neighbours may point back at n
		for _, nb := range n.Neighbors {
			c.Neighbors = append(c.Neighbors, clone(nb))
		}
		return c
	}
	if node == nil {
		return nil
	}
	return clone(node)
}

// CloneGraphBFS deep-copies iteratively: create a copy when a node is first
// discovered, and wire up its neighbour list when it is dequeued.
func CloneGraphBFS(node *GraphNode) *GraphNode {
	if node == nil {
		return nil
	}
	copies := map[*GraphNode]*GraphNode{node: {Val: node.Val}}
	queue := []*GraphNode{node}
	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		for _, nb := range n.Neighbors {
			if _, ok := copies[nb]; !ok {
				copies[nb] = &GraphNode{Val: nb.Val}
				queue = append(queue, nb)
			}
			copies[n].Neighbors = append(copies[n].Neighbors, copies[nb])
		}
	}
	return copies[node]
}
