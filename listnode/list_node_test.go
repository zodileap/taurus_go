package listnode

import (
	"fmt"
	"testing"
)

func Test1(t *testing.T) {
	l1 := &ListNode{2, nil}
	l1.Next = &ListNode{4, nil}
	l1.Next.Next = &ListNode{3, nil}

	l2 := &ListNode{5, nil}
	l2.Next = &ListNode{6, nil}
	l2.Next.Next = &ListNode{4, nil}

	l3 := AddTwoNumbers(l1, l2)
	for node := l3; node != nil; node = node.Next {
		fmt.Print(node.Val)
	}

}

func Test2(t *testing.T) {
	l1 := &ListNode{9, nil}
	l1.Next = &ListNode{9, nil}
	l1.Next.Next = &ListNode{9, nil}
	l1.Next.Next.Next = &ListNode{9, nil}
	l1.Next.Next.Next.Next = &ListNode{9, nil}
	l1.Next.Next.Next.Next.Next = &ListNode{9, nil}
	l1.Next.Next.Next.Next.Next.Next = &ListNode{9, nil}

	l2 := &ListNode{9, nil}
	l2.Next = &ListNode{9, nil}
	l2.Next.Next = &ListNode{9, nil}
	l2.Next.Next.Next = &ListNode{9, nil}

	l3 := AddTwoNumbers(l1, l2)
	for node := l3; node != nil; node = node.Next {
		fmt.Print(node.Val)
	}
}

func Test3(t *testing.T) {
	l1 := &ListNode{9, nil}
	l1.Next = &ListNode{9, nil}
	l1.Next.Next = &ListNode{9, nil}
	l1.Next.Next.Next = &ListNode{9, nil}
	l1.Next.Next.Next.Next = &ListNode{9, nil}
	l1.Next.Next.Next.Next.Next = &ListNode{9, nil}
	l1.Next.Next.Next.Next.Next.Next = &ListNode{9, nil}

	l2 := &ListNode{9, nil}
	l2.Next = &ListNode{9, nil}
	l2.Next.Next = &ListNode{9, nil}
	l2.Next.Next.Next = &ListNode{9, nil}

	l3 := AddTwoNumbers(l2, l1)
	for node := l3; node != nil; node = node.Next {
		fmt.Print(node.Val)
	}
}
