// 4. Median of Two Sorted Arrays
// Hard

// 4596

// 642

// Favorite

// Share
// There are two sorted arrays nums1 and nums2 of size m and n respectively.
// Find the median of the two sorted arrays. The overall run time complexity should be O(log (m+n)).
// You may assume nums1 and nums2 cannot be both empty.

// Example 1:

// nums1 = [1, 3]
// nums2 = [2]

// The median is 2.0
// Example 2:

// nums1 = [1, 2]
// nums2 = [3, 4]

// The median is (2 + 3)/2 = 2.5

package leetcode

import (
	"fmt"
	"math"
	"math/rand"
)

// Approach 1: Recursive Approach
func findMedianSortedArrays_1(nums1 []int, nums2 []int) float64 {
	m := len(nums1)
	n := len(nums2)

	if m > n { // to ensure m <= n
		nums1, nums2 = nums2, nums1
		m, n = n, m
	}

	iMin := 0
	iMax := m
	halfLen := (m + n + 1) / 2
	for iMin <= iMax {
		i := (iMin + iMax) / 2
		j := halfLen - i

		if i < iMax && nums2[j-1] > nums1[i] {
			iMin = i + 1 // i is too small
		} else if i > iMin && nums1[i-1] > nums2[j] {
			iMax = i - 1 // i is too big
		} else { // i is perfect
			maxLeft := 0
			if i == 0 {
				maxLeft = nums2[j-1]
			} else if j == 0 {
				maxLeft = nums1[i-1]
			} else {
				maxLeft = int(math.Max(float64(nums1[i-1]), float64(nums2[j-1])))
			}

			if (m+n)%2 == 1 {
				return float64(maxLeft)
			}

			minRight := 0
			if i == m {
				minRight = nums2[j]
			} else if j == n {
				minRight = nums1[i]
			} else {
				minRight = int(math.Min(float64(nums2[j]), float64(nums1[i])))
			}

			return float64(maxLeft+minRight) / 2.0
		}
	}
	return 0
}

// 単純に考えると2つの配列を結合して中央値をとればいい。
// [1] 合成配列が偶数の場合
//

func FindMedianSortedArrays_2(nums1 []int, nums2 []int) float64 {
	A := nums1
	i := rand.Intn(len(A) + 1)
	leftA := A[:i]
	rightA := A[i:]
	fmt.Println(leftA)
	fmt.Println(rightA)
	m := len(nums1)
	n := len(nums2)

	if m > n { // to ensure m <= n
		nums1, nums2 = nums2, nums1
		m, n = n, m
	}

	iMin := 0
	iMax := m
	halfLen := (m + n + 1) / 2
	for iMin <= iMax {
		i := (iMin + iMax) / 2
		j := halfLen - i

		if i < iMax && nums2[j-1] > nums1[i] {
			iMin = i + 1 // i is too small
		} else if i > iMin && nums1[i-1] > nums2[j] {
			iMax = i - 1 // i is too big
		} else { // i is perfect
			maxLeft := 0
			if i == 0 {
				maxLeft = nums2[j-1]
			} else if j == 0 {
				maxLeft = nums1[i-1]
			} else {
				maxLeft = int(math.Max(float64(nums1[i-1]), float64(nums2[j-1])))
			}

			if (m+n)%2 == 1 {
				return float64(maxLeft)
			}

			minRight := 0
			if i == m {
				minRight = nums2[j]
			} else if j == n {
				minRight = nums1[i]
			} else {
				minRight = int(math.Min(float64(nums2[j]), float64(nums1[i])))
			}

			return float64(maxLeft+minRight) / 2.0
		}
	}
	return 0
}

func findMedianSortedArrays_3(nums1 []int, nums2 []int) float64 {
	m := len(nums1)
	n := len(nums2)

	KthNumber := func(k int, nums1, nums2 []int) int {
		// k番目に小さい数字を見つける条件は
		// min(nums1[i], nums2[j])
		// かつ i+j=k+1
		// k番目に小さい数字を見つけるにはk個の数があれば良いが、
		// 2つの配列それぞれに少なくとも1つの数字が必要なため+1
		// j = k+1-i => 0
		// i >= 0

		// 求める条件は
		// 0 <= i < len(nums1)
		// 0 <= j < len(nums2)
		// (i + 1) + (j + 1) = k + 1 <=> i + j = k - 1 <=> j = k - i - 1
		// のとき
		// min(nums1[i], nums2[j])

		m := len(nums1)
		n := len(nums2)
		// nums1のサイズを必ず大きくする
		if m < n {
			nums1, nums2 = nums2, nums1
			m, n = n, m
		}

		i := k / 2
		j := k - i - 1
		if i == 0 {

		} else if j == 0 {

		} else {

		}

		return 0
	}

	// ビット演算で奇数の場合
	if (m+n)&1 == 1 {
		// 奇数の場合は真ん中( (m+n+1)/2 )の値
		mid := KthNumber((m+n+1)/2, nums1, nums2)
		return float64(mid)
	}
	// 偶数の場合は中央二つの値の( (m+n)/2, (m+n)/2+1 )の平均値
	left := KthNumber((m+n)/2, nums1, nums2)
	right := KthNumber(((m+n)/2 + 1), nums1, nums2)
	return float64(left+right) / 2.0
}
