from walkthrough import KV, Array, Frame, Text

TITLE = "Encode and Decode Strings"
PATTERN = "Length prefix"


def styles_for(encoded, spans, upto=None):
    """Colour each frame: digits special, '#' compare, payload match."""
    hl = {}
    for k, (d0, h, p0, p1) in enumerate(spans):
        if upto is not None and k > upto:
            break
        for j in range(d0, h):
            hl[j] = "special"
        hl[h] = "compare"
        for j in range(p0, p1):
            hl[j] = "match"
    return hl


def encode(strs):
    return "".join(f"{len(s)}#{s}" for s in strs)


def frames():
    strs = ["lint", "code", "love", "you"]
    expected = "4#lint4#code4#love3#you"

    yield Frame(
        Text("Encode([]string) string", "Decode(string) []string",
             "input: " + str(strs).replace("'", '"')),
        speaker="Interviewer",
        say="Turn a list of strings into one string, and back again exactly. "
            "The strings may contain any character, including your delimiter.",
    )
    yield Frame(
        Text('join with "#":  ["a#b"] → "a#b"  and  ["a","b"] → "a#b"',
             ("same encoding for different inputs → can't decode", "reject")),
        say="The obvious idea is a delimiter. But if a string can contain the delimiter, two different lists "
            "can encode to the same text. Escaping works but is fiddly.",
    )

    # Encode: "<byte length>#<bytes>" per string, recording where each part lands.
    out, spans = "", []
    for k, s in enumerate(strs):
        d0 = len(out)
        out += str(len(s))
        h = len(out)
        out += "#"
        p0 = len(out)
        out += s
        spans.append((d0, h, p0, len(out)))
        if k in (0, len(strs) - 1):
            yield Frame(
                Array(list(out), label="encoded", hl=styles_for(out, spans, k), cell=29),
                Text(f'write len("{s}") = {len(s)}, then "#", then "{s}"'),
                say="Instead I write each string's byte length, a '#', then the bytes. "
                    "The decoder never has to look inside the payload." if k == 0 else
                    " and ".join(f'"{w}"' for w in strs[1:-1]) + f' go the same way, then "{s}". '
                    'Purple is the length, orange the "#", green the payload.',
            )
    encoded = out
    assert encoded == expected == encode(strs)

    # Decode
    res, i = [], 0
    while i < len(encoded):
        hash_at = encoded.index("#", i) - i
        if hash_at <= 0:
            raise ValueError("malformed")
        n = int(encoded[i:i + hash_at])
        start = i + hash_at + 1
        assert start + n <= len(encoded)
        res.append(encoded[start:start + n])
        k = len(res) - 1
        hl = {j: "visited" for j in range(i)}
        hl.update({j: "special" for j in range(i, i + hash_at)})
        hl[i + hash_at] = "compare"
        hl.update({j: "match" for j in range(start, start + n)})
        if k == 0:
            say = (f"To decode, I read digits up to the first '#' from i. That gives n = {n}, "
                   f"so the next {n} bytes are the string, whatever they contain.")
        elif k == 1:
            say = f"Then i jumps to start + n = {i}. Same again: read the length, slice out {n} bytes."
        else:
            say = f'Read {n}, slice "{res[-1]}".'
        yield Frame(
            Array(list(encoded), label="encoded", hl=hl, pointers={"i": i}, cell=29),
            Text(f'digits "{encoded[i:i + hash_at]}" → n = {n}, slice [{start}:{start + n}] = "{res[-1]}"'),
            KV(res, label="out"),
            say=say,
        )
        i = start + n

    assert res == strs
    yield Frame(
        Text(("decoded == input ✓", "match"), "Encode: O(total length) time and space",
             "Decode: O(total length) time and space",
             f'["", "#"] → "{encode(["", "#"])}" also round-trips'),
        say="Decoding gives back the original list. Both directions touch each byte once: O(total length). "
            "Empty strings and '#' inside a string work too.",
    )
