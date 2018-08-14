package main

import "fmt"

// fibonacci is a function that returns
// a function that returns an int.
func fibonacci() func() int {
	a := 0
	b := 0
	return func() int {
		b = a + b
		a = b - a
		
		if a == 0 && b == 0 {
			a = 1 // Just at the first iteration
		}
		
		return b
	}
}

func main() {
	f := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println(f())
	}
}

