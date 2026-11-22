from walkthrough import Array, Frame, Text

TITLE = "Find the Index of the First Occurrence"
PATTERN = "KMP"


def prefix_function(s):
    pi = [0] * len(s)
    for i in range(1, len(s)):
        k = pi[i - 1]
        while k > 0 and s[i] != s[k]:
            k = pi[k - 1]
        if s[i] == s[k]:
            k += 1
        pi[i] = k
    return pi


def strstr_kmp(hay, needle):
    m = len(needle)
    if m == 0:
        return 0
    pi = prefix_function(needle)
    k = 0
    for i in range(len(hay)):
        while k > 0 and hay[i] != needle[k]:
            k = pi[k - 1]
        if hay[i] == needle[k]:
            k += 1
        if k == m:
            return i - m + 1
    return -1


def frames():
    # The doc's first example ("sadbutsad", "sad") matches at once; KMP's
    # fallback only shows when a partial match fails, so we trace the
    # module README's example and still check the doc examples below.
    assert strstr_kmp("sadbutsad", "sad") == 0
    assert strstr_kmp("leetcode", "leeto") == -1
    hay, needle = "abababc", "ababc"
    m = len(needle)
    pi = prefix_function(needle)

    def view(i=None, k=0, hl_h=None, hl_n=None, extra=()):
        ptr_h = {"i": i} if i is not None else {}
        ptr_n = {"j": k} if i is not None else {}
        return (
            Array(list(hay), label="haystack", hl=hl_h or {}, pointers=ptr_h),
            Array(list(needle), label="needle", hl=hl_n or {}, pointers=ptr_n, cell=44),
            Array(pi, label="pi (longest border of needle[:j+1])", cell=44, show_index=False),
            Text(*extra) if extra else Text(""),
        )

    yield Frame(
        Array(list(hay), label="haystack"),
        Array(list(needle), label="needle", cell=44),
        Text(f'find "{needle}" in "{hay}"'),
        speaker="Interviewer",
        say="Return the first index where needle appears in haystack, or -1. Can you do it in linear time?",
    )
    yield Frame(
        Array(list(hay), label="haystack"),
        Array(list(needle), label="needle", cell=44),
        Text("brute force: try every start, compare up to m chars → O(n·m)"),
        say="The simple way tries every start and compares up to m characters. That's O(n·m). "
            "After a mismatch it throws away everything it just matched and starts over.",
    )
    yield Frame(
        *view(extra=("pi[j] = longest proper prefix of needle[:j+1] that is also its suffix",
                     ("pi[3] = 2: \"abab\" ends with its own prefix \"ab\"", "special"))),
        say="KMP first builds pi for the needle. pi[3] = 2 says \"abab\" ends with \"ab\", its own start. "
            "So after a failure I can keep those 2 chars matched instead of restarting.",
    )

    k = 0
    for i in range(len(hay)):
        fell = []
        while k > 0 and hay[i] != needle[k]:
            fell.append((k, pi[k - 1]))
            k = pi[k - 1]
        if fell:
            old_k, new_k = fell[0]
            yield Frame(
                *view(i, old_k, hl_h={**{j: "visited" for j in range(i - old_k, i - new_k)},
                                      **{j: "match" for j in range(i - new_k, i)}, i: "reject"},
                      hl_n={**{j: "match" for j in range(new_k)}, old_k: "reject"},
                      extra=(f"haystack[{i}] = '{hay[i]}' ≠ needle[{old_k}] = '{needle[old_k]}'",
                             (f"fall back: k = pi[{old_k - 1}] = {new_k}  (i stays at {i})", "reject"))),
                say=f"Mismatch after {old_k} matched chars. Instead of moving i back, I set k = pi[{old_k - 1}] = {new_k}: "
                    f"the last {new_k} chars I read already match the needle's first {new_k}.",
            )
        if hay[i] == needle[k]:
            k += 1
            hl_n = {j: "match" for j in range(k)}
            hl_h = {j: "match" for j in range(i - k + 1, i + 1)}
            if k == m:
                ans = i - m + 1
                assert ans == strstr_kmp(hay, needle) == 2
                yield Frame(
                    *view(i, k - 1, hl_h=hl_h, hl_n=hl_n,
                          extra=(f"k = {k} = len(needle) → match at i − m + 1 = {i} − {m} + 1",
                                 (f"answer: {ans}", "match"))),
                    say=f"k reached {m}, so the whole needle matched, ending at i = {i}. The answer is {ans}. "
                        "i only ever moved forward: O(n + m) time and O(m) space for pi.",
                )
                return
            yield Frame(
                *view(i, k - 1, hl_h=hl_h, hl_n=hl_n,
                      extra=(f"haystack[{i}] = '{hay[i]}' = needle[{k - 1}]", (f"k = {k}", "match"))),
                say=f"'{hay[i]}' matches needle[{k - 1}], so {k} char{'s' if k > 1 else ''} of the needle now match. "
                    "Move i forward.",
            )
