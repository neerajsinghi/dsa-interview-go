from walkthrough import Array, Frame, Text

TITLE = "Reverse Words in a String"
PATTERN = "Reverse twice"

SP = "␣"


def show(buf):
    return [SP if c == " " else c for c in buf]


def reverse_in_place(buf, lo, hi):  # reverses buf[lo:hi]
    i, j = lo, hi - 1
    while i < j:
        buf[i], buf[j] = buf[j], buf[i]
        i, j = i + 1, j - 1


def reverse_words_in_place(buf):
    w = 0
    r = 0
    while r < len(buf):
        if buf[r] == " ":
            r += 1
            continue
        if w > 0:
            buf[w] = " "
            w += 1
        while r < len(buf) and buf[r] != " ":
            buf[w] = buf[r]
            w, r = w + 1, r + 1
    buf = buf[:w]
    reverse_in_place(buf, 0, len(buf))
    start = 0
    while start < len(buf):
        end = start
        while end < len(buf) and buf[end] != " ":
            end += 1
        reverse_in_place(buf, start, end)
        start = end + 1
    return buf


def frames():
    # The second doc example has extra spaces, so all three phases do work.
    assert "".join(reverse_words_in_place(list("the sky is blue"))) == "blue is sky the"
    s = "  hello   world  "
    buf = list(s)

    yield Frame(
        Array(show(buf), label="buf"),
        Text(f'"{s}"'),
        speaker="Interviewer",
        say="Reverse the order of the words. Collapse runs of spaces to one, and drop leading and trailing spaces.",
    )
    yield Frame(
        Array(show(buf), label="buf"),
        Text("strings.Fields + reverse the slice + Join → O(n) time, O(n) extra"),
        say="In Go I'd call strings.Fields, reverse the slice and Join. That's fine, but the follow-up is: "
            "do it in place on a byte buffer with O(1) extra space.",
    )

    # Phase 1: squeeze, one frame per copied word.
    w = r = 0
    while r < len(buf):
        if buf[r] == " ":
            r += 1
            continue
        if w > 0:
            buf[w] = " "
            w += 1
        start_r, start_w = r, w
        while r < len(buf) and buf[r] != " ":
            buf[w] = buf[r]
            w, r = w + 1, r + 1
        word = s[start_r:r]
        hl = {i: "match" for i in range(start_w, w)}
        yield Frame(
            Array(show(buf), label="buf", hl=hl, pointers={"w": w, "r": r}),
            Text(f'copied "{word}" from index {start_r} to index {start_w}', ("w never passes r", "special")),
            say=f"Step 1 squeezes spaces. A read pointer r skips spaces; a write pointer w copies \"{word}\" left"
                + (", with one space before it." if start_w > 0 else ".")
                + " w never passes r, so I never overwrite unread text.",
        )
    buf = buf[:w]
    yield Frame(
        Array(show(buf), label="buf[:w]"),
        Text(f'cut to buf[:{w}] = "{"".join(buf)}"'),
        say=f"After the copy, the clean text is the first {w} bytes. I keep buf[:w] and forget the rest.",
    )

    # Phase 2: reverse everything.
    reverse_in_place(buf, 0, len(buf))
    yield Frame(
        Array(show(buf), label="buf", hl={i: "window" for i in range(len(buf))}),
        Text(f'reverse all → "{"".join(buf)}"'),
        say="Step 2 reverses the whole buffer. Now the words are in the right order, but each word is backwards.",
    )

    # Phase 3: reverse each word back.
    start = 0
    done = set()
    while start < len(buf):
        end = start
        while end < len(buf) and buf[end] != " ":
            end += 1
        before = "".join(buf[start:end])
        reverse_in_place(buf, start, end)
        done |= set(range(start, end))
        yield Frame(
            Array(show(buf), label="buf", hl={i: "match" for i in done}, pointers={"start": start, "end": end}),
            Text(f'reverse buf[{start}:{end}]: "{before}" → "{"".join(buf[start:end])}"'),
            say=f"Step 3 walks the words. For each one I find its end at the next space and reverse it in place, "
                f"so \"{before}\" reads \"{''.join(buf[start:end])}\" again.",
        )
        start = end + 1

    out = "".join(buf)
    assert out == "world hello"
    yield Frame(
        Array(show(buf), label="result", hl={i: "match" for i in range(len(buf))}),
        Text((f'answer: "{out}"', "match"), "time O(n) · extra space O(1)"),
        say="The result is \"world hello\". Each phase touches every byte a constant number of times, "
            "so it's O(n) time with O(1) extra space.",
    )
