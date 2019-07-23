package main

import (
	"fmt"
)

func main() {
	// fmt.Println(leetcode.LengthOfLongestSubstring("abcabcbb"))
	// fmt.Println(leetcode.FindMedianSortedArrays_2([]int{1, 2, 3, 4, 5, 6, 7}, []int{1}))
	maxInt := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	minInt := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}

	findMedianSortedArrays := func(nums1 []int, nums2 []int) float64 {
		// m := len(nums1)
		// n := len(nums2)

		// if m > n { // to ensure m <= n
		// 	nums1, nums2 = nums2, nums1
		// 	m, n = n, m
		// }

		// iMin := 0
		// iMax := m
		// halfLen := (m + n + 1) / 2
		// for iMin <= iMax {
		// 	i := (iMin + iMax) / 2
		// 	j := halfLen - i

		// 	if i < iMax && nums2[j-1] > nums1[i] {
		// 		iMin = i + 1 // i is too small
		// 	} else if i > iMin && nums1[i-1] > nums2[j] {
		// 		iMax = i - 1 // i is too big
		// 	} else { // i is perfect
		// 		maxLeft := 0
		// 		if i == 0 {
		// 			maxLeft = nums2[j-1]
		// 		} else if j == 0 {
		// 			maxLeft = nums1[i-1]
		// 		} else {
		// 			maxLeft = int(math.Max(float64(nums1[i-1]), float64(nums2[j-1])))
		// 		}

		// 		if (m+n)%2 == 1 {
		// 			return float64(maxLeft)
		// 		}

		// 		minRight := 0
		// 		if i == m {
		// 			minRight = nums2[j]
		// 		} else if j == n {
		// 			minRight = nums1[i]
		// 		} else {
		// 			minRight = int(math.Min(float64(nums2[j]), float64(nums1[i])))
		// 		}

		// 		return float64(maxLeft+minRight) / 2.0
		// 	}
		// }
		// return 0
		A := nums1
		B := nums2
		m := len(A)
		n := len(B)

		if m > n {
			A, B = B, A
			m, n = n, m
		}

		iMin, iMax := 0, m
		halfLen := (m + n + 1) / 2

		for iMin <= iMax {
			i := (iMin + iMax) / 2
			j := halfLen - i

			switch {
			case i < iMax && B[j-1] > A[i]:
				iMin = i + 1
			case i > iMin && A[i-1] > B[j]:
				iMax = i - 1
			default:
				maxLeft := 0
				switch {
				case i == 0:
					maxLeft = B[j-1]
				case j == 0:
					maxLeft = A[i-1]
				default:
					maxLeft = maxInt(A[i-1], B[j-1])
				}

				if (m+n)%2 == 1 {
					return float64(maxLeft)
				}

				minRight := 0
				switch {
				case i == m:
					minRight = B[j]
				case j == n:
					minRight = A[i]
				default:
					minRight = minInt(B[j], A[i])
				}

				return float64(maxLeft+minRight) / 2
			}
		}

		return 0
	}
	A := []int{2, 3, 4, 5, 6, 7, 9}
	B := []int{1}
	fmt.Println(A)
	fmt.Println(B)
	fmt.Println(findMedianSortedArrays(A, B))
}
