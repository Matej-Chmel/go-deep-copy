# Deep copy
Create a deep copy of any variable.

## Installation
```bash
go get github.com/Matej-Chmel/go-deep-copy@v1.0.5
```

## Features
- Uses 2 stacks to copy values instead of recursion
- Copy built-in and composite types
- Copy both exported and unexported fields of a struct (thanks [brunoga](https://github.com/brunoga/deep/tree/92c699d4e2e304e7c4a4d4138817a8c96e8abb72))

## Example
```go
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
```

### Output
```none
Original at  0xc00009a040
&{flag:true Node1:{Next:0xc00008a030 val:1} Node2:{Next:<nil> val:2} IntSlice:[3 4 5]}

Deep copy at 0xc00009a080
&{flag:true Node1:{Next:0xc00008a070 val:1} Node2:{Next:<nil> val:2} IntSlice:[3 4 5]}

Copy and original live on different memory addresses

Copy and original are detached, changing one doesn't affect the other

Original
&{flag:true Node1:{Next:<nil> val:1} Node2:{Next:<nil> val:2} IntSlice:[3 4 5 100]}
Deep copy
&{flag:true Node1:{Next:0xc00008a070 val:1} Node2:{Next:<nil> val:2} IntSlice:[3 4 5]}
```
