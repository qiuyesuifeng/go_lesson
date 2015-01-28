package main

import (
	"fmt"
	"sync"
)

func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func merge(cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	output := func(c <-chan int) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func demo_case_one() {
	c := gen(2, 3)
	out := square(c)

	fmt.Println("---------demo_case_one---------")
	fmt.Println(<-out) // 4
	fmt.Println(<-out) // 9
}

func demo_case_two() {
	fmt.Println("---------demo_case_two---------")

	for n := range square(square(gen(2, 3))) {
		fmt.Println(n) // 16 then 81
	}
}

func demo_case_three() {
	in := gen(2, 3)

	c1 := square(in)
	c2 := square(in)

	fmt.Println("---------demo_case_three---------")

	for n := range merge(c1, c2) {
		fmt.Println(n) // 4 then 9, or 9 then 4
	}
}

func main() {
	demo_case_one()
	demo_case_two()
	demo_case_three()
}
