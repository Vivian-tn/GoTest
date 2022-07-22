package main

import (
	"fmt"
)

type ListNode struct {
	Val  int
	Next *ListNode
}

func main() {
	list1 := &ListNode{
		Val: 1,
		Next: &ListNode{
			Val: 2,
			Next: &ListNode{
				Val: 4,
				// Next:&ListNode{},
			},
		},
	}
	list2 := &ListNode{
		Val: 1,
		Next: &ListNode{
			Val: 3,
			Next: &ListNode{
				Val: 4,
				// Next:&ListNode{},
			},
		},
	}
	list := mergeTwoLists(list1, list2)
	fmt.Print(*list)

}
func mergeTwoLists(list1 *ListNode, list2 *ListNode) *ListNode {
	if list1.Next == nil {
		return list2
	} else if list2.Next == nil {
		return list1
	} else if list1.Val <= list2.Val {
		list1 = mergeTwoLists(list1.Next, list2)
		return list1
	} else {
		list2.Next = mergeTwoLists(list1, list2.Next)
		return list2
	}
}
