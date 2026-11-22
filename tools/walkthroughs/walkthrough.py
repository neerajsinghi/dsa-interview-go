"""Draw step-by-step animated GIFs for the course walkthroughs.

A tracer script describes an algorithm run as a list of Frames. Each Frame
stacks a few widgets (an array, a grid, a tree, a graph, some text) above a
caption in which the candidate explains the step. render() lays every frame
out at the same size and writes an animated GIF.

Only Pillow is needed: pip install pillow
"""

import math
import os
import textwrap

from PIL import Image, ImageDraw, ImageFont

WIDTH = 760
PAD = 22

INK = "#1f2328"
MUTED = "#656d76"
BORDER = "#d0d7de"
CELL = "#f6f8fa"
BG = "#ffffff"
BAR = "#0d1117"

# Highlight styles: (fill, border, text). Tracers refer to them by name.
STYLES = {
    None: (CELL, BORDER, INK),
    "active": ("#ddf4ff", "#0969da", "#0a3069"),
    "compare": ("#fff1e5", "#bc4c00", "#762c00"),
    "match": ("#dafbe1", "#1a7f37", "#044f1e"),
    "reject": ("#ffebe9", "#cf222e", "#82071e"),
    "window": ("#fff8c5", "#9a6700", "#4d2d00"),
    "visited": ("#eaeef2", "#8c959f", "#57606a"),
    "special": ("#fbefff", "#8250df", "#512a97"),
    "empty": (BG, "#eaeef2", "#8c959f"),
}

POINTER_COLORS = ["#0969da", "#bc4c00", "#8250df", "#1a7f37", "#cf222e", "#9a6700"]

_FONT_CANDIDATES = {
    "sans": ["/System/Library/Fonts/Helvetica.ttc", "/System/Library/Fonts/SFNS.ttf",
             "/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", "DejaVuSans.ttf"],
    "bold": ["/System/Library/Fonts/Helvetica.ttc", "/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf",
             "DejaVuSans-Bold.ttf"],
    "mono": ["/System/Library/Fonts/Menlo.ttc", "/System/Library/Fonts/SFNSMono.ttf",
             "/usr/share/fonts/truetype/dejavu/DejaVuSansMono.ttf", "DejaVuSansMono.ttf"],
}
_font_cache = {}


def font(kind="sans", size=16):
    key = (kind, size)
    if key in _font_cache:
        return _font_cache[key]
    for path in _FONT_CANDIDATES[kind]:
        try:
            # Helvetica.ttc index 1 is bold.
            index = 1 if kind == "bold" and path.endswith("Helvetica.ttc") else 0
            f = ImageFont.truetype(path, size, index=index)
            break
        except OSError:
            continue
    else:
        f = ImageFont.load_default()
    _font_cache[key] = f
    return f


_FALLBACK = ["/System/Library/Fonts/Menlo.ttc", "/Library/Fonts/Arial Unicode.ttf",
             "/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", "DejaVuSans.ttf"]
_glyph_cache = {}


def _fallback(f):
    size = getattr(f, "size", 16)
    key = ("fallback", size)
    if key not in _font_cache:
        for path in _FALLBACK:
            try:
                _font_cache[key] = ImageFont.truetype(path, size)
                break
            except OSError:
                continue
        else:
            _font_cache[key] = f
    return _font_cache[key]


def _has_glyph(f, ch):
    key = (id(f), ch)
    if key not in _glyph_cache:
        if ord(ch) < 128 or not hasattr(f, "getmask"):
            _glyph_cache[key] = True
        else:
            _glyph_cache[key] = bytes(f.getmask(ch)) != bytes(f.getmask("\uffff"))
    return _glyph_cache[key]


def _runs(s, f):
    """Split s into (text, font) runs, using a fallback font for missing glyphs."""
    runs = []
    for ch in str(s):
        use = f if _has_glyph(f, ch) else _fallback(f)
        if runs and runs[-1][1] is use:
            runs[-1][0] += ch
        else:
            runs.append([ch, use])
    return runs


def text_w(s, f):
    return sum(ff.getlength(part) for part, ff in _runs(s, f))


