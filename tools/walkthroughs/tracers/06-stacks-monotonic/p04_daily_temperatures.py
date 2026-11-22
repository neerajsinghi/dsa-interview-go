from walkthrough import Array, Frame, Text, stack_widget

TITLE = "Daily Temperatures"
PATTERN = "Monotonic stack"


def frames():
    temps = [73, 74, 75, 71, 69, 72, 76, 73]
    want = [1, 1, 4, 2, 1, 1, 0, 0]
    n = len(temps)

    yield Frame(
        Array(temps, label="temps"),
        Array([None] * n, label="ans"),
        speaker="Interviewer",
        say="For each day, how many days until a strictly warmer day? Store 0 if none comes.",
    )
    yield Frame(
        Array(temps, label="temps", hl={2: "active", 3: "compare", 4: "compare", 5: "compare", 6: "match"}),
        Array([None] * n, label="ans"),
        Text("brute force: scan forward from every day → O(n²)"),
        say="Scanning forward from each day is O(n²) on a falling sequence. Instead I'll keep the days "
            "still waiting for a warmer day on a stack.",
    )

    ans = [0] * n
    shown = [None] * n
    st = []
    for i, t in enumerate(temps):
        popped = []
        while st and temps[st[-1]] < t:
            j = st.pop()
            ans[j] = i - j
            shown[j] = ans[j]
            popped.append(j)
        st.append(i)
        stack_vals = [f"{j}:{temps[j]}" for j in st]
        hl_t = {i: "active"}
        hl_t.update({j: "match" for j in popped})
        if popped:
            res = ", ".join(f"ans[{j}] = {i} − {j} = {ans[j]}" for j in popped)
            if len(popped) == 1:
                say = (f"{t} is warmer than day {popped[0]}'s {temps[popped[0]]}, so day {i} answers it: "
                       f"wait {ans[popped[0]]}. Then I push day {i}.")
            else:
                say = (f"{t} is warmer than {len(popped)} waiting days. I pop each and record its wait, "
                       f"stopping at a day that isn't colder. Then I push day {i}.")
        else:
            res = "nothing colder on top → just push"
            say = (f"Day {i} ({t}) can't answer anyone yet, so it waits on the stack. "
                   "Temperatures on the stack only go down from bottom to top.")
        yield Frame(
            Array(temps, label="temps", hl=hl_t, pointers={"i": i}),
            stack_widget(stack_vals, label="stack (index:temp), decreasing"),
            Array(shown, label="ans", hl={j: "match" for j in popped}),
            Text(res, size=14),
            say=say,
        )

    assert ans == want, ans
    yield Frame(
        Array(temps, label="temps"),
        stack_widget([f"{j}:{temps[j]}" for j in st], label="left on stack (they stay 0)"),
        Array(ans, label="ans", hl={j: "match" for j in range(n)}),
        Text(f"answer: {ans}", ("time O(n) · space O(n)", "match")),
        say="Days 6 and 7 never found a warmer day, so they keep 0. Each index is pushed and popped at "
            "most once, so it's O(n) time and O(n) space.",
    )
