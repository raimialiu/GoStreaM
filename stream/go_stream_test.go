package stream

import (
	"github.com/raimialiu/gostream/stream/collectors"
	"github.com/raimialiu/gostream/stream/iterators"
	"testing"
)

func TestFrom(t *testing.T) {
	s := From([]int{1, 2, 3})
	result := s.ToList()
	assertSliceEqual(t, []int{1, 2, 3}, result)
}

func TestOf(t *testing.T) {
	result := Of(1, 2, 3).ToList()
	assertSliceEqual(t, []int{1, 2, 3}, result)
}

func TestEmpty(t *testing.T) {
	result := Empty[int]().ToList()
	if len(result) != 0 {
		t.Errorf("expected empty slice, got %v", result)
	}
}

func TestFilter(t *testing.T) {
	result := From([]int{1, 2, 3, 4, 5, 6}).
		Filter(func(n int) bool { return n%2 == 0 }).
		ToList()
	assertSliceEqual(t, []int{2, 4, 6}, result)
}

func TestFilterEmpty(t *testing.T) {
	result := From([]int{1, 3, 5}).
		Filter(func(n int) bool { return n%2 == 0 }).
		ToList()
	if len(result) != 0 {
		t.Errorf("expected empty slice, got %v", result)
	}
}

func TestDistinct(t *testing.T) {
	result := From([]int{1, 2, 2, 3, 3, 3, 4}).Distinct().ToList()
	assertSliceEqual(t, []int{1, 2, 3, 4}, result)
}

func TestDistinctBy(t *testing.T) {
	type item struct {
		name string
		val  int
	}
	items := []item{{"a", 1}, {"b", 2}, {"a", 3}}
	result := From(items).DistinctBy(func(i item) interface{} { return i.name }).ToList()
	if len(result) != 2 {
		t.Errorf("expected 2 items, got %d", len(result))
	}
}

func TestTake(t *testing.T) {
	result := From([]int{1, 2, 3, 4, 5}).Take(3).ToList()
	assertSliceEqual(t, []int{1, 2, 3}, result)
}

func TestTakeMoreThanAvailable(t *testing.T) {
	result := From([]int{1, 2}).Take(5).ToList()
	assertSliceEqual(t, []int{1, 2}, result)
}

func TestSkip(t *testing.T) {
	result := From([]int{1, 2, 3, 4, 5}).Skip(2).ToList()
	assertSliceEqual(t, []int{3, 4, 5}, result)
}

func TestSkipMoreThanAvailable(t *testing.T) {
	result := From([]int{1, 2}).Skip(5).ToList()
	if len(result) != 0 {
		t.Errorf("expected empty, got %v", result)
	}
}

func TestTakeWhile(t *testing.T) {
	result := From([]int{1, 2, 3, 4, 5, 1}).
		TakeWhile(func(n int) bool { return n < 4 }).
		ToList()
	assertSliceEqual(t, []int{1, 2, 3}, result)
}

func TestSkipWhile(t *testing.T) {
	result := From([]int{1, 2, 3, 4, 5, 1}).
		SkipWhile(func(n int) bool { return n < 4 }).
		ToList()
	assertSliceEqual(t, []int{4, 5, 1}, result)
}

func TestReverse(t *testing.T) {
	result := From([]int{1, 2, 3, 4, 5}).Reverse().ToList()
	assertSliceEqual(t, []int{5, 4, 3, 2, 1}, result)
}

func TestConcat(t *testing.T) {
	a := From([]int{1, 2})
	b := From([]int{3, 4})
	result := a.Concat(*b).ToList()
	assertSliceEqual(t, []int{1, 2, 3, 4}, result)
}

func TestConcatDistinct(t *testing.T) {
	a := From([]int{1, 2, 3})
	b := From([]int{2, 3, 4})
	result := a.Concat(*b).Distinct().ToList()
	assertSliceEqual(t, []int{1, 2, 3, 4}, result)
}

