package main

import (
	"context"
	"fmt"
	"gostream/stream"
)

type Person struct {
	ID     int
	Name   string
	Age    int
	Email  string
	Salary float64
}

type Employee struct {
	Person
	Department string
	Manager    string
}

func main() {

	people := []Person{
		{1, "Alice", 30, "alice@example.com", 75000},
		{2, "Bob", 25, "bob@example.com", 65000},
		{3, "Charlie", 35, "charlie@example.com", 85000},
		{4, "Charlie", 35, "charlie@example.com", 85000},
	}

	// numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 2, 3, 7}

	names := stream.From(people).
		CountBy(func(person Person) bool {
			return len(person.Name) > 4
		})

	fmt.Println("List:", names)

	/*
			ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
			for value := range repeatValue(ctx, 2) {
				fmt.Println("Repeat value:", value)
			}


		for value := range FilterFunc(MultiplyFunc(EnumerableRange(1, 30), 2), func(value int) bool {
			return value%2 == 0
		}) {
			fmt.Println(value)
		}

	*/
}

// write a function that generate numbers like enumerable.range in c#
// multiply it by a factor
// filter it by prime
// get the result

func MultiplyFunc(input <-chan int, factor int) <-chan int {
	values := make(chan int)
	go func() {
		defer close(values)
		for value := range input {
			values <- (value * factor)
		}
	}()

	return values
}

func FilterFunc(input <-chan int, predicate func(value int) bool) <-chan int {
	result := make(chan int)
	go func() {
		defer close(result)
		for value := range input {
			if predicate(value) {
				result <- value
			}
		}
	}()

	return result
}

func EnumerableRange(start, stop int) <-chan int {
	values := make(chan int)
	go func() {
		defer close(values)
		for i := start; i <= stop; i++ {
			values <- i
		}
	}()

	return values
}

func repeatValue(ctx context.Context, value int) <-chan int {
	valueStream := make(chan int)
	go func() {
		defer close(valueStream)
		for {
			select {
			case <-ctx.Done():
				return
			case valueStream <- value:
			}
		}
	}()

	return valueStream
}
