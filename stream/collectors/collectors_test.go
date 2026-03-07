package collectors

import (
	"testing"
)

func TestToList(t *testing.T) {
	result := Collect([]int{1, 2, 3}, ToList[int]())
	assertSlice(t, []int{1, 2, 3}, result)
}

func TestToMap(t *testing.T) {
	type item struct {
		id   int
		name string
	}
	items := []item{{1, "a"}, {2, "b"}, {3, "c"}}
	result := Collect(items, ToMap[item, int, string](
		func(i item) int { return i.id },
		func(i item) string { return i.name },
	))
	if len(result) != 3 {
		t.Errorf("expected 3 entries, got %d", len(result))
	}
	if result[2] != "b" {
		t.Errorf("expected result[2]='b', got '%s'", result[2])
	}
}

func TestGroupingBy(t *testing.T) {
	items := []int{1, 2, 3, 4, 5, 6}
	result := Collect(items, GroupingBy[int, int](func(n int) int { return n % 2 }))
	if len(result) != 2 {
		t.Errorf("expected 2 groups, got %d", len(result))
	}
	assertSlice(t, []int{1, 3, 5}, result[1])
	assertSlice(t, []int{2, 4, 6}, result[0])
}

func TestPartitioningBy(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}
	result := Collect(items, PartitioningBy[int](func(n int) bool { return n%2 == 0 }))
	assertSlice(t, []int{2, 4}, result.Matching)
	assertSlice(t, []int{1, 3, 5}, result.NonMatching)
}

func TestJoining(t *testing.T) {
	result := Collect([]string{"a", "b", "c"}, Joining(", "))
	if result != "a, b, c" {
		t.Errorf("expected 'a, b, c', got '%s'", result)
	}
}

func TestJoiningWithPrefixSuffix(t *testing.T) {
	result := Collect([]string{"a", "b", "c"}, JoiningWithPrefixSuffix(", ", "[", "]"))
	if result != "[a, b, c]" {
		t.Errorf("expected '[a, b, c]', got '%s'", result)
	}
}

func TestCounting(t *testing.T) {
	result := Collect([]int{1, 2, 3, 4, 5}, Counting[int]())
	if result != 5 {
		t.Errorf("expected 5, got %d", result)
	}
}

func TestSumming(t *testing.T) {
	result := Collect([]int{1, 2, 3, 4, 5}, Summing[int](func(n int) float64 { return float64(n) }))
	if result != 15.0 {
		t.Errorf("expected 15.0, got %f", result)
	}
}

func TestAveraging(t *testing.T) {
	result := Collect([]int{2, 4, 6}, Averaging[int](func(n int) float64 { return float64(n) }))
	if result != 4.0 {
		t.Errorf("expected 4.0, got %f", result)
	}
}

func TestAveragingEmpty(t *testing.T) {
	result := Collect([]int{}, Averaging[int](func(n int) float64 { return float64(n) }))
	if result != 0 {
		t.Errorf("expected 0, got %f", result)
	}
}

func TestMinBy(t *testing.T) {
	result := Collect([]int{3, 1, 4, 1, 5}, MinBy[int](func(a, b int) bool { return a < b }))
	if result != 1 {
		t.Errorf("expected 1, got %d", result)
	}
}

func TestMaxBy(t *testing.T) {
	result := Collect([]int{3, 1, 4, 1, 5}, MaxBy[int](func(a, b int) bool { return a < b }))
	if result != 5 {
		t.Errorf("expected 5, got %d", result)
	}
}

func TestReducing(t *testing.T) {
	result := Collect([]int{1, 2, 3, 4}, Reducing(0, func(a, b int) int { return a + b }))
	if result != 10 {
		t.Errorf("expected 10, got %d", result)
	}
}

func TestToSet(t *testing.T) {
	result := Collect([]int{1, 2, 2, 3, 3, 3}, ToSet[int]())
	if len(result) != 3 {
		t.Errorf("expected 3 unique, got %d", len(result))
	}
}

func TestMapping(t *testing.T) {
	result := Collect([]int{1, 2, 3},
		Mapping(func(n int) string {
			return string(rune('a' - 1 + n))
		}, ToList[string]()),
	)
	assertSlice(t, []string{"a", "b", "c"}, result)
}

func TestFiltering(t *testing.T) {
	result := Collect([]int{1, 2, 3, 4, 5},
		Filtering(func(n int) bool { return n%2 == 0 }, ToList[int]()),
	)
	assertSlice(t, []int{2, 4}, result)
}

func TestFlatMapping(t *testing.T) {
	result := Collect([]int{1, 2, 3},
		FlatMapping(func(n int) []int { return []int{n, n * 10} }, ToList[int]()),
	)
	assertSlice(t, []int{1, 10, 2, 20, 3, 30}, result)
}

func TestNewCollector(t *testing.T) {
	// Custom collector: collect into a comma-separated string of ints
	collector := NewCollector[int, *[]int, int](
		func() *[]int { s := make([]int, 0); return &s },
		func(acc *[]int, n int) *[]int { *acc = append(*acc, n*n); return acc },
		func(acc *[]int) int {
			sum := 0
			for _, v := range *acc {
				sum += v
			}
			return sum
		},
	)
	result := Collect([]int{1, 2, 3}, collector) // 1+4+9 = 14
	if result != 14 {
		t.Errorf("expected 14, got %d", result)
	}
}

func TestToStringCollector(t *testing.T) {
	type item struct{ name string }
	result := Collect([]item{{"Alice"}, {"Bob"}},
		ToStringCollector[item](" & ", func(i item) string { return i.name }),
	)
	if result != "Alice & Bob" {
		t.Errorf("expected 'Alice & Bob', got '%s'", result)
	}
}

func assertSlice[T comparable](t *testing.T, expected, actual []T) {
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
