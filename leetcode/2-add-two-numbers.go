package leetcode

// You are given two non-empty linked lists representing two non-negative integers.
// The digits are stored in reverse order and each of their nodes contain a single digit.
// Add the two numbers and return it as a linked list.
// You may assume the two numbers do not contain any leading zero, except the number 0 itself.

// Example:

// Input: (2 -> 4 -> 3) + (5 -> 6 -> 4)
// Output: 7 -> 0 -> 8
// Explanation: 342 + 465 = 807.

/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */

type ListNode struct {
	Val  int
	Next *ListNode
}

// Runtime: 8 ms
// Memory Usage: 4.9 MB
// ポイント、先頭に番兵を置いて、返却時に先頭から2番目を返す。
func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	var node *ListNode
	sentinel := &ListNode{}
	curr := sentinel
	carry := 0
	for l1 != nil || l2 != nil {
		var n1, n2 int
		if l1 != nil {
			n1 = l1.Val
			l1 = l1.Next
		}
		if l2 != nil {
			n2 = l2.Val
			l2 = l2.Next
		}

		sum := n1 + n2 + carry
		carry = 0
		if sum >= 10 {
			sum = sum - 10
			carry = 1
		}
		node = &ListNode{Val: sum}
		curr.Next = node
		curr = node
	}
	if carry == 1 {
		node = &ListNode{Val: 1}
		curr.Next = node
	}

	return sentinel.Next
}

// Runtime: 8 ms
// Memory Usage: 4.9 MB
func addTwoNumbers_2(l1 *ListNode, l2 *ListNode) *ListNode {
	var node *ListNode
	head := &ListNode{}
	prevNode := head
	carry := 0
	next := func(node *ListNode) (currVal int, nextNode *ListNode) {
		if node == nil {
			return 0, nil
		}
		return node.Val, node.Next
	}
	exists := func(nodes ...*ListNode) bool {
		for _, node := range nodes {
			if node != nil {
				return true
			}
		}
		return false
	}

	for exists(l1, l2) {
		var n1, n2 int
		n1, l1 = next(l1)
		n2, l2 = next(l2)

		sum := n1 + n2 + carry
		carry = 0
		// sum, carry = sum/10, sum%10
		if sum >= 10 {
			sum = sum - 10
			carry = 1
		}
		node = &ListNode{Val: sum}
		prevNode.Next = node
		prevNode = node
	}
	if carry == 1 {
		node = &ListNode{Val: 1}
		prevNode.Next = node
	}

	return head.Next
}
