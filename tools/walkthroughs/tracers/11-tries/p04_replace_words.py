import os
import sys

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

from trieview import TNode, TrieView  # noqa: E402
from walkthrough import KV, Frame, Text  # noqa: E402

TITLE = "Replace Words"
PATTERN = "Trie shortest prefix"


def insert(root, w):
    cur = root
    for ch in w:
        cur = cur.children.setdefault(ch, TNode())
    cur.end = True


def shortest_root(root, w):
    """Port of shortestRoot. Returns (answer, walked path, how it stopped)."""
    cur = root
    for i, ch in enumerate(w):
        cur = cur.children.get(ch)
        if cur is None:
            return w, w[:i], "fell"
        if cur.end:
            return w[: i + 1], w[: i + 1], "root"
    return w, w, "end"


def chips(words):
    # KV highlights by key, so make repeated words distinct with trailing spaces.
    seen, out = {}, []
    for w in words:
        k = w + " " * seen.get(w, 0)
        seen[w] = seen.get(w, 0) + 1
        out.append(k)
    return out


def frames():
    roots = ["cat", "bat", "rat"]
    sentence = "the cattle was rattled by the battery"
    expected = "the cat was rat by the bat"
    words = sentence.split()
    keys = chips(words)

    yield Frame(
        KV(roots, label="roots"),
        Text(sentence),
        speaker="Interviewer",
        say="Replace every word that starts with a root by the shortest such root. "
            "Words with no root stay as they are.",
    )
    yield Frame(
        KV(roots, label="roots"),
        Text(sentence, "brute force: HasPrefix(word, root) for every word × every root",
             ("O(words · roots · L)", "reject")),
        say="I could test every root against every word with HasPrefix and keep the shortest match. "
            "That's O(words × roots × L), and it re-reads the same letters many times.",
    )

    root = TNode()
    for r in roots:
        insert(root, r)
    yield Frame(
        TrieView(root, label="trie of roots"),
        KV(keys, label="sentence"),
        say="I put the roots in a trie. Now each word needs one walk from the root, "
            "and the first word-end I pass is the shortest root.",
    )

    out = []
    for i, w in enumerate(words):
        ans, path, how = shortest_root(root, w)
        out.append(ans)
        hl = {path[:j]: "active" for j in range(len(path) + 1)}
        if how == "root":
            hl[path] = "match"
            say = (f"\"{w}\": I walk {' → '.join(path)} and hit a word end after {len(path)} letters. "
                   f"I stop there and use \"{ans}\". No shorter root was passed on the way.")
            line = (f'"{w}" -> "{ans}"  (IsEnd at "{path}")', "match")
        else:
            hl[path] = "reject"
            missing = w[len(path)]
            if path:
                say = (f"\"{w}\": '{path}' exists, but no root continues with '{missing}'. "
                       "The walk falls off, so the word stays as it is.")
            else:
                say = f"\"{w}\": the root has no '{missing}' child, so no root can be a prefix. The word stays."
            line = (f'"{w}" -> "{ans}"  (no child \'{missing}\')', "reject")
        if i == 5:
            say = "\"the\" again: same walk as before, no 't' child, so it stays."
        yield Frame(
            TrieView(root, label="trie of roots", hl=hl),
            KV(keys, label="sentence", hl={keys[i]: "active"}),
            Text(line, "output: " + " ".join(out)),
            say=say,
        )

    result = " ".join(out)
    assert result == expected, result
    yield Frame(
        TrieView(root, label="trie of roots"),
        Text((f'answer: "{result}"', "match"),
             "time: O(R + S)   R = total root length, S = sentence length",
             "space: O(R) for the trie"),
        say="Building the trie costs the total root length. Each word walk stops at its first root, "
            "so the sentence is read at most once: O(R + S).",
    )
