from fractions import Fraction

from walkthrough import Frame, Grid, Text, stack_widget

TITLE = "Car Fleet"
PATTERN = "Sort + monotonic stack"


def fmt_t(dist, speed):
    f = Fraction(dist, speed)
    return str(f.numerator) if f.denominator == 1 else f"{f.numerator}/{f.denominator}"


def frames():
    target = 12
    position = [10, 8, 0, 5, 3]
    speed = [2, 4, 1, 1, 3]
    want = 3
    n = len(position)

    yield Frame(
        Grid([position, speed], label=f"cars (target = {target})", row_labels=["pos", "speed"],
             col_labels=[f"car {i}" for i in range(n)], cell=62),
        speaker="Interviewer",
        say="Cars drive toward the target on one lane and can't pass. A faster car that catches a slower "
            "one joins it as a fleet. How many fleets arrive?",
    )
    yield Frame(
        Grid([position, speed], label=f"cars (target = {target})", row_labels=["pos", "speed"],
             col_labels=[f"car {i}" for i in range(n)], cell=62),
        Text("brute force: simulate time steps and merge cars → slow, and ties are fiddly"),
        say="Simulating the road step by step is slow and error-prone. What matters is only when each "
            "car would arrive alone, and who is in front of it.",
    )

    idx = sorted(range(n), key=lambda i: -position[i])  # closest to target first
    times = [fmt_t(target - position[i], speed[i]) for i in idx]
    cols = [f"car {i}" for i in idx]
    rows = [[position[i] for i in idx], [speed[i] for i in idx], times]
    labels = ["pos", "speed", "time"]

    def table(hl=None):
        return Grid(rows, label="sorted: closest to target first", row_labels=labels, col_labels=cols,
                    hl=hl, cell=62)

    yield Frame(
        table(hl={(2, c): "special" for c in range(n)}),
        Text("time = (target − pos) / speed   (kept as an exact fraction)"),
        say="I sort by position, closest to the target first, and compute each car's solo arrival time. "
            "A car can only be blocked by cars ahead of it.",
    )

    fleets = []  # stack of (dist, speed); arrival times grow bottom to top
    for c, i in enumerate(idx):
        t = (target - position[i], speed[i])
        ts = times[c]
        hl = {(r, k): "visited" for r in range(3) for k in range(c)}
        if fleets:
            ahead = fleets[-1]
            slower = t[0] * ahead[1] > ahead[0] * t[1]  # cross-multiply: no floats
        else:
            slower = True
        if slower:
            fleets.append(t)
            hl.update({(r, c): "active" for r in range(3)})
            if len(fleets) == 1:
                say = f"Car {i} is closest to the target and nothing blocks it. It arrives at time {ts} and starts the first fleet."
            else:
                say = (f"Car {i} would arrive at {ts}, later than the fleet ahead at "
                       f"{fmt_t(*fleets[-2])}. It never catches up, so it starts a new fleet.")
            line = (f"car {i}: time {ts} > top {fmt_t(*fleets[-2])} → push (new fleet)" if len(fleets) > 1 else
                    f"car {i}: stack empty → push (first fleet)", "match")
        else:
            hl.update({(r, c): "reject" for r in range(3)})
            say = (f"Car {i} would arrive at {ts}, no later than the fleet ahead at {fmt_t(*fleets[-1])}. "
                   "It catches up before or at the target, so it joins. No push.")
            line = (f"car {i}: time {ts} ≤ top {fmt_t(*fleets[-1])} → merges", "reject")
        yield Frame(
            table(hl=hl),
            stack_widget([fmt_t(*f) for f in fleets], label="fleet arrival times"),
            Text(line),
            say=say,
        )

    assert len(fleets) == want
    yield Frame(
        table(),
        stack_widget([fmt_t(*f) for f in fleets], label="fleet arrival times",
                     hl={k: "match" for k in range(len(fleets))}),
        Text(f"answer: {len(fleets)} fleets", ("time O(n log n) for the sort · space O(n)", "match")),
        say="Three fleets are left on the stack, so the answer is 3. Sorting dominates: O(n log n) time "
            "and O(n) space.",
    )
