from walkthrough import BORDER, MUTED, POINTER_COLORS, STYLES, Frame, Text, Widget, centered, font, arrow

TITLE = "Container With Most Water"
PATTERN = "Two pointers"


class Bars(Widget):
    """Vertical bars with an optional water rectangle between two bars."""

    def __init__(self, heights, hl=None, pointers=None, water=None, best=None, h=150):
        self.heights = heights
        self.hl = hl or {}
        self.pointers = pointers or {}
        self.water = water  # (l, r): fill between the bars up to the shorter one
        self.best = best    # (l, r): dashed outline of the best container so far
        self.h = h

    def height(self, width):
        return self.h + 16 + 42 + 14

    def draw(self, d, x, y, width):
        y += 14
        n = len(self.heights)
        c = min(56, width / n)
        top = max(self.heights)
        unit = (self.h - 8) / top
        base = y + self.h

        def bx(i):
            return x + i * c

        if self.water:
            l, r = self.water
            level = min(self.heights[l], self.heights[r])
            d.rectangle([bx(l) + c / 2, base - level * unit, bx(r) + c / 2, base], fill="#cfe8ff")
        if self.best:
            l, r = self.best
            level = min(self.heights[l], self.heights[r])
            x1, x2, yt = bx(l) + c / 2, bx(r) + c / 2, base - level * unit
            for xx in range(int(x1), int(x2), 10):
                d.line([(xx, yt), (min(xx + 5, x2), yt)], fill="#1a7f37", width=2)
        for i, v in enumerate(self.heights):
            fill, border, ink = STYLES[self.hl.get(i)]
            if self.hl.get(i) is None:
                fill, border = "#d0d7de", "#8c959f"
            d.rectangle([bx(i) + c * 0.3, base - v * unit, bx(i) + c * 0.7, base], fill=fill, outline=border,
                        width=2 if i in self.hl else 1)
            centered(d, bx(i) + c / 2, base - v * unit - 9, v, font("mono", 12), ink if i in self.hl else MUTED)
        d.line([(x, base), (x + n * c, base)], fill=BORDER, width=2)
        for i in range(n):
            centered(d, bx(i) + c / 2, base + 9, i, font("mono", 11), MUTED)
        for k, (name, i) in enumerate(self.pointers.items()):
            col = POINTER_COLORS[k % len(POINTER_COLORS)]
            cx = bx(i) + c / 2
            arrow(d, cx, base + 28, cx, base + 19, col, head=5)
            centered(d, cx, base + 38, name, font("bold", 14), col)


def frames():
    heights = [1, 8, 6, 2, 5, 4, 8, 3, 7]

    yield Frame(
        Bars(heights),
        Text("area(i, j) = (j - i) × min(h[i], h[j])"),
        speaker="Interviewer",
        say="Each number is the height of a vertical line. Pick two lines that, with the x-axis, "
            "hold the most water, and return that area.",
    )
    yield Frame(
        Bars(heights),
        Text("brute force: every pair (i, j)", "n(n-1)/2 = 36 pairs here → O(n²)"),
        say="Trying every pair is O(n²). Instead I'll start with the widest container and only ever drop "
            "a line that can't be part of anything better.",
    )

    best, best_pair = 0, None
    l, r = 0, len(heights) - 1
    step = 0
    while l < r:
        area = (r - l) * min(heights[l], heights[r])
        new_best = area > best
        if new_best:
            best, best_pair = area, (l, r)
        short = "l" if heights[l] < heights[r] else "r"
        hl = {l: "active", r: "active"}
        hl[l if short == "l" else r] = "reject"
        why = (f"h[l] = {heights[l]} < h[r] = {heights[r]} → l++" if short == "l"
               else f"h[r] = {heights[r]} <= h[l] = {heights[l]} → r--")
        if step == 0:
            say = (f"Start at both ends: width {r - l}, height {min(heights[l], heights[r])}, area {area}. "
                   f"The left line is shorter, so every narrower container using it holds less. Drop it.")
        elif new_best:
            say = (f"Width {r - l} times height {min(heights[l], heights[r])} is {area}, a new best. "
                   f"Then I drop the shorter line, {'left' if short == 'l' else 'right'}, since it caps everything narrower.")
        elif heights[l] == heights[r]:
            say = (f"Area {area} doesn't beat {best}. Both lines are {heights[l]}; on a tie either can go, "
                   "because neither can do better with a narrower partner. The code moves r.")
        else:
            say = (f"Area {area} doesn't beat {best}. The shorter line is "
                   f"{'left' if short == 'l' else 'right'}, so I move that pointer in.")
        yield Frame(
            Bars(heights, hl=hl, pointers={"l": l, "r": r}, water=(l, r), best=best_pair),
            Text(f"area = ({r} - {l}) × min({heights[l]}, {heights[r]}) = {area}" + ("   new best!" if new_best else ""),
                 (f"best = {best}", "match"), (why, "reject")),
            say=say,
        )
        if heights[l] < heights[r]:
            l += 1
        else:
            r -= 1
        step += 1

    assert best == 49  # doc comment: [1, 8, 6, 2, 5, 4, 8, 3, 7] → 49
    bl, br = best_pair
    yield Frame(
        Bars(heights, hl={bl: "match", br: "match"}, water=best_pair),
        Text((f"answer: {best}  (lines {bl} and {br}: {br - bl} × {min(heights[bl], heights[br])})", "match"),
             "time O(n): each step moves one pointer inward", "space O(1)"),
        say=f"The pointers met, so the answer is {best}, from lines {bl} and {br}. "
            "Each step discards one line for good, so it's O(n) time and O(1) space.",
    )
