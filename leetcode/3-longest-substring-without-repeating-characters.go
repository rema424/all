package leetcode

import (
	"math"
)

// Given a string, find the length of the longest substring without repeating characters.

// Example 1:

// Input: "abcabcbb"
// Output: 3
// Explanation: The answer is "abc", with the length of 3.
// Example 2:

// Input: "bbbbb"
// Output: 1
// Explanation: The answer is "b", with the length of 1.
// Example 3:

// Input: "pwwkew"
// Output: 3
// Explanation: The answer is "wke", with the length of 3.
//              Note that the answer must be a substring, "pwke" is a subsequence and not a substring.

// Approach 1: Brute Force
func LengthOfLongestSubstring_1(s string) int {
	r := []rune(s)
	n := len(r)
	ans := 0
	unique := func(r []rune, start int, end int) bool {
		var c string
		m := map[string]bool{}
		for i := start; i < end; i++ {
			c = string(r[i])
			if m[c] {
				return false
			} else {
				m[c] = true
			}
		}
		return true
	}

	for i := 0; i < n; i++ {
		for j := i + 1; j <= n; j++ {
			if unique(r, i, j) {
				ans = int(math.Max(float64(ans), float64(j-i)))
			}
		}
	}
	return ans
}

// Approach 2: Sliding Window
func LengthOfLongestSubstring_2(s string) int {
	r := []rune(s)
	n := len(r)
	var ans, i, j int
	m := map[string]bool{}
	// iは左端、jは右端
	for i < n && j < n {
		c := string(r[j])
		if !m[c] {
			m[c] = true
			j++
			ans = int(math.Max(float64(ans), float64(j-i)))
		} else {
			delete(m, string(r[i]))
			i++
		}
	}
	return ans
}

// Approach 3: Sliding Window Optimized
func LengthOfLongestSubstring_3(s string) int {
	r := []rune(s)
	n := len(r)
	ans := 0
	start := 0
	nextStart := map[string]int{}
	for end := 0; end < n; end++ {
		c := string(r[end])
		if _, ok := nextStart[c]; ok {
			start = int(math.Max(float64(nextStart[c]), float64(start)))
		}
		ans = int(math.Max(float64(ans), float64(end+1-start)))
		nextStart[c] = end + 1 // 次cが出現したら、開始位置は今のcを取り除いた一個右になる
	}
	return ans
}
