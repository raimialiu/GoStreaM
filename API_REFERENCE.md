# GoStream API Reference

A Java Streams / C# LINQ-style lazy stream processing library for Go with generics.

```
go get github.com/raimialiu/gostream
```

```go
import "github.com/raimialiu/gostream/stream"
```

---

## Table of Contents

- [Stream Creation](#stream-creation)
- [Intermediate Operations (Lazy)](#intermediate-operations-lazy)
- [Terminal Operations (Eager)](#terminal-operations-eager)
- [Collectors](#collectors)
- [Builders](#builders)
- [Parallel Streams](#parallel-streams)
- [Custom Extensions](#custom-extensions)

---

## Stream Creation

### `stream.From(slice []T) *GoStream[T]`
Creates a stream from a slice.

```go
s := stream.From([]int{1, 2, 3, 4, 5})
```

### `stream.Of(items ...T) *GoStream[T]`
Creates a stream from variadic arguments.

```go
s := stream.Of(1, 2, 3, 4, 5)
```

### `stream.Empty[T]() *GoStream[T]`
Creates an empty stream.

```go
s := stream.Empty[int]()
```

### `stream.IntRange(start, stop, step int) *GoStream[int]`
Creates a stream of integers from `start` (inclusive) to `stop` (exclusive) with the given step.

```go
s := stream.IntRange(0, 10, 2) // [0, 2, 4, 6, 8]
```

### `stream.Range[T Numeric](start, stop, step int) *GoStream[T]`
Generic range for any numeric type.

```go
s := stream.Range[int64](0, 100, 5)
```

### `stream.LongRange(start, stop, step int) *GoStream[int64]`
Shorthand for `Range[int64]`.

### `stream.DoubleRange(start, stop, step float32) *GoStream[float32]`
Range for float32 values.

### `stream.Repeat[T](value T, count int) *GoStream[T]`
Creates a stream that repeats `value` exactly `count` times.

```go
s := stream.Repeat("hello", 3) // ["hello", "hello", "hello"]
```

### `stream.RepeatForever[T](value T) *GoStream[T]`
Creates an infinite stream of the same value. **Must use `Take` to limit.**

```go
s := stream.RepeatForever(42).Take(5) // [42, 42, 42, 42, 42]
```

### `stream.Iterate[T](seed T, hasNext func(T) bool, next func(T) T) *GoStream[T]`
Creates a stream by iteratively applying `next` starting from `seed`, while `hasNext` returns true.

```go
// Powers of 2: [1, 2, 4, 8, 16]
s := stream.Iterate(1,
    func(n int) bool { return n <= 16 },
    func(n int) int { return n * 2 },
)
```

### `stream.Generate[T](generator func(args ...interface{}) T) *GoStream[T]`
Creates an infinite stream from a generator function.

### `stream.FromMap[K comparable, V any](m map[K]V) *GoStream[KeyValue[K, V]]`
Creates a stream of key-value pairs from a map.

```go
m := map[string]int{"a": 1, "b": 2}
s := stream.FromMap(m) // stream of KeyValue[string, int]
```

### `stream.FromMapValues[K comparable, V any](m map[K]V) *GoStream[V]`
Creates a stream of just the values from a map.

### `stream.FromChannel[T](ch <-chan T) *GoStream[T]`
Creates a stream from a channel.

```go
ch := make(chan int, 3)
ch <- 1; ch <- 2; ch <- 3
close(ch)
s := stream.FromChannel(ch)
```

### `stream.FromFile(path string) *GoStream[string]`
Creates a stream where each element is a line from the file.

### `stream.FromReader[T](r io.Reader, bufferSize int) *GoStream[string]`
Creates a stream that reads chunks from an `io.Reader`.

### `stream.FromIterator[T](iter iterators.Iterator[T]) *GoStream[T]`
Creates a stream from a custom `Iterator` implementation. See [Custom Extensions](#custom-extensions).

### `stream.Concat[T](streams ...*GoStream[T]) *GoStream[T]`
Concatenates multiple streams into one (standalone function).

```go
s := stream.Concat(stream.Of(1, 2), stream.Of(3, 4)) // [1, 2, 3, 4]
```

### `stream.Chunk[T](s *GoStream[T], size int) *GoStream[[]T]`
Splits a stream into chunks of the given size.

```go
chunks := stream.Chunk(stream.Of(1,2,3,4,5), 2).ToList()
// [[1,2], [3,4], [5]]
```

### `stream.Window[T](s *GoStream[T], size int) *GoStream[[]T]`
Creates sliding windows of the given size.

```go
windows := stream.Window(stream.Of(1,2,3,4,5), 3).ToList()
// [[1,2,3], [2,3,4], [3,4,5]]
```

---

## Intermediate Operations (Lazy)

These operations return a new `*GoStream[T]` and are not evaluated until a terminal operation is called.

### `Filter(predicate func(T) bool) *GoStream[T]`
Keeps only elements matching the predicate.

```go
evens := stream.From([]int{1,2,3,4,5,6}).
    Filter(func(n int) bool { return n%2 == 0 }).
    ToList() // [2, 4, 6]
```

### `Map(mapper func(T) interface{}) *GoStream[T]`
Transforms each element. The mapper returns `interface{}` which is type-asserted back to `T`.

```go
doubled := stream.From([]int{1,2,3}).
    Map(func(n int) interface{} { return n * 2 }).
    ToList() // [2, 4, 6]
```

### `FlatMap(mapper func(T) []T) *GoStream[T]`
Transforms each element into a slice and flattens all results into a single stream.

```go
result := stream.From([]int{1, 2, 3}).
    FlatMap(func(n int) []int { return []int{n, n * 10} }).
    ToList() // [1, 10, 2, 20, 3, 30]
```

### `Distinct() *GoStream[T]`
Removes duplicate elements (based on value equality).

```go
result := stream.From([]int{1,2,2,3,3,3}).Distinct().ToList() // [1, 2, 3]
```

### `DistinctBy(keySelector func(T) interface{}) *GoStream[T]`
Removes duplicates based on a key extracted from each element.

```go
type Person struct { Name string; Age int }
people := []Person{{"Alice", 30}, {"Bob", 25}, {"Alice", 35}}
unique := stream.From(people).
    DistinctBy(func(p Person) interface{} { return p.Name }).
    ToList() // [Alice(30), Bob(25)]
```

### `Take(count int) *GoStream[T]`
Takes at most `count` elements from the beginning.

```go
result := stream.From([]int{1,2,3,4,5}).Take(3).ToList() // [1, 2, 3]
```

### `Skip(count int) *GoStream[T]`
Skips the first `count` elements.

```go
result := stream.From([]int{1,2,3,4,5}).Skip(2).ToList() // [3, 4, 5]
```

### `TakeWhile(predicate func(T) bool) *GoStream[T]`
Takes elements while the predicate is true, stops at the first false.

```go
result := stream.From([]int{1,2,3,4,5}).
    TakeWhile(func(n int) bool { return n < 4 }).
    ToList() // [1, 2, 3]
```

### `SkipWhile(predicate func(T) bool) *GoStream[T]`
Skips elements while the predicate is true, then takes the rest.

```go
result := stream.From([]int{1,2,3,4,5,1}).
    SkipWhile(func(n int) bool { return n < 4 }).
    ToList() // [4, 5, 1]
```

### `Peek(action func(T)) *GoStream[T]`
Performs a side-effect action on each element without modifying the stream. Useful for debugging.

```go
stream.From([]int{1,2,3}).
    Peek(func(n int) { fmt.Println("Processing:", n) }).
    Filter(func(n int) bool { return n > 1 }).
    ToList()
```

### `Reverse() *GoStream[T]`
Returns a stream with elements in reverse order. **Materializes the stream.**

```go
result := stream.From([]int{1,2,3,4,5}).Reverse().ToList() // [5, 4, 3, 2, 1]
```

### `Concat(others ...GoStream[T]) *GoStream[T]`
Concatenates other streams onto this one (method form).

```go
a := stream.From([]int{1, 2})
b := stream.From([]int{3, 4})
result := a.Concat(*b).ToList() // [1, 2, 3, 4]
```

### `Union(other GoStream[T]) *GoStream[T]`
Returns the set union (concat + distinct).

```go
a := stream.From([]int{1, 2, 3})
b := *stream.From([]int{2, 3, 4, 5})
result := a.Union(b).ToList() // [1, 2, 3, 4, 5]
```

### `Intersect(other GoStream[T]) *GoStream[T]`
Returns elements present in both streams.

```go
a := stream.From([]int{1, 2, 3, 4})
b := *stream.From([]int{3, 4, 5, 6})
result := a.Intersect(b).ToList() // [3, 4]
```

### `Except(other GoStream[T]) *GoStream[T]`
Returns elements in this stream but not in the other.

```go
a := stream.From([]int{1, 2, 3, 4})
b := *stream.From([]int{3, 4, 5, 6})
result := a.Except(b).ToList() // [1, 2]
```

### `Cache() *GoStream[T]`
Materializes the stream into memory for reuse. **Materializes the stream.**

### `Zip(other *GoStream[T], zipper func(T, T) T) *GoStream[T]`
Combines elements from two streams pairwise. Stops at the shorter stream.

```go
a := stream.From([]int{1, 2, 3})
b := stream.From([]int{10, 20, 30})
result := a.Zip(b, func(x, y int) int { return x + y }).ToList() // [11, 22, 33]
```

### `OrderBy(less func(a, b T) bool) *GoStream[T]`
Sorts elements in ascending order. **Materializes the stream.**

```go
result := stream.From([]int{5,3,1,4,2}).
    OrderBy(func(a, b int) bool { return a < b }).
    ToList() // [1, 2, 3, 4, 5]
```

### `OrderByDesc(less func(a, b T) bool) *GoStream[T]`
Sorts elements in descending order. **Materializes the stream.**

```go
result := stream.From([]int{5,3,1,4,2}).
    OrderByDesc(func(a, b int) bool { return a < b }).
    ToList() // [5, 4, 3, 2, 1]
```

### `Select(selector func(T) interface{}) *GoStream[T]`
Alias for `Map`. LINQ-style naming.

---

## Terminal Operations (Eager)

These operations trigger evaluation and return a result.

### `ToList() []T` / `ToSlice() []T`
Collects all elements into a slice.

```go
result := stream.Of(1, 2, 3).ToList() // [1, 2, 3]
```

### `ForEach(action func(T))`
Executes an action on each element.

```go
stream.Of(1, 2, 3).ForEach(func(n int) {
    fmt.Println(n)
})
```

### `Count() int`
Returns the number of elements.

### `CountBy(predicate func(T) bool) int`
Returns the number of elements matching the predicate.

### `Sum() float64`
Returns the sum of all numeric elements (uses reflection).

```go
sum := stream.Of(1, 2, 3, 4, 5).Sum() // 15.0
```

### `Average() float64`
Returns the average of all numeric elements.

### `Min() (T, error)` / `Max() (T, error)`
Returns the min/max element. Returns `ErrEmptyStream` if empty.

```go
min, err := stream.Of(3, 1, 4, 1, 5).Min() // 1, nil
max, err := stream.Of(3, 1, 4, 1, 5).Max() // 5, nil
```

### `First() (T, error)` / `Last() (T, error)`
Returns the first/last element. Returns `ErrEmptyStream` if empty.

### `FirstOrDefault(defaultValue T) T`
Returns the first element, or `defaultValue` if empty.

### `Single() (T, error)`
Returns the single element. Returns `ErrEmptyStream` if empty, `ErrMultipleElements` if more than one.

### `Any() bool`
Returns true if the stream has at least one element.

### `AnyMatch(predicate func(T) bool) bool`
Returns true if any element matches the predicate.

### `AllMatch(predicate func(T) bool) bool`
Returns true if all elements match the predicate.

### `NoneMatch(predicate func(T) bool) bool`
Returns true if no element matches the predicate.

### `Reduce(identity T, accumulator func(T, T) T) T`
Reduces elements to a single value using the accumulator.

```go
sum := stream.Of(1, 2, 3, 4).Reduce(0, func(a, b int) int { return a + b }) // 10
product := stream.Of(1, 2, 3, 4).Reduce(1, func(a, b int) int { return a * b }) // 24
```

### `Contains(value T) bool`
Returns true if the stream contains the value (uses `reflect.DeepEqual`).

### `ToMap(keySelector func(T) interface{}) map[interface{}]T`
Collects elements into a map using the key selector.

```go
type User struct { ID int; Name string }
users := stream.Of(User{1, "Alice"}, User{2, "Bob"})
m := users.ToMap(func(u User) interface{} { return u.ID })
// map[1:Alice, 2:Bob]
```

### `GroupBy(keySelector func(T) interface{}) map[interface{}][]T`
Groups elements by key.

```go
groups := stream.Of(1,2,3,4,5,6).GroupBy(func(n int) interface{} { return n % 2 })
// map[0:[2,4,6], 1:[1,3,5]]
```

### `Partition(predicate func(T) bool) (matching []T, nonMatching []T)`
Splits elements into two slices based on the predicate.

```go
evens, odds := stream.Of(1,2,3,4,5).Partition(func(n int) bool { return n%2 == 0 })
// evens=[2,4], odds=[1,3,5]
```

### `OrderBy(less func(a, b T) bool) *GoStream[T]`
See [Intermediate Operations](#orderbyless-funca-b-t-bool-gostreamt).

### `Iterator() iterators.Iterator[T]`
Returns an iterator for manual iteration.

### `Close()`
Releases resources. Any subsequent operation will panic.

---

## Collectors

```go
import "github.com/raimialiu/gostream/stream/collectors"
```

Collectors provide a reusable, composable way to accumulate stream elements into results.

### Using Collectors with Streams

```go
result := stream.Collect(myStream, collectors.ToList[int]())
```

### Available Collectors

#### `collectors.ToList[T]() Collector`
Collects into a `[]T`.

#### `collectors.ToSet[T comparable]() Collector`
Collects into a `map[T]struct{}` (unique elements).

#### `collectors.ToMap[T, K comparable, V any](keySelector, valueSelector) Collector`
Collects into a `map[K]V`.

```go
type User struct { ID int; Name string }
result := stream.Collect(users,
    collectors.ToMap[User, int, string](
        func(u User) int { return u.ID },
        func(u User) string { return u.Name },
    ),
)
```

#### `collectors.GroupingBy[T, K comparable](keySelector) Collector`
Groups elements by key into `map[K][]T`.

```go
result := stream.Collect(numbers,
    collectors.GroupingBy[int, int](func(n int) int { return n % 2 }),
)
```

#### `collectors.PartitioningBy[T](predicate) Collector`
Partitions into `*Partition[T]` with `.Matching` and `.NonMatching` fields.

#### `collectors.Joining(separator string) Collector`
Joins strings with a separator.

```go
result := stream.Collect(names, collectors.Joining(", ")) // "Alice, Bob, Charlie"
```

#### `collectors.JoiningWithPrefixSuffix(separator, prefix, suffix string) Collector`
Joins with separator, prefix, and suffix.

```go
result := stream.Collect(names, collectors.JoiningWithPrefixSuffix(", ", "[", "]"))
// "[Alice, Bob, Charlie]"
```

#### `collectors.Counting[T]() Collector`
Counts elements.

#### `collectors.Summing[T](toFloat func(T) float64) Collector`
Sums elements using a conversion function.

```go
result := stream.Collect(people,
    collectors.Summing[Person](func(p Person) float64 { return p.Salary }),
)
```

#### `collectors.Averaging[T](toFloat func(T) float64) Collector`
Averages elements.

#### `collectors.MinBy[T](less func(a, b T) bool) Collector`
Finds the minimum element.

#### `collectors.MaxBy[T](less func(a, b T) bool) Collector`
Finds the maximum element.

#### `collectors.Reducing[T](identity T, accumulator func(T, T) T) Collector`
General-purpose reduction.

#### `collectors.Mapping[T, U, A, R](mapper, downstream) Collector`
Applies a mapper before collecting with the downstream collector.

```go
result := stream.Collect(people,
    collectors.Mapping(
        func(p Person) string { return p.Name },
        collectors.ToList[string](),
    ),
)
```

#### `collectors.Filtering[T, A, R](predicate, downstream) Collector`
Filters before collecting with the downstream collector.

```go
result := stream.Collect(numbers,
    collectors.Filtering(
        func(n int) bool { return n > 3 },
        collectors.ToList[int](),
    ),
)
```

#### `collectors.FlatMapping[T, U, A, R](mapper, downstream) Collector`
FlatMaps before collecting with the downstream collector.

#### `collectors.ToStringCollector[T](separator string, formatter func(T) string) Collector`
Converts each element to string and joins.

#### `collectors.FormatCollector[T](format string, separator string) Collector`
Uses `fmt.Sprintf` to format each element and joins.

### Creating Custom Collectors

#### Using `NewCollector`

```go
sumOfSquares := collectors.NewCollector[int, *float64, float64](
    func() *float64 { v := 0.0; return &v },                    // supplier
    func(acc *float64, n int) *float64 { *acc += float64(n*n); return acc }, // accumulator
    func(acc *float64) float64 { return *acc },                  // finisher
)
result := stream.Collect(myStream, sumOfSquares)
```

#### Implementing the Interface

```go
type Collector[T, A, R any] interface {
    Supplier() A
    Accumulator(accumulator A, element T) A
    Finisher(accumulator A) R
}
```

You can also use `collectors.Collect(slice, collector)` to apply collectors directly to slices without creating a stream.

---

## Builders

```go
import "github.com/raimialiu/gostream/stream/builders"
```

Builders provide a fluent way to construct streams with conditional logic.

### Creating Builders

```go
b := builders.NewBuilder[int]()
b := builders.NewBuilderWithCapacity[int](100)
b := builders.NewBuilderFrom(1, 2, 3)
b := builders.NewBuilderFromSlice(mySlice)
```

### Basic Addition

```go
b.Add(1)                    // single element
b.AddAll(1, 2, 3)           // multiple elements
b.AddSlice([]int{4, 5, 6})  // from slice
b.AddFirst(0)               // prepend
b.AddAt(2, 99)              // insert at index
```

### Conditional Addition

```go
isAdmin := true
b.AddIf(isAdmin, adminItem)
b.AddIfElse(isAdmin, adminItem, guestItem)
b.AddAllIf(hasData, items...)
b.AddSliceIf(hasData, items)
b.AddUnless(isEmpty, defaultItem)
```

### Functional Addition

```go
b.AddWith(func() int { return computeValue() })
b.AddWithIf(condition, func() int { return expensiveComputation() })
b.AddRange(0, 10, func(i int) int { return i * i })
```

### Stream Integration

```go
b.AddStream(existingStream)
b.Merge(otherBuilder)
```

### Conditional Logic Blocks

#### When/Unless

```go
result := builders.NewBuilder[string]().
    Add("always").
    When(userIsAdmin).
        Add("admin-panel").
        AddAll("settings", "users").
    End().
    Unless(isMobile).
        Add("desktop-widget").
    End().
    Build()
```

#### When/Else

```go
result := builders.NewBuilder[string]().
    When(isPremium).
        Add("premium-feature").
    ElseAdd("free-feature").
    End().
    Build()
```

#### Switch

```go
result := builders.NewBuilder[string]().
    Switch(userRole).
        Case("admin", "admin-dashboard", "user-management").
        Case("editor", "content-editor").
        CaseWith("viewer", func() string { return loadViewerConfig() }).
        Default("guest-page").
    End().
    Build()
```

### Inspection

```go
b.Count()          // number of elements
b.IsEmpty()        // true if empty
b.Contains(value)  // true if element exists
b.Last()           // (element, bool)
b.Get(index)       // (element, bool)
```

### Modification

```go
b.Clear()                                          // remove all
b.RemoveLast()                                     // remove last element
b.RemoveAt(index)                                  // remove at index
b.RemoveIf(func(n int) bool { return n < 0 })     // remove matching
```

### Terminal

```go
s := b.Build()          // *GoStream[T] (builder unchanged)
s := b.BuildAndClear()  // *GoStream[T] (builder cleared)
s := b.ToStream()       // alias for Build
items := b.ToSlice()    // []T directly
```

---

## Parallel Streams

```go
import "github.com/raimialiu/gostream/stream/parallel"
```

Parallel streams execute operations concurrently across multiple goroutines.

### Creating Parallel Streams

```go
ps := parallel.From(myGoStream)
ps := parallel.FromSlice([]int{1, 2, 3, 4, 5})
ps := parallel.FromSlice(items).WithWorkers(4)  // default: runtime.NumCPU()
```

### Operations

```go
// ForEach (parallel execution)
ps.ForEach(func(item Item) {
    process(item) // runs concurrently
})

// Filter (order-preserving)
result := ps.Filter(func(n int) bool { return n > 5 }).ToSlice()

// Map (order-preserving, type-changing)
results := parallel.Map(ps, func(n int) string { return fmt.Sprint(n) })

// AnyMatch / AllMatch (short-circuiting)
hasNegative := ps.AnyMatch(func(n int) bool { return n < 0 })
allPositive := ps.AllMatch(func(n int) bool { return n > 0 })

// Reduce (parallel reduction with combiner)
sum := ps.Reduce(0, func(a, b int) int { return a + b })

// Count / ToSlice
count := ps.Count()
items := ps.ToSlice()
```

### Parallel Collect

```go
result := parallel.Collect(ps,
    func(chunk []int) int {         // per-chunk collector
        sum := 0
        for _, v := range chunk { sum += v }
        return sum
    },
    func(a, b int) int { return a + b },  // combiner
    0,                                     // identity
)
```

### Converting Back

```go
sequential := ps.Sequential() // back to *GoStream[T]
```

**Note:** Each terminal call on a `ParallelStream` consumes the underlying stream. Create new instances for multiple operations.

---

## Custom Extensions

GoStream provides three levels of extensibility for custom behavior.

### Level 1: Transform (Reusable Pipeline Fragments)

Compose reusable stream transformations as functions:

```go
func TopN[T any](n int, less func(a, b T) bool) func(*stream.GoStream[T]) *stream.GoStream[T] {
    return func(s *stream.GoStream[T]) *stream.GoStream[T] {
        return s.OrderBy(less).Take(n)
    }
}

result := stream.From(items).Transform(TopN[int](5, func(a, b int) bool { return a < b })).ToList()
```

### Level 2: ApplyIterator (Custom Iterator Logic)

Wrap the pipeline with arbitrary iterator transformation logic:

```go
result := stream.From([]int{1, 2, 3}).
    ApplyIterator(func(iter iterators.Iterator[int]) iterators.Iterator[int] {
        var items []int
        for iter.HasNext() {
            items = append(items, iter.Next() * 2)
        }
        return iterators.AsListIterator(items...)
    }).
    ToList() // [2, 4, 6]
```

### Level 3: Custom StreamOperation

Implement the `StreamOperation` interface for fully lazy custom operations:

```go
import "github.com/raimialiu/gostream/stream/operations"

type DoubleOperation struct{}

func (d DoubleOperation) CanFuse(next operations.StreamOperation[int]) bool { return true }

func (d DoubleOperation) Apply(source iterators.Iterator[int]) iterators.Iterator[int] {
    return &doublingIterator{source: source}
}

// Custom iterator
type doublingIterator struct {
    source iterators.Iterator[int]
}

func (it *doublingIterator) HasNext() bool { return it.source.HasNext() }
func (it *doublingIterator) Next() int     { return it.source.Next() * 2 }
func (it *doublingIterator) Close() error  { return it.source.Close() }

// Usage
result := stream.From([]int{1, 2, 3}).Apply(DoubleOperation{}).ToList()
```

### Level 4: Custom Iterator + FromIterator

Create entirely custom data sources by implementing the `Iterator` interface:

```go
import "github.com/raimialiu/gostream/stream/iterators"

type Iterator[T any] interface {
    HasNext() bool
    Next() T
    Close() error
}

// Example: Fibonacci iterator
type fibIterator struct {
    a, b  int
    limit int
    count int
}

func (f *fibIterator) HasNext() bool { return f.count < f.limit }
func (f *fibIterator) Next() int {
    f.count++
    result := f.a
    f.a, f.b = f.b, f.a+f.b
    return result
}
func (f *fibIterator) Close() error { return nil }

// Usage
fib := stream.FromIterator[int](&fibIterator{a: 0, b: 1, limit: 10})
result := fib.ToList() // [0, 1, 1, 2, 3, 5, 8, 13, 21, 34]
```

### Level 5: Custom Collectors

Create reusable terminal operations. See [Creating Custom Collectors](#creating-custom-collectors).

---

## Errors

```go
stream.ErrEmptyStream      // returned by First, Last, Single, Min, Max on empty streams
stream.ErrMultipleElements // returned by Single when stream has >1 element
stream.ErrIndexOutOfRange  // index out of range
stream.ErrStreamClosed     // operation on closed stream
```

---

## Type Definitions

```go
// stream/types
type Numeric     interface { ~int | ~int8 | ... | ~float64 }
type Predicate[T]   func(T) bool
type Consumer[T]    func(T)
type Mapper[T, R]   func(T) R
type Comparator[T]  func(T) int
type BiFunction[T, U, R] func(T, U) R

// stream/delegates
type Predicate[T]       func(T) bool
type KeySelector[K]     func() K
```

---

## Chaining Examples

```go
// Complex pipeline
result := stream.From(employees).
    Filter(func(e Employee) bool { return e.Department == "Engineering" }).
    DistinctBy(func(e Employee) interface{} { return e.Email }).
    OrderBy(func(a, b Employee) bool { return a.Salary > b.Salary }).
    Take(10).
    Peek(func(e Employee) { log.Println("Top earner:", e.Name) }).
    ToList()

// Sum of squares of even numbers from 1 to 100
result := stream.IntRange(1, 101, 1).
    Filter(func(n int) bool { return n%2 == 0 }).
    Reduce(0, func(a, b int) int { return a + b*b })

// Build a stream conditionally
config := builders.NewBuilder[string]().
    Add("base-module").
    When(isPremium).
        AddAll("analytics", "reports").
    End().
    Switch(region).
        Case("US", "us-compliance").
        Case("EU", "gdpr-module").
        Default("standard").
    End().
    Build()

// Parallel processing
parallel.FromSlice(largeDataset).
    WithWorkers(8).
    Filter(func(r Record) bool { return r.IsValid() }).
    ForEach(func(r Record) { saveToDatabase(r) })

// Collect with custom collector
avgSalary := stream.Collect(
    stream.From(employees),
    collectors.Averaging[Employee](func(e Employee) float64 { return e.Salary }),
)
```
