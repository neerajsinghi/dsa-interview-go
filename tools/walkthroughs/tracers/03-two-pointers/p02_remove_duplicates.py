from walkthrough import Array, Frame, Text

TITLE = "Remove Duplicates from Sorted Array"
PATTERN = "Read / write pointers"


def frames():
    nums = [1, 1, 2]
    original = list(nums)

    yield Frame(
        Array(nums, label="nums (sorted)"),
        Text("rearrange in place: first k slots = distinct values, return k"),
        speaker="Interviewer",
        say="nums is sorted. Rearrange it in place so the first k elements are the distinct values "
            "in order, and return k. What's after index k-1 doesn't matter.",
    )
    yield Frame(
        Array(nums, label="nums (sorted)"),
        Text("brute force: copy distinct values into a new slice, copy back",
             "O(n) time, O(n) extra space"),
        say="I could collect the distinct values in a second slice and copy them back. "
            "That's O(n) time but O(n) extra memory. I want O(1) space, so I'll compact in place.",
    )

    write = 1
    yield Frame(
        Array(nums, label="nums", hl={0: "match"}, pointers={"write": write}, window=(0, 0)),
        Text("nums[:write] = [1] is already deduplicated"),
        say="nums[0] is always kept, so I start write at 1. Everything left of write is the clean prefix. "
            "A second pointer, read, scans the rest.",
    )
    for read in range(1, len(nums)):
        if nums[read] != nums[write - 1]:
            before = list(nums)
            nums[write] = nums[read]
            hl = {i: "match" for i in range(write + 1)}
            hl[read] = hl.get(read, "active")
            yield Frame(
                Array(before, label="before", hl={read: "active", write: "compare"},
                      pointers={"write": write, "read": read}),
                Array(nums, label="after nums[write] = nums[read]", hl=hl, window=(0, write)),
                Text(f"nums[read] = {before[read]} != nums[write-1] = {before[write - 1]}",
                     (f"new value → copy to slot {write}, write = {write + 1}", "match")),
                say=f"{before[read]} differs from the last kept value {before[write - 1]}. "
                    "Because the array is sorted, that means it's new. Copy it to the write slot and advance write.",
            )
            write += 1
        else:
            hl = {i: "match" for i in range(write)}
            hl[read] = "reject"
            yield Frame(
                Array(nums, label="nums", hl=hl, pointers={"write": write, "read": read}, window=(0, write - 1)),
                Text(f"nums[read] = {nums[read]} == nums[write-1] = {nums[write - 1]}",
                     ("duplicate → skip, write stays", "reject")),
                say=f"{nums[read]} equals the last value I kept. In a sorted array, duplicates sit next to each other, "
                    "so one comparison is enough. Skip it; write doesn't move.",
            )

    k = write
    assert k == 2 and nums[:k] == [1, 2], (k, nums, original)  # doc comment: [1, 1, 2] → k = 2, [1, 2]
    yield Frame(
        Array(nums, label="nums", hl={i: "match" for i in range(k)}, pointers={"write": k}, window=(0, k - 1)),
        Text((f"answer: k = {k}, nums[:{k}] = {nums[:k]}", "match"),
             "time O(n): read visits each index once", "space O(1): only two indices"),
        say=f"read reached the end, so I return k = {k}. The prefix is [1, 2]. One pass is O(n) time, "
            "and two integers are O(1) extra space.",
    )
