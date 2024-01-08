package listnode

type ListNode struct {
	Val  int
	Next *ListNode
}

// 两数相加
//
// 例如：输入：(2 -> 4 -> 3) + (5 -> 6 -> 4) 输出：7 -> 0 -> 8
//
// 参数：
//   - l1: 链表1
//   - l2: 链表2
func AddTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	newNode := l1
	temL2 := l2
	count := 0
	for i := 1; i <= 100; i++ {
		var val = count + newNode.Val + temL2.Val
		if val >= 10 {
			count = 1
			val = val - 10
		} else {
			count = 0
		}
		newNode.Val = val
		if newNode.Next == nil {
			if temL2.Next == nil {
				if count == 1 {
					newNode.Next = &ListNode{1, nil}
				}
				break
			} else {
				newNode.Next = &ListNode{0, nil}
			}
		}
		if temL2.Next == nil {
			temL2.Next = &ListNode{0, nil}
		}
		newNode = newNode.Next
		temL2 = temL2.Next
	}

	return l1
}
