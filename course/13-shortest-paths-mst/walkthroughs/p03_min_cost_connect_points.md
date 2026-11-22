# Walkthrough: Min Cost to Connect All Points

**Module:** [13 · Shortest paths & minimum spanning trees](../README.md) · **Difficulty:** Medium · **Solution:** [`p03_min_cost_connect_points.go`](../p03_min_cost_connect_points.go)

![Animated walkthrough of Min Cost to Connect All Points](p03_min_cost_connect_points.gif)

## The problem in plain words

You have points on a grid. Joining two points costs their Manhattan distance, `|x1 - x2| + |y1 - y2|`. Join them so every point can reach every other point, as cheaply as possible, and return the total cost.

```
[[0,0],[2,2],[3,10],[5,2],[7,0]]  →  20
```

## Mock interview

**Interviewer:** Given points on a 2-D grid, connecting two points costs their Manhattan distance. Return the minimum total cost to connect all of them.

**Candidate:** Can points repeat? How many points are there?

**Interviewer:** Points are distinct, and there are up to 1,000 of them.

**Candidate:** "Connect everything with minimum total cost" is a minimum spanning tree. In the example, the cheapest tree uses edges 0–1 (4), 1–3 (3), 3–4 (4) and 1–2 (9), for a total of 20. One point gives a cost of 0.

**Candidate:** A brute force would try every spanning tree, and there are exponentially many. The usual MST algorithms are Kruskal and Prim. Here every pair of points is a possible edge, so there are about n²/2 edges. Kruskal has to build and sort all of them: O(n² log n) time and O(n²) memory.

**Interviewer:** Can you avoid building the edge list?

**Candidate:** Yes, with array-based Prim. I grow one tree from point 0. I keep `minDist[v]`, the cheapest edge from the tree to point `v`. Each round I scan for the closest point outside the tree, add it, and add its `minDist` to the total. Then I let the new point offer cheaper edges: for every outside point, `minDist[v] = min(minDist[v], manhattan(u, v))`. Each round is O(n), and there are n rounds.

**Interviewer:** Why is it safe to always take the closest point?

**Candidate:** That's the cut property. Split the points into "in the tree" and "outside". The cheapest edge crossing that split belongs to some minimum spanning tree. If an MST didn't use it, I could swap it in for whatever edge crosses the split on that tree's path, and the total wouldn't go up.

**Candidate:** In code: `inTree` and `minDist` slices, `minDist[0] = 0`, and a loop of `n` rounds. Each round does one scan to pick `u` and one scan to update the others. I don't need a heap. With a complete graph, a heap would cost O(n² log n) because every pair is an edge.

**Candidate:** Testing it: round 1 takes point 0 for free and sets `minDist` to `[0, 4, 13, 7, 7]`. Round 2 adds point 1 at 4, which lowers 2 to 9 and 3 to 3. Round 3 adds point 3 at 3, which lowers 4 to 4. Round 4 adds point 4 at 4. Round 5 adds point 2 at 9. The total is 4 + 3 + 4 + 9 = 20.

**Interviewer:** Complexity?

**Candidate:** **O(n²) time and O(n) space** for Prim. Kruskal, `MinCostConnectPointsKruskal`, is O(n² log n) time and O(n²) space.

## What the animation shows

| Step | Point added | Cost | `minDist` updates | Total |
|---|---|---|---|---|
| 1 | Interviewer states the problem | | | |
| 2 | Complete graph; Kruskal would sort O(n²) edges | | | |
| 3 | Start: `minDist[0] = 0` | | | |
| 4 | 0 | 0 | 1→4, 2→13, 3→7, 4→7 | 0 |
| 5 | 1 (via 0) | 4 | 2: 13→9, 3: 7→3 | 4 |
| 6 | 3 (via 1) | 3 | 4: 7→4 | 7 |
| 7 | 4 (via 3) | 4 | none | 11 |
| 8 | 2 (via 1) | 9 | none | 20 |
| 9 | Answer 20, O(n²) time | | | |

## Detailed explanation

Prim keeps one invariant: *the chosen points form a subtree of some minimum spanning tree, and `minDist[v]` is the cheapest edge from that subtree to each outside point `v`.* Picking the smallest `minDist` takes the cheapest edge across the cut between the tree and the rest. By the cut property that edge is safe, so the invariant still holds after adding it. Updating `minDist` from the new point restores the second half of the invariant, because the only new edges leaving the tree start at `u`.

Why not a heap? Heap-based Prim is O(E log V). On a complete graph E ≈ n², so that becomes O(n² log n). The array version finds the minimum with a linear scan, O(n) per round, and never stores edges. That makes it O(n²) overall, which is optimal: you must look at every pair at least once.

The Go code doesn't store which tree point each edge comes from. The animation tracks it only to draw the edges.

Edge cases:

- **No points:** the function returns 0.
- **One point:** one round adds it at cost 0.
- **Ties in `minDist`:** any choice gives the same total.

## Complexity

| Approach | Time | Space |
|---|---|---|
| Dense Prim (`MinCostConnectPoints`) | O(n²) | O(n) |
| Kruskal on all pairs (`MinCostConnectPointsKruskal`) | O(n² log n) | O(n²) |

## Common mistakes

- Using Euclidean distance instead of Manhattan.
- Forgetting to skip points already in the tree when updating `minDist`.
- Starting every `minDist` at 0, so the first scan picks the wrong point.
- Using a heap on a complete graph and paying an extra log factor for nothing.
- Summing all `n·(n-1)/2` edges in Kruskal without checking that the two points are in different components.

## Try it yourself

```sh
go test ./course/13-shortest-paths-mst/ -run TestMinCostConnectPoints -v
```