func TestFlatMap(t *testing.T) {
	result := From([]int{1, 2, 3}).
		FlatMap(func(n int) []int { return []int{n, n * 10} }).
		ToList()
	assertSliceEqual(t, []int{1, 10, 2, 20, 3, 30}, result)
}

func TestPeek(t *testing.T) {
	var peeked []int
	result := From([]int{1, 2, 3}).
		Peek(func(n int) { peeked = append(peeked, n) }).
		ToList()
	assertSliceEqual(t, []int{1, 2, 3}, result)
	assertSliceEqual(t, []int{1, 2, 3}, peeked)
}

func TestUnion(t *testing.T) {
	a := From([]int{1, 2, 3})
	b := *From([]int{2, 3, 4, 5})
	result := a.Union(b).ToList()
	assertSliceEqual(t, []int{1, 2, 3, 4, 5}, result)
}

func TestIntersect(t *testing.T) {
	a := From([]int{1, 2, 3, 4})
	b := *From([]int{3, 4, 5, 6})
	result := a.Intersect(b).ToList()
	assertSliceEqual(t, []int{3, 4}, result)
}

func TestExcept(t *testing.T) {
	a := From([]int{1, 2, 3, 4})
	b := *From([]int{3, 4, 5, 6})
	result := a.Except(b).ToList()
	assertSliceEqual(t, []int{1, 2}, result)
}

// ========== Terminal Operations ==========

func TestCount(t *testing.T) {
	count := From([]int{1, 2, 3, 4, 5}).Count()
	if count != 5 {
		t.Errorf("expected 5, got %d", count)
	}
}

func TestCountBy(t *testing.T) {
	count := From([]int{1, 2, 3, 4, 5}).CountBy(func(n int) bool { return n > 3 })
	if count != 2 {
		t.Errorf("expected 2, got %d", count)
	}
}

func TestSum(t *testing.T) {
	sum := From([]int{1, 2, 3, 4, 5}).Sum()
	if sum != 15 {
		t.Errorf("expected 15, got %f", sum)
	}
}

func TestAverage(t *testing.T) {
	avg := From([]int{2, 4, 6}).Average()
	if avg != 4.0 {
		t.Errorf("expected 4.0, got %f", avg)
	}
}

func TestAverageEmpty(t *testing.T) {
	avg := Empty[int]().Average()
	if avg != 0 {
		t.Errorf("expected 0, got %f", avg)
	}
}

func TestMin(t *testing.T) {
	min, err := From([]int{3, 1, 4, 1, 5}).Min()
	if err != nil || min != 1 {
		t.Errorf("expected 1, got %d (err: %v)", min, err)
	}
}

func TestMinEmpty(t *testing.T) {
	_, err := Empty[int]().Min()
	if err != ErrEmptyStream {
		t.Errorf("expected ErrEmptyStream, got %v", err)
	}
}

func TestMax(t *testing.T) {
	max, err := From([]int{3, 1, 4, 1, 5}).Max()
	if err != nil || max != 5 {
		t.Errorf("expected 5, got %d (err: %v)", max, err)
	}
}

func TestFirst(t *testing.T) {
	first, err := From([]int{10, 20, 30}).First()
	if err != nil || first != 10 {
		t.Errorf("expected 10, got %d (err: %v)", first, err)
	}
}

func TestFirstEmpty(t *testing.T) {
	_, err := Empty[int]().First()
	if err != ErrEmptyStream {
		t.Errorf("expected ErrEmptyStream, got %v", err)
	}
}

func TestFirstOrDefault(t *testing.T) {
	val := Empty[int]().FirstOrDefault(42)
	if val != 42 {
		t.Errorf("expected 42, got %d", val)
	}
}

func TestLast(t *testing.T) {
	last, err := From([]int{10, 20, 30}).Last()
	if err != nil || last != 30 {
		t.Errorf("expected 30, got %d (err: %v)", last, err)
	}
}

