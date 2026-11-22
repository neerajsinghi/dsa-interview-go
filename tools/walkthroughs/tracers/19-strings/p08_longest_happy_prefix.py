from walkthrough import Array, Frame, Text

TITLE = "Longest Happy Prefix"
PATTERN = "KMP prefix function"


def frames():
    s = "level"
    n = len(s)
    pi = [None] * n

    yield Frame(
        Array(list(s), label="s"),
        Array([""] * n, label="pi", show_index=False),
        Text('"level"'),
        speaker="Interviewer",
        say="A happy prefix is a non-empty prefix that is also a suffix, but not the whole string. "
            "Return the longest one, or an empty string.",
    )
    yield Frame(
        Array(list(s), label="s"),
        Array([""] * n, label="pi", show_index=False),
        Text("brute force: for k = n−1 down to 1, compare s[:k] with s[n−k:] → O(n²)"),
        say="I could try every length k from longest to shortest and compare prefix with suffix. "
            "Each comparison is O(k), so that's O(n²) in the worst case.",
    )
    pi[0] = 0
    yield Frame(
        Array(list(s), label="s", hl={0: "active"}),
        Array(pi, label="pi", show_index=False, hl={0: "active"}),
        Text("pi[i] = longest proper border of s[:i+1]", "pi[0] = 0 (one char has no proper border)",
             ("answer = s[:pi[n−1]]", "special")),
        say="This is exactly the KMP prefix function. pi[i] is the longest proper prefix of s[:i+1] that is also "
            "its suffix, so the answer is just s[:pi[n−1]].",
    )

    for i in range(1, n):
        k = pi[i - 1]
        steps = []
        while k > 0 and s[i] != s[k]:
            steps.append(f"s[{i}]='{s[i]}' ≠ s[{k}]='{s[k]}' → k = pi[{k - 1}] = {pi[k - 1]}")
            k = pi[k - 1]
        ok = s[i] == s[k]
        cmp_line = f"try to extend border k = {pi[i - 1]}: s[{i}]='{s[i]}' vs s[{k}]='{s[k]}'"
        if ok:
            k += 1
        pi[i] = k
        hl_s = {i: "match" if ok else "reject", k - 1 if ok else k: "match" if ok else "compare"}
        lines = [cmp_line] + steps + [(f"pi[{i}] = {k}", "match" if ok else None)]
        if ok:
            say = (f"'{s[i]}' equals s[{k - 1}], so the border grows to {k}: the prefix \"{s[:k]}\" "
                   f"is also a suffix of \"{s[:i + 1]}\".")
        else:
            say = (f"'{s[i]}' doesn't match s[{k}] = '{s[k]}', and k is already 0, so there's no border to fall "
                   f"back to. pi[{i}] = 0.")
        yield Frame(
            Array(list(s), label="s", hl=hl_s, pointers={"i": i, "k": k - 1 if ok else k}),
            Array(pi, label="pi", show_index=False, hl={i: "match" if ok else "active"}),
            Text(*lines),
            say=say,
        )

    ans = s[:pi[n - 1]]
    assert ans == "l"
    yield Frame(
        Array(list(s), label="s", hl={0: "match", n - 1: "match"}),
        Array(pi, label="pi", show_index=False, hl={n - 1: "match"}),
        Text((f'answer: s[:pi[{n - 1}]] = s[:{pi[n - 1]}] = "{ans}"', "match"), "time O(n) · space O(n)"),
        say="pi[4] = 1, so the answer is \"l\". k rises by at most one per step, so the fall-backs total O(n). "
            "That's O(n) time and O(n) space for pi.",
    )
