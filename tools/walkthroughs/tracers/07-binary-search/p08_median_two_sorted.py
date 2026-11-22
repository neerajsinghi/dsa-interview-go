from walkthrough import STYLES, MUTED, Array, Frame, Row, Text, Widget, centered, fit_font, font, _fmt

TITLE = "Median of Two Sorted Arrays"
PATTERN = "Binary search"

INF = float("inf")


class Cut(Widget):
    """Two sorted arrays drawn one above the other, each cut by a red bar.

    The cells left of the bars form the "left half" of the merged order.
    hl: {(row, index): style} with row 0 = a, 1 = b.
    """

    CELL = 50

    def __init__(self, a, b, i, j, hl=None, names=("a", "b")):
        self.rows = [list(a), list(b)]
        self.cuts = [i, j]
        self.hl = hl or {}
        self.names = names

    def height(self, width):
        return 22 + 2 * (self.CELL + 26) + 6

    def draw(self, d, x, y, width):
        c = self.CELL
        d.text((x + 60, y), "left half", font=font("bold", 13), fill=STYLES["active"][1])
        y += 22
        cut_x = x + 60 + c * max(self.cuts)  # line all cuts up for easy reading
        d.text((cut_x + 14, y - 22), "right half", font=font("bold", 13), fill=STYLES["compare"][1])
        for r, (vals, cut, name) in enumerate(zip(self.rows, self.cuts, self.names)):
            yy = y + r * (c + 26)
            d.text((x, yy + c / 2), f"{name}", font=font("bold", 16), fill=MUTED, anchor="lm")
            # Shift the row so its cut sits on the shared cut line.
            x0 = cut_x - cut * c
            for k, v in enumerate(vals):
                st = self.hl.get((r, k), "active" if k < cut else "compare")
                fill, border, ink = STYLES[st]
                d.rounded_rectangle([x0 + k * c + 3, yy, x0 + (k + 1) * c - 3, yy + c - 6], 6,
                                    fill=fill, outline=border, width=2 if (r, k) in self.hl else 1)
                s = _fmt(v)
                centered(d, x0 + k * c + c / 2, yy + (c - 6) / 2, s, fit_font(s, "mono", 17, c - 8), ink)
            d.text((cut_x + 8, yy + c + 6), f"{'i' if r == 0 else 'j'} = {cut}", font=font("mono", 13), fill="#cf222e", anchor="lm")
        d.line([(cut_x, y - 4), (cut_x, y + 2 * (c + 26) - 26)], fill="#cf222e", width=3)


def at(xs, k):
    if k < 0:
        return -INF
    if k >= len(xs):
        return INF
    return xs[k]


def median(a, b):
    """Port of FindMedianSortedArrays. Yields (event, info) for each step."""
    if len(a) > len(b):
        a, b = b, a
    m, n = len(a), len(b)
    half = (m + n + 1) // 2
    yield "setup", dict(a=a, b=b, m=m, n=n, half=half)
    lo, hi = 0, m
    while lo <= hi:
        i = lo + (hi - lo) // 2
        j = half - i
        aL, aR, bL, bR = at(a, i - 1), at(a, i), at(b, j - 1), at(b, j)
        info = dict(a=a, b=b, lo=lo, hi=hi, i=i, j=j, aL=aL, aR=aR, bL=bL, bR=bR)
        if aL > bR:
            yield "too_many", info
            hi = i - 1
        elif bL > aR:
            yield "too_few", info
            lo = i + 1
        else:
            left_max = max(aL, bL)
            if (m + n) % 2 == 1:
                info["result"] = float(left_max)
            else:
                info["result"] = (left_max + min(aR, bR)) / 2
            info["left_max"] = left_max
            info["right_min"] = min(aR, bR)
            yield "done", info
            return


def checks(f):
    v = lambda x: _fmt(x)
    return (f"a[i-1] = {v(f['aL'])}  a[i] = {v(f['aR'])}   b[j-1] = {v(f['bL'])}  b[j] = {v(f['bR'])}")


