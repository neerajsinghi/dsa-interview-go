"""Shared trie node and trie drawing widget for the 11-tries tracers.

This file is not a tracer (build.py only picks up p*.py); the tracers import it.
"""

from walkthrough import BORDER, MUTED, STYLES, Widget, _label, centered, fit_font, font


class TNode:
    """Python port of the Go trie Node: children keyed by letter (or bit)."""

    def __init__(self):
        self.children = {}
        self.end = False
        self.passes = 0
        self.word = ""  # used by Word Search II

    def kids(self):
        return [(k, self.children[k]) for k in sorted(self.children)]


class TrieView(Widget):
    """Draws a trie top-down. Nodes are addressed by their path string.

    hl:    {path: style}   e.g. {"": "active", "c": "active", "ca": "match"}
    notes: {path: text}    drawn under the node
    End-of-word nodes get a double ring.
    """

    def __init__(self, root, label=None, hl=None, notes=None, level_h=38, r=13, h=None):
        self.root = root
        self.label = label
        self.hl = hl or {}
        self.notes = notes or {}
        self.level_h = level_h
        self.r = r
        self.fixed_h = h
        self._layout()

    def _layout(self):
        self.pos = {}  # path -> (slot, depth)
        self.edges = []
        self.char = {}
        self.ends = set()
        counter = [0]

        def walk(n, path, depth):
            if n.end or n.word:
                self.ends.add(path)
            kids = n.kids()
            if not kids:
                slot = counter[0]
                counter[0] += 1
            else:
                slots = []
                for k, c in kids:
                    p = path + str(k)
                    self.char[p] = str(k)
                    self.edges.append((path, p))
                    slots.append(walk(c, p, depth + 1))
                slot = (slots[0] + slots[-1]) / 2
            self.pos[path] = (slot, depth)
            return slot

        walk(self.root, "", 0)
        self.leaves = max(1, counter[0])
        self.depth = max(d for _, d in self.pos.values())

    def height(self, width):
        if self.fixed_h:
            return self.fixed_h
        return (20 if self.label else 0) + self.depth * self.level_h + 2 * self.r + 22

    def draw(self, d, x, y, width):
        _label(d, x, y, self.label)
        y += 20 if self.label else 0
        r = self.r
        gap = min(52, (width - 2 * r - 8) / max(1, self.leaves - 1)) if self.leaves > 1 else 0
        total = gap * (self.leaves - 1)
        x0 = x + (width - total) / 2

        def xy(p):
            s, dep = self.pos[p]
            return x0 + s * gap, y + r + 2 + dep * self.level_h

        for a, b in self.edges:
            (x1, y1), (x2, y2) = xy(a), xy(b)
            st = self.hl.get(b) if a in self.hl else None
            col = STYLES[st][1] if st else "#8c959f"
            d.line([(x1, y1 + r), (x2, y2 - r)], fill=col, width=3 if st else 2)
        for p in self.pos:
            px, py = xy(p)
            st = self.hl.get(p)
            fill, border, ink = STYLES[st]
            if p in self.ends:
                d.ellipse([px - r - 3, py - r - 3, px + r + 3, py + r + 3], outline=border if st else "#57606a", width=2)
            d.ellipse([px - r, py - r, px + r, py + r], fill=fill, outline=border if st else BORDER,
                      width=2 if st else 1)
            s = "·" if p == "" else self.char[p]
            centered(d, px, py, s, fit_font(s, "mono", 14, 2 * r - 2), ink if p else MUTED)
            if p in self.notes:
                s = str(self.notes[p])
                leaf = not any(a == p for a, _ in self.edges)
                if leaf:
                    centered(d, px, py + r + 11, s, font("bold", 11), "#8250df")
                else:
                    d.text((px + r + 8, py), s, font=font("bold", 11), fill="#8250df", anchor="lm")
