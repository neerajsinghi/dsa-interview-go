from walkthrough import Array, Frame, Text

TITLE = "String Compression"
PATTERN = "Read/write pointers"


def frames():
    s = "aabbccc"
    chars = list(s)
    n = len(chars)

    yield Frame(
        Array(chars, label="chars"),
        Text('"aabbccc" → "a2b2c3", return 6'),
        speaker="Interviewer",
        say="Compress the byte slice in place with run-length encoding. A run of one keeps just the letter. "
            "Return the new length, using O(1) extra space.",
    )
    yield Frame(
        Array(chars, label="chars"),
        Text("naive: build a new []byte, then copy back → O(n) extra space"),
        say="Building a fresh output slice is easy and is what CompressNaive does, but it needs O(n) extra memory. "
            "I want to write the answer over the input.",
    )

    write = 0
    read = 0
    while read < n:
        ch, start = chars[read], read
        while read < n and chars[read] == ch:
            read += 1
        run = read - start
        hl = {i: "window" for i in range(start, read)}
        yield Frame(
            Array(chars, label="chars", hl=hl, pointers={"w": write, "start": start, "read": read}),
            Text(f"run of '{ch}' = chars[{start}:{read}], length {run}"),
            say=f"read scans to the end of the run of '{ch}'. The run has length {run}.",
        )
        w0 = write
        chars[write] = ch
        write += 1
        digits = str(run) if run > 1 else ""
        for d in digits:
            chars[write] = d
            write += 1
        hl = {i: "match" for i in range(w0, write)}
        yield Frame(
            Array(chars, label="chars", hl=hl, pointers={"w": write, "read": read}),
            Text(f"write '{ch}'" + (f" + digits \"{digits}\"" if digits else " (run of 1: no count)"),
                 (f"w = {write} ≤ read = {read}", "special")),
            say=f"I write '{ch}' and the digits \"{digits}\" at w. A run of length k needs 1 + digits(k) ≤ k bytes, "
                "so w can never overtake read and destroy unread input.",
        )

    assert write == 6 and "".join(chars[:write]) == "a2b2c3"
    yield Frame(
        Array(chars, label="chars", hl={i: "match" for i in range(write)}, window=(0, write - 1)),
        Text((f'answer: {write}, chars[:{write}] = "{"".join(chars[:write])}"', "match"),
             "time O(n) · space O(1)"),
        say="The compressed form is chars[:6] = \"a2b2c3\" and I return 6. Each byte is read once and written "
            "at most once: O(n) time, O(1) extra space.",
    )
