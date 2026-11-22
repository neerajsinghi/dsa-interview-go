import os
import sys

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

from _trees09_util import Node, level  # noqa: E402
from walkthrough import Array, Frame, Text, Tree  # noqa: E402

TITLE = "Construct from Preorder + Inorder"
PATTERN = "Divide & conquer"


def frames():
    preorder = [3, 9, 20, 15, 7]
    inorder = [9, 3, 15, 20, 7]
    built = {"root": None}

    def widgets(nxt=None, lo=None, hi=None, mid=None, new=None, pre_done=0):
        pre_hl = {k: "visited" for k in range(pre_done)}
        if nxt is not None and nxt < len(preorder):
            pre_hl[nxt] = "active"
        in_hl = {}
        if lo is not None and lo < hi:
            for k in range(lo, hi):
                in_hl[k] = "window"
        if mid is not None:
            in_hl[mid] = "match"
        ptr = {"next": nxt} if nxt is not None and nxt < len(preorder) else {}
        window = (lo, hi - 1) if lo is not None and lo < hi else None
        if built["root"] is None:
            tree = Tree([], label="tree so far")
        else:
            lv, idx = level(built["root"])
            tree = Tree(lv, label="tree so far", hl={idx[id(new)]: "match"} if new else {})
        return (Array(preorder, label="preorder", hl=pre_hl, pointers=ptr, cell=50),
                Array(inorder, label="inorder  (window = values of this subtree)", hl=in_hl, window=window, cell=50),
                tree)

    yield Frame(
        *widgets(),
        speaker="Interviewer",
        say="Here are the preorder and inorder traversals of a binary tree with distinct values. Rebuild the tree.",
    )
    yield Frame(
        *widgets(),
        Text("brute force: find the root in inorder with a scan, slice, recurse → O(n²)"),
        say="preorder[0] is the root, and its spot in inorder splits left from right. Scanning for it at every "
            "level costs O(n²) on a skewed tree. A value→index map makes each lookup O(1).",
    )

    pos = {v: i for i, v in enumerate(inorder)}
    nxt = 0
    events = []

    def build(lo, hi, parent=None, side=None):
        nonlocal nxt
        if lo >= hi:
            return None
        root = Node(preorder[nxt])
        if parent is None:
            built["root"] = root
        else:
            setattr(parent, side, root)  # attach now so the picture grows as we go
        used = nxt
        nxt += 1
        mid = pos[root.val]
        events.append((root, lo, hi, mid, used))
        yield root, lo, hi, mid, used
        root.left = yield from build(lo, mid, root, "left")
        root.right = yield from build(mid + 1, hi, root, "right")
        return root

    gen = build(0, len(inorder))
    first = True
    for root, lo, hi, mid, used in gen:
        left_n, right_n = mid - lo, hi - mid - 1
        if first:
            say = (f"pos = {{9:0, 3:1, 15:2, 20:3, 7:4}}. next points at {root.val}, so {root.val} is the root. "
                   f"It sits at inorder index {mid}: {left_n} value to its left, {right_n} to its right.")
            first = False
        elif left_n == 0 and right_n == 0:
            say = (f"The window holds only {root.val}, so the next preorder value {root.val} becomes a leaf. "
                   "Both child windows are empty and return nil.")
        else:
            say = (f"Left side done. next now points at {root.val}, the root of the right window. "
                   f"Its inorder index {mid} splits that window into {left_n} left and {right_n} right.")
        yield Frame(
            *widgets(used, lo, hi, mid, root, pre_done=used),
            Text(f"build(lo={lo}, hi={hi}): root = preorder[{used}] = {root.val}, mid = pos[{root.val}] = {mid}"),
            say=say,
        )

    lv, _ = level(built["root"])
    assert lv == [3, 9, 20, None, None, 15, 7], lv
    yield Frame(
        *widgets(pre_done=len(preorder)),
        Text(("answer: [3,9,20,null,null,15,7]", "match"), "time O(n), space O(n) for the map + O(h) stack"),
        say="Every preorder value became a node exactly once, with an O(1) map lookup each. That's O(n) time, "
            "and O(n) space for the map. The left call must run before the right one, because preorder is root, left, right.",
    )
