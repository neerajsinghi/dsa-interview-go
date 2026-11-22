import math

from walkthrough import Array, Frame, Text

TITLE = "Coin Change I & II"
PATTERN = "1-D DP · unbounded knapsack"


def frames():
    coins, amount = [1, 2, 5], 11
    dp = [None] * (amount + 1)

    def arr(hl=None):
        return Array(dp, label="fewest[a] = fewest coins that sum to a", hl=hl)

    yield Frame(
        arr(),
        Text(f"coins = {coins}, amount = {amount}", "unlimited coins of each kind"),
        speaker="Interviewer",
        say="You have unlimited coins of 1, 2 and 5. What's the fewest number of coins that adds up to 11? "
            "Return -1 if it can't be done.",
    )
    yield Frame(
        arr(),
        Text("recurse: 1 + min over coins of fewest(a - c)", ("fewest(6), fewest(4)... recomputed → exponential", "reject")),
        say="Recursion tries every coin as the last one. That's exponential, because fewest(6) is reached "
            "from 11-5, from 7-1, from 8-2 and so on.",
    )

    dp[0] = 0
    yield Frame(
        arr({0: "match"}),
        Text("fewest[0] = 0  (no coins needed)", "fill a = 1, 2, ..., 11 in increasing order"),
        say="So I fill a table over amounts. Making 0 needs no coins. Every bigger amount reads smaller ones, "
            "so I go upward.",
    )

    shown = 0
    for a in range(1, amount + 1):
        best = math.inf
        parts, hl = [], {a: "active"}
        for c in coins:
            if c <= a:
                parts.append(f"{dp[a-c]}+1")
                hl[a - c] = "compare"
                best = min(best, dp[a - c] + 1)
        dp[a] = best
        terms = ", ".join(f"fewest[{a-c}]" for c in coins if c <= a)
        detail = f"fewest[{a}] = 1 + min({terms}) = min({', '.join(parts)}) = {dp[a]}"
        if a <= 5 or a == amount:
            shown += 1
            if a == 1:
                say = "Only the 1-coin fits under 1, so fewest[1] = fewest[0] + 1 = 1."
            elif a == 2:
                say = "For 2, the last coin is a 1 (after fewest[1]) or a 2 (after fewest[0]). The 2-coin wins: 1 coin."
            elif a == 5:
                say = "At 5 the 5-coin fits for the first time, and fewest[0] + 1 = 1 beats everything else."
            elif a == amount:
                say = (f"For 11 I look at fewest[10], fewest[9] and fewest[6]. "
                       f"Both 10 and 6 give {dp[a]} coins. That's 5 + 5 + 1.")
            else:
                say = f"Same rule for {a}: try each coin as the last one and keep the smallest count."
            yield Frame(arr(hl), Text(detail), say=say)
        elif a == amount - 1:
            yield Frame(
                arr({k: "visited" for k in range(6, a + 1)}),
                Text("fewest[6..10] = " + ", ".join(str(dp[k]) for k in range(6, a + 1))),
                say="Same again for 6 through 10. Each cell does one check per coin, so the whole table costs "
                    "coins × amount.",
            )

    assert dp[amount] == 3
    yield Frame(
        arr({amount: "match", 10: "compare", 5: "compare", 0: "compare"}),
        Text(f"answer: fewest[{amount}] = {dp[amount]}  (5 + 5 + 1)",
             ("O(len(coins) · amount) time · O(amount) space", "match"),
             "unreachable amounts stay at a huge value → return -1"),
        say="The answer is 3. Amounts that can't be made keep a huge sentinel, and I return -1 for those. "
            "Time is O(coins × amount), space O(amount).",
    )
