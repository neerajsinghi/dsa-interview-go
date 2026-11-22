from walkthrough import Array, Frame, Text

TITLE = "Range Sum Query — Immutable"
PATTERN = "Prefix sums"


def frames():
    nums = [-2, 0, 3, -5, 2, -1]
    queries = [((0, 2), 1), ((2, 5), -1)]
    n = len(nums)

    yield Frame(
        Array(nums, label="nums"),
        Text("SumRange(0, 2) → ?", "SumRange(2, 5) → ?"),
        speaker="Interviewer",
        say="nums never changes, but I'll ask many range-sum queries. Build something that answers "
            "SumRange(l, r) fast.",
    )
    yield Frame(
        Array(nums, label="nums", hl={i: "compare" for i in range(2, 6)}, window=(2, 5)),
        Text("brute force: add nums[l..r] for every query", "→ O(n) per query, O(n·q) total"),
        say="Adding the range each time is O(n) per query. With many queries the same additions "
            "repeat again and again, so I'd rather pay once up front.",
    )

    p = [0] * (n + 1)
    shown = [0] + [None] * n
    yield Frame(
        Array(nums, label="nums"),
        Array(shown, label="P  (P[i] = sum of nums[0..i-1])", hl={0: "special"}),
        Text("P[0] = 0  (the empty prefix)"),
        say="I build a prefix array P with one extra slot. P[0] is 0, the sum of nothing. "
            "That leading zero removes the special case for ranges starting at 0.",
    )
    for i, x in enumerate(nums):
        p[i + 1] = p[i] + x
        shown[i + 1] = p[i + 1]
        if i < 2 or i == n - 1:
            say = (f"P[{i + 1}] = P[{i}] + nums[{i}] = {p[i]} + {x} = {p[i + 1]}. "
                   "Each entry reuses the one before it, so the build is one pass.")
            if i == n - 1:
                say = f"Last entry: P[{n}] = {p[n]}, the total of the whole array. The build cost O(n) once."
            yield Frame(
                Array(nums, label="nums", hl={i: "active"}),
                Array(shown, label="P  (P[i] = sum of nums[0..i-1])", hl={i: "visited", i + 1: "active"}),
                Text(f"P[{i + 1}] = P[{i}] + nums[{i}] = {p[i]} + ({x}) = {p[i + 1]}"),
                say=say,
            )
        if i == n - 2:  # entries 3..5 are summarised in one frame
            yield Frame(
                Array(nums, label="nums", hl={2: "active", 3: "active", 4: "active"}),
                Array(shown, label="P  (P[i] = sum of nums[0..i-1])", hl={3: "active", 4: "active", 5: "active"}),
                Text(f"P[3] = {p[3]},  P[4] = {p[4]},  P[5] = {p[5]}"),
                say="Same rule for the next three entries: add the current number to the previous total.",
            )

    for (l, r), want in queries:
        got = p[r + 1] - p[l]
        assert got == want, (l, r, got, want)
        yield Frame(
            Array(nums, label="nums", hl={i: "window" for i in range(l, r + 1)}, window=(l, r)),
            Array(p, label="P", hl={r + 1: "match", l: "reject"}),
            Text(f"SumRange({l}, {r}) = P[{r + 1}] − P[{l}] = {p[r + 1]} − ({p[l]}) = {got}", (f"answer: {got}", "match")),
            say=(f"P[{r + 1}] covers nums[0..{r}] and P[{l}] covers nums[0..{l - 1}]. "
                 f"Subtracting leaves exactly nums[{l}..{r}], which is {got}." if l else
                 f"P[{r + 1}] covers nums[0..{r}] and P[0] is the empty prefix, so the answer is "
                 f"{p[r + 1]} − 0 = {got}. No special case for l = 0."),
        )

    yield Frame(
        Array(p, label="P"),
        Text("SumRange(0, 2) = 1    SumRange(2, 5) = -1",
             ("build O(n) · each query O(1) · space O(n)", "match")),
        say="Each query is now one subtraction. Building P is O(n) time and O(n) space, "
            "and every query after that is O(1).",
    )
