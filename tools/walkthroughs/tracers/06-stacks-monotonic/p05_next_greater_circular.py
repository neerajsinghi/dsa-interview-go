from walkthrough import Array, Frame, Text, stack_widget

TITLE = "Next Greater Element (circular)"
PATTERN = "Monotonic stack"


def frames():
    nums = [1, 2, 1]
    want = [2, -1, 2]
    n = len(nums)
    doubled = nums + nums  # display only: the 2n indices the loop visits

    yield Frame(
        Array(nums, label="nums (circular)"),
        speaker="Interviewer",
        say="The array is circular. For each element, find the first strictly greater element walking "
            "forward, wrapping around at most once. Use -1 if there's none.",
    )
    yield Frame(
        Array(nums, label="nums (circular)", hl={2: "active", 0: "compare", 1: "match"}),
        Text("brute force: from each i, step up to n − 1 times with (i + step) % n → O(n²)"),
        say="From each element I could walk forward up to n − 1 steps with modulo. That's O(n²). "
            "The last 1 needs to wrap around to find the 2.",
    )
    yield Frame(
        Array(doubled, label="i = 0 .. 2n−1, value nums[i % n]", hl={j: "visited" for j in range(n, 2 * n)}),
        Text("first pass: push and resolve    second pass: resolve only"),
        say="Trick: loop i from 0 to 2n − 1 and read nums[i % n]. The second lap lets late elements see "
            "the ones before them. Only the first lap pushes.",
    )

    ans = [-1] * n
    st = []
    for i in range(2 * n):
        v = nums[i % n]
        popped = []
        while st and nums[st[-1]] < v:
            j = st.pop()
            ans[j] = v
            popped.append(j)
        pushed = i < n
        if pushed:
            st.append(i)
        hl = {i: "active"}
        hl.update({j: "match" for j in popped})
        hl.update({j: "visited" for j in range(n, 2 * n) if j != i})
        lap = "1st lap" if i < n else "2nd lap"
        if popped:
            act = f"{v} > " + ", ".join(f"nums[{j}]" for j in popped) + " → " + \
                  ", ".join(f"ans[{j}] = {v}" for j in popped)
        else:
            act = "nothing smaller on top"
        act += "; push " + str(i) if pushed else "; no push (2nd lap)"
        if popped and pushed:
            say = f"{v} is greater than the {nums[popped[0]]} waiting on top, so index {popped[0]}'s answer is {v}. Then I push index {i}."
        elif popped:
            say = (f"Second lap: index {i} is nums[{i % n}] = {v} again. It's greater than the {nums[popped[0]]} "
                   f"at index {popped[0]}, so ans[{popped[0]}] = {v}. That's the wrap-around.")
        elif pushed:
            say = f"Index {i} holds {v}. Nothing smaller is waiting, so I push it. It waits for something bigger."
        else:
            say = (f"Second lap, value {v}. It isn't greater than the top of the stack, and the second lap "
                   "never pushes, so nothing happens.")
        yield Frame(
            Array(doubled, label=f"i = {i} ({lap}),  value nums[{i % n}] = {v}", hl=hl, pointers={"i": i}),
            stack_widget([f"{j}:{nums[j]}" for j in st], label="stack (index:value)"),
            Array(ans, label="ans", hl={j: "match" for j in popped}),
            Text(act, size=14),
            say=say,
        )

    assert ans == want, ans
    yield Frame(
        Array(nums, label="nums"),
        stack_widget([f"{j}:{nums[j]}" for j in st], label="left on stack (no greater element)"),
        Array(ans, label="ans", hl={j: "match" for j in range(n)}),
        Text(f"answer: {ans}", ("time O(n) · space O(n)", "match")),
        say="Index 1 holds the maximum, so it stays on the stack with -1. The loop runs 2n times and each "
            "index is pushed and popped at most once: O(n) time, O(n) space.",
    )
