import os
import sys

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

from _heaps10_util import GoHeap  # noqa: E402
from walkthrough import Array, Frame, KV, Row, Text  # noqa: E402

TITLE = "Task Scheduler"
PATTERN = "Greedy + max-heap"


def frames():
    tasks, n = "AAABBB", 2
    want = 8

    yield Frame(
        Array(list(tasks), label="tasks"),
        Text(f"n = {n}: two runs of the same task must be at least {n} units apart"),
        speaker="Interviewer",
        say="A CPU runs one task per time unit. The same task needs at least n = 2 units between runs. "
            "What's the minimum time to finish everything?",
    )
    yield Frame(
        Array(list(tasks), label="tasks"),
        Text("brute force: try every ordering and pick the shortest → exponential"),
        say="Trying every order is exponential. Greedy works instead: at each tick, run the available task "
            "with the most copies left, because it's the one most likely to force idling later.",
    )

    counts = {}
    for t in tasks:
        counts[t] = counts.get(t, 0) + 1
    # The Go heap holds plain counts. Labels ride along here only so the timeline can show them;
    # the comparison uses the count alone, exactly like NewMax[int]().
    ready = GoHeap(lambda a, b: a[0] > b[0])
    for t in sorted(counts):
        ready.push((counts[t], t))
    queue = []  # (remaining, readyAt, label)
    timeline = []

    def widgets(hl_t=None):
        return (
            Array(timeline + [""] * (want - len(timeline)), label="timeline", hl=hl_t or {}, cell=50),
            Row(KV({f"{t}": c for c, t in ready.items}, label="ready (max-heap of counts)", empty="(empty)"),
                KV({f"{lab} left {r}": f"ready@{at}" for r, at, lab in queue}, label="cooling queue (FIFO)",
                   empty="(empty)"),
                weights=[1, 1.4]),
        )

    yield Frame(
        *widgets(),
        Text("counts: A=3, B=3 → ready heap = [3, 3]"),
        say="I count each task and push the counts into a max-heap. Tasks that just ran wait in a FIFO "
            "queue with the time they become ready again.",
    )

    time = 0
    while len(ready) or queue:
        time += 1
        if len(ready):
            c, lab = ready.pop()
            left = c - 1
            if left > 0:
                queue.append((left, time + n, lab))
            timeline.append(lab)
            what = f"t={time}: run {lab} (count {c} → {left})" + (f", cools until t={time + n}" if left else ", done")
        else:
            timeline.append("idle")
            what = f"t={time}: nothing ready → idle"
            lab = None
        back = None
        if queue and queue[0][1] == time:
            r, _, blab = queue.pop(0)
            ready.push((r, blab))
            back = blab
        line2 = f"{back} is ready again → back into the heap" if back else ""
        if time == 1:
            say = "At t=1 both have 3 left. I pop one, run A, and it cools until t=3."
        elif time == 2:
            say = "At t=2 only B is ready. It runs and cools until t=4."
        elif time == 3 and lab is None:
            say = "At t=3 the heap is empty: A is still cooling. The CPU idles. At the end of t=3, A's wait is over and it goes back in the heap."
        elif lab is None:
            say = "Idle again: both tasks are cooling. A comes back at the end of this tick."
        elif time >= 7:
            say = f"t={time}: run {lab}, its last copy." + (" Nothing is left anywhere, so the loop ends." if not len(ready) and not queue else "")
        else:
            say = f"t={time}: {lab} is the only ready task, so it runs."
        yield Frame(
            *widgets({time - 1: "reject" if lab is None else "active"}),
            Text(what, line2),
            say=say,
        )

    assert time == want
    maxc = max(counts.values())
    num_max = sum(1 for c in counts.values() if c == maxc)
    formula = max(len(tasks), (maxc - 1) * (n + 1) + num_max)
    assert formula == want
    yield Frame(
        *widgets({k: ("reject" if v == "idle" else "match") for k, v in enumerate(timeline)}),
        Text(("answer: 8   (A B idle A B idle A B)", "match"),
             f"formula: max(len={len(tasks)}, ({maxc}-1)·({n}+1) + {num_max}) = {formula}"),
        say="The answer is 8. LeastInterval gets it directly: the top count forces (3-1) frames of n+1 slots, "
            "plus one slot per task tied at the top. That's O(len(tasks)) time and O(26) space.",
    )
