import os
import sys

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

from _trees09_util import from_level, level  # noqa: E402
from walkthrough import Frame, Text, Tree  # noqa: E402

TITLE = "Lowest Common Ancestor (BST + general)"
PATTERN = "BST descent"


def lca_bst(root, p, q):
    """LowestCommonAncestorBST as a generator: yields (node, decision)."""
    n = root
    while n is not None:
        if p < n.val and q < n.val:
            yield n, "left"
            n = n.left
        elif p > n.val and q > n.val:
            yield n, "right"
            n = n.right
        else:
            yield n, "split"
            return n
    return None


def frames():
    vals = [6, 2, 8, 0, 4, 7, 9, None, None, 3, 5]
    root = from_level(vals)
    lo, idx = level(root)

    def tree(hl):
        return Tree(lo, label="BST", hl=hl)

    yield Frame(
        tree({}),
        Text("LCA(2, 8) = ?   LCA(3, 5) = ?"),
        speaker="Interviewer",
        say="Given two nodes p and q in this BST, return their lowest common ancestor: "
            "the deepest node that has both in its subtree. A node counts as its own descendant.",
    )
    yield Frame(
        tree({}),
        Text("general approach: post-order search of the whole tree → O(n)",
             "or store both root-to-node paths and compare → extra O(h) space"),
        say="The general-tree method searches every node, O(n). "
            "But this is a BST, so the values tell me which side p and q are on without searching.",
    )

    results = {}
    for p, q in ((3, 5), (2, 8)):
        hl = {}
        gen = lca_bst(root, p, q)
        steps = []
        try:
            while True:
                steps.append(next(gen))
        except StopIteration as e:
            results[(p, q)] = e.value.val
        for n, decision in steps:
            i = idx[id(n)]
            if decision == "split":
                hl[i] = "match"
                text = Text(f"p={p}, q={q}, node={n.val}: {p} <= {n.val} <= {q} → split here",
                            (f"LCA({p}, {q}) = {n.val}", "match"))
                if (p, q) == (2, 8):
                    say = "LCA(2, 8): at the root, 2 is smaller and 8 is larger, so they split right away. The answer is 6."
                else:
                    say = (f"At {n.val}, {p} is on the left and {q} is on the right. They split here, "
                           f"so {n.val} is the lowest node holding both.")
            else:
                hl[i] = "visited"
                if decision == "left":
                    text = Text(f"p={p}, q={q}, node={n.val}: both < {n.val} → go left")
                    say = f"Both {p} and {q} are smaller than {n.val}, so both live in the left subtree. The LCA must be down there too."
                else:
                    text = Text(f"p={p}, q={q}, node={n.val}: both > {n.val} → go right")
                    say = f"Both are larger than {n.val}, so both are in its right subtree. I move right and ignore everything on the left."
            pq = {lo.index(p): "special", lo.index(q): "special"}  # p and q themselves
            yield Frame(tree({**pq, **hl, i: hl[i] if decision == "split" else "active"}), text, say=say)

    assert results[(3, 5)] == 4 and results[(2, 8)] == 6
    yield Frame(
        tree({idx[id(root)]: "match", idx[id(root.left.right)]: "match"}),
        Text(("LCA(2, 8) = 6,  LCA(3, 5) = 4", "match"), "BST: O(h) time, O(1) space · general tree: O(n), O(h)"),
        say="The BST walk follows one path, so it's O(h) time and O(1) space. For a general tree I'd use the "
            "post-order search: the first node that finds p on one side and q on the other.",
    )
