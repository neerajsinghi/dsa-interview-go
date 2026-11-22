from walkthrough import Array, Frame, Grid, KV, Row, Text

TITLE = "Palindrome Partitioning"
PATTERN = "Backtracking over cuts"


def fmt(parts):
    return "[" + " ".join(parts) + "]"


def frames():
    s = "aab"
    n = len(s)

    # isPal table, filled with i descending so isPal[i+1] is ready.
    is_pal = [[False] * n for _ in range(n)]
    for i in range(n - 1, -1, -1):
        for j in range(i, n):
            is_pal[i][j] = s[i] == s[j] and (j - i < 2 or is_pal[i + 1][j - 1])

    def table(hl=None):
        rows = [[("T" if is_pal[i][j] else "F") if j >= i else "" for j in range(n)] for i in range(n)]
        base = {(i, j): ("match" if is_pal[i][j] else "reject") for i in range(n) for j in range(i, n)}
        base.update({(i, j): "empty" for i in range(n) for j in range(i)})
        if hl:
            base.update(hl)
        return Grid(rows, label="isPal[i][j]", row_labels=[f"i={i}" for i in range(n)],
                    col_labels=[f"j={j}" for j in range(n)], hl=base, cell=40)

    def chars(start, end=None, style="active"):
        hl = {k: "visited" for k in range(start)}
        if end is not None:
            for k in range(start, end + 1):
                hl[k] = style
        return Array(list(s), label="s", hl=hl, window=(start, end) if end is not None else None,
                     pointers={"start": start} if start < n else {}, cell=52)

    def view(start, end, path, res, style="active", new=None, cell_hl=None):
        return [
            Row(chars(start, end, style), table(cell_hl), weights=[1, 1]),
            Row(KV(path, label="path (pieces so far)", empty="[ ]"),
                KV([fmt(r) for r in res], label=f"results ({len(res)})", hl={new: "match"} if new else None, empty="[ ]"),
                weights=[1, 1]),
        ]

    yield Frame(
        chars(0),
        Text(f's = "{s}"', "split s so that every piece is a palindrome; return all splits"),
        speaker="Interviewer",
        say='Split the string into pieces so that every piece reads the same both ways. Return every such split.',
    )
    yield Frame(
        chars(0),
        Text(f"brute force: 2^(n-1) = {2 ** (n - 1)} ways to place cuts, check each piece",
             ("re-checks the same substrings again and again", "reject")),
        say="I could try every set of cut positions and check each piece. "
            "Instead I'll cut left to right, only take pieces that are palindromes, and look them up in a table.",
    )
    yield Frame(
        Row(chars(0), table(), weights=[1, 1]),
        Text("isPal[i][j] = s[i] == s[j] && (j-i < 2 || isPal[i+1][j-1])", "filled with i going down so row i+1 is ready"),
        say="First I precompute isPal. A substring is a palindrome if its ends match and the inside is one. "
            "Then every check during the search is O(1).",
    )

    res, path = [], []
    events = []

    def backtrack(start):
        if start == n:
            res.append(list(path))
            events.append(("leaf", start, None, list(path), list(res)))
            return
        for end in range(start, n):
            if not is_pal[start][end]:
                events.append(("reject", start, end, list(path), list(res)))
                continue
            path.append(s[start:end + 1])
            events.append(("take", start, end, list(path), list(res)))
            backtrack(end + 1)
            path.pop()

    backtrack(0)

    for ev, start, end, p, r in events:
        piece = s[start:end + 1] if end is not None else None
        if ev == "take":
            say = (f'From index {start}, "{piece}" is a palindrome, so I take it and recurse from {end + 1}.')
            yield Frame(*view(start, end, p, r, cell_hl={(start, end): "active"}),
                        Text(f'isPal[{start}][{end}] = T → take "{piece}", recurse from {end + 1}'), say=say)
        elif ev == "reject":
            say = f'"{piece}" is not a palindrome, so no valid partition can use it as the next piece. I skip this cut.'
            yield Frame(*view(start, end, p, r, style="reject", cell_hl={(start, end): "compare"}),
                        Text(f'isPal[{start}][{end}] = F', (f'"{piece}" rejected', "reject")), say=say)
        else:
            f = fmt(p)
            yield Frame(*view(start, None, p, r, new=f),
                        Text(f"start == len(s) → record {f}"),
                        say=f"start reached the end, so every character is covered. {f} is a full partition. "
                            "I copy it and unchoose the last piece.")

    assert res == [["a", "a", "b"], ["aa", "b"]]
    yield Frame(
        KV([fmt(x) for x in res], label=f"results ({len(res)})", hl={fmt(x): "match" for x in res}),
        Text("answer = [a a b] [aa b]", "O(n · 2^n) time, O(n²) space for isPal"),
        say="Two partitions. There can be up to 2 to the n minus 1 of them and each copy is O(n): "
            "O(n times 2 to the n) time. The isPal table takes O(n squared) space.",
    )
