from walkthrough import (BORDER, MUTED, POINTER_COLORS, STYLES, Frame, Text, Widget, arrow, centered, font,
                         stack_widget, _label)

TITLE = "Largest Rectangle in Histogram"
PATTERN = "Monotonic stack"


class Bars(Widget):
    """A histogram: one bar per height, with highlights, an optional rectangle and an index pointer.

    hl:   {index: style}
    rect: (left, right, height) inclusive bar range drawn as an outlined rectangle
    ptr:  index of the "i" pointer; len(heights) points at the sentinel slot
    """

    PLOT_H = 150

    def __init__(self, heights, label=None, hl=None, rect=None, rect_style="match", ptr=None, sentinel=True):
        self.h = list(heights)
        self.label = label
        self.hl = hl or {}
        self.rect = rect
        self.rect_style = rect_style
        self.ptr = ptr
        self.sentinel = sentinel

    def height(self, width):
        return (20 if self.label else 0) + 18 + self.PLOT_H + 18 + 30

    def draw(self, d, x, y, width):
        _label(d, x, y, self.label)
        y += 20 if self.label else 0
        slots = len(self.h) + (1 if self.sentinel else 0)
        c = min(64, width / slots)
        top = y + 18
        base = top + self.PLOT_H
        unit = self.PLOT_H / max(1, max(self.h))
        d.line([(x, base), (x + slots * c, base)], fill=BORDER, width=2)
        for i, v in enumerate(self.h):
            fill, border, ink = STYLES[self.hl.get(i)]
            x0, x1 = x + i * c + 4, x + (i + 1) * c - 4
            d.rectangle([x0, base - v * unit, x1, base], fill=fill, outline=border, width=2 if i in self.hl else 1)
            centered(d, (x0 + x1) / 2, base - v * unit - 9, v, font("mono", 13), ink)
            centered(d, (x0 + x1) / 2, base + 10, i, font("mono", 11), MUTED)
        if self.sentinel:
            i = len(self.h)
            centered(d, x + i * c + c / 2, base - 9, "0", font("mono", 13), MUTED)
            centered(d, x + i * c + c / 2, base + 10, f"{i} (end)", font("mono", 11), MUTED)
        if self.rect:
            l, r, h = self.rect
            col = STYLES[self.rect_style][1]
            d.rectangle([x + l * c + 1, base - h * unit, x + (r + 1) * c - 1, base], outline=col, width=4)
            area = h * (r - l + 1)
            cx, cy = x + (l + r + 1) * c / 2, base - h * unit / 2
            label = f"area {area}"
            tw = font("bold", 15).getlength(label)
            d.rounded_rectangle([cx - tw / 2 - 6, cy - 12, cx + tw / 2 + 6, cy + 12], 6, fill="#ffffff", outline=col)
            centered(d, cx, cy, label, font("bold", 15), col)
        if self.ptr is not None:
            px = x + self.ptr * c + c / 2
            col = POINTER_COLORS[0]
            arrow(d, px, base + 26, px, base + 19, col, head=5)
            centered(d, px, base + 36, "i", font("bold", 14), col)


def frames():
    heights = [2, 1, 5, 6, 2, 3]
    want = 10
    n = len(heights)

    yield Frame(
        Bars(heights, label="heights", sentinel=False),
        speaker="Interviewer",
        say="Each bar has width 1. Find the area of the largest rectangle that fits entirely inside the "
            "histogram.",
    )
    yield Frame(
        Bars(heights, label="heights", sentinel=False, hl={2: "compare", 3: "compare"}, rect=(2, 3, 5),
             rect_style="compare"),
        Text("brute force: every left edge, extend right, track the min height → O(n²)"),
        say="Trying every left edge and extending right while tracking the min height is O(n²). Better: "
            "for each bar, find how far its full height can stretch.",
    )
    yield Frame(
        Bars(heights, label="heights  (+ sentinel 0 at the end)"),
        stack_widget([], label="stack of indices (increasing heights)"),
        Text("bar's rectangle: from previous shorter bar + 1 to next shorter bar − 1"),
        say="A bar's rectangle stretches left to the previous shorter bar and right to the next shorter "
            "one. An increasing stack finds both. A 0 at the end flushes the stack.",
    )

    best = 0
    st = []
    for i in range(n + 1):
        h = heights[i] if i < n else 0
        while st and heights[st[-1]] >= h:
            top = st.pop()
            left = st[-1] if st else -1
            width = i - left - 1
            area = heights[top] * width
            improved = area > best
            best = max(best, area)
            why_left = f"new top is index {left}" if left >= 0 else "stack is empty, so it reaches index 0"
            yield Frame(
                Bars(heights, label="heights  (+ sentinel 0 at the end)", hl={top: "match", **({i: "reject"} if i < n else {})},
                     rect=(left + 1, i - 1, heights[top]), rect_style="match" if improved else "compare", ptr=i),
                stack_widget(st, label="stack of indices (increasing heights)"),
                Text(f"pop {top} (h={heights[top]}): width = {i} − ({left}) − 1 = {width}, area = {area}",
                     (f"best = {best}", "match" if improved else None)),
                say=(f"Height {h} at index {i} is not taller than bar {top}, so bar {top} can't extend past {i - 1}. "
                     f"On the left, the {why_left}. Area {heights[top]}×{width} = {area}.")
                if i < n else
                f"The sentinel 0 is shorter than everything. Bar {top} pops: width {width}, area {area}. "
                f"The {why_left}.",
            )
        st.append(i)
        if i < n:
            yield Frame(
                Bars(heights, label="heights  (+ sentinel 0 at the end)", hl={i: "active"}, ptr=i),
                stack_widget(st, label="stack of indices (increasing heights)", hl={len(st) - 1: "active"}),
                Text(f"push {i} (h={h})", f"best = {best}"),
                say=f"Bar {i} (height {h}) is taller than the top of the stack, so I don't know its right "
                    "edge yet. I push it." if len(st) > 1 else
                    f"The stack is empty, so I push bar {i}. Its right edge isn't known yet.",
            )

    assert best == want, best
    yield Frame(
        Bars(heights, label="heights", sentinel=False, hl={2: "match", 3: "match"}, rect=(2, 3, 5)),
        Text(f"answer: {best}   (height 5 × width 2)", ("time O(n) · space O(n)", "match")),
        say="The largest area is 10, bars 2 and 3 at height 5. Every index is pushed and popped once, "
            "so it's O(n) time and O(n) space.",
    )
