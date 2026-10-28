package stacks

// Problem: Daily Temperatures
// Difficulty: Medium
//
// Given a list of daily temperatures, produce a list where position i holds
// how many days you must wait after day i to see a strictly warmer day. If no
// warmer day ever comes, store 0.
//
// Examples:
//
//	[73,74,75,71,69,72,76,73] -> [1,1,4,2,1,1,0,0]
//	[30,20,10]                -> [0,0,0]
//
// Reference: https://leetcode.com/problems/daily-temperatures/
//
// Pattern: monotonic decreasing stack of indices. Each index waits on the
// stack until a warmer day arrives; that day "resolves" every colder index on
// top of the stack at once. Every index is pushed and popped at most once.
//
// Time: brute O(n^2); optimal O(n). Space: O(n).

// DailyTemperaturesBrute scans forward from every day.
func DailyTemperaturesBrute(temps []int) []int {
	ans := make([]int, len(temps))
	for i := range temps {
		for j := i + 1; j < len(temps); j++ {
			if temps[j] > temps[i] {
				ans[i] = j - i
				break
			}
		}
	}
	return ans
}

// DailyTemperatures uses a monotonic stack of unresolved day indices.
func DailyTemperatures(temps []int) []int {
	ans := make([]int, len(temps))
	var st Stack[int] // indices; their temperatures are decreasing bottom->top
	for i, t := range temps {
		for {
			j, ok := st.Peek()
			if !ok || temps[j] >= t {
				break
			}
			st.Pop()
			ans[j] = i - j // day i is the first warmer day after j
		}
		st.Push(i)
	}
	return ans // indices still on the stack keep their default 0
}
