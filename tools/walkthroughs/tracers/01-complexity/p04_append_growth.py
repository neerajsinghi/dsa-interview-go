from walkthrough import Array, Frame, Grid, Text

TITLE = "Amortized Append"
PATTERN = "Amortized O(1)"

# Go runtime malloc size classes (bytes), enough for the sizes used here.
SIZE_CLASSES = [8, 16, 24, 32, 48, 64, 80, 96, 112, 128, 144, 160, 176, 192, 208, 224, 240, 256,
                288, 320, 352, 384, 416, 448, 480, 512, 576, 640, 704, 768, 896, 1024, 1152, 1280,
                1408, 1536, 1792, 2048, 2304, 2688, 3072, 3200, 3456, 4096, 4864, 5376, 6144, 6528,
                6784, 6912, 8192, 9472, 9728, 10240, 10880, 12288, 13568, 14336, 16384, 18432,
                19072, 20480, 21760, 24576, 27264, 28672, 32768]


def grow(old_cap, need, elem=8):
    """Port of the runtime's nextslicecap + roundupsize for an append of one int."""
    if old_cap == 0:
        newcap = need
    elif old_cap < 256:
        newcap = 2 * old_cap
    else:
        newcap = old_cap
        while newcap < need:
            newcap += (newcap + 3 * 256) // 4
    size = newcap * elem
    for c in SIZE_CLASSES:
        if c >= size:
            size = c
            break
    return size // elem


def append(s, cap, x):
    """Returns (new slice, new cap, elements copied)."""
    copied = 0
    if len(s) == cap:
        cap = grow(cap, len(s) + 1)
        copied = len(s)  # old elements move to the new array
    return s + [x], cap, copied


def count_growths(n):
    s, cap, growths = [], 0, 0
    for i in range(n):
        before = cap
        s, cap, _ = append(s, cap, i)
        if cap != before:
            growths += 1
    return growths


def view(s, cap, hl=None):
    vals = s + [None] * (cap - len(s))
    styles = {j: "empty" for j in range(len(s), cap)}
    styles.update(hl or {})
    return Array(vals, label=f"backing array   len={len(s)}  cap={cap}", hl=styles, cell=40)


def table(rows_l, rows_c, rows_k, hl_col=None):
    hl = {}
    if hl_col is not None:
        hl = {(r, hl_col): "active" for r in range(3)}
    return Grid([rows_l, rows_c, rows_k], label="history", row_labels=["len", "cap", "copies"],
                col_labels=None, hl=hl, cell=34)


def frames():
    yield Frame(
        Text("var s []int", "for i := range n { s = append(s, i) }", "CountGrowths(1) → 1    CountGrowths(1000) → ~11"),
        speaker="Interviewer",
        say="append sometimes copies the whole slice. So why do we call append O(1)? Show me with CountGrowths.",
    )
    yield Frame(
        Text("naive idea: grow the array by exactly 1 each time",
             "append #k copies k-1 elements → 0+1+2+…+(n-1) ≈ n²/2",
             ("n appends = O(n²)   (same cost as PrependN shifting every time)", "reject")),
        say="If the array grew by one slot each time, every append would copy everything. "
            "That's n²/2 copies in total. Prepending has that cost too.",
    )

    s, cap = [], 0
    L, C, K = [], [], []
    total_copies, growths = 0, 0
    narrated = {0, 1, 2, 3, 4, 8}
    for i in range(9):
        before = cap
        s, cap, copied = append(s, cap, i)
        total_copies += copied
        grew = cap != before
        growths += grew
        L.append(len(s))
        C.append(cap)
        K.append(copied)
        if i not in narrated:
            if i == 7:
                yield Frame(
                    view(s, cap, {j: "match" for j in range(5, 8)}),
                    table(L, C, K, hl_col=i),
                    Text(f"append 5, 6, 7: room left, copied 0   total copies so far = {total_copies}"),
                    say="The next three appends also fit in the spare room. No copying at all.",
                )
            continue
        if grew:
            hl = {len(s) - 1: "active"}
            hl.update({j: "compare" for j in range(len(s) - 1)})
            if i == 0:
                say = "The first append allocates an array of capacity 1. Nothing to copy yet."
            elif i == 8:
                say = f"Capacity runs out at 8, so it doubles to {cap} and copies 8 values. " \
                      "Big copies happen, but each one buys as much free room as it cost."
            else:
                say = f"Full again at {before}, so the runtime doubles capacity to {cap} and copies the {copied} old " \
                      f"value{'s' if copied > 1 else ''}."
            yield Frame(
                view(s, cap, hl),
                table(L, C, K, hl_col=i),
                Text(f"append {i}: cap {before} → {cap}, copied {copied}", ("grow + copy", "compare")),
                say=say,
            )
        else:
            yield Frame(
                view(s, cap, {len(s) - 1: "match"}),
                table(L, C, K, hl_col=i),
                Text(f"append {i}: cap {cap} has room, copied 0", ("cheap append", "match")),
                say="This time there's spare capacity, so append just writes the value. O(1), no copy.",
            )

    g1, g1000 = count_growths(1), count_growths(1000)
    assert g1 == 1
    assert 10 <= g1000 <= 12, g1000  # doc comment: "~11"; the exact count depends on the runtime policy
    caps, cap, copies = [], 0, 0
    for k in range(1000):
        if k == cap:
            copies += k
            cap = grow(cap, k + 1)
            caps.append(cap)
    yield Frame(
        Text(f"CountGrowths(1000) = {g1000} growths", "caps: " + " ".join(map(str, caps)),
             f"total copies = {copies}  < 2n = 2000",
             ("doubling below 256, then ~1.25× plus size-class rounding", "special")),
        say=f"For 1000 appends there are only {g1000} growths, and all the copying adds up to {copies}, "
            "less than 2n. Growth is geometric, so the copies form a halving series.",
    )
    yield Frame(
        Text(("total cost of n appends: O(n)  →  O(1) amortized each", "match"),
             "AppendN              O(n) total",
             "AppendNPreallocated  O(n), zero regrowth",
             ("PrependN             O(n²): shifts every element each time", "reject")),
        say="So n appends cost O(n) total, O(1) amortized each. If I know n, make([]int, 0, n) skips regrowth. "
            "Inserting at the front stays O(n²).",
    )
