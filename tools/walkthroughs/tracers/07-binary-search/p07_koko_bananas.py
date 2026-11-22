from walkthrough import Array, Frame, Row, Text

TITLE = "Koko Eating Bananas"
PATTERN = "Search on the answer"


def hours_at(piles, k):
    return sum((p + k - 1) // k for p in piles)  # ceil(p / k) without floats


def speeds(lo, hi, known, mid=None, mid_style="compare"):
    """Speeds 1..max drawn as cells; known = {speed: True/False}. hi is exclusive."""
    top = max(known.get("_max", 0), 1)
    vals = list(range(1, top + 1))
    hl = {}
    for k in vals:
        if k in known:
            hl[k - 1] = "match" if known[k] else "reject"
        elif k < lo or k >= hi:
            hl[k - 1] = "visited"
    ptr = {"lo": lo - 1, "hi": hi - 1}
    if mid is not None:
        hl[mid - 1] = mid_style
        ptr["mid"] = mid - 1
    return Array(vals, label="speed k (cells are speeds, hi is exclusive)", hl=hl, pointers=ptr, show_index=False,
                 window=(lo - 1, hi - 2) if lo < hi else None)


def frames():
    piles, h = [3, 6, 7, 11], 8
    mx = max(piles)
    known = {"_max": mx}

    yield Frame(
        Array(piles, label="piles"),
        Text(f"h = {h} hours", "each hour: eat up to k bananas from one pile"),
        speaker="Interviewer",
        say="Koko eats up to k bananas an hour from one pile. Find the smallest whole speed k that finishes all piles in 8 hours.",
    )
    yield Frame(
        Array(piles, label="piles"),
        Text("try k = 1, 2, 3, ... until hours(k) <= h", "→ O(max · n) = O(11 · 4) here, huge for big piles"),
        say="Brute force tries k = 1, 2, 3 and so on. Each check costs O(n), and k can go up to the largest pile, "
            "so it's O(max · n).",
    )
    yield Frame(
        speeds(1, mx + 1, known),
        Text("hours(k) = Σ ceil(pile / k)   only goes down as k grows",
             "can finish at k?  F F F ... T T T   → first true", "k = max(piles) = 11 always works (1 hour per pile)"),
        say="Faster never takes more hours, so 'can finish at speed k' is false, then true. I binary-search the "
            "answer range [1, 11] for the first true.",
    )

    lo, hi = 1, mx + 1  # FirstTrue(1, max+1, pred): hi is exclusive
    while lo < hi:
        mid = lo + (hi - lo) // 2
        per = [(p + mid - 1) // mid for p in piles]
        hrs = hours_at(piles, mid)
        ok = hrs <= h
        detail = Row(Array(piles, label="piles", show_index=False),
                     Array(per, label=f"hours at k = {mid}", show_index=False, hl={i: "active" for i in range(len(per))}))
        calc = f"hours({mid}) = " + " + ".join(map(str, per)) + f" = {hrs}"
        if ok:
            frame = Frame(speeds(lo, hi, known, mid, "match"), detail,
                          Text(calc, (f"{hrs} <= {h} → feasible → hi = {mid}", "match")),
                          say=f"At speed {mid} she needs {hrs} hours, within {h}. {mid} works, but a slower speed might too, "
                              f"so I keep {mid} as a candidate: hi = {mid}.")
            known[mid] = True
            hi = mid
        else:
            frame = Frame(speeds(lo, hi, known, mid, "reject"), detail,
                          Text(calc, (f"{hrs} > {h} → too slow → lo = {mid + 1}", "reject")),
                          say=f"At speed {mid} she needs {hrs} hours, more than {h}. Every slower speed is worse, "
                              f"so the answer is above {mid}: lo = {mid + 1}.")
            known[mid] = False
            lo = mid + 1
        yield frame

    answer = lo
    assert answer == 4
    yield Frame(
        speeds(lo, lo, known),
        Text(f"lo == hi == {answer}", (f"minimum speed = {answer}", "match"),
             "O(n · log max) time, O(1) space"),
        say=f"lo meets hi at {answer}: speed {answer} works and {answer - 1} doesn't. Each check is O(n) and there are "
            "log(max) checks, so O(n log max) time and O(1) space.",
    )
