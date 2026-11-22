from walkthrough import Frame, Grid, Text

TITLE = "Reverse Bits"
PATTERN = "Bits · masks"

W = 32
M32 = 0xFFFFFFFF
COLS = list(range(W - 1, -1, -1))
COL_LABELS = [k if k % 8 in (0, 7) else "" for k in COLS]


def bits(x):
    x &= M32
    return [(x >> k) & 1 for k in COLS]


def reverse_loop(n):
    out = 0
    for _ in range(32):
        out = (out << 1 | n & 1) & M32
        n >>= 1
    return out


def block_hl(row, size, flipped=False):
    """Colour each bit by the block it belongs to (before) or came from (after)."""
    hl = {}
    for j, k in enumerate(COLS):
        src = k ^ size if flipped else k
        hl[(row, j)] = "active" if (src // size) % 2 else "compare"
    return hl


def frames():
    n0 = 0b00000010100101000001111010011100
    expected = 0b00111001011110000010100101000000

    yield Frame(
        Grid([bits(n0)], row_labels=["n"], col_labels=COL_LABELS),
        Text(f"n = {n0}"),
        speaker="Interviewer",
        say="Reverse all 32 bits of an unsigned 32-bit integer: bit 0 swaps with bit 31, bit 1 with bit 30, and so on.",
    )
    yield Frame(
        Grid([bits(n0)], row_labels=["n"], col_labels=COL_LABELS, hl={(0, W - 1): "active"}),
        Text("loop: out = out<<1 | n&1;  n >>= 1   (32 times)"),
        say="The loop pops the low bit of n and pushes it onto out, 32 times. It must run all 32, "
            "or leading zeros end up in the wrong place. Masks need only 5 steps.",
    )

    steps = [
        (16, "n = n>>16 | n<<16", lambda n: (n >> 16 | n << 16) & M32, "16-bit halves"),
        (8, "n = (n&0xFF00FF00)>>8 | (n&0x00FF00FF)<<8",
         lambda n: ((n & 0xFF00FF00) >> 8 | (n & 0x00FF00FF) << 8) & M32, "bytes within each half"),
        (4, "n = (n&0xF0F0F0F0)>>4 | (n&0x0F0F0F0F)<<4",
         lambda n: ((n & 0xF0F0F0F0) >> 4 | (n & 0x0F0F0F0F) << 4) & M32, "nibbles within each byte"),
        (2, "n = (n&0xCCCCCCCC)>>2 | (n&0x33333333)<<2",
         lambda n: ((n & 0xCCCCCCCC) >> 2 | (n & 0x33333333) << 2) & M32, "pairs within each nibble"),
        (1, "n = (n&0xAAAAAAAA)>>1 | (n&0x55555555)<<1",
         lambda n: ((n & 0xAAAAAAAA) >> 1 | (n & 0x55555555) << 1) & M32, "neighbouring bits"),
    ]
    n = n0
    for size, code, fn, what in steps:
        new = fn(n)
        hl = block_hl(0, size) | block_hl(1, size, flipped=True)
        if size == 16:
            say = ("First I swap the two 16-bit halves. Reversing 32 bits means reversing the order of the halves, "
                   "then reversing inside each half.")
        elif size == 8:
            say = ("Now each half needs reversing, so I swap its two bytes. One mask picks the high bytes, "
                   "the other the low bytes, and I shift each 8 places.")
        elif size == 1:
            say = "Last, I swap neighbouring bits. Every block is now one bit, so the whole word is reversed."
        else:
            say = f"Same again at half the size: swap the {what}, using masks and a shift of {size}."
        yield Frame(
            Grid([bits(n), bits(new)], row_labels=["before", "after"], col_labels=COL_LABELS, hl=hl),
            Text(code, f"swap {what} (blocks of {size})"),
            say=say,
        )
        n = new

    assert n == expected == reverse_loop(n0)
    yield Frame(
        Grid([bits(n0), bits(n)], row_labels=["n", "result"], col_labels=COL_LABELS,
             hl={(1, j): "match" for j in range(W)}),
        Text((f"answer: {n}  (0x{n:08X})", "match"), "5 mask steps vs 32 loop steps, O(1) space"),
        say="The result is the input read backwards. Five mask steps instead of 32 shifts, both O(1) space. "
            "In real code I'd call bits.Reverse32.",
    )
