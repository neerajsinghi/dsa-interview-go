from walkthrough import BORDER, MUTED, POINTER_COLORS, STYLES, Frame, Text, Widget, arrow, centered, font

TITLE = "Trapping Rain Water"
PATTERN = "Two pointers"


class Bars(Widget):
    """Elevation bars with trapped water drawn on top of settled bars."""

    def __init__(self, heights, hl=None, pointers=None, water=None, levels=None, h=130):
        self.heights = heights
        self.hl = hl or {}
        self.pointers = pointers or {}
        self.water = water or {}    # {i: units of water above bar i}
        self.levels = levels or {}  # {name: (height, colour index)} dashed running-max lines
        self.h = h

    def height(self, width):
        return 14 + self.h + 16 + 42

    def draw(self, d, x, y, width):
        y += 14
        n = len(self.heights)
        c = min(56, width / n)
        top = max(self.heights)
        unit = (self.h - 8) / top
        base = y + self.h

        def bx(i):
            return x + i * c

        for i, w in self.water.items():
            if w > 0:
                v = self.heights[i]
                d.rectangle([bx(i) + 1, base - (v + w) * unit, bx(i) + c - 1, base - v * unit], fill="#9fd0ff")
                centered(d, bx(i) + c / 2, base - (v + w / 2) * unit, f"+{w}", font("bold", 12), "#0a3069")
        for i, v in enumerate(self.heights):
            st = self.hl.get(i)
            fill, border, ink = STYLES[st]
            if st is None:
                fill, border = "#d0d7de", "#8c959f"
            if v > 0:
                d.rectangle([bx(i) + 1, base - v * unit, bx(i) + c - 1, base], fill=fill, outline=border,
                            width=2 if st else 1)
            else:
                d.line([(bx(i) + 3, base - 1), (bx(i) + c - 3, base - 1)], fill=border, width=3)
        d.line([(x, base), (x + n * c, base)], fill=BORDER, width=2)
        for i in range(n):
            centered(d, bx(i) + c / 2, base + 9, self.heights[i], font("mono", 12), MUTED)
        for k, (name, i) in enumerate(self.pointers.items()):
            col = POINTER_COLORS[k % len(POINTER_COLORS)]
            cx = bx(i) + c / 2
            arrow(d, cx, base + 28, cx, base + 19, col, head=5)
            centered(d, cx, base + 38, name, font("bold", 14), col)


def frames():
    heights = [0, 1, 0, 2, 1, 0, 1, 3, 2, 1, 2, 1]

    yield Frame(
        Bars(heights),
        Text("heights = " + str(heights)),
        speaker="Interviewer",
        say="Each number is a bar one unit wide. After it rains, how many units of water are trapped between the bars?",
    )
    yield Frame(
        Bars(heights),
        Text("water above i = min(max left of i, max right of i) - h[i]",
             "brute force: rescan both sides for every i → O(n²)"),
        say="Water above a bar is the smaller of the tallest bar to its left and to its right, minus its height. "
            "Rescanning both sides for every bar is O(n²).",
    )

    l, r = 0, len(heights) - 1
    left_max = right_max = 0
    total = 0
    water = {}
    step = 0
    prev = None
    while l <= r:
        left_max = max(left_max, heights[l])
        right_max = max(right_max, heights[r])
        ptr = {"l": l, "r": r}
        settled = {i: "visited" for i in water}
        if left_max <= right_max:
            w = left_max - heights[l]
            water[l] = w
            total += w
            side, i, m = "left", l, left_max
            l += 1
        else:
            w = right_max - heights[r]
            water[r] = w
            total += w
            side, i, m = "right", r, right_max
            r -= 1
        hl = dict(settled)
        hl[i] = "match" if w else "active"
        line = (f"leftMax = {left_max}, rightMax = {right_max} → settle the {side} bar (index {i})")
        calc = (f"water[{i}] = {m} - {heights[i]} = {w}   total = {total}", "match" if w else None)
        smaller = "smaller" if left_max != right_max else "equal"
        mname = "leftMax" if side == "left" else "rightMax"
        if step == 0:
            say = (f"leftMax is {left_max} and rightMax is {right_max}. The smaller max is the bottleneck, "
                   f"because the other side already has a wall at least that tall. Bar {i} holds {m} - {heights[i]} = {w}.")
        elif w and prev == (side, smaller, True):
            say = f"Same situation: {mname} is still {m}, so bar {i} holds {m} - {heights[i]} = {w}. Running total: {total}."
        elif w and smaller == "equal":
            say = (f"Both maxima are {m}, and on a tie the code settles the left side. The right side has a wall "
                   f"of at least {m}, so bar {i} holds {m} - {heights[i]} = {w}.")
        elif w:
            say = (f"{mname} = {m} is the smaller max, and the other side has a wall at least that tall. "
                   f"So bar {i} holds {m} - {heights[i]} = {w}.")
        elif smaller == "equal":
            say = (f"Both maxima are {m}, so on the tie I settle the left bar {i}. "
                   f"It is as tall as leftMax, so it holds nothing.")
        else:
            say = (f"The {side} side has the smaller max, so I settle bar {i}. It is as tall as {mname}, "
                   "so it holds nothing; it only raises the wall.")
        prev = (side, smaller, bool(w))
        yield Frame(
            Bars(heights, hl=hl, pointers=ptr, water=dict(water)),
            Text(line, calc),
            say=say,
        )
        step += 1

    assert total == 6  # doc comment: [0, 1, 0, 2, 1, 0, 1, 3, 2, 1, 2, 1] → 6
    yield Frame(
        Bars(heights, water=water),
        Text((f"answer: {total} units", "match"), "time O(n): every bar is settled exactly once",
             "space O(1): two pointers and two running maxima"),
        say=f"Every bar is settled once, and the total is {total}. That's O(n) time, "
            "and unlike the prefix-max version it needs only O(1) extra space.",
    )