def _text(d, xy, s, font=None, fill=None, anchor=None):
    """ImageDraw.text with per-character font fallback (Helvetica has no arrows)."""
    runs = _runs(s, font)
    if len(runs) <= 1:
        d.text(xy, str(s), font=font, fill=fill, anchor=anchor)
        return
    anchor = anchor or "la"
    x, y = xy
    total = sum(ff.getlength(part) for part, ff in runs)
    if anchor[0] == "m":
        x -= total / 2
    elif anchor[0] == "r":
        x -= total
    for part, ff in runs:
        d.text((x, y), part, font=ff, fill=fill, anchor="l" + anchor[1])
        x += ff.getlength(part)


def fit_font(s, kind, size, max_w, min_size=9):
    """Shrink the font until s fits in max_w."""
    while size > min_size and text_w(s, font(kind, size)) > max_w:
        size -= 1
    return font(kind, size)


def centered(d, cx, cy, s, f, fill):
    _text(d, (cx, cy), str(s), font=f, fill=fill, anchor="mm")


def arrow(d, x1, y1, x2, y2, color, width=2, head=8):
    d.line([(x1, y1), (x2, y2)], fill=color, width=width)
    ang = math.atan2(y2 - y1, x2 - x1)
    for side in (-0.45, 0.45):
        d.line([(x2, y2), (x2 - head * math.cos(ang + side), y2 - head * math.sin(ang + side))],
               fill=color, width=width)


def _fmt(v):
    if v is None:
        return ""
    if isinstance(v, bool):
        return "T" if v else "F"
    if isinstance(v, float):
        if v == math.inf:
            return "∞"
        if v == -math.inf:
            return "-∞"
        return f"{v:g}"
    return str(v)


# --------------------------------------------------------------------------- widgets


class Widget:
    def height(self, width):
        raise NotImplementedError

    def draw(self, d, x, y, width):
        raise NotImplementedError


def _label(d, x, y, label):
    if label:
        _text(d, (x, y), label, font=font("bold", 14), fill=MUTED)


class Array(Widget):
    """A row of cells.

    hl:       {index: style}
    pointers: {name: index}; several names may share an index
    window:   (lo, hi) inclusive range drawn as a bracket under the cells
    """

    def __init__(self, values, label=None, hl=None, pointers=None, window=None,
                 show_index=True, cell=None):
        self.values = list(values)
        self.label = label
        self.hl = hl or {}
        self.pointers = pointers or {}
        self.window = window
        self.show_index = show_index
        self.cell = cell

    def _cell(self, width):
        n = max(1, len(self.values))
        return self.cell or int(max(26, min(56, width / n)))

    def _rows(self):
        stack = {}
        for name, i in self.pointers.items():
            stack.setdefault(i, []).append(name)
        return stack, max((len(v) for v in stack.values()), default=0)

    def height(self, width):
        _, depth = self._rows()
        h = (20 if self.label else 0) + (16 if self.show_index else 0) + self._cell(width)
        h += 20 * depth + (6 if depth else 0) + (14 if self.window else 0)
        return h + 8

    def draw(self, d, x, y, width):
        c = self._cell(width)
        _label(d, x, y, self.label)
        y += 20 if self.label else 0
        if self.show_index:
            for i in range(len(self.values)):
                centered(d, x + i * c + c / 2, y + 7, i, font("mono", 11), MUTED)
            y += 16
        for i, v in enumerate(self.values):
            fill, border, ink = STYLES[self.hl.get(i)]
            d.rounded_rectangle([x + i * c + 2, y, x + (i + 1) * c - 2, y + c - 4], 6,
                                fill=fill, outline=border, width=2 if i in self.hl else 1)
            s = _fmt(v)
            centered(d, x + i * c + c / 2, y + (c - 4) / 2, s,
                     fit_font(s, "mono", 16 if c >= 40 else 13, c - 8), ink)
        y += c
        if self.window:
            lo, hi = self.window
            if lo <= hi:
                x1, x2 = x + lo * c + 4, x + (hi + 1) * c - 4
                col = STYLES["window"][1]
                d.line([(x1, y), (x1, y + 8), (x2, y + 8), (x2, y)], fill=col, width=2)
            y += 14
        stack, _ = self._rows()
        names = list(self.pointers)
        for i, group in stack.items():
            for k, name in enumerate(group):
                col = POINTER_COLORS[names.index(name) % len(POINTER_COLORS)]
                if 0 <= i < len(self.values):
                    cx = x + i * c + c / 2
                elif i < 0:
                    cx = x - c / 2 + 6
                else:
                    cx = x + len(self.values) * c + c / 2 - 6
                yy = y + 2 + k * 20
                if k == 0:
                    arrow(d, cx, yy + 5, cx, yy - 3, col, width=2, head=5)
                centered(d, cx, yy + 16, name, font("bold", 14), col)


