from walkthrough import KV, Array, Frame, Text

TITLE = "Group Anagrams"
PATTERN = "Canonical key"


def letter_counts(w):
    c = [0] * 26
    for ch in w:
        if "a" <= ch <= "z":
            c[ord(ch) - 97] += 1
    return tuple(c)


def show_key(key):
    """Print a [26]int compactly: only the non-zero counters, e.g. a1 e1 t1."""
    return " ".join(f"{chr(97 + k)}{v}" for k, v in enumerate(key) if v) or "all zero"


def frames():
    words = ["eat", "tea", "tan", "ate", "nat", "bat"]
    expected = [["eat", "tea", "ate"], ["tan", "nat"], ["bat"]]

    yield Frame(
        Array(words, label="words", cell=70),
        speaker="Interviewer",
        say="Group the words so that each group holds anagrams of each other. "
            "Keep groups in order of their first word.",
    )
    yield Frame(
        Array(words, label="words", cell=70, hl={0: "compare", 1: "compare"}),
        Text("brute force: compare each word with every group's first word",
             ("O(n²·k) with n words of length k", "reject")),
        say="I could test each word against every existing group with IsAnagram. That's O(n²·k). "
            "Better: give all anagrams the same key and bucket by it in a map.",
    )

    group_idx = {}
    groups = []
    for wi, w in enumerate(words):
        key = letter_counts(w)
        new = key not in group_idx
        if new:
            group_idx[key] = len(groups)
            groups.append([])
        i = group_idx[key]
        groups[i].append(w)
        keys_view = {show_key(k): v for k, v in group_idx.items()}
        groups_view = {g: "[" + ", ".join(ws) + "]" for g, ws in enumerate(groups)}
        if wi == 0:
            say = ("The key is the [26]int of letter counts. Arrays are comparable in Go, so the array itself "
                   "is a map key. 'eat' gets a new key, so it starts group 0.")
        elif new:
            say = f"'{w}' has a key I haven't seen, so it opens group {i}."
        else:
            say = f"'{w}' has the same counts as group {i}'s key, so it joins group {i}."
        yield Frame(
            Array(words, label="words", cell=70, hl={k: "visited" for k in range(wi)} | {wi: "active"},
                  pointers={"w": wi}),
            Text(f"letterCounts(\"{w}\") = {show_key(key)}  → " + (f"new group {i}" if new else f"group {i}")),
            KV(keys_view, label="groupIdx (key: group index)", hl={show_key(key): "active" if new else "match"}),
            KV(groups_view, label="groups", hl={i: "active" if new else "match"}),
            say=say,
        )

    assert groups == expected
    yield Frame(
        KV({g: "[" + ", ".join(ws) + "]" for g, ws in enumerate(groups)}, label="groups",
           hl={g: "match" for g in range(len(groups))}),
        Text(("answer: [" + " ".join("[" + " ".join(g) + "]" for g in groups) + "]", "match"),
             "GroupAnagramsSorted: O(n·k log k)   sorted-word key",
             "GroupAnagrams:       O(n·k)         [26]int key",
             "space O(n·k) for both"),
        say="Done. The count key costs O(k) per word, so O(n·k) overall. The group-index map keeps the output "
            "in input order, since Go map order is random.",
    )
