from walkthrough import Array, Frame, Text

TITLE = "Partition Equal Subset Sum"
PATTERN = "1-D DP · 0/1 knapsack"


def frames():
    nums = [1, 5, 11, 5]
    total = sum(nums)

    def widgets(hl_n=None, can=None, hl_c=None):
        return (
            Array(nums, label="nums", hl=hl_n),
            Array(can if can is not None else [None] * 12, label="can[s] = some subset of the items so far sums to s",
                  hl=hl_c),
        )

    yield Frame(
        *widgets(),
        speaker="Interviewer",
        say="Can [1, 5, 11, 5] be split into two groups with equal sums?",
    )
    yield Frame(
        *widgets(),
        Text("try every subset: 2^n", ("too slow", "reject")),
        say="Trying every subset is 2^n. But two equal halves just means one subset sums to half the total. "
            "That's subset sum, a 0/1 knapsack.",
    )

    assert total % 2 == 0
    target = total // 2
    can = [False] * (target + 1)
    can[0] = True
    yield Frame(
        *widgets(None, can, {0: "match"}),
        Text(f"total = {total} (even) → target = {target}", "can[0] = T  (the empty subset)"),
        say=f"The total is {total}, which is even, so the target is {target}. "
            "The empty subset makes 0. Every other sum starts false.",
    )

    for idx, x in enumerate(nums):
        before = can[:]
        new = []
        for s in range(target, x - 1, -1):
            if not can[s] and can[s - x]:
                new.append(s)
            can[s] = can[s] or can[s - x]
        hl = {s: "match" for s in new}
        for s in new:
            hl.setdefault(s - x, "compare")
        sums = [s for s in range(target + 1) if can[s]]
        lines = [f"item {x}: s from {target} down to {x}: can[s] |= can[s-{x}]"]
        lines.append("new sums: " + (", ".join(f"{s} (= {s - x} + {x})" for s in new) if new else "none"))
        lines.append(f"reachable: {sums}")
        if idx == 0:
            say = ("For each item I sweep s downward. Going down means can[s-x] still holds the value from before "
                   "this item, so each item is used at most once.")
        elif idx == 1:
            say = "With 5 I can reach 5 and 6. Going upward would wrongly reuse 5 to also make 10, like Coin Change."
        else:
            say = f"With {x}, can[{target}] becomes true. The Go code exits early here."
        yield Frame(*widgets({idx: "active"}, can, hl), Text(*lines), say=say)
        if can[target]:
            break

    assert can[target] is True
    yield Frame(
        *widgets({2: "match"}, can, {target: "match"}),
        Text(f"answer: can[{target}] = T  ({{11}} and {{1, 5, 5}})", ("O(n · target) time · O(target) space", "match")),
        say="So the answer is true: {11} and {1, 5, 5}. Each item sweeps the array once: "
            "O(n × target) time and O(target) space.",
    )