def stack_widget(values, label="stack", hl=None):
    """A stack drawn left to right, top on the right."""
    ptr = {"top": len(values) - 1} if values else {}
    return Array(values, label=label, hl=hl, pointers=ptr, show_index=False)


class Grid(Widget):
    """A 2-D table: DP tables, boards, bit rows.

    hl:         {(r, c): style}
    row_labels, col_labels: optional headers
    """

    def __init__(self, rows, label=None, hl=None, row_labels=None, col_labels=None, cell=None):
        self.rows = [list(r) for r in rows]
        self.label = label
        self.hl = hl or {}
        self.row_labels = row_labels
        self.col_labels = col_labels
        self.cell = cell

    def _geom(self, width):
        ncols = max((len(r) for r in self.rows), default=1)
        lw = 0
        if self.row_labels:
            lw = max(text_w(_fmt(s), font("mono", 13)) for s in self.row_labels) + 12
        c = self.cell or int(max(20, min(46, (width - lw) / max(1, ncols))))
        return c, lw

    def height(self, width):
        c, _ = self._geom(width)
        return (20 if self.label else 0) + (18 if self.col_labels else 0) + c * len(self.rows) + 8

    def draw(self, d, x, y, width):
        c, lw = self._geom(width)
        _label(d, x, y, self.label)
        y += 20 if self.label else 0
        if self.col_labels:
            for j, s in enumerate(self.col_labels):
                s = _fmt(s)
                centered(d, x + lw + j * c + c / 2, y + 8, s, fit_font(s, "mono", 13, c - 2), MUTED)
            y += 18
        for i, row in enumerate(self.rows):
            if self.row_labels:
                s = _fmt(self.row_labels[i])
                _text(d, (x + lw - 8, y + i * c + c / 2), s, font=font("mono", 13), fill=MUTED, anchor="rm")
            for j, v in enumerate(row):
                st = self.hl.get((i, j))
                fill, border, ink = STYLES[st]
                x0, y0 = x + lw + j * c, y + i * c
                d.rectangle([x0, y0, x0 + c, y0 + c], fill=fill, outline=border, width=1)
                if st:
                    d.rectangle([x0 + 1, y0 + 1, x0 + c - 1, y0 + c - 1], outline=border, width=2)
                s = _fmt(v)
                centered(d, x0 + c / 2, y0 + c / 2, s, fit_font(s, "mono", 15 if c >= 34 else 12, c - 4), ink)


class LinkedList(Widget):
    """Nodes joined by arrows.

    links:    {i: j or None}; default i -> i+1 and the last node -> nil
    pointers: {name: index}
    """

    def __init__(self, values, label=None, hl=None, pointers=None, links=None):
        self.values = list(values)
        self.label = label
        self.hl = hl or {}
        self.pointers = pointers or {}
        self.links = links

    def _geom(self, width):
        n = max(1, len(self.values))
        step = min(90, (width - 40) / n)
        box = max(26, min(46, step - 22))
        return step, box

    def height(self, width):
        stack = {}
        for name, i in self.pointers.items():
            stack.setdefault(i, []).append(name)
        depth = max((len(v) for v in stack.values()), default=0)
        return (20 if self.label else 0) + 34 + 46 + 20 * depth + 12

    def draw(self, d, x, y, width):
        step, box = self._geom(width)
        _label(d, x, y, self.label)
        y += 20 if self.label else 0
        top = y + 34  # room for back-arcs
        n = len(self.values)
        links = self.links if self.links is not None else {i: (i + 1 if i + 1 < n else None) for i in range(n)}

        def cx(i):
            return x + i * step + box / 2

        cy = top + 23
        for i, j in links.items():
            if j is None:
                if i == n - 1 or links.get(i, -2) is None:
                    pass
                continue
            col = "#57606a"
            if j == i + 1:
                arrow(d, cx(i) + box / 2, cy, cx(j) - box / 2 - 1, cy, col)
            else:
                # curved arc above the nodes
                x1, x2 = cx(i), cx(j)
                h = 14 + 6 * min(4, abs(j - i))
                pts = []
                for t in range(21):
                    tt = t / 20
                    px = x1 + (x2 - x1) * tt
                    py = top - 2 - h * math.sin(math.pi * tt)
                    pts.append((px, py))
                d.line(pts, fill="#8250df", width=2)
                arrow(d, pts[-2][0], pts[-2][1], x2, top - 1, "#8250df", head=7)
        for i, v in enumerate(self.values):
            fill, border, ink = STYLES[self.hl.get(i)]
            d.rounded_rectangle([cx(i) - box / 2, top, cx(i) + box / 2, top + 46], 8,
                                fill=fill, outline=border, width=2 if i in self.hl else 1)
            s = _fmt(v)
            centered(d, cx(i), top + 23, s, fit_font(s, "mono", 16, box - 6), ink)
            if links.get(i, "x") is None:
                _text(d, (cx(i) + box / 2 + 4, top + 23), "∅", font=font("mono", 14), fill=MUTED, anchor="lm")
        stack = {}
        for name, i in self.pointers.items():
            stack.setdefault(i, []).append(name)
        names = list(self.pointers)
        for i, group in stack.items():
            for k, name in enumerate(group):
                col = POINTER_COLORS[names.index(name) % len(POINTER_COLORS)]
                px = cx(i) if i is not None and 0 <= i < n else (x + n * step + 10 if i is not None and i >= n else x - 6)
                yy = top + 48 + k * 20
                if k == 0:
                    arrow(d, px, yy + 5, px, yy - 1, col, head=5)
                centered(d, px, yy + 16, name, font("bold", 14), col)


