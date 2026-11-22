from walkthrough import Frame, Grid, Text

TITLE = "Gas Station"
PATTERN = "Greedy restart"


def table(gas, cost, hl=None):
    diff = [g - c for g, c in zip(gas, cost)]
    return Grid([gas, cost, diff], row_labels=["gas", "cost", "gas-cost"],
                col_labels=[f"st {i}" for i in range(len(gas))], hl=hl or {}, cell=62)


def frames():
    gas, cost = [1, 2, 3, 4, 5], [3, 4, 5, 1, 2]
    expected = 3
    n = len(gas)

    yield Frame(
        table(gas, cost),
        speaker="Interviewer",
        say="Stations sit on a circle. You fill gas[i] at station i and burn cost[i] to drive to the next one. "
            "Starting empty, which station lets you finish a lap?",
    )
    yield Frame(
        table(gas, cost),
        Text("brute force: simulate a lap from every start → O(n²)"),
        say="I could simulate a full lap from every station. That's n starts times n steps, O(n²).",
    )
    yield Frame(
        table(gas, cost, hl={(2, j): "special" for j in range(n)}),
        Text("only gas[i] - cost[i] matters", "total = sum of the last row"),
        say="Only the net gain at each station matters. If the total is negative, no start works. "
            "If my tank goes negative, I restart after that station.",
    )

    total, tank, start = 0, 0, 0
    for i in range(n):
        diff = gas[i] - cost[i]
        total += diff
        old = tank
        tank += diff
        hl = {(r, j): ("visited" if j < start else "window") for r in range(3) for j in range(0, i)}
        for r in range(3):
            hl[(r, i)] = "reject" if tank < 0 else "active"
        line1 = f"i = {i}: tank = {old} + ({diff}) = {tank}    total = {total}"
        if tank < 0:
            say = (f"The tank drops to {tank}, so start {start} can't reach station {i + 1}. "
                   f"No station from {start} to {i} can either, so the next candidate is {i + 1}.")
            if i > 0:
                say = (f"Tank is {tank} again. Every station up to {i} fails, "
                       f"so I move the candidate to {i + 1} and reset the tank.")
            start, tank = i + 1, 0
            yield Frame(table(gas, cost, hl), Text(line1, (f"tank < 0 → start = {start}, tank = 0", "reject")),
                        say=say)
        else:
            yield Frame(table(gas, cost, hl), Text(line1, (f"tank ≥ 0 → keep start = {start}", "match")),
                        say=f"Starting from {start}, the tank is {tank} after station {i}. Still fine, so I keep going.")

    answer = -1 if total < 0 else start
    assert answer == expected
    hl = {(r, answer): "match" for r in range(3)}
    yield Frame(
        table(gas, cost, hl),
        Text(f"total = {total} ≥ 0", (f"answer: station {answer}", "match"), "O(n) time, O(1) space"),
        say=f"The total is {total}, not negative, so a lap exists, and it must start at {answer}. "
            "One pass, three variables: O(n) time, O(1) space.",
    )
