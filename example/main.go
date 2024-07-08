package main

import (
	"fmt"

	gd "github.com/Matej-Chmel/go-deep-copy"
)

type Example struct {
	flag     bool
	Node1    Node
	Node2    Node
	IntSlice []int
}

type Node struct {
	Next *Node
	val  int
}

func main() {
	node1 := Node{Next: nil, val: 1}
	node2 := Node{Next: nil, val: 2}
	node1.Next = &node2

	original := &Example{
		flag:     true,
		Node1:    node1,
		Node2:    node2,
		IntSlice: []int{3, 4, 5},
	}

	aCopy := gd.DeepCopy(original)

	fmt.Printf("Original at  %p\n%s\n\n", original, fmt.Sprintf("%+v", original))
	fmt.Printf("Deep copy at %p\n%s\n\n", aCopy, fmt.Sprintf("%+v", aCopy))

	if aCopy != original {
		fmt.Println("Copy and original live on different memory addresses")
	}

	if fmt.Sprintf("%+v", aCopy) == fmt.Sprintf("%+v", original) {
		fmt.Println("String representations match")
	}

	original.Node1.Next = nil
	original.IntSlice = append(original.IntSlice, 100)

	if aCopy.Node1.Next != nil && len(aCopy.IntSlice) == 3 {
		fmt.Println("\nCopy and original are detached,",
			"changing one doesn't affect the other")

		fmt.Printf("\nOriginal\n%s\n", fmt.Sprintf("%+v", original))
		fmt.Printf("Deep copy\n%s\n", fmt.Sprintf("%+v", aCopy))
	}
}
