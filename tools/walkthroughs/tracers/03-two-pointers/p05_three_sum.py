from walkthrough import Array, Frame, KV, Text

TITLE = "3Sum"
PATTERN = "Sort + two pointers"


def frames():
    nums = [-1, 0, 1, 2, -1, -4]

    yield Frame(
        Array(nums, label="nums"),
        Text("find every distinct triplet with a + b + c = 0", "no duplicate triplets in the output"),
        speaker="Interviewer",
        say="Return every distinct triplet of values from different positions that sums to zero. "
            "Each triplet sorted, no triplet twice.",
    )
    yield Frame(
        Array(nums, label="nums"),
        Text("brute force: three nested loops + a set to drop repeats", "O(n³) time"),
        say="Three nested loops try every triple, and a set removes repeats. That's O(n³). "
            "If I fix one value, the other two are just Two Sum, which two pointers solve in O(n).",
    )

    s = sorted(nums)
    yield Frame(
        Array(nums, label="nums"),
        Array(s, label="s = sorted copy"),
        Text("sorting: O(n log n), and equal values end up next to each other"),
        say="First I sort a copy. Sorting lets two pointers discard values, "
            "and it puts equal values side by side, so skipping duplicates is easy.",
    )

    res = []

    def view(i, l, r, hl_extra=None):
        hl = {i: "special"}
        if hl_extra:
            hl.update(hl_extra)
        ptr = {"i": i}
        if l < len(s):
            ptr["l"] = l
        if r >= 0:
            ptr["r"] = r
        return Array(s, label="s", hl=hl, pointers=ptr)

    def found():
        return KV([str(t) for t in res], label="result", empty="[ ]")

    for i in range(len(s) - 2):
        if s[i] > 0:
            break
        if i > 0 and s[i] == s[i - 1]:
            yield Frame(
                view(i, len(s), -1, {i - 1: "visited", i: "reject"}),
                found(),
                Text((f"s[{i}] = {s[i]} == s[{i - 1}] → skip", "reject")),
                say=f"s[{i}] is {s[i]} again. Every triplet starting with {s[i]} is already found, "
                    "so I skip it to avoid duplicates.",
            )
            continue
        l, r = i + 1, len(s) - 1
        steps = []
        while l < r:
            total = s[i] + s[l] + s[r]
            if total < 0:
                steps.append(("small", l, r, total))
                l += 1
            elif total > 0:
                steps.append(("big", l, r, total))
                r -= 1
            else:
                res.append([s[i], s[l], s[r]])
                steps.append(("hit", l, r, total))
                l += 1
                r -= 1
                while l < r and s[l] == s[l - 1]:
                    l += 1

        # For i = -4 every step is "too small"; show one, then summarise the run.
        if steps and all(k == "small" for k, *_ in steps) and len(steps) > 2:
            k, l0, r0, t0 = steps[0]
            yield Frame(
                view(i, l0, r0, {l0: "reject", r0: "compare"}),
                found(),
                Text(f"fix s[i] = {s[i]}: {s[i]} + {s[l0]} + {s[r0]} = {t0} < 0", ("too small → l++", "reject")),
                say=f"Fix {s[i]}. Now I need two values summing to {-s[i]}. {s[i]} + {s[l0]} + {s[r0]} = {t0}, "
                    "too small, so l moves right.",
            )
            sums = ", ".join(str(t) for *_, t in steps[1:])
            yield Frame(
                view(i, l, r, {j: "visited" for j in range(i + 1, len(s))}),
                found(),
                Text(f"next sums: {sums} — all < 0", ("l meets r → no triplet starts with " + str(s[i]), "reject")),
                say=f"Even the largest pair, {s[-2]} + {s[-1]} = {s[-2] + s[-1]}, can't reach {-s[i]}. Every sum stays negative, "
                    f"l meets r, and no triplet starts with {s[i]}.",
            )
            continue

        for k, l0, r0, t0 in steps:
            if k == "hit":
                yield Frame(
                    view(i, l0, r0, {l0: "match", r0: "match", i: "match"}),
                    found(),
                    Text(f"{s[i]} + {s[l0]} + {s[r0]} = 0 ✓", ("record it, then l++ and r--", "match")),
                    say=f"{s[i]} + {s[l0]} + {s[r0]} is zero, so I record it. Then I move both pointers "
                        "and skip any repeat of the same second value.",
                )
            elif k == "big":
                yield Frame(
                    view(i, l0, r0, {l0: "compare", r0: "reject"}),
                    found(),
                    Text(f"{s[i]} + {s[l0]} + {s[r0]} = {t0} > 0", ("too big → r--", "reject")),
                    say=f"Fix {s[i]}. {s[i]} + {s[l0]} + {s[r0]} = {t0}, too big, so r moves left and drops {s[r0]}."
                        + (" Now l meets r, so this i is done." if l0 >= r0 - 1 else ""),
                )
            else:
                yield Frame(
                    view(i, l0, r0, {l0: "reject", r0: "compare"}),
                    found(),
                    Text(f"{s[i]} + {s[l0]} + {s[r0]} = {t0} < 0", ("too small → l++", "reject")),
                    say=f"{t0} is too small, so l moves right.",
                )

    assert res == [[-1, -1, 2], [-1, 0, 1]]  # doc comment: → [[-1, -1, 2], [-1, 0, 1]]
    yield Frame(
        Array(s, label="s", hl={j: "visited" for j in range(len(s))}),
        found(),
        Text((f"answer: {res}", "match"), "time O(n²): n choices of i × O(n) two-pointer scan",
             "space O(1) beyond the output (the sorted copy is O(n))"),
        say="The answer is [[-1, -1, 2], [-1, 0, 1]]. For each of n first values I run an O(n) scan, so it's O(n²). "
            "Sorting a copy costs O(n log n) and O(n) space.",
    )
