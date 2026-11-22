from walkthrough import Array, Frame, KV, Text

TITLE = "Partition Labels"
PATTERN = "Greedy · intervals"

PIECE_STYLES = ["match", "special", "active"]


def frames():
    s = "ababcbacadefegdehijhklij"
    expected = [9, 7, 8]
    chars = list(s)

    yield Frame(
        Array(chars, label="s"),
        speaker="Interviewer",
        say="Cut the string into as many pieces as possible so each letter appears in only one piece. "
            "Return the piece lengths.",
    )
    yield Frame(
        Array(chars, label="s"),
        Text("brute force: try every set of cut points → O(2ⁿ · n)"),
        say="I could try every set of cut points and keep the valid one with the most pieces. "
            "With 23 gaps that's 2²³ options. Too slow.",
    )

    last = {}
    for i, c in enumerate(s):
        last[c] = i
    yield Frame(
        Array(chars, label="s"),
        KV(dict(sorted(last.items())), label="last index of each letter"),
        say="Each letter spans from its first to its last index, so the pieces are merged intervals. "
            "First pass: record each letter's last index.",
    )

    sizes = []
    start, end = 0, 0
    done_hl = {}
    for i, c in enumerate(s):
        old = end
        end = max(end, last[c])
        cut = i == end
        hl = dict(done_hl)
        for k in range(start, i + 1):
            hl[k] = "window"
        hl[i] = "compare"
        if last[c] not in done_hl:
            hl[last[c]] = "compare"
        shown = i == 0 or i == 1 or end != old or cut
        if cut:
            sizes.append(end - start + 1)
            style = PIECE_STYLES[(len(sizes) - 1) % 3]
            for k in range(start, end + 1):
                hl[k] = style
                done_hl[k] = style
            yield Frame(
                Array(chars, label="s", hl=hl, pointers={"i = end": i},
                      window=(start, end)),
                Text(f"i = {i} == end → cut", (f"piece length = {end} - {start} + 1 = {sizes[-1]}", "match"),
                     f"sizes = {sizes}"),
                say=(f"i reached end. Every letter in this piece has its last copy inside it, so I cut. "
                     f"Length {sizes[-1]}. The next piece starts at {i + 1}."),
            )
            start = i + 1
            continue
        if not shown:
            continue
        if end != old:
            say = (f"'{c}' last appears at {last[c]}, so this piece must reach at least {last[c]}. "
                   f"end becomes {end}.")
            if i >= 16:
                say = f"'{c}' last appears at {last[c]}. end grows to {end}."
        else:
            say = (f"'{c}' last appears at {last[c]}, which is inside the current end {end}. "
                   "Nothing changes. Letters 2 to 7 are the same, so I skip ahead.")
        yield Frame(
            Array(chars, label="s", hl=hl, pointers={"i": i, "end": end}, window=(start, end)),
            Text(f"last['{c}'] = {last[c]}", f"end = max({old}, {last[c]}) = {end}"),
            say=say,
        )

    assert sizes == expected
    yield Frame(
        Array(chars, label="s", hl=done_hl),
        Text((f"answer: {sizes}", "match"), "O(n) time, O(1) space (26 counters)"),
        say="The pieces are ababcbaca, defegde and hijhklij: lengths 9, 7, 8. "
            "Two passes over the string, O(n) time, and only 26 counters: O(1) space.",
    )
