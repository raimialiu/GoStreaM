package parallel

import (
	"github.com/raimialiu/gostream/stream"
	"sync/atomic"
	"testing"
)

func TestParallelForEach(t *testing.T) {
	var count int64
	ps := From(stream.From([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}))
	ps.ForEach(func(n int) {
		atomic.AddInt64(&count, 1)
	})
	if count != 10 {
		t.Errorf("expected 10, got %d", count)
	}
}

func TestParallelFilter(t *testing.T) {
	ps := FromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	result := ps.Filter(func(n int) bool { return n%2 == 0 }).ToSlice()
	if len(result) != 5 {
		t.Errorf("expected 5 evens, got %d", len(result))
	}
	for _, v := range result {
		if v%2 != 0 {
			t.Errorf("expected even, got %d", v)
		}
	}
}

func TestParallelMap(t *testing.T) {
	ps := FromSlice([]int{1, 2, 3, 4, 5})
	result := Map(ps, func(n int) int { return n * n })
	expected := []int{1, 4, 9, 16, 25}
	if len(result) != len(expected) {
		t.Fatalf("expected len %d, got %d", len(expected), len(result))
	}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("at index %d: expected %d, got %d", i, expected[i], v)
		}
	}
}

func TestParallelAnyMatch(t *testing.T) {
	ps := FromSlice([]int{1, 2, 3, 4, 5})
	if !ps.AnyMatch(func(n int) bool { return n == 3 }) {
		t.Error("expected true for AnyMatch 3")
	}
	if ps.AnyMatch(func(n int) bool { return n == 99 }) {
		t.Error("expected false for AnyMatch 99")
	}
}

func TestParallelAllMatch(t *testing.T) {
	if !FromSlice([]int{2, 4, 6, 8}).AllMatch(func(n int) bool { return n%2 == 0 }) {
		t.Error("expected true for AllMatch even")
	}
	if FromSlice([]int{2, 4, 6, 8}).AllMatch(func(n int) bool { return n > 5 }) {
		t.Error("expected false for AllMatch >5")
	}
}

func TestParallelReduce(t *testing.T) {
	ps := FromSlice([]int{1, 2, 3, 4, 5})
	result := ps.Reduce(0, func(a, b int) int { return a + b })
	if result != 15 {
		t.Errorf("expected 15, got %d", result)
	}
}

func TestParallelReduceProduct(t *testing.T) {
	ps := FromSlice([]int{1, 2, 3, 4})
	result := ps.Reduce(1, func(a, b int) int { return a * b })
	if result != 24 {
		t.Errorf("expected 24, got %d", result)
	}
}

func TestParallelWithWorkers(t *testing.T) {
	ps := FromSlice([]int{1, 2, 3, 4, 5}).WithWorkers(2)
	result := ps.Reduce(0, func(a, b int) int { return a + b })
	if result != 15 {
		t.Errorf("expected 15, got %d", result)
	}
}

func TestParallelCount(t *testing.T) {
	ps := FromSlice([]int{1, 2, 3, 4, 5})
	if ps.Count() != 5 {
		t.Errorf("expected 5, got %d", ps.Count())
	}
}

func TestParallelSequential(t *testing.T) {
	ps := FromSlice([]int{1, 2, 3})
	s := ps.Sequential()
	result := s.ToList()
	if len(result) != 3 {
		t.Errorf("expected 3 elements, got %d", len(result))
	}
}

func TestParallelEmpty(t *testing.T) {
	ps := FromSlice[int](nil)
	if ps.Count() != 0 {
		t.Errorf("expected 0, got %d", ps.Count())
	}
	result := ps.Reduce(0, func(a, b int) int { return a + b })
	if result != 0 {
		t.Errorf("expected 0, got %d", result)
	}
}

func TestParallelCollect(t *testing.T) {
	ps := FromSlice([]int{1, 2, 3, 4, 5})
	result := Collect(ps,
		func(items []int) int {
			sum := 0
			for _, v := range items {
				sum += v
			}
			return sum
		},
		func(a, b int) int { return a + b },
		0,
	)
	if result != 15 {
		t.Errorf("expected 15, got %d", result)
	}
}

func TestParallelMapTypeChange(t *testing.T) {
	ps := FromSlice([]int{1, 2, 3})
	result := Map(ps, func(n int) string {
		return string(rune('a' - 1 + n))
	})
	expected := []string{"a", "b", "c"}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("at index %d: expected %s, got %s", i, expected[i], v)
		}
	}
}

func TestParallelFilterChain(t *testing.T) {
	result := FromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}).
		Filter(func(n int) bool { return n%2 == 0 }).
		Filter(func(n int) bool { return n > 5 }).
		ToSlice()
	expected := []int{6, 8, 10}
	if len(result) != len(expected) {
		t.Fatalf("expected %d elements, got %d", len(expected), len(result))
	}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("at index %d: expected %d, got %d", i, expected[i], v)
		}
	}
}
