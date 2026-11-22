from walkthrough import Frame, Grid, Row, Text

TITLE = "Range Sum Query 2D — Immutable"
PATTERN = "2-D prefix sums"


def mat_grid(m, hl=None, label="matrix"):
    return Grid(m, label=label, hl=hl, row_labels=list(range(len(m))), col_labels=list(range(len(m[0]))), cell=44)


def p_grid(p, hl=None, label="P (padded)"):
    return Grid(p, label=label, hl=hl, row_labels=list(range(len(p))), col_labels=list(range(len(p[0]))), cell=44)


def frames():
    matrix = [[3, 0, 1],
              [5, 6, 3],
              [1, 2, 0]]
    queries = [((1, 1, 2, 2), 11), ((0, 0, 0, 0), 3)]
    rows, cols = len(matrix), len(matrix[0])

    region = {(r, c): "window" for r in range(1, 3) for c in range(1, 3)}
    yield Frame(
        Row(mat_grid(matrix, hl=region), Text("SumRegion(1, 1, 2, 2) → ?", "SumRegion(0, 0, 0, 0) → ?", size=14),
            weights=[1, 1.4]),
        speaker="Interviewer",
        say="The matrix never changes. Answer many queries for the sum of a rectangle given its "
            "top-left and bottom-right corners.",
    )
    yield Frame(
        Row(mat_grid(matrix, hl={k: "compare" for k in region}),
            Text("brute force: add every cell", "→ O(rows·cols) per query", size=14), weights=[1, 1.4]),
        say="Adding every cell in the rectangle costs up to rows times cols per query. With many "
            "queries I'd rather precompute, like 1-D prefix sums but in two directions.",
    )

    p = [[0] * (cols + 1) for _ in range(rows + 1)]
    pad = {(0, j): "special" for j in range(cols + 1)}
    pad.update({(i, 0): "special" for i in range(rows + 1)})
    yield Frame(
        Row(mat_grid(matrix), p_grid(p, hl=pad), weights=[1, 1.3]),
        Text("P[i][j] = sum of matrix[0..i-1][0..j-1]   (row 0 and column 0 are zeros)"),
        say="P has one extra row and column of zeros. P[i][j] is the sum of everything above and to "
            "the left of (i, j). The padding removes edge cases.",
    )

    for i in range(1, rows + 1):
        for j in range(1, cols + 1):
            p[i][j] = matrix[i - 1][j - 1] + p[i - 1][j] + p[i][j - 1] - p[i - 1][j - 1]
        # show the rule on the last cell of the row
        j = 2 if i == 1 else cols
        hl = {(i, c): "match" for c in range(1, cols + 1)}
        hl.update({(i - 1, j): "active", (i, j - 1): "active", (i - 1, j - 1): "reject"})
        hl[(i, j)] = "match"
        formula = (f"P[{i}][{j}] = m[{i - 1}][{j - 1}] + P[{i - 1}][{j}] + P[{i}][{j - 1}] − P[{i - 1}][{j - 1}]"
                   f" = {matrix[i - 1][j - 1]} + {p[i - 1][j]} + {p[i][j - 1]} − {p[i - 1][j - 1]} = {p[i][j]}")
        if i == 1:
            say = ("Each cell is its matrix value plus the box above plus the box to the left, minus the "
                   "corner, which both of those already counted.")
        elif i == 2:
            say = f"Row {i} uses the same rule. Each cell needs only three neighbours already filled in, so the build is one pass."
        else:
            say = f"Row {i} finishes the table. P[3][3] = {p[3][3]} is the sum of the whole matrix. Build cost: O(rows·cols)."
        yield Frame(
            Row(mat_grid(matrix, hl={(i - 1, j - 1): "active"}), p_grid(p_partial(p, i), hl=hl), weights=[1, 1.3]),
            Text(formula),
            say=say,
        )

    for (r1, c1, r2, c2), want in queries:
        got = p[r2 + 1][c2 + 1] - p[r1][c2 + 1] - p[r2 + 1][c1] + p[r1][c1]
        assert got == want, (r1, c1, r2, c2, got)
        mhl = {(r, c): "window" for r in range(r1, r2 + 1) for c in range(c1, c2 + 1)}
        phl = {(r2 + 1, c2 + 1): "match", (r1, c2 + 1): "reject", (r2 + 1, c1): "reject", (r1, c1): "active"}
        expr = (f"P[{r2 + 1}][{c2 + 1}] − P[{r1}][{c2 + 1}] − P[{r2 + 1}][{c1}] + P[{r1}][{c1}]"
                f" = {p[r2 + 1][c2 + 1]} − {p[r1][c2 + 1]} − {p[r2 + 1][c1]} + {p[r1][c1]} = {got}")
        if (r1, c1) == (1, 1):
            say = ("Whole box P[3][3], minus the strip above, P[1][3], minus the strip to the left, P[3][1], "
                   "plus the corner P[1][1], which I subtracted twice. That's 11.")
        else:
            say = "A single cell works too: P[1][1] minus two zero padding cells plus a zero corner gives 3."
        yield Frame(
            Row(mat_grid(matrix, hl=mhl), p_grid(p, hl=phl), weights=[1, 1.3]),
            Text(f"SumRegion({r1}, {c1}, {r2}, {c2})", expr, (f"answer: {got}", "match")),
            say=say,
        )

    yield Frame(
        Row(mat_grid(matrix), p_grid(p), weights=[1, 1.3]),
        Text("SumRegion(1,1,2,2) = 11    SumRegion(0,0,0,0) = 3",
             ("build O(rows·cols) · query O(1) · space O(rows·cols)", "match")),
        say="Every query is four lookups. Building P is O(rows·cols) time and space, then each query "
            "is O(1).",
    )


def p_partial(p, upto):
    """P with rows below `upto` blanked out, so the table fills in row by row."""
    return [row if i <= upto else [0] + [None] * (len(row) - 1) for i, row in enumerate(p)]
