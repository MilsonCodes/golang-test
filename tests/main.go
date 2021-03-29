// A concurrent prime sieve
// source: https://play.golang.org/p/9U22NfrXeq
package main

import (
	"fmt"
)

var target = 1000

// Send the sequence 2, 3, 4, ... to channel 'ch'.
func generate(ch chan<- int) {
	for i := 2; ; i++ {
		ch <- i // Send 'i' to channel 'ch'.
	}
}

// Copy the values from channel 'in' to channel 'out',
// removing those divisible by 'prime'.
func filter(in <-chan int, out chan<- int, prime int) {
	for {
		i := <-in // Receive value from 'in'.
		if i%prime != 0 {
			out <- i // Send 'i' to 'out'.
		}
	}
}

// The prime sieve: Daisy-chain Filter processes.
func main() {
	var prime int
	ch := make(chan int) // Create a new channel.
	go generate(ch)      // Launch Generate goroutine.
	for i := 0; i < target; i++ {
		prime = <-ch
		print(prime, "\n")
		ch1 := make(chan int)
		go filter(ch, ch1, prime) // Filter ch into ch1
		ch = ch1
	}
	fmt.Printf("%dth prime number: %d\n", target, prime)
}
