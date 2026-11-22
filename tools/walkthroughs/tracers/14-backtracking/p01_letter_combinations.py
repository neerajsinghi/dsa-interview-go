from walkthrough import Array, Frame, KV, Text

TITLE = "Letter Combinations of a Phone Number"
PATTERN = "Backtracking"

KEYPAD = ["", "", "abc", "def", "ghi", "jkl", "mno", "pqrs", "tuv", "wxyz"]


def trace(digits):
    """Run LetterCombinations and yield (event, buf, res, i) at every choose and leaf."""
    res = []
    buf = []

    def backtrack(i):
        if i == len(digits):
            res.append("".join(buf))
            yield "leaf", list(buf), list(res), i
            return
        for ch in KEYPAD[int(digits[i])]:
            buf.append(ch)  # choose
            yield "choose", list(buf), list(res), i
            yield from backtrack(i + 1)  # explore
            buf.pop()  # unchoose

    if digits:
        yield from backtrack(0)
    trace.result = res


def frames():
    digits = "23"
    keys = {d: KEYPAD[int(d)] for d in digits}

    def path_arr(buf, hl=None):
        cells = buf + [""] * (len(digits) - len(buf))
        return Array(cells, label=f"buf: one letter per digit of \"{digits}\"", hl=hl or {}, cell=56)

    def results(res, new=None):
        return KV(res, label=f"results ({len(res)})", hl={new: "match"} if new else None, empty="[ ]")

    yield Frame(
        KV(keys, label="keypad"),
        Text(f'digits = "{digits}"'),
        speaker="Interviewer",
        say='Each digit maps to a few letters on a phone keypad. Return every string you could type for "23", '
            "picking one letter per digit.",
    )
    yield Frame(
        KV(keys, label="keypad"),
        Text("brute force: one nested loop per digit", ("but the number of digits isn't fixed", "reject"),
             "every answer must be output anyway: 3 × 3 = 9 strings"),
        say="With two digits I could write two nested loops, but the input length varies. "
            "Recursion gives me one loop per digit. There's no faster way: I must output every combination.",
    )

    events = list(trace(digits))
    res_final = trace.result
    shown_leaves = 0
    for ev, buf, res, i in events:
        if ev == "choose":
            if buf[0] != "a":
                continue
            hl = {len(buf) - 1: "active"}
            if len(buf) == 1:
                say = (f"Depth 0 picks the letter for '{digits[0]}'. I choose '{buf[0]}', then recurse to choose "
                       f"a letter for '{digits[1]}'.")
                yield Frame(path_arr(buf, hl), KV(keys, label="keypad", hl={digits[0]: "active"}), results(res),
                            Text(f"choose '{buf[0]}' from \"{keys[digits[0]]}\""), say=say)
            continue
        # leaf
        shown_leaves += 1
        if buf[0] == "a":
            word = "".join(buf)
            say = (f"Depth {len(digits)} means every digit has a letter, so '{word}' is an answer. "
                   "Then I unchoose the last letter and try the next one.")
            if shown_leaves > 1:
                say = f"Unchoose, choose '{buf[-1]}': '{word}' is another answer."
            yield Frame(path_arr(buf, {len(buf) - 1: "match"}), KV(keys, label="keypad", hl={digits[1]: "active"}),
                        results(res, word), Text(f"i == len(digits) → append \"{word}\""), say=say)
        elif buf[-1] == keys[digits[-1]][-1]:
            first = buf[0]
            new = [w for w in res if w[0] == first]
            yield Frame(path_arr(buf, {0: "active"}), KV(keys, label="keypad", hl={digits[0]: "active"}),
                        KV(res, label=f"results ({len(res)})", hl={w: "match" for w in new}),
                        Text(f"unchoose back to depth 0, choose '{first}'", f"same again: {' '.join(new)}"),
                        say=f"Back at depth 0, I swap '{keys[digits[0]][keys[digits[0]].index(first) - 1]}' for '{first}' "
                            "and repeat the same three choices for the second digit.")

    assert res_final == ["ad", "ae", "af", "bd", "be", "bf", "cd", "ce", "cf"]
    yield Frame(
        KV(res_final, label=f"results ({len(res_final)})", hl={w: "match" for w in res_final}),
        Text("answer = " + " ".join(res_final), "O(n · 4^n) time, O(n) recursion besides the output"),
        say="Nine strings, in keypad order. Each digit has up to 4 letters, and each answer takes n bytes to build: "
            "O(n times 4 to the n) time, O(n) extra space.",
    )
