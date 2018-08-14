package main

import (
	"fmt"
	"golang.org/x/tour/tree"
)

// Walk walks the tree t sending all values
// from the tree to the channel ch. 
func Walk(t *tree.Tree, ch chan int) {
	WalkAux(t,ch)
	close(ch) // Avoids panic!
}

// Recursive auxiliar functions
func WalkAux(t *tree.Tree, ch chan int) {
	if t == nil {
		return
	}
	
	// Left node
	if t.Left != nil {
		WalkAux(t.Left, ch)
	}

	ch <- t.Value
	
	// Right node	
	if t.Right != nil {
		WalkAux(t.Right, ch)
	}
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int)
	ch2 := make(chan int)
	
	go Walk(t1,ch1)
	go Walk(t2,ch2)
	
	// Compare channels
    for i := 0; i < 10; i++ {
        x, y := <-ch1, <-ch2
        if x != y {
            return false
        }
    }
    return true
}
	

func main() {
	ch := make(chan int)
	go Walk(tree.New(1),ch)
	
	for i := range ch {
		fmt.Println(i)
	}
	
	fmt.Println(Same(tree.New(1),tree.New(1))) // True
	fmt.Println(Same(tree.New(1),tree.New(10))) // False
}

