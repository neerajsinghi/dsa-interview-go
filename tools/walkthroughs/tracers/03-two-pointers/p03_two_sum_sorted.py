from walkthrough import Array, Frame, Text

TITLE = "Two Sum II (sorted)"
PATTERN = "Two pointers"


def frames():
    nums, target = [2, 7, 11, 15], 9

    yield Frame(
        Array(nums, label="nums (sorted)"),
        Text(f"target = {target}"),
        speaker="Interviewer",
        say="The array is sorted. Find two numbers that add up to 9 and return their indices, using O(1) extra space.",
    )
    yield Frame(
        Array(nums, label="nums (sorted)"),
        Text(f"target = {target}", "brute force: try every pair → O(n²)"),
        say="Trying every pair works, but it's O(n²) and ignores that the array is sorted. "
            "Sorted data lets me throw values away.",
    )

    l, r = 0, len(nums) - 1
    yield Frame(
        Array(nums, label="nums (sorted)", hl={l: "active", r: "active"}, pointers={"l": l, "r": r}),
        Text(f"target = {target}"),
        say="I'll put one pointer at each end: l on the smallest value and r on the largest.",
    )
    while l < r:
        s = nums[l] + nums[r]
        if s == target:
            yield Frame(
                Array(nums, label="nums (sorted)", hl={l: "match", r: "match"}, pointers={"l": l, "r": r}),
                Text(f"nums[l] + nums[r] = {nums[l]} + {nums[r]} = {s} = target ✓", ("answer: [%d, %d]" % (l, r), "match")),
                say=f"{nums[l]} + {nums[r]} is exactly {target}, so the answer is [{l}, {r}]. "
                    "Each pointer moved at most n times, so it's O(n) time and O(1) space.",
            )
            return
        if s > target:
            yield Frame(
                Array(nums, label="nums (sorted)", hl={l: "compare", r: "reject"}, pointers={"l": l, "r": r}),
                Text(f"nums[l] + nums[r] = {nums[l]} + {nums[r]} = {s} > {target}", ("too big → move r left", "reject")),
                say=f"{s} is too big. Even with the smallest value, {nums[r]} overshoots, "
                    f"so {nums[r]} can't be in any answer. Move r left.",
            )
            r -= 1
        else:
            yield Frame(
                Array(nums, label="nums (sorted)", hl={l: "reject", r: "compare"}, pointers={"l": l, "r": r}),
                Text(f"nums[l] + nums[r] = {nums[l]} + {nums[r]} = {s} < {target}", ("too small → move l right", "reject")),
                say=f"{s} is too small. Even with the largest value, {nums[l]} falls short, "
                    f"so {nums[l]} can't be in any answer. Move l right.",
            )
            l += 1
