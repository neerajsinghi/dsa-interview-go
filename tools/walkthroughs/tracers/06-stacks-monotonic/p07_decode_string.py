from walkthrough import KV, Array, Frame, Text

TITLE = "Decode String"
PATTERN = "Stack of contexts"


def frame_label(f):
    prefix, k = f
    return f'("{prefix}", {k})'


def stack_view(st):
    """Stack of saved contexts as chips, top on the right (wide enough for ("aaa", 2))."""
    vals = [frame_label(f) for f in st]
    return KV(vals, label="stack of (prefix, k), top on the right", hl={vals[-1]: "special"} if vals else None,
              empty="(empty)")


def frames():
    s = "3[a]2[bc]"
    want = "aaabcbc"
    chars = list(s)

    yield Frame(
        Array(chars, label="s"),
        speaker="Interviewer",
        say="Decode a string where k[text] means text repeated k times. Encodings can nest and k can "
            "have several digits. Assume the input is well-formed.",
    )
    yield Frame(
        Array(chars, label="s", hl={0: "compare", 1: "compare", 2: "compare", 3: "compare"}),
        Text("brute force: find an innermost k[...] , expand it, rescan → O(n · output)"),
        say="I could repeatedly find an innermost bracket, expand it in place and rescan. That rebuilds "
            "the string many times. Nested brackets suggest a stack.",
    )
    yield Frame(
        Array(chars, label="s"),
        stack_view([]),
        Text('cur = ""   k = 0'),
        say="I keep cur, the text built at the current level, and k, the number being read. On '[' I save "
            "both and start fresh. On ']' I restore and repeat.",
    )

    st = []
    cur = ""
    k = 0
    for i, c in enumerate(s):
        hl = {**{j: "visited" for j in range(i)}, i: "active"}
        if c.isdigit():
            k = k * 10 + int(c)
            line = f"digit → k = {k}"
            say = f"'{c}' is a digit, so k becomes {k}. Multi-digit counts work because I do k*10 + digit."
            style = "active"
        elif c == "[":
            st.append((cur, k))
            line = f'push ("{cur}", {k}); cur = "", k = 0'
            say = (f"'[' opens a group. I push what I had, prefix \"{cur}\" and count {k}, then reset cur and k "
                   "for the inside.")
            cur, k = "", 0
            style = "active"
        elif c == "]":
            prefix, rep = st.pop()
            inner = cur
            cur = prefix + inner * rep
            line = f'pop ("{prefix}", {rep}): cur = "{prefix}" + "{inner}"×{rep} = "{cur}"'
            say = (f"']' closes the group. I pop the saved prefix \"{prefix}\" and count {rep}, and set cur to "
                   f"the prefix plus \"{inner}\" repeated {rep} times: \"{cur}\".")
            style = "match"
        else:
            cur += c
            run = [j for j in range(i, -1, -1) if all(s[t].isalpha() for t in range(j, i + 1))]
            first = run[-1]  # start of this run of letters
            if first < i:
                for j in range(first, i):
                    hl[j] = "active"
                letters = " and ".join(f"'{s[j]}'" for j in range(first, i + 1))
                line = f'letters → cur = "{cur}"'
                say = f"{letters} are letters, so they go straight onto cur, which is now \"{cur}\"."
            else:
                line = f'letter → cur = "{cur}"'
                say = f"'{c}' is a letter, so it goes straight onto cur, which is now \"{cur}\"."
            style = "active"
        hl[i] = style
        if c.isalpha() and i + 1 < len(s) and s[i + 1].isalpha():
            continue  # a run of letters is shown in one frame, at its last letter
        yield Frame(
            Array(chars, label="s", hl=hl, pointers={"i": i}),
            stack_view(st),
            Text(line, f'cur = "{cur}"   k = {k}'),
            say=say,
        )
    assert cur == want, cur
    yield Frame(
        Array(chars, label="s", hl={j: "visited" for j in range(len(s))}),
        stack_view([]),
        Text(f'answer: "{cur}"', ("time O(output) · space O(output)", "match")),
        say="The stack is empty and cur holds the answer, \"aaabcbc\". The work and memory are "
            "proportional to the output length.",
    )
