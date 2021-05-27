package main

import (
	"fmt"
)

func Chann(ch chan int) {
	for j := 0; j < 10; j++ {
		ch <- j
		//time.Sleep(time.Second)
	}
}
func main() {
	var  c2, c3 chan int
	c1 := make(chan int, 10)
	var i1, i2 int
	go Chann(c1)
	for {
		select {
		case i1 = <-c1:
			fmt.Printf("received ", i1, " from c1\n")
		case c2 <- i2:
			fmt.Printf("sent ", i2, " to c2\n")
		case i3, ok := (<-c3): // same as: i3, ok := <-c3
			if ok {
				fmt.Printf("received ", i3, " from c3\n")
			} else {
				fmt.Printf("c3 is closed\n")
			}
		//default:
		//	fmt.Printf("no communication\n")
		}
		//time.Sleep(2*time.Second)
	}
}