class Tree(Widget):
    """A binary tree from a level-order list (None = missing child).

    hl:     {level-order index: style}
    notes:  {level-order index: short text drawn under the node}
    """

    def __init__(self, level_order, label=None, hl=None, notes=None):
        self.vals = list(level_order)
        self.label = label
        self.hl = hl or {}
        self.notes = notes or {}
        self._build()

    def _build(self):
        vals = self.vals
        self.children = {}
        self.depth = {}
        if not vals or vals[0] is None:
            self.order = []
            return
        self.depth[0] = 0
        queue, i = [0], 1
        while queue and i < len(vals):
            p = queue.pop(0)
            kids = []
            for _ in range(2):
                if i < len(vals) and vals[i] is not None:
                    kids.append(i)
                    self.depth[i] = self.depth[p] + 1
                    queue.append(i)
                else:
                    kids.append(None)
                i += 1
            self.children[p] = kids
        self.order = []

        def inorder(n):
            if n is None:
                return
            l, r = self.children.get(n, [None, None])
            inorder(l)
            self.order.append(n)
            inorder(r)

        inorder(0)

    def height(self, width):
        levels = (max(self.depth.values()) + 1) if self.depth else 1
        return (20 if self.label else 0) + levels * 58 + 14

    def draw(self, d, x, y, width):
        _label(d, x, y, self.label)
        y += 20 if self.label else 0
        if not self.order:
            _text(d, (x, y + 10), "(empty tree)", font=font("sans", 14), fill=MUTED)
            return
        n = len(self.order)
        gap = min(56, (width - 40) / max(1, n))
        total = gap * (n - 1)
        x0 = x + (width - total) / 2
        pos = {node: (x0 + k * gap, y + 22 + self.depth[node] * 58) for k, node in enumerate(self.order)}
        r = max(13, min(20, gap / 2 - 2))
        for p, kids in self.children.items():
            for k in kids:
                if k is not None:
                    d.line([pos[p], pos[k]], fill="#8c959f", width=2)
        for node, (px, py) in pos.items():
            fill, border, ink = STYLES[self.hl.get(node)]
            d.ellipse([px - r, py - r, px + r, py + r], fill=fill, outline=border,
                      width=3 if node in self.hl else 1)
            s = _fmt(self.vals[node])
            centered(d, px, py, s, fit_font(s, "mono", 15, 2 * r - 4), ink)
            if node in self.notes:
                centered(d, px, py + r + 9, self.notes[node], font("mono", 11), "#8250df")


