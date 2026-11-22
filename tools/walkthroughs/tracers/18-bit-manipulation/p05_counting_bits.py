from walkthrough import Array, Frame, Grid, Text

TITLE = "Counting Bits"
PATTERN = "Bits · DP"

W = 8
COLS = list(range(W - 1, -1, -1))


def bits(x):
    x &= (1 << W) - 1
    return [(x >> k) & 1 for k in COLS]


def frames():
    # The doc comment's second example (n=5): n=2 is too short to show the recurrence.
    n = 5
    expected = [0, 1, 1, 2, 1, 2]

    yield Frame(
        Array([None] * (n + 1), label="ans"),
        speaker="Interviewer",
        say=f"For every i from 0 to n, report how many 1 bits i has. Take n = {n}. Can you do it in one linear pass?",
    )
    yield Frame(
        Array([None] * (n + 1), label="ans"),
        Text("naive: popcount each i separately", "O(n · bits) unless popcount is one instruction"),
        say="I could popcount every number on its own. That costs up to one step per bit for each i. "
            "Instead I want to reuse answers I already have.",
    )

    ans = [0] * (n + 1)
    yield Frame(
        Array(ans[:1] + [None] * n, label="ans", hl={0: "match"}),
        Text("ans[0] = 0", "ans[i] = ans[i & (i-1)] + 1"),
        say="i & (i-1) is i with its lowest 1 cleared. It's smaller than i, so its answer is ready. "
            "i has exactly one more 1 bit than it.",
    )
    for i in range(1, n + 1):
        j = i & (i - 1)
        ans[i] = ans[j] + 1
        hl = {k: "visited" for k in range(i)}
        hl[j] = "compare"
        hl[i] = "match"
        low = (i & -i).bit_length() - 1
        ghl = {(0, COLS.index(low)): "reject"}
        yield Frame(
            Array(ans[:i + 1] + [None] * (n - i), label="ans", hl=hl,
                  pointers={"i": i, "i&(i-1)": j} if j != i else {"i": i}),
            Grid([bits(i), bits(j)], row_labels=[f"i = {i}", f"i & (i-1) = {j}"], col_labels=COLS, hl=ghl),
            Text(f"ans[{i}] = ans[{j}] + 1 = {ans[j]} + 1 = {ans[i]}"),
            say=(f"{i} with its lowest 1 removed is {j}. ans[{j}] is {ans[j]}, so ans[{i}] = {ans[i]}."
                 if j != 0 else
                 f"{i} is a power of two, so clearing its only 1 gives 0. ans[{i}] = ans[0] + 1 = 1."),
        )

    assert ans == expected
    yield Frame(
        Array(ans, label="ans", hl={k: "match" for k in range(n + 1)}),
        Text((f"answer: {ans}", "match"), "O(n) time, O(n) space for the output"),
        say=f"The answer is {ans}. Each entry is one AND and one lookup: O(n) time. "
            "The only space is the output itself.",
    )