func TestLastEmpty(t *testing.T) {
	_, err := Empty[int]().Last()
	if err != ErrEmptyStream {
		t.Errorf("expected ErrEmptyStream, got %v", err)
	}
}

func TestSingle(t *testing.T) {
	val, err := Of(42).Single()
	if err != nil || val != 42 {
		t.Errorf("expected 42, got %d (err: %v)", val, err)
	}
}

func TestSingleEmpty(t *testing.T) {
	_, err := Empty[int]().Single()
	if err != ErrEmptyStream {
		t.Errorf("expected ErrEmptyStream, got %v", err)
	}
}

func TestSingleMultiple(t *testing.T) {
	_, err := Of(1, 2).Single()
	if err != ErrMultipleElements {
		t.Errorf("expected ErrMultipleElements, got %v", err)
	}
}

func TestReduce(t *testing.T) {
	result := From([]int{1, 2, 3, 4}).Reduce(0, func(a, b int) int { return a + b })
	if result != 10 {
		t.Errorf("expected 10, got %d", result)
	}
}

func TestReduceProduct(t *testing.T) {
	result := From([]int{1, 2, 3, 4}).Reduce(1, func(a, b int) int { return a * b })
	if result != 24 {
		t.Errorf("expected 24, got %d", result)
	}
}

func TestContains(t *testing.T) {
	if !From([]int{1, 2, 3}).Contains(2) {
		t.Error("expected true")
	}
	if From([]int{1, 2, 3}).Contains(4) {
		t.Error("expected false")
	}
}

func TestAny(t *testing.T) {
	if !From([]int{1}).Any() {
		t.Error("expected true")
	}
	if Empty[int]().Any() {
		t.Error("expected false")
	}
}

func TestAnyMatch(t *testing.T) {
	if !From([]int{1, 2, 3}).AnyMatch(func(n int) bool { return n == 2 }) {
		t.Error("expected true")
	}
	if From([]int{1, 2, 3}).AnyMatch(func(n int) bool { return n == 5 }) {
		t.Error("expected false")
	}
}

func TestAllMatch(t *testing.T) {
	if !From([]int{2, 4, 6}).AllMatch(func(n int) bool { return n%2 == 0 }) {
		t.Error("expected true")
	}
	if From([]int{2, 3, 6}).AllMatch(func(n int) bool { return n%2 == 0 }) {
		t.Error("expected false")
	}
}

func TestNoneMatch(t *testing.T) {
	if !From([]int{1, 3, 5}).NoneMatch(func(n int) bool { return n%2 == 0 }) {
		t.Error("expected true")
	}
	if From([]int{1, 2, 5}).NoneMatch(func(n int) bool { return n%2 == 0 }) {
		t.Error("expected false")
	}
}

func TestForEach(t *testing.T) {
	var collected []int
	From([]int{1, 2, 3}).ForEach(func(n int) { collected = append(collected, n) })
	assertSliceEqual(t, []int{1, 2, 3}, collected)
}

func TestToMap(t *testing.T) {
	type item struct {
		id   int
		name string
	}
	items := []item{{1, "a"}, {2, "b"}}
	result := From(items).ToMap(func(i item) interface{} { return i.id })
	if len(result) != 2 {
		t.Errorf("expected 2 entries, got %d", len(result))
	}
}

func TestGroupBy(t *testing.T) {
	result := From([]int{1, 2, 3, 4, 5, 6}).GroupBy(func(n int) interface{} { return n % 2 })
	if len(result) != 2 {
		t.Errorf("expected 2 groups, got %d", len(result))
	}
}

func TestPartition(t *testing.T) {
	evens, odds := From([]int{1, 2, 3, 4, 5}).Partition(func(n int) bool { return n%2 == 0 })
	assertSliceEqual(t, []int{2, 4}, evens)
	assertSliceEqual(t, []int{1, 3, 5}, odds)
}

