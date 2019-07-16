package leetcode

// Given an array of integers, return indices of the two numbers such that they add up to a specific target.

// You may assume that each input would have exactly one solution, and you may not use the same element twice.

// Example:

// Given nums = [2, 7, 11, 15], target = 9,

// Because nums[0] + nums[1] = 2 + 7 = 9,
// return [0, 1].

// Approach 1: Brute Force
func twoSum_1(nums []int, target int) []int {
	for i := 0; i < len(nums)-1; i++ {
		for j := i + 1; j < len(nums); j++ {
			if nums[i]+nums[j] == target {
				return []int{i, j}
			}
		}
	}
	return nil
}

// Approach 2: One-pass Hash Table
func twoSum_2(nums []int, target int) []int {
	m := make(map[int]int)
	for i, n := range nums {
		if val, ok := m[n]; ok {
			return []int{val, i}
		} else {
			m[target-n] = i
		}
	}
	return nil
}
