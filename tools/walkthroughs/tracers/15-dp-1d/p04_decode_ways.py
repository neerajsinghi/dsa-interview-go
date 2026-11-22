from walkthrough import Array, Frame, Text

TITLE = "Decode Ways"
PATTERN = "1-D DP"


def frames():
    s = "226"
    n = len(s)
    dp = [None] * (n + 1)

    def widgets(hl_s=None, hl_dp=None):
        return (
            Array(list(s), label="s (digits)", hl=hl_s),
            Array(dp, label="ways[i] = decodings of the prefix s[:i]", hl=hl_dp),
        )

    yield Frame(
        *widgets(),
        Text("A=1, B=2, ..., Z=26 · \"0\" and \"06\" are not valid codes"),
        speaker="Interviewer",
        say="Letters are encoded as 1 to 26. Given the digits \"226\", how many letter strings could have produced it?",
    )
    yield Frame(
        *widgets(),
        Text("branch on a 1-digit or 2-digit piece at every position", ("same suffixes re-solved → O(2^n)", "reject")),
        say="The brute force branches on a one- or two-digit piece every time. That's exponential, "
            "and the same prefixes get counted again and again.",
    )

    dp[0] = 1
    yield Frame(
        *widgets(hl_dp={0: "match"}),
        Text("ways[0] = 1  (the empty prefix decodes one way: nothing)"),
        say="This is Climbing Stairs with checks. ways[i] counts decodings of the first i digits. "
            "The empty prefix has exactly one decoding, so ways[0] = 1.",
    )
    for i in range(1, n + 1):
        dp[i] = 0
        lines = []
        one = s[i - 1] != "0"
        two = i >= 2 and s[i - 2] != "0" and int(s[i - 2:i]) <= 26
        hl_dp = {i: "active"}
        if one:
            dp[i] += dp[i - 1]
            hl_dp[i - 1] = "compare"
            lines.append(f"last digit \"{s[i-1]}\" is 1-9 ✓  → + ways[{i-1}] = {dp[i-1]}")
        else:
            lines.append((f"last digit \"{s[i-1]}\" is 0 ✗", "reject"))
        if i >= 2:
            if two:
                dp[i] += dp[i - 2]
                hl_dp[i - 2] = "compare"
                lines.append(f"last two \"{s[i-2:i]}\" is 10-26 ✓  → + ways[{i-2}] = {dp[i-2]}")
            else:
                lines.append((f"last two \"{s[i-2:i]}\" is not 10-26 ✗", "reject"))
        lines.append((f"ways[{i}] = {dp[i]}", "active"))
        hl_s = {k: "window" for k in range(max(0, i - 2) if two else i - 1, i)}
        says = {
            1: "The prefix \"2\" can only end with the single letter B, so it has ways[0] = 1 decoding.",
            2: "For \"22\" the last letter is either \"2\" alone, after ways[1], or \"22\" as V, after ways[0]. "
               "Both are valid, so I add them: 2.",
            3: "For \"226\": \"6\" alone is valid and adds ways[2] = 2. \"26\" is Z and adds ways[1] = 1. Total 3.",
        }
        yield Frame(*widgets(hl_s, hl_dp), Text(*lines), say=says[i])

    assert dp[n] == 3
    yield Frame(
        *widgets(hl_dp={n: "match"}),
        Text(f"answer: ways[{n}] = {dp[n]}  (BBF, BZ, VF)", ("O(n) time · O(1) space", "match")),
        say="The answer is 3: BBF, BZ and VF. Each cell reads only the two before it, "
            "so the Go code keeps two variables: O(n) time, O(1) space.",
    )
