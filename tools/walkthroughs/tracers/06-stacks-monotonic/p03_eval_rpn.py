from walkthrough import Array, Frame, Text, stack_widget

TITLE = "Evaluate Reverse Polish Notation"
PATTERN = "Stack"


def trunc_div(a, b):
    """Go's integer division truncates toward zero; Python's // floors."""
    q = abs(a) // abs(b)
    return q if (a >= 0) == (b >= 0) else -q


def frames():
    tokens = ["2", "1", "+", "3", "*"]
    want = 9

    yield Frame(
        Array(tokens, label="tokens"),
        speaker="Interviewer",
        say="Evaluate this postfix expression. Each operator applies to the two values just before it. "
            "Division truncates toward zero. Report an error if it's malformed.",
    )
    yield Frame(
        Array(tokens, label="tokens", hl={0: "compare", 1: "compare", 2: "compare"}),
        Text("brute force: find the first operator, replace 'a b op' with its value, repeat",
             "→ O(n) rescan per operator → O(n²)"),
        say="I could repeatedly find the first operator, collapse it with its two operands and rescan. "
            "That's O(n²). The values I need are always the most recent ones: a stack.",
    )

    st = []
    for i, tok in enumerate(tokens):
        hl = {**{j: "visited" for j in range(i)}, i: "active"}
        if tok in "+-*/":
            b = st.pop()
            a = st.pop()
            r = {"+": a + b, "-": a - b, "*": a * b, "/": trunc_div(a, b) if b else None}[tok]
            st.append(r)
            yield Frame(
                Array(tokens, label="tokens", hl=hl, pointers={"i": i}),
                stack_widget(st, label="operand stack", hl={len(st) - 1: "match"}),
                Text(f"pop b = {b}, pop a = {a}", (f"push a {tok} b = {a} {tok} {b} = {r}", "match")),
                say=f"'{tok}' pops two values. The top is the right operand, so b = {b} and a = {a}. "
                    f"I push {a} {tok} {b} = {r}. Order matters for - and /.",
            )
        else:
            st.append(int(tok))
            yield Frame(
                Array(tokens, label="tokens", hl=hl, pointers={"i": i}),
                stack_widget(st, label="operand stack", hl={len(st) - 1: "active"}),
                Text(f"'{tok}' is a number → push {tok}"),
                say=f"'{tok}' is a number. It's waiting for an operator, so I push it." if i == 0 else
                    f"'{tok}' is another number: push it on top.",
            )

    assert len(st) == 1 and st[0] == want
    yield Frame(
        Array(tokens, label="tokens", hl={j: "visited" for j in range(len(tokens))}),
        stack_widget(st, label="operand stack", hl={0: "match"}),
        Text("exactly one value left → answer 9   ((2 + 1) * 3)", ("time O(n) · space O(n)", "match")),
        say="Exactly one value is left, so the expression was well-formed and the answer is 9. "
            "One pass, stack up to n values: O(n) time and space.",
    )
