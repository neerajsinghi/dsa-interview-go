package graphs

// Problem: Course Schedule I & II
// Difficulty: Medium
//
// There are numCourses courses labelled 0..numCourses-1. Each prerequisite
// pair {a, b} means "take b before a". Part I: can every course be finished?
// Part II: return one valid order to take all courses, or an empty result if
// that is impossible.
//
// Examples:
//
//	2 courses, [[1,0]]        -> I: true,  II: [0 1]
//	2 courses, [[1,0],[0,1]]  -> I: false, II: []   (they require each other)
//	4 courses, [[1,0],[2,0],[3,1],[3,2]] -> II: [0 1 2 3] or [0 2 1 3]
//
// References: https://leetcode.com/problems/course-schedule/
//
//	https://leetcode.com/problems/course-schedule-ii/
//
// Pattern: topological sort. Draw an edge b -> a for each pair; a valid
// schedule is exactly a topological order, which exists iff the graph has no
// cycle. Kahn's algorithm (BFS on in-degree 0) answers both parts at once:
// if it emits every node there is no cycle and the emission order is the
// schedule. The DFS three-colour method is the other standard answer.
//
// Time: O(V + E). Space: O(V + E).

// prereqGraph builds the edge b -> a for every pair {a, b}.
func prereqGraph(numCourses int, prerequisites [][]int) [][]int {
	adj := make([][]int, numCourses)
	for _, p := range prerequisites {
		adj[p[1]] = append(adj[p[1]], p[0])
	}
	return adj
}

// CanFinish answers part I with Kahn's algorithm.
func CanFinish(numCourses int, prerequisites [][]int) bool {
	_, ok := TopoSortKahn(prereqGraph(numCourses, prerequisites))
	return ok
}

// CanFinishDFS answers part I by looking for a back edge (a cycle).
func CanFinishDFS(numCourses int, prerequisites [][]int) bool {
	_, ok := TopoSortDFS(prereqGraph(numCourses, prerequisites))
	return ok
}

// FindOrder answers part II; it returns an empty, non-nil slice when no order
// exists so callers can distinguish "no courses" from a bug.
func FindOrder(numCourses int, prerequisites [][]int) []int {
	order, ok := TopoSortKahn(prereqGraph(numCourses, prerequisites))
	if !ok {
		return []int{}
	}
	return order
}