class Graph(Widget):
    """Nodes and edges.

    edges:    [(u, v)] or [(u, v, w)]
    pos:      {node: (fx, fy)} with fx, fy in 0..1; default is a circle
    hl_nodes: {node: style}; hl_edges: {(u, v): style}
    notes:    {node: short text drawn beside the node (a distance, a colour)}
    """

    def __init__(self, nodes, edges, label=None, pos=None, hl_nodes=None, hl_edges=None,
                 directed=False, notes=None, h=230):
        self.nodes = list(nodes)
        self.edges = list(edges)
        self.label = label
        self.pos = pos
        self.hl_nodes = hl_nodes or {}
        self.hl_edges = hl_edges or {}
        self.directed = directed
        self.notes = notes or {}
        self.h = h

    def height(self, width):
        return (20 if self.label else 0) + self.h + 10

    def _edge_style(self, u, v):
        st = self.hl_edges.get((u, v))
        if st is None and not self.directed:
            st = self.hl_edges.get((v, u))
        return st

    def draw(self, d, x, y, width):
        _label(d, x, y, self.label)
        y += 20 if self.label else 0
        n = len(self.nodes)
        r = 18
        if self.pos:
            pos = {k: (x + 30 + fx * (width - 60), y + 24 + fy * (self.h - 48)) for k, (fx, fy) in self.pos.items()}
        else:
            cxm, cym = x + width / 2, y + self.h / 2
            rad = min(width / 2 - 40, self.h / 2 - 26)
            pos = {k: (cxm + rad * math.cos(2 * math.pi * i / n - math.pi / 2),
                       cym + rad * math.sin(2 * math.pi * i / n - math.pi / 2)) for i, k in enumerate(self.nodes)}
        for e in self.edges:
            u, v = e[0], e[1]
            st = self._edge_style(u, v)
            col = STYLES[st][1] if st else "#8c959f"
            wdt = 3 if st else 2
            (x1, y1), (x2, y2) = pos[u], pos[v]
            ang = math.atan2(y2 - y1, x2 - x1)
            sx, sy = x1 + r * math.cos(ang), y1 + r * math.sin(ang)
            ex, ey = x2 - r * math.cos(ang), y2 - r * math.sin(ang)
            if self.directed:
                arrow(d, sx, sy, ex, ey, col, width=wdt, head=9)
            else:
                d.line([(sx, sy), (ex, ey)], fill=col, width=wdt)
            if len(e) > 2:
                mx, my = (x1 + x2) / 2, (y1 + y2) / 2
                ox, oy = -math.sin(ang) * 11, math.cos(ang) * 11
                s = _fmt(e[2])
                f = font("mono", 12)
                tw = text_w(s, f)
                d.rectangle([mx + ox - tw / 2 - 3, my + oy - 8, mx + ox + tw / 2 + 3, my + oy + 8], fill=BG)
                centered(d, mx + ox, my + oy, s, f, col if st else MUTED)
        for k in self.nodes:
            px, py = pos[k]
            st = self.hl_nodes.get(k)
            fill, border, ink = STYLES[st]
            d.ellipse([px - r, py - r, px + r, py + r], fill=fill, outline=border, width=3 if st else 1)
            s = _fmt(k)
            centered(d, px, py, s, fit_font(s, "mono", 14, 2 * r - 4), ink)
            if k in self.notes:
                s = _fmt(self.notes[k])
                f = font("bold", 12)
                tw = text_w(s, f)
                d.rounded_rectangle([px + r - 4, py - r - 12, px + r + tw + 4, py - r + 4], 4,
                                    fill="#fbefff", outline="#8250df")
                _text(d, (px + r, py - r - 4), s, font=f, fill="#512a97", anchor="lm")


class KV(Widget):
    """A map or set drawn as chips: {key: value} or a list of keys."""

    def __init__(self, mapping, label=None, hl=None, empty="{ }"):
        if isinstance(mapping, dict):
            self.items = [f"{_fmt(k)}: {_fmt(v)}" for k, v in mapping.items()]
            self.keys = list(mapping)
        else:
            self.items = [_fmt(k) for k in mapping]
            self.keys = list(mapping)
        self.label = label
        self.hl = hl or {}
        self.empty = empty

    def _layout(self, width):
        f = font("mono", 14)
        rows, cur, cx = [], [], 0
        for k, s in zip(self.keys, self.items):
            w = text_w(s, f) + 18
            if cur and cx + w > width:
                rows.append(cur)
                cur, cx = [], 0
            cur.append((k, s, w))
            cx += w + 8
        if cur:
            rows.append(cur)
        return rows

    def height(self, width):
        return (20 if self.label else 0) + max(1, len(self._layout(width))) * 32 + 6

    def draw(self, d, x, y, width):
        _label(d, x, y, self.label)
        y += 20 if self.label else 0
        rows = self._layout(width)
        if not rows:
            _text(d, (x, y + 6), self.empty, font=font("mono", 14), fill=MUTED)
            return
        f = font("mono", 14)
        for r, row in enumerate(rows):
            cx = x
            for k, s, w in row:
                fill, border, ink = STYLES[self.hl.get(k)]
                d.rounded_rectangle([cx, y + r * 32, cx + w, y + r * 32 + 26], 13, fill=fill, outline=border,
                                    width=2 if k in self.hl else 1)
                _text(d, (cx + 9, y + r * 32 + 13), s, font=f, fill=ink, anchor="lm")
                cx += w + 8


