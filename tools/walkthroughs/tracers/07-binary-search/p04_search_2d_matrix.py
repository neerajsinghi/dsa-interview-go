from walkthrough import Array, Frame, Grid, Text

TITLE = "Search a 2D Matrix"
PATTERN = "Binary search"

MATRIX = [[1, 3, 5, 7],
          [10, 11, 16, 20],
          [23, 30, 34, 60]]


def views(lo, hi, mid=None, mid_style="compare"):
    m, n = len(MATRIX), len(MATRIX[0])
    flat = [v for row in MATRIX for v in row]
    fhl, ghl = {}, {}
    for k in range(m * n):
        if k < lo or k > hi:
            fhl[k] = ghl[(k // n, k % n)] = "visited"
    ptr = {"lo": lo, "hi": hi}
    if mid is not None:
        fhl[mid] = ghl[(mid // n, mid % n)] = mid_style
        ptr["mid"] = mid
    grid = Grid(MATRIX, label="matrix (m = 3, n = 4)", hl=ghl, row_labels=[0, 1, 2], col_labels=[0, 1, 2, 3], cell=44)
    arr = Array(flat, label="the same values as one flat array (k = row·n + col)", hl=fhl, pointers=ptr,
                window=(lo, hi) if lo <= hi else None)
    return grid, arr


def frames():
    target = 3
    m, n = len(MATRIX), len(MATRIX[0])

    yield Frame(
        *views(0, m * n - 1),
        Text(f"target = {target}"),
        speaker="Interviewer",
        say="Each row is sorted and each row starts above the previous row's last value. Is 3 in the matrix? Aim for O(log(m·n)).",
    )
    yield Frame(
        *views(0, m * n - 1),
        Text("check every cell → O(m·n)"),
        say="Checking every cell is O(m·n). But read row by row, the matrix is one sorted list of 12 values, "
            "so I can binary-search it without copying it.",
    )
    yield Frame(
        *views(0, m * n - 1),
        Text("flat index k  →  row = k / n,  col = k % n", "lo = 0, hi = m·n - 1 = 11"),
        say="I search flat indices 0 to 11. A flat index k lives at row k / 4, column k % 4. That's the only new idea.",
    )

    found = False
    lo, hi = 0, m * n - 1
    while lo <= hi:
        mid = lo + (hi - lo) // 2
        r, c = mid // n, mid % n
        v = MATRIX[r][c]
        where = f"mid = {mid} → ({mid}/{n}, {mid}%{n}) = ({r}, {c}) → {v}"
        if v == target:
            found = True
            yield Frame(
                *views(lo, hi, mid, "match"),
                Text(where, (f"{v} == {target} → true", "match"), "O(log(m·n)) time, O(1) space"),
                say=f"Flat index {mid} maps to row {r}, column {c}, which holds {v}. Found it, so I return true. "
                    "That's O(log(m·n)) time and O(1) space.",
            )
            break
        if v < target:
            yield Frame(*views(lo, hi, mid), Text(where, (f"{v} < {target} → lo = {mid + 1}", "reject")),
                        say=f"Index {mid} is row {r}, column {c}: value {v}, too small. Everything before it is smaller too, so lo = {mid + 1}.")
            lo = mid + 1
        else:
            yield Frame(*views(lo, hi, mid), Text(where, (f"{v} > {target} → hi = {mid - 1}", "reject")),
                        say=f"Index {mid} is row {r}, column {c}: value {v}, too big. Everything after it is bigger, so hi = {mid - 1}.")
            hi = mid - 1
    assert found is True
