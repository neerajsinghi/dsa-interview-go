import os
import sys

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

from trieview import TNode, TrieView  # noqa: E402
from walkthrough import KV, Frame, Text  # noqa: E402

TITLE = "Longest Common Prefix"
PATTERN = "Trie trunk"


def insert(root, w):
    cur = root
    for ch in w:
        cur = cur.children.setdefault(ch, TNode())
    cur.end = True


def frames():
    strs = ["flower", "flow", "flight"]
    expected = "fl"

    yield Frame(
        KV(strs, label="strs"),
        speaker="Interviewer",
        say="Given a list of words, return the longest string that every word starts with. "
            "If they share nothing, return the empty string.",
    )
    yield Frame(
        KV(strs, label="strs"),
        Text("brute force: for each prefix of strs[0], check HasPrefix on every word",
             ("O(L · S): the same characters get compared again and again", "reject")),
        say="The naive way tries each prefix of the first word and checks every word for it. "
            "That re-reads the same characters for every prefix length, so it's O(L·S).",
    )

    root = TNode()
    for k, w in enumerate(strs):
        insert(root, w)
        hl = {w[:i]: "active" for i in range(len(w) + 1)}
        yield Frame(
            KV(strs, label="strs", hl={w: "active"}),
            TrieView(root, label="trie", hl=hl),
            Text(f'Insert("{w}")'),
            say=(f"I insert every word into a trie. Shared beginnings share nodes, so \"{w}\" "
                 + ("starts a single chain." if k == 0 else "reuses the nodes it has in common and then branches off.")),
        )

    # Follow the trunk: one child and no word ending here.
    prefix = ""
    cur = root
    while not cur.end:
        kids = cur.kids()
        hl = {prefix[:i]: "match" for i in range(len(prefix) + 1)}
        if len(kids) > 1:
            hl[prefix] = "reject"
            for ch, _ in kids:
                hl[prefix + ch] = "compare"
            yield Frame(
                KV(strs, label="strs"),
                TrieView(root, label="trie", hl=hl),
                Text(f'node "{prefix}" has {len(kids)} children: {", ".join(c for c, _ in kids)}',
                     ("the words split here -> stop", "reject")),
                say=f"The node for \"{prefix}\" has two children, {kids[0][0]} and {kids[1][0]}. "
                    "The words disagree on the next letter, so the common prefix ends here.",
            )
            break
        if not kids:
            break
        ch, nxt = kids[0]
        hl[prefix + ch] = "active"
        yield Frame(
            KV(strs, label="strs"),
            TrieView(root, label="trie", hl=hl),
            Text(f'node "{prefix}" has one child "{ch}" and no word ends here',
                 (f"prefix = \"{prefix + ch}\"", "match")),
            say=f"From {'the root' if not prefix else 'node ' + repr(prefix)} there is exactly one child, '{ch}', "
                "and no word ends here. Every word continues with that letter, so I add it.",
        )
        prefix += ch
        cur = nxt

    assert prefix == expected, prefix
    yield Frame(
        KV(strs, label="strs"),
        TrieView(root, label="trie", hl={prefix[:i]: "match" for i in range(len(prefix) + 1)}),
        Text((f'answer: "{prefix}"', "match"),
             "trie: O(S) time, O(S) space",
             "vertical scan: O(S) time, O(1) space"),
        say="The answer is \"fl\". Building the trie is O(S) for S total characters. "
            "For a one-off query the vertical scan is just as fast with O(1) space.",
    )