class Intervals(Widget):
    """Intervals on a number line, one per row.

    hl: {index: style}; marks: {name: position} vertical marker lines
    """

    def __init__(self, intervals, label=None, hl=None, axis=None, marks=None):
        self.iv = [tuple(i) for i in intervals]
        self.label = label
        self.hl = hl or {}
        lo = min((a for a, _ in self.iv), default=0)
        hi = max((b for _, b in self.iv), default=1)
        self.axis = axis or (lo, hi)
        self.marks = marks or {}

    def height(self, width):
        return (20 if self.label else 0) + len(self.iv) * 24 + 34

    def draw(self, d, x, y, width):
        _label(d, x, y, self.label)
        y += 20 if self.label else 0
        lo, hi = self.axis
        span = max(1, hi - lo)
        x0, x1 = x + 70, x + width - 20

        def px(v):
            return x0 + (v - lo) / span * (x1 - x0)

        for i, (a, b) in enumerate(self.iv):
            fill, border, ink = STYLES[self.hl.get(i)]
            yy = y + i * 24
            _text(d, (x, yy + 9), f"[{_fmt(a)},{_fmt(b)}]", font=font("mono", 12), fill=MUTED, anchor="lm")
            d.rounded_rectangle([px(a), yy + 2, max(px(b), px(a) + 4), yy + 18], 5, fill=fill, outline=border,
                                width=2 if i in self.hl else 1)
        ya = y + len(self.iv) * 24 + 6
        d.line([(x0, ya), (x1, ya)], fill=BORDER, width=1)
        step = max(1, round(span / 10))
        v = lo
        while v <= hi:
            d.line([(px(v), ya - 3), (px(v), ya + 3)], fill=BORDER)
            centered(d, px(v), ya + 12, _fmt(v), font("mono", 11), MUTED)
            v += step
        for k, (name, v) in enumerate(self.marks.items()):
            col = POINTER_COLORS[k % len(POINTER_COLORS)]
            d.line([(px(v), y - 4), (px(v), ya)], fill=col, width=2)
            centered(d, px(v), y - 10, name, font("bold", 11), col)


class Text(Widget):
    """Status lines such as 'sum = 2 + 15 = 17 > 9'. Lines may be (text, style)."""

    def __init__(self, *lines, mono=True, size=15):
        self.lines = [l if isinstance(l, tuple) else (l, None) for l in lines]
        self.kind = "mono" if mono else "sans"
        self.size = size

    def height(self, width):
        return len(self.lines) * (self.size + 9) + 4

    def draw(self, d, x, y, width):
        for i, (s, st) in enumerate(self.lines):
            col = STYLES[st][1] if st else INK
            f = fit_font(s, self.kind, self.size, width)
            _text(d, (x, y + i * (self.size + 9)), s, font=f, fill=col)


class Row(Widget):
    """Widgets side by side, splitting the width by the given weights."""

    def __init__(self, *widgets, weights=None, gap=24):
        self.widgets = widgets
        self.weights = weights or [1] * len(widgets)
        self.gap = gap

    def _widths(self, width):
        total = sum(self.weights)
        free = width - self.gap * (len(self.widgets) - 1)
        return [free * w / total for w in self.weights]

    def height(self, width):
        return max(w.height(ww) for w, ww in zip(self.widgets, self._widths(width)))

    def draw(self, d, x, y, width):
        for w, ww in zip(self.widgets, self._widths(width)):
            w.draw(d, x, y, ww)
            x += ww + self.gap


# --------------------------------------------------------------------------- frames


