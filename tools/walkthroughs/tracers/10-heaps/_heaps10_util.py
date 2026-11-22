"""Heap helpers shared by the 10-heaps tracers (not a tracer itself).

GoHeap is a line-by-line port of Go's container/heap (Init, Push, Pop with
the same up/down sifting), wrapped like the module's Heap[T]. Using the same
sift order means the array shown in each frame is exactly what the Go code
holds at that point.
"""

from walkthrough import Array, Row, Tree


class GoHeap:
    def __init__(self, less, items=None):
        self.less = less
        self.items = list(items or [])
        n = len(self.items)
        for i in range(n // 2 - 1, -1, -1):  # heap.Init
            self._down(i, n)

    def __len__(self):
        return len(self.items)

    def _lt(self, i, j):
        return self.less(self.items[i], self.items[j])

    def _swap(self, i, j):
        self.items[i], self.items[j] = self.items[j], self.items[i]

    def _up(self, j):
        while True:
            i = (j - 1) // 2 if j > 0 else 0
            if i == j or not self._lt(j, i):
                break
            self._swap(i, j)
            j = i

    def _down(self, i0, n):
        i = i0
        while True:
            j1 = 2 * i + 1
            if j1 >= n:
                break
            j = j1
            j2 = j1 + 1
            if j2 < n and self._lt(j2, j1):
                j = j2
            if not self._lt(j, i):
                break
            self._swap(i, j)
            i = j
        return i > i0

    def push(self, v):
        self.items.append(v)
        self._up(len(self.items) - 1)

    def pop(self):
        n = len(self.items) - 1
        self._swap(0, n)
        self._down(0, n)
        return self.items.pop()

    def peek(self):
        return self.items[0]


def heap_view(items, label, fmt=str, hl=None, tree_fmt=None, weights=(1, 1)):
    """The heap drawn twice: as a tree (level order = the array) and as the array."""
    tree_fmt = tree_fmt or fmt
    hl = hl or {}
    return Row(
        Tree([tree_fmt(v) for v in items], label=f"{label} as a tree", hl=hl),
        Array([fmt(v) for v in items], label=f"{label} as an array", hl=hl),
        weights=list(weights),
    )
