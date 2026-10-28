package slidingwindow

// Problem 02 — Best Time to Buy and Sell Stock (Easy)
//
// prices[i] is a stock's price on day i. Buy on one day and sell on a LATER
// day to maximise profit. Return the best profit, or 0 if no trade makes money.
//
//	[7, 1, 5, 3, 6, 4] → 5   (buy at 1, sell at 6)
//	[7, 6, 4, 3, 1]    → 0   (prices only fall)
//
// Reference: https://leetcode.com/problems/best-time-to-buy-and-sell-stock/
//
// Pattern: a window whose left edge is the cheapest day seen so far. When
// today's price is lower than the left edge, no future sale would prefer the
// old buy day, so jump the left edge to today.
//
//	MaxProfitBrute: Time O(n²), Space O(1)
//	MaxProfit:      Time O(n),  Space O(1)

// MaxProfitBrute tries every (buy, sell) pair.
func MaxProfitBrute(prices []int) int {
	best := 0
	for buy := range prices {
		for sell := buy + 1; sell < len(prices); sell++ {
			best = max(best, prices[sell]-prices[buy])
		}
	}
	return best
}

// MaxProfit tracks the minimum price to the left of each day.
func MaxProfit(prices []int) int {
	if len(prices) == 0 {
		return 0
	}
	best, minPrice := 0, prices[0]
	for _, p := range prices[1:] {
		if p < minPrice {
			minPrice = p // cheaper buy day: move the window's left edge
		} else {
			best = max(best, p-minPrice)
		}
	}
	return best
}