class Frame:
    """One step: widgets top to bottom, then the caption.

    say:     what the speaker says during this step (keep it to one or two sentences)
    speaker: "Candidate" by default; use "Interviewer" for questions
    """

    def __init__(self, *widgets, say="", speaker="Candidate", ms=None):
        self.widgets = widgets
        self.say = say
        self.speaker = speaker
        self.ms = ms


def _caption_lines(frame, width):
    f = font("sans", 16)
    words = frame.say.split()
    lines, cur = [], ""
    prefix_w = text_w(frame.speaker + ":  ", font("bold", 16))
    for w in words:
        trial = (cur + " " + w).strip()
        limit = width - (prefix_w if not lines else 0)
        if text_w(trial, f) > limit and cur:
            lines.append(cur)
            cur = w
        else:
            cur = trial
    if cur:
        lines.append(cur)
    return lines or [""]


def render(title, pattern, frames, out_path, ms=2200, last_ms=4000, per_word_ms=170):
    """Write frames as an animated GIF. Durations grow with caption length."""
    frames = list(frames)
    if not frames:
        raise ValueError("no frames")
    inner = WIDTH - 2 * PAD
    title_h = 54
    body_h = max(sum(w.height(inner) + 14 for w in fr.widgets) for fr in frames)
    cap_lines = max(len(_caption_lines(fr, inner - 24)) for fr in frames)
    cap_h = 22 * cap_lines + 26
    foot_h = 26
    H = title_h + 18 + body_h + cap_h + foot_h + 10

    images, durations = [], []
    for k, fr in enumerate(frames):
        img = Image.new("RGB", (WIDTH, int(H)), BG)
        d = ImageDraw.Draw(img)
        # title bar
        d.rectangle([0, 0, WIDTH, title_h], fill=BAR)
        _text(d, (PAD, title_h / 2), title, font=fit_font(title, "bold", 20, WIDTH - 260), fill="#ffffff", anchor="lm")
        pf = font("bold", 13)
        pw = text_w(pattern, pf)
        d.rounded_rectangle([WIDTH - PAD - pw - 20, 15, WIDTH - PAD, title_h - 15], 12, fill="#1f6feb")
        _text(d, (WIDTH - PAD - pw - 10, title_h / 2), pattern, font=pf, fill="#ffffff", anchor="lm")
        # widgets
        y = title_h + 18
        for w in fr.widgets:
            w.draw(d, PAD, y, inner)
            y += w.height(inner) + 14
        # caption
        cy = title_h + 18 + body_h
        is_q = fr.speaker.lower().startswith("interviewer")
        fill, border = ("#fff8c5", "#d4a72c") if is_q else ("#f6f8fa", "#d0d7de")
        d.rounded_rectangle([PAD, cy, WIDTH - PAD, cy + cap_h - 8], 10, fill=fill, outline=border)
        lines = _caption_lines(fr, inner - 24)
        sp = fr.speaker + ":  "
        _text(d, (PAD + 12, cy + 10), sp, font=font("bold", 16), fill="#9a6700" if is_q else "#0969da")
        sw = text_w(sp, font("bold", 16))
        for i, line in enumerate(lines):
            _text(d, (PAD + 12 + (sw if i == 0 else 0), cy + 10 + i * 22), line, font=font("sans", 16), fill=INK)
        # progress
        fy = H - foot_h - 4
        n = len(frames)
        bw = inner - 90
        d.rounded_rectangle([PAD, fy + 8, PAD + bw, fy + 14], 3, fill="#eaeef2")
        d.rounded_rectangle([PAD, fy + 8, PAD + bw * (k + 1) / n, fy + 14], 3, fill="#1f6feb")
        _text(d, (WIDTH - PAD, fy + 11), f"step {k + 1}/{n}", font=font("mono", 12), fill=MUTED, anchor="rm")

        images.append(img.quantize(colors=96, method=Image.Quantize.MEDIANCUT, dither=Image.Dither.NONE))
        words = len(fr.say.split())
        base = fr.ms or max(ms, 900 + per_word_ms * words)
        durations.append(last_ms + base // 2 if k == n - 1 else base)

    os.makedirs(os.path.dirname(out_path), exist_ok=True)
    images[0].save(out_path, save_all=True, append_images=images[1:], duration=durations,
                   loop=0, optimize=True, disposal=1)
    return out_path


def wrap(s, width=80):
    return "\n".join(textwrap.wrap(s, width))
