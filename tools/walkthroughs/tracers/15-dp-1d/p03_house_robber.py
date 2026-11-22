from walkthrough import Array, Frame, Text

TITLE = "House Robber I & II"
PATTERN = "1-D DP"


def rob(nums):
    skip, take = 0, 0
    for x in nums:
        skip, take = take, max(take, skip + x)
    return take


def frames():
    nums = [1, 2, 3, 1]
    n = len(nums)
    dp = [None] * (n + 1)

    def widgets(hl_n=None, hl_dp=None):
        return (
            Array(nums, label="nums (money in each house)", hl=hl_n),
            Array(dp, label="best[i] = most money from the first i houses", hl=hl_dp),
        )

    yield Frame(
        *widgets(),
        Text("no two adjacent houses can both be robbed"),
        speaker="Interviewer",
        say="Each house holds some money. You can't rob two neighbouring houses. "
            "What's the most you can take? Then: what if the street is a circle?",
    )
    yield Frame(
        *widgets(),
        Text("try every subset with no two neighbours", ("2^n subsets → too slow", "reject")),
        say="The brute force tries every set of houses with no two side by side. That's 2^n subsets. "
            "But each choice only depends on the best answer for a shorter street.",
    )

    dp[0], dp[1] = 0, nums[0]
    yield Frame(
        *widgets(hl_n={0: "match"}, hl_dp={0: "match", 1: "match"}),
        Text("best[0] = 0  (no houses)", f"best[1] = nums[0] = {nums[0]}"),
        say="State: best[i] is the most money from the first i houses. With no houses it's 0. "
            "With one house I just rob it.",
    )
    for i in range(2, n + 1):
        skip, take = dp[i - 1], dp[i - 2] + nums[i - 1]
        dp[i] = max(skip, take)
        took = take > skip
        yield Frame(
            *widgets(hl_n={i - 1: "match" if took else "reject"},
                     hl_dp={i: "active", i - 1: "compare", i - 2: "compare"}),
            Text(f"skip house {i-1}: best[{i-1}] = {skip}",
                 f"rob house {i-1}:  best[{i-2}] + nums[{i-1}] = {dp[i-2]} + {nums[i-1]} = {take}",
                 (f"best[{i}] = max({skip}, {take}) = {dp[i]}", "active")),
            say=("For house 1 I either skip it and keep best[1], or rob it and add its money to best[0], "
                 "which never touches house 0. Robbing wins." if i == 2 else
                 f"Same take-or-skip choice for house {i-1}. "
                 + ("Robbing it on top of best[%d] gives more." % (i - 2) if took else
                    "Skipping it is better, so best carries over.")),
        )

    assert dp[n] == 4 and rob(nums) == 4
    a, b = rob(nums[:-1]), rob(nums[1:])
    circ = max(a, b)
    assert circ == 4
    yield Frame(
        *widgets(hl_n={0: "match", 2: "match"}, hl_dp={n: "match"}),
        Text(f"part I answer: best[{n}] = {dp[n]}  (houses 0 and 2)",
             f"part II (circle): max(Rob(houses 0..2), Rob(houses 1..3)) = max({a}, {b}) = {circ}",
             ("O(n) time · O(1) space (two variables)", "match")),
        say=f"The line answer is {dp[n]}. On a circle, houses 0 and 3 can't both be robbed, so I solve two lines, "
            "without the last house and without the first, and keep the better one.",
    )
