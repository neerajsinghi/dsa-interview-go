from walkthrough import Array, Frame, Text

TITLE = "String to Integer (atoi)"
PATTERN = "State machine"

MAX32 = 2**31 - 1
SP = "␣"


def frames():
    # "42" (first doc example) has no spaces or sign; the second shows every state.
    s = "   -042abc"
    cells = [SP if c == " " else c for c in s]
    n = len(s)

    def arr(hl=None, i=None):
        return Array(cells, label="s", hl=hl or {}, pointers={"i": i} if i is not None else {})

    yield Frame(
        arr(),
        Text(f'"{s}"'),
        speaker="Interviewer",
        say="Parse a 32-bit signed integer from the start of s: skip spaces, read an optional sign, read digits, "
            "stop at the first non-digit, and clamp to the int32 range.",
    )
    yield Frame(
        arr(),
        Text("naive: collect digits, then strconv.Atoi / ParseInt", ("a long digit string overflows even int64", "reject")),
        say="I could cut out the digit run and call strconv. But a hundred digits overflow even int64, "
            "so I must check for overflow while I read, before each multiply.",
    )

    i = 0
    while i < n and s[i] == " ":
        i += 1
    yield Frame(
        arr({k: "visited" for k in range(i)}, i),
        Text(f"state 1: skipped {i} spaces → i = {i}"),
        say="State one: skip leading spaces. i moves to the first non-space character.",
    )
    negative = False
    if i < n and s[i] in "+-":
        negative = s[i] == "-"
        i += 1
    limit = MAX32 + (1 if negative else 0)
    yield Frame(
        arr({**{k: "visited" for k in range(i - 1)}, i - 1: "special"}, i),
        Text(f"state 2: sign '{s[i - 1]}' → negative = {str(negative).lower()}",
             (f"limit = {limit} (|MinInt32| is one more than MaxInt32)", "special")),
        say="State two: an optional sign. It's '-', so the magnitude may go up to 2147483648, "
            "one more than MaxInt32.",
    )

    result = 0
    start = i
    while i < n and "0" <= s[i] <= "9":
        d = ord(s[i]) - 48
        bound = (limit - d) // 10
        hl = {**{k: "visited" for k in range(start)}, **{k: "match" for k in range(start, i)}, i: "active"}
        if result > bound:
            result = limit
            break
        new = result * 10 + d
        yield Frame(
            arr(hl, i),
            Text(f"digit {d}: is result {result} > (limit − {d}) / 10 = {bound}?  no",
                 (f"result = {result}·10 + {d} = {new}", "match")),
            say=f"State three reads digit {d}. Before multiplying I check result > (limit − d) / 10. "
                "That's result·10 + d > limit rearranged, so the check itself can't overflow.",
        )
        result = new
        i += 1

    yield Frame(
        arr({**{k: "visited" for k in range(start)}, **{k: "match" for k in range(start, i)}, i: "reject"}, i),
        Text(f"s[{i}] = '{s[i]}' is not a digit → stop", f"result = {result}"),
        say=f"'{s[i]}' isn't a digit, so parsing stops here. Whatever follows is ignored.",
    )
    ans = -result if negative else result
    assert ans == -42
    yield Frame(
        arr({k: "match" for k in range(start - 1, i)}),
        Text((f"answer: {ans}", "match"), "time O(n) · space O(1)"),
        say="Apply the sign and return -42. One pass over the string: O(n) time and O(1) space.",
    )
