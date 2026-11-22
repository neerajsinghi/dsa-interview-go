"""Tree helpers shared by the 09-trees tracers (not a tracer itself).

Node mirrors the Go TreeNode. from_level / level convert to and from the
LeetCode level-order list that the Tree widget draws, and level() also
returns where each node sits in that list so tracers can highlight it.
"""


class Node:
    def __init__(self, val, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


def from_level(vals):
    """Port of FromLevelOrder; None marks a missing child."""
    if not vals or vals[0] is None:
        return None
    root = Node(vals[0])
    queue, i = [root], 1
    while queue and i < len(vals):
        node = queue.pop(0)
        if i < len(vals) and vals[i] is not None:
            node.left = Node(vals[i])
            queue.append(node.left)
        i += 1
        if i < len(vals) and vals[i] is not None:
            node.right = Node(vals[i])
            queue.append(node.right)
        i += 1
    return root


def level(root):
    """Port of ToLevelOrder. Returns (values, {id(node): index})."""
    out, idx = [], {}
    queue = [root]
    while queue:
        node = queue.pop(0)
        if node is None:
            out.append(None)
            continue
        idx[id(node)] = len(out)
        out.append(node.val)
        queue += [node.left, node.right]
    while out and out[-1] is None:
        out.pop()
    return out, idx
