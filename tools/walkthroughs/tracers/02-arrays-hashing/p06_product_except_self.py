from walkthrough import Array, Frame, Text

TITLE = "Product of Array Except Self"
PATTERN = "Prefix / suffix"


def frames():
    nums = [1, 2, 3, 4]
    n = len(nums)

    yield Frame(
        Array(nums, label="nums"),
        speaker="Interviewer",
        say="Return out where out[i] is the product of every element except nums[i]. No division.",
    )
    yield Frame(
        Array(nums, label="nums", hl={0: "visited", 1: "compare", 2: "compare", 3: "compare"}),
        Text("brute force: for each i, multiply the other n-1 values → O(n²)",
             ("division would fail on zeros, and it's not allowed", "reject")),
        say="The brute force multiplies the other n − 1 values for each i, which is O(n²). "
            "Instead: out[i] is the product left of i times the product right of i.",
    )

    out = [None] * n
    left = 1
    for i in range(n):
        out[i] = left
        before = left
        left *= nums[i]
        if i == 0:
            say = "Pass 1 goes left to right. Nothing is left of index 0, so out[0] gets the empty product, 1."
        elif i == n - 1:
            say = f"out[{i}] = {before}. After this pass, out holds every prefix product."
        else:
            say = f"out[{i}] = {before}, the product of everything before i. Then left absorbs nums[{i}]."
        yield Frame(
            Array(nums, label="nums", hl={k: "window" for k in range(i)} | {i: "active"}, pointers={"i": i}),
            Array(out, label="out", hl={i: "active"}),
            Text(f"out[{i}] = left = {before}", f"left *= nums[{i}] → {left}"),
            say=say,
        )

    right = 1
    for i in range(n - 1, -1, -1):
        before_out, before = out[i], right
        out[i] *= right
        right *= nums[i]
        if i == n - 1:
            say = (f"Pass 2 goes right to left with a running product, right. Nothing is right of the end, "
                   f"so out[{i}] stays {out[i]}.")
        elif i == 0:
            say = f"Finally out[0] = 1 × {before} = {out[0]}. Every slot now has left × right."
        else:
            say = f"out[{i}] = {before_out} × {before} = {out[i]}: what was left of i, times what is right of i."
        yield Frame(
            Array(nums, label="nums", hl={k: "window" for k in range(i + 1, n)} | {i: "active"}, pointers={"i": i}),
            Array(out, label="out", hl={k: "match" for k in range(i + 1, n)} | {i: "active"}),
            Text(f"out[{i}] = {before_out} × right({before}) = {out[i]}", f"right *= nums[{i}] → {right}"),
            say=say,
        )

    assert out == [24, 12, 8, 6]
    yield Frame(
        Array(out, label="out", hl={k: "match" for k in range(n)}),
        Text((f"answer: {out}", "match"),
             "ProductExceptSelfBrute: O(n²) time, O(1) extra",
             "ProductExceptSelf:      O(n) time,  O(1) extra (output not counted)"),
        say=f"The answer is {out}. Two passes is O(n) time. Reusing out for the prefixes "
            "leaves only two ints of extra space.",
    )
