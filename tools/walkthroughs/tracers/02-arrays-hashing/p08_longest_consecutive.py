from walkthrough import KV, Array, Frame, Text

TITLE = "Longest Consecutive Sequence"
PATTERN = "Hash set"


def frames():
    nums = [100, 4, 200, 1, 3, 2]

    yield Frame(
        Array(nums, label="nums"),
        speaker="Interviewer",
        say="Return the length of the longest run of consecutive values in this unsorted slice, in O(n).",
    )
    yield Frame(
        Array(sorted(nums), label="sorted copy"),
        Text("sort, then scan for neighbours that differ by 1", ("O(n log n): too slow for the O(n) ask", "reject")),
        say="Sorting puts runs next to each other, but that's O(n log n). "
            "Checking 'is x + 1 there?' for every x needs a set.",
    )

    # NewSet(nums); Go iterates the set in random order, here insertion order.
    s = list(dict.fromkeys(nums))
    yield Frame(
        Array(nums, label="nums"),
        KV(s, label="set"),
        Text("only count from a run START: x with x-1 absent"),
        say="I put everything in a set. Then I only start counting at a run's beginning: "
            "a value whose x − 1 is missing. Other values get counted from their start.",
    )

    best, best_run = 0, []
    for x in s:
        if x - 1 in s:
            yield Frame(
                KV(s, label="set", hl={x: "visited", x - 1: "compare"}),
                Text(f"x = {x}: {x - 1} is in the set → skip", f"best = {best}"),
                say=f"{x - 1} is present, so {x} is inside a run that starts lower. Skip it; "
                    "that run gets counted once, from its start.",
            )
            continue
        length = 1
        steps = [x]
        while x + length in s:
            length += 1
            steps.append(x + length - 1)
        if length > best:
            best, best_run = length, steps
        hl = {v: "match" for v in steps}
        hl[x] = "active"
        if length == 1:
            say = f"{x - 1} is missing, so {x} starts a run. {x + 1} isn't there, so the run has length 1."
        else:
            say = (f"{x - 1} is missing, so {x} is a start. I walk up: " + ", ".join(map(str, steps[1:]))
                   + f" are all there, {x + length} isn't. Length {length}.")
        yield Frame(
            KV(s, label="set", hl=hl),
            Text(f"x = {x}: {x - 1} absent → start; count " + " → ".join(map(str, steps)),
                 (f"length = {length}, best = {best}", "match")),
            say=say,
        )

    assert best == 4
    yield Frame(
        KV(s, label="set", hl={v: "match" for v in best_run}),
        Text((f"answer: {best}", "match"),
             "each value is touched at most twice: once by the outer loop,",
             "once while extending its run → O(n) time, O(n) space"),
        say=f"The answer is {best}. The inner loop only runs from starts, so each value is visited at most twice. "
            "That's O(n) time with O(n) space for the set.",
    )
