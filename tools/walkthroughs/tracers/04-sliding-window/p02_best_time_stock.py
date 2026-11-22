from walkthrough import Array, Frame, Text

TITLE = "Best Time to Buy and Sell Stock"
PATTERN = "Sliding window"


def frames():
    prices = [7, 1, 5, 3, 6, 4]

    yield Frame(
        Array(prices, label="prices (day i)"),
        Text("buy on one day, sell on a later day", "return the best profit, or 0"),
        speaker="Interviewer",
        say="prices[i] is the stock price on day i. Buy once and sell on a later day. "
            "Return the best profit, or 0 if no trade makes money.",
    )
    yield Frame(
        Array(prices, label="prices (day i)"),
        Text("brute force: every (buy, sell) pair with buy < sell", "O(n²)"),
        say="Trying every buy day with every later sell day is O(n²). But for a given sell day, "
            "the only buy day that matters is the cheapest one before it.",
    )

    best, min_i = 0, 0
    best_pair = None
    yield Frame(
        Array(prices, label="prices", hl={0: "special"}, pointers={"buy": 0}),
        Text(f"minPrice = {prices[0]}, best = {best}"),
        say=f"So I keep the cheapest price seen so far as the window's left edge. Day 0 starts it at {prices[0]}.",
    )
    for day in range(1, len(prices)):
        p = prices[day]
        if p < prices[min_i]:
            old = prices[min_i]
            min_i = day
            yield Frame(
                Array(prices, label="prices", hl={day: "special", **({best_pair[0]: "visited"} if best_pair else {})},
                      pointers={"buy": min_i, "sell": day}),
                Text(f"{p} < minPrice {old}", (f"cheaper buy day → minPrice = {p}", "special"), f"best = {best}"),
                say=f"{p} is cheaper than {old}. Any later sale makes more money buying here, "
                    "so I jump the left edge to today.",
            )
        else:
            profit = p - prices[min_i]
            improved = profit > best
            if improved:
                best, best_pair = profit, (min_i, day)
            hl = {i: "window" for i in range(min_i, day + 1)}
            hl[min_i] = "special"
            hl[day] = "match" if improved else "compare"
            yield Frame(
                Array(prices, label="prices", hl=hl, window=(min_i, day), pointers={"buy": min_i, "sell": day}),
                Text(f"profit = {p} - {prices[min_i]} = {profit}",
                     (f"best = {best}" + ("  (new best)" if improved else ""), "match" if improved else None)),
                say=(f"Selling today at {p} after buying at {prices[min_i]} gives {profit}. "
                     + ("That's a new best." if improved else f"That doesn't beat {best}.")),
            )

    assert best == 5  # doc comment: [7, 1, 5, 3, 6, 4] → 5
    b, s = best_pair
    yield Frame(
        Array(prices, label="prices", hl={b: "special", s: "match"}, window=(b, s), pointers={"buy": b, "sell": s}),
        Text((f"answer: {best}  (buy at {prices[b]}, sell at {prices[s]})", "match"), "time O(n): one pass",
             "space O(1): minPrice and best"),
        say=f"The answer is {best}: buy at {prices[b]} on day {b}, sell at {prices[s]} on day {s}. "
            "One pass with two variables, so O(n) time and O(1) space.",
    )
