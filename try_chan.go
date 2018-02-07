package main

import (
	"fmt"
)

func main() {
	ch  := make(chan int, 3)
	ch <- 1
	ch <- 2
	close(ch)
	e, ok := <-ch
	fmt.Println(e, ok)
	e, ok = <-ch
	fmt.Println(e, ok)
	e, ok = <-ch
	fmt.Println(e, ok)
	e, ok = <-ch
	fmt.Println(e, ok)
}

func testAAa(in chan int) {
	fmt.Println("in function")
	<- in
	fmt.Println("output")
}