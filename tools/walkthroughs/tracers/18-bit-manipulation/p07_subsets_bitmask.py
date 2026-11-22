from walkthrough import Array, Frame, Grid, KV, Text

TITLE = "Subsets via Bitmask"
PATTERN = "Bits · bitmask"


def fmt(sub):
    return "[" + " ".join(map(str, sub)) + "]"


def frames():
    nums = [1, 2, 3]
    expected = [[], [1], [2], [1, 2], [3], [1, 3], [2, 3], [1, 2, 3]]
    n = len(nums)
    cols = list(range(n - 1, -1, -1))  # bit n-1 on the left, bit 0 on the right
    col_labels = [f"nums[{i}]={nums[i]}" for i in cols]

    def mask_grid(mask, hl=True):
        row = [(mask >> i) & 1 for i in cols]
        return Grid([row], row_labels=[f"mask = {mask}"], col_labels=col_labels, cell=86,
                    hl={(0, j): "match" for j, b in enumerate(row) if b} if hl else {})

    yield Frame(
        Array(nums, label="nums"),
        speaker="Interviewer",
        say="Given distinct integers, return every subset. Elements inside a subset keep their original order.",
    )
    yield Frame(
        Array(nums, label="nums"),
        Text("backtracking: include / skip each element → also O(n · 2ⁿ)"),
        say="Backtracking works: include or skip each element. But each subset is just n yes/no choices, "
            "and an n-bit integer stores exactly that. So I can count instead of recurse.",
    )

    out = []
    for mask in range(1 << n):
        subset = [nums[i] for i in range(n) if (mask >> i) & 1 == 1]
        out.append(subset)
        if mask == 0:
            say = "mask 0 has no bits set, so it's the empty subset."
        elif mask == 3:
            say = "mask 3 is 011: bits 0 and 1 are set, so I take nums[0] and nums[1]. That gives [1 2]."
        elif mask == 5:
            say = "mask 5 is 101: bits 0 and 2, so [1 3]. Scanning i from 0 up keeps the original order."
        elif mask < 3:
            say = f"mask {mask} has only bit {mask - 1} set, so the subset is [{nums[mask - 1]}]."
        elif mask == 4:
            say = "mask 4 sets only bit 2, so the subset is [3]."
        elif mask == 6:
            say = "Same again: mask 6 is 110, so [2 3]."
        else:
            say = f"mask 7 = 2ⁿ - 1 has every bit set: the full set. That's the last mask."
        keys = [fmt(s) for s in out]
        yield Frame(
            Array(nums, label="nums", hl={i: "match" for i in range(n) if (mask >> i) & 1}),
            mask_grid(mask),
            KV(keys, label=f"out ({len(out)} of {1 << n})", hl={keys[-1]: "active"}),
            say=say,
        )

    assert out == expected
    keys = [fmt(s) for s in out]
    yield Frame(
        Array(nums, label="nums"),
        KV(keys, label=f"out ({len(out)} subsets)", hl={k: "match" for k in keys}),
        Text((f"answer: {len(out)} subsets", "match"), "O(n · 2ⁿ) time and space (that's the output size)"),
        say="Masks 0 to 2ⁿ-1 give all 8 subsets. Each takes n bit checks, so O(n · 2ⁿ) time and space, "
            "which is the size of the output. n must stay small.",
    )
