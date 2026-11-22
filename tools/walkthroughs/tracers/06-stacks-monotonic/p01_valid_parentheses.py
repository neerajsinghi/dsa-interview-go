from walkthrough import Array, Frame, Text, stack_widget

TITLE = "Valid Parentheses"
PATTERN = "Stack"


def frames():
    s = "([]{})"
    want = True
    chars = list(s)
    opener = {")": "(", "]": "[", "}": "{"}

    yield Frame(
        Array(chars, label="s"),
        speaker="Interviewer",
        say="The string holds only ( ) [ ] { }. Is it balanced? Every opener must be closed by the "
            "same kind, in the reverse order it was opened.",
    )
    yield Frame(
        Array(chars, label="s", hl={1: "compare", 2: "compare"}),
        Text("brute force: delete \"()\", \"[]\", \"{}\" pairs until nothing changes", "→ up to n/2 rounds of O(n) → O(n²)"),
        say="I could keep deleting adjacent matching pairs until the string stops changing. Each round "
            "is O(n) and there can be n/2 rounds, so O(n²).",
    )
    yield Frame(
        Array(chars, label="s"),
        stack_widget([], label="stack of openers"),
        Text(f"len(s) = {len(s)} is even → keep going"),
        say="The most recent unclosed opener is the only one allowed to close next. That's the top of "
            "a stack. First, an odd length can never pair up, but 6 is even.",
    )

    st = []
    ok = len(s) % 2 == 0
    for i, c in enumerate(s):
        if c in "([{":
            st.append(c)
            yield Frame(
                Array(chars, label="s", hl={**{j: "visited" for j in range(i)}, i: "active"}, pointers={"i": i}),
                stack_widget(st, label="stack of openers", hl={len(st) - 1: "active"}),
                Text(f"'{c}' is an opener → push"),
                say=f"'{c}' opens something. I can't check it yet, so I push it and wait for its closer."
                if i < 3 else f"'{c}' is another opener, so it goes on top. It must close before anything under it.",
            )
        else:
            top = st[-1] if st else None
            if top != opener[c]:
                ok = False
                break
            yield Frame(
                Array(chars, label="s", hl={**{j: "visited" for j in range(i)}, i: "match"}, pointers={"i": i}),
                stack_widget(st, label="stack of openers", hl={len(st) - 1: "match"}),
                Text(f"'{c}' closes: top is '{top}' = opener['{c}'] → pop"),
                say=f"'{c}' is a closer. The top of the stack is '{top}', its partner, so this pair is "
                    "matched and I pop it.",
            )
            st.pop()
    ok = ok and not st

    assert ok == want
    yield Frame(
        Array(chars, label="s", hl={j: "match" for j in range(len(s))}),
        stack_widget(st, label="stack of openers"),
        Text("end of string, stack empty → true", ("time O(n) · space O(n)", "match")),
        say="At the end the stack is empty, so every opener was closed: true. One pass, and the stack "
            "holds at most n characters: O(n) time and space.",
    )