def frames():
    a0, b0 = [1, 3], [2]

    yield Frame(
        Row(Array(a0, label="a"), Array(b0, label="b")),
        Text("a = [1, 3], b = [2]   → median 2.0", "a = [1, 2], b = [3, 4] → median 2.5"),
        speaker="Interviewer",
        say="Two sorted arrays. Return the median of all their values together, in O(log(min(m, n))).",
    )
    yield Frame(
        Row(Array(a0, label="a"), Array(b0, label="b")),
        Text("merge both, pick the middle → O(m + n) time and space"),
        say="Merging like merge sort and taking the middle is O(m + n). I don't need the whole merge, "
            "only where the merged order splits into a left and right half.",
    )

    result = None
    for ev, f in median(a0, b0):
        if ev == "setup":
            a, b, half = f["a"], f["b"], f["half"]
            yield Frame(
                Cut(a, b, 0, half),
                Text(f"swap so a is the shorter: a = {a}, b = {b}",
                     f"left half size = (m + n + 1) / 2 = {half};  i from a, j = {half} - i from b",
                     "valid cut: a[i-1] <= b[j]  and  b[j-1] <= a[i]"),
                say=f"I take i values from a and j = {half} - i from b as the left half. "
                    "The cut is right when every left value is <= every right value. I binary-search i over the shorter array.",
            )
            continue
        a, b, i, j = f["a"], f["b"], f["i"], f["j"]
        head = f"lo = {f['lo']}, hi = {f['hi']}  →  i = {i}, j = {j}"
        if ev == "too_few":
            yield Frame(
                Cut(a, b, i, j, hl={(1, j - 1): "reject"} if j - 1 >= 0 else {}),
                Text(head, checks(f), (f"b[j-1] = {_fmt(f['bL'])} > a[i] = {_fmt(f['aR'])} → too few from a → lo = {i + 1}", "reject")),
                say=f"b[j-1] = {_fmt(f['bL'])} is on the left but is bigger than a[i] = {_fmt(f['aR'])} on the right. "
                    f"I took too few from a, so lo = {i + 1}.",
            )
        elif ev == "too_many":
            yield Frame(
                Cut(a, b, i, j, hl={(0, i - 1): "reject"}),
                Text(head, checks(f), (f"a[i-1] = {_fmt(f['aL'])} > b[j] = {_fmt(f['bR'])} → too many from a → hi = {i - 1}", "reject")),
                say=f"a[i-1] = {_fmt(f['aL'])} is on the left but bigger than b[j] = {_fmt(f['bR'])}. I took too many from a, so hi = {i - 1}.",
            )
        else:
            result = f["result"]
            hl = {}
            lm = f["left_max"]
            if i - 1 >= 0 and a[i - 1] == lm:
                hl[(0, i - 1)] = "match"
            elif j - 1 >= 0:
                hl[(1, j - 1)] = "match"
            yield Frame(
                Cut(a, b, i, j, hl=hl),
                Text(head, checks(f), (f"both checks pass; total {len(a) + len(b)} is odd → median = max(left) = {_fmt(lm)}", "match")),
                say=f"Now {_fmt(f['aL'])} <= {_fmt(f['bR'])} and {_fmt(f['bL'])} <= {_fmt(f['aR'])}, so the cut is valid. "
                    f"The total is odd, so the median is the largest left value, {_fmt(lm)}.",
            )
    assert result == 2.0

    # Second doc example: the even case uses both sides of the cut.
    a1, b1 = [1, 2], [3, 4]
    steps = list(median(a1, b1))
    lines = []
    for ev, f in steps[1:]:
        if ev == "too_few":
            lines.append(f"i={f['i']}, j={f['j']}: b[j-1]={f['bL']} > a[i]={f['aR']} → lo={f['i'] + 1}")
        elif ev == "too_many":
            lines.append(f"i={f['i']}, j={f['j']}: a[i-1]={f['aL']} > b[j]={f['bR']} → hi={f['i'] - 1}")
    ev, f = steps[-1]
    res2 = f["result"]
    assert res2 == 2.5
    lines.append(f"i={f['i']}, j={f['j']}: valid → (max({_fmt(f['aL'])}, {_fmt(f['bL'])}) + min({_fmt(f['aR'])}, {_fmt(f['bR'])})) / 2 = {res2}")
    yield Frame(
        Row(Array(a1, label="a"), Array(b1, label="b")),
        Text("total = 4 values, even"),
        speaker="Interviewer",
        say="What about an even total, like a = [1, 2] and b = [3, 4]? And what if a cut leaves one side empty?",
    )
    yield Frame(
        Cut(f["a"], f["b"], f["i"], f["j"]),
        Text(*lines, (f"median = {res2};  O(log min(m, n)) time, O(1) space", "match")),
        say=f"Even total: average the largest left value and the smallest right value, here {res2}. Empty sides count as "
            "-∞ or +∞ so they never fail a check. O(log min(m, n)) time.",
    )