func TestOrderBy(t *testing.T) {
	result := From([]int{5, 3, 1, 4, 2}).
		OrderBy(func(a, b int) bool { return a < b }).
		ToList()
	assertSliceEqual(t, []int{1, 2, 3, 4, 5}, result)
}

func TestOrderByDesc(t *testing.T) {
	result := From([]int{5, 3, 1, 4, 2}).
		OrderByDesc(func(a, b int) bool { return a < b }).
		ToList()
	assertSliceEqual(t, []int{5, 4, 3, 2, 1}, result)
}

// ========== Source Constructors ==========

func TestRange(t *testing.T) {
	result := IntRange(0, 5, 1).ToList()
	assertSliceEqual(t, []int{0, 1, 2, 3, 4}, result)
}

func TestRangeWithStep(t *testing.T) {
	result := IntRange(0, 10, 2).ToList()
	assertSliceEqual(t, []int{0, 2, 4, 6, 8}, result)
}

func TestRepeat(t *testing.T) {
	result := Repeat("x", 3).ToList()
	assertSliceEqual(t, []string{"x", "x", "x"}, result)
}

func TestRepeatForeverWithTake(t *testing.T) {
	result := RepeatForever(7).Take(4).ToList()
	assertSliceEqual(t, []int{7, 7, 7, 7}, result)
}

func TestIterate(t *testing.T) {
	result := Iterate(1, func(n int) bool { return n <= 16 }, func(n int) int { return n * 2 }).ToList()
	assertSliceEqual(t, []int{1, 2, 4, 8, 16}, result)
}

func TestFromChannel(t *testing.T) {
	ch := make(chan int, 3)
	ch <- 10
	ch <- 20
	ch <- 30
	close(ch)
	result := FromChannel(ch).ToList()
	assertSliceEqual(t, []int{10, 20, 30}, result)
}

func TestCache(t *testing.T) {
	result := From([]int{3, 1, 2}).Cache().ToList()
	assertSliceEqual(t, []int{3, 1, 2}, result)
}

func TestZip(t *testing.T) {
	a := From([]int{1, 2, 3})
	b := From([]int{10, 20, 30})
	result := a.Zip(b, func(x, y int) int { return x + y }).ToList()
	assertSliceEqual(t, []int{11, 22, 33}, result)
}

func TestZipDifferentLengths(t *testing.T) {
	a := From([]int{1, 2, 3, 4})
	b := From([]int{10, 20})
	result := a.Zip(b, func(x, y int) int { return x + y }).ToList()
	assertSliceEqual(t, []int{11, 22}, result)
}

func TestChunk(t *testing.T) {
	result := Chunk(From([]int{1, 2, 3, 4, 5}), 2).ToList()
	if len(result) != 3 {
		t.Errorf("expected 3 chunks, got %d", len(result))
	}
	assertSliceEqual(t, []int{1, 2}, result[0])
	assertSliceEqual(t, []int{3, 4}, result[1])
	assertSliceEqual(t, []int{5}, result[2])
}

func TestWindow(t *testing.T) {
	result := Window(From([]int{1, 2, 3, 4, 5}), 3).ToList()
	if len(result) != 3 {
		t.Errorf("expected 3 windows, got %d", len(result))
	}
	assertSliceEqual(t, []int{1, 2, 3}, result[0])
	assertSliceEqual(t, []int{2, 3, 4}, result[1])
	assertSliceEqual(t, []int{3, 4, 5}, result[2])
}

// ========== Chaining ==========

func TestChainedOperations(t *testing.T) {
	result := From([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}).
		Filter(func(n int) bool { return n%2 == 0 }).
		Take(3).
		ToList()
	assertSliceEqual(t, []int{2, 4, 6}, result)
}

func TestFilterThenDistinct(t *testing.T) {
	result := From([]int{1, 2, 2, 3, 3, 4}).
		Filter(func(n int) bool { return n > 1 }).
		Distinct().
		ToList()
	assertSliceEqual(t, []int{2, 3, 4}, result)
}

