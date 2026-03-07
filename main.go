package main

import (
	"fmt"
	"github.com/raimialiu/gostream/stream"
)

type Person struct {
	ID     int
	Name   string
	Age    int
	Email  string
	Salary float64
}

func main() {
	people := []Person{
		{1, "Alice", 30, "alice@example.com", 75000},
		{2, "Bob", 25, "bob@example.com", 65000},
		{3, "Charlie", 35, "charlie@example.com", 85000},
		{4, "Charlie", 35, "charlie@example.com", 85000},
	}

	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 2, 3, 7}
	nums := []int{3, 4, 5}

	// CountBy
	count := stream.From(people).
		CountBy(func(person Person) bool {
			return len(person.Name) > 4
		})
	fmt.Println("People with name > 4 chars:", count)

	// Concat + Distinct
	goStream1 := stream.From(numbers)
	goStream2 := stream.From(nums).
		Concat(*goStream1).
		Distinct().
		ToList()
	fmt.Println("Concat+Distinct:", goStream2)

	// Filter + ToList
	evens := stream.From(numbers).
		Filter(func(n int) bool { return n%2 == 0 }).
		ToList()
	fmt.Println("Evens:", evens)

	// Take + Skip
	taken := stream.From(numbers).Take(5).ToList()
	skipped := stream.From(numbers).Skip(8).ToList()
	fmt.Println("Take 5:", taken)
	fmt.Println("Skip 8:", skipped)

	// TakeWhile + SkipWhile
	tw := stream.From(numbers).TakeWhile(func(n int) bool { return n < 5 }).ToList()
	sw := stream.From(numbers).SkipWhile(func(n int) bool { return n < 5 }).ToList()
	fmt.Println("TakeWhile <5:", tw)
	fmt.Println("SkipWhile <5:", sw)

	// Range + Sum
	rangeSum := stream.IntRange(1, 101, 1).Sum()
	fmt.Println("Sum 1..100:", rangeSum)

	// First, Last
	first, _ := stream.From(people).First()
	last, _ := stream.From(people).Last()
	fmt.Println("First:", first.Name, "Last:", last.Name)

	// Reduce
	sum := stream.From(numbers).Reduce(0, func(a, b int) int { return a + b })
	fmt.Println("Reduce sum:", sum)

	// AnyMatch, AllMatch, NoneMatch
	fmt.Println("Any >5:", stream.From(numbers).AnyMatch(func(n int) bool { return n > 5 }))
	fmt.Println("All >0:", stream.From(numbers).AllMatch(func(n int) bool { return n > 0 }))
	fmt.Println("None <0:", stream.From(numbers).NoneMatch(func(n int) bool { return n < 0 }))

	// Reverse
	rev := stream.From([]int{1, 2, 3, 4, 5}).Reverse().ToList()
	fmt.Println("Reverse:", rev)

	// FlatMap
	flat := stream.From([]int{1, 2, 3}).FlatMap(func(n int) []int {
		return []int{n, n * 10}
	}).ToList()
	fmt.Println("FlatMap:", flat)

	// OrderBy
	sorted := stream.From([]int{5, 3, 1, 4, 2}).OrderBy(func(a, b int) bool {
		return a < b
	}).ToList()
	fmt.Println("OrderBy:", sorted)

	// Chunk
	chunks := stream.Chunk(stream.From(numbers), 4).ToList()
	fmt.Println("Chunks:", chunks)

	// Repeat
	repeated := stream.Repeat("hello", 3).ToList()
	fmt.Println("Repeat:", repeated)

	// Contains
	fmt.Println("Contains 7:", stream.From(numbers).Contains(7))
	fmt.Println("Contains 99:", stream.From(numbers).Contains(99))

	// DistinctBy
	uniqueByName := stream.From(people).DistinctBy(func(p Person) interface{} {
		return p.Name
	}).ToList()
	fmt.Println("DistinctBy Name:")
	for _, p := range uniqueByName {
		fmt.Printf("  %s (age %d)\n", p.Name, p.Age)
	}

	// Peek
	fmt.Print("Peek: ")
	stream.From([]int{1, 2, 3}).
		Peek(func(n int) { fmt.Printf("[%d] ", n) }).
		ToList()
	fmt.Println()

	// Min, Max, Average
	min, _ := stream.From(numbers).Min()
	max, _ := stream.From(numbers).Max()
	avg := stream.From(numbers).Average()
	fmt.Printf("Min: %d, Max: %d, Average: %.2f\n", min, max, avg)

	// Partition
	evensP, oddsP := stream.From(numbers).Partition(func(n int) bool { return n%2 == 0 })
	fmt.Println("Partition evens:", evensP)
	fmt.Println("Partition odds:", oddsP)

	// GroupBy
	groups := stream.From(people).GroupBy(func(p Person) interface{} { return p.Age })
	fmt.Println("GroupBy Age:")
	for age, group := range groups {
		fmt.Printf("  Age %v: ", age)
		for _, p := range group {
			fmt.Printf("%s ", p.Name)
		}
		fmt.Println()
	}
}
