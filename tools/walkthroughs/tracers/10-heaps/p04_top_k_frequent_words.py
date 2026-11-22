import os
import sys

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

from _heaps10_util import GoHeap, heap_view  # noqa: E402
from walkthrough import Array, Frame, KV, Text  # noqa: E402

TITLE = "Top K Frequent Words"
PATTERN = "Size-k heap"


def better(a, b):
    """a ranks ahead of b: higher count, then alphabetically first."""
    if a[1] != b[1]:
        return a[1] > b[1]
    return a[0] < b[0]


def wc(x):
    return f"{x[0]}:{x[1]}"


def frames():
    words, k = ["i", "love", "go", "i", "love", "code"], 2

    yield Frame(
        Array(words, label="words"),
        Text(f"k = {k}; order: most frequent first, ties alphabetical"),
        speaker="Interviewer",
        say="Return the 2 most frequent words, most frequent first. Words with the same count go in alphabetical order.",
    )

    counts = {}
    for w in words:
        counts[w] = counts.get(w, 0) + 1
    yield Frame(
        Array(words, label="words"),
        KV(counts, label="counts"),
        Text("brute force: sort all u distinct words by (count desc, word asc) → O(u log u)"),
        say="First I count with a map. Sorting all distinct words then works, in O(u log u). "
            "With a heap of size k I only pay log k per word.",
    )

    # Go ranges over the map in random order; the result doesn't depend on it.
    h = GoHeap(lambda a, b: better(b, a))  # worst candidate on top
    for w, c in counts.items():
        h.push((w, c))
        evicted = h.pop() if len(h) > k else None
        if evicted is None:
            line = f"push {w}:{c} → size {len(h)}"
            say = ("The heap's top is the WORST of my current picks: lowest count, and for a tie the "
                   "alphabetically later word." if len(h) == 1 else
                   f"Push {w}:{c}. It ties with i:2, and 'love' comes after 'i', so love is worse and sits on top.")
        else:
            line = f"push {w}:{c} → size {k + 1} > k, evict worst = {wc(evicted)}"
            say = (f"Push {w}:{c}. Now there are 3, so I evict the worst. {wc(evicted)} has the lowest count, so it goes."
                   if evicted[0] == w else f"Evict the worst: {wc(evicted)}.")
        yield Frame(
            Array(words, label="words"),
            KV(counts, label="counts", hl={w: "active"}),
            heap_view(h.items, "heap (worst on top)", fmt=wc, tree_fmt=lambda x: x[0]),
            Text(line),
            say=say,
        )

    out = [None] * len(h)
    pops = []
    for i in range(len(out) - 1, -1, -1):
        x = h.pop()
        pops.append(wc(x))
        out[i] = x[0]
    assert out == ["i", "love"], out
    yield Frame(
        Array(words, label="words"),
        KV(counts, label="counts"),
        Array(out, label="answer (filled from the back)", hl={0: "match", 1: "match"}, show_index=True),
        Text(f"pops come out worst-first: {', '.join(pops)} → write from the end"),
        say="Pops come out worst first: love, then i. So I fill the answer from the back, which gives [i, love] "
            "without a separate reverse.",
    )
    yield Frame(
        Array(words, label="words"),
        KV(counts, label="counts"),
        Array(out, label="answer", hl={0: "match", 1: "match"}),
        Text(("answer: [i, love]", "match"), "time O(n + u log k), space O(u)"),
        say="The answer is [i, love]. Counting is O(n), the heap work is O(u log k) for u distinct words, "
            "and the map takes O(u) space.",
    )