func TestSkipThenTake(t *testing.T) {
	result := From([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}).
		Skip(3).
		Take(4).
		ToList()
	assertSliceEqual(t, []int{4, 5, 6, 7}, result)
}

func TestZeroValuesPreserved(t *testing.T) {
	result := From([]int{0, 1, 0, 2, 0}).ToList()
	assertSliceEqual(t, []int{0, 1, 0, 2, 0}, result)
}

func TestFilterWithZeroValues(t *testing.T) {
	result := From([]int{0, 1, 2, 3}).
		Filter(func(n int) bool { return n >= 0 }).
		ToList()
	assertSliceEqual(t, []int{0, 1, 2, 3}, result)
}

func TestSumRange(t *testing.T) {
	sum := IntRange(1, 101, 1).Sum()
	if sum != 5050 {
		t.Errorf("expected 5050, got %f", sum)
	}
}

func TestFromMap(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	result := FromMap(m).ToList()
	if len(result) != 2 {
		t.Errorf("expected 2 pairs, got %d", len(result))
	}
}

func TestFromMapValues(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	result := FromMapValues(m)
	sum := result.Sum()
	if sum != 6 {
		t.Errorf("expected sum 6, got %f", sum)
	}
}

func TestClose(t *testing.T) {
	s := From([]int{1, 2, 3})
	s.Close()

	defer func() {
		r := recover()
		if r == nil {
			t.Error("expected panic on closed stream")
		}
	}()
	s.ToList()
}

// ========== Collect + Custom Extensibility ==========

func TestCollectWithCollector(t *testing.T) {
	s := From([]int{1, 2, 3, 4, 5})
	result := Collect(s, collectors.Counting[int]())
	if result != 5 {
		t.Errorf("expected 5, got %d", result)
	}
}

func TestCollectJoining(t *testing.T) {
	s := From([]string{"a", "b", "c"})
	result := Collect(s, collectors.Joining(", "))
	if result != "a, b, c" {
		t.Errorf("expected 'a, b, c', got '%s'", result)
	}
}

func TestCollectGroupingBy(t *testing.T) {
	s := From([]int{1, 2, 3, 4, 5, 6})
	result := Collect(s, collectors.GroupingBy[int, int](func(n int) int { return n % 2 }))
	if len(result) != 2 {
		t.Errorf("expected 2 groups, got %d", len(result))
	}
}

func TestApplyCustomOperation(t *testing.T) {
	// Custom operation: double each element using ApplyIterator
	result := From([]int{1, 2, 3}).
		ApplyIterator(func(iter iterators.Iterator[int]) iterators.Iterator[int] {
			var items []int
			for iter.HasNext() {
				items = append(items, iter.Next()*2)
			}
			return iterators.AsListIterator(items...)
		}).
		ToList()
	assertSliceEqual(t, []int{2, 4, 6}, result)
}

func TestTransform(t *testing.T) {
	// Reusable pipeline fragment
	topThreeEvens := func(s *GoStream[int]) *GoStream[int] {
		return s.Filter(func(n int) bool { return n%2 == 0 }).Take(3)
	}
	result := From([]int{1, 2, 3, 4, 5, 6, 7, 8}).Transform(topThreeEvens).ToList()
	assertSliceEqual(t, []int{2, 4, 6}, result)
}

func TestFromIterator(t *testing.T) {
	iter := iterators.AsListIterator(10, 20, 30)
	result := FromIterator[int](iter).ToList()
	assertSliceEqual(t, []int{10, 20, 30}, result)
}

// ========== Helpers ==========

func assertSliceEqual[T comparable](t *testing.T, expected, actual []T) {
	t.Helper()
	if len(expected) != len(actual) {
		t.Errorf("length mismatch: expected %v (len %d), got %v (len %d)", expected, len(expected), actual, len(actual))
		return
	}
	for i := range expected {
		if expected[i] != actual[i] {
			t.Errorf("mismatch at index %d: expected %v, got %v", i, expected[i], actual[i])
			return
		}
	}
}
