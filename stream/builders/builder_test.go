package builders

import (
	"github.com/raimialiu/gostream/stream"
	"testing"
)

func TestBuilderAdd(t *testing.T) {
	result := NewBuilder[int]().Add(1).Add(2).Add(3).ToSlice()
	assertSlice(t, []int{1, 2, 3}, result)
}

func TestBuilderAddAll(t *testing.T) {
	result := NewBuilder[int]().AddAll(1, 2, 3).ToSlice()
	assertSlice(t, []int{1, 2, 3}, result)
}

func TestBuilderAddSlice(t *testing.T) {
	result := NewBuilder[int]().AddSlice([]int{4, 5, 6}).ToSlice()
	assertSlice(t, []int{4, 5, 6}, result)
}

func TestBuilderAddFirst(t *testing.T) {
	result := NewBuilder[int]().Add(2).Add(3).AddFirst(1).ToSlice()
	assertSlice(t, []int{1, 2, 3}, result)
}

func TestBuilderAddAt(t *testing.T) {
	result := NewBuilder[int]().AddAll(1, 3).AddAt(1, 2).ToSlice()
	assertSlice(t, []int{1, 2, 3}, result)
}

func TestBuilderAddIf(t *testing.T) {
	result := NewBuilder[int]().
		AddIf(true, 1).
		AddIf(false, 2).
		AddIf(true, 3).
		ToSlice()
	assertSlice(t, []int{1, 3}, result)
}

func TestBuilderAddIfElse(t *testing.T) {
	result := NewBuilder[string]().
		AddIfElse(true, "yes", "no").
		AddIfElse(false, "yes", "no").
		ToSlice()
	assertSlice(t, []string{"yes", "no"}, result)
}

func TestBuilderAddUnless(t *testing.T) {
	result := NewBuilder[int]().
		AddUnless(false, 1).
		AddUnless(true, 2).
		ToSlice()
	assertSlice(t, []int{1}, result)
}

func TestBuilderAddWith(t *testing.T) {
	counter := 0
	result := NewBuilder[int]().
		AddWith(func() int { counter++; return counter }).
		AddWith(func() int { counter++; return counter }).
		ToSlice()
	assertSlice(t, []int{1, 2}, result)
}

func TestBuilderAddWithIf(t *testing.T) {
	result := NewBuilder[int]().
		AddWithIf(true, func() int { return 1 }).
		AddWithIf(false, func() int { return 2 }).
		ToSlice()
	assertSlice(t, []int{1}, result)
}

func TestBuilderAddRange(t *testing.T) {
	result := NewBuilder[int]().
		AddRange(0, 5, func(i int) int { return i * 2 }).
		ToSlice()
	assertSlice(t, []int{0, 2, 4, 6, 8}, result)
}

func TestBuilderAddStream(t *testing.T) {
	s := stream.Of(10, 20, 30)
	result := NewBuilder[int]().Add(1).AddStream(s).ToSlice()
	assertSlice(t, []int{1, 10, 20, 30}, result)
}

func TestBuilderMerge(t *testing.T) {
	b1 := NewBuilder[int]().AddAll(1, 2)
	b2 := NewBuilder[int]().AddAll(3, 4)
	result := b1.Merge(b2).ToSlice()
	assertSlice(t, []int{1, 2, 3, 4}, result)
}

func TestBuilderWhen(t *testing.T) {
	result := NewBuilder[int]().
		Add(1).
		When(true).Add(2).End().
		When(false).Add(3).End().
		ToSlice()
	assertSlice(t, []int{1, 2}, result)
}

func TestBuilderWhenElse(t *testing.T) {
	result := NewBuilder[string]().
		When(false).Add("yes").ElseAdd("no").End().
		ToSlice()
	assertSlice(t, []string{"no"}, result)
}

func TestBuilderUnless(t *testing.T) {
	result := NewBuilder[int]().
		Unless(false).Add(1).End().
		Unless(true).Add(2).End().
		ToSlice()
	assertSlice(t, []int{1}, result)
}

func TestBuilderSwitch(t *testing.T) {
	result := NewBuilder[string]().
		Switch("b").
		Case("a", "alpha").
		Case("b", "beta").
		Case("c", "gamma").
		Default("unknown").
		End().
		ToSlice()
	assertSlice(t, []string{"beta"}, result)
}

func TestBuilderSwitchDefault(t *testing.T) {
	result := NewBuilder[string]().
		Switch("z").
		Case("a", "alpha").
		Default("unknown").
		End().
		ToSlice()
	assertSlice(t, []string{"unknown"}, result)
}

func TestBuilderInspection(t *testing.T) {
	b := NewBuilder[int]().AddAll(1, 2, 3)

	if b.Count() != 3 {
		t.Errorf("expected count 3, got %d", b.Count())
	}
	if b.IsEmpty() {
		t.Error("expected not empty")
	}
	if !b.Contains(2) {
		t.Error("expected contains 2")
	}
	if b.Contains(5) {
		t.Error("expected not contains 5")
	}
	last, ok := b.Last()
	if !ok || last != 3 {
		t.Errorf("expected last=3, got %d (ok=%v)", last, ok)
	}
	val, ok := b.Get(1)
	if !ok || val != 2 {
		t.Errorf("expected get(1)=2, got %d", val)
	}
}

func TestBuilderRemoveLast(t *testing.T) {
	result := NewBuilder[int]().AddAll(1, 2, 3).RemoveLast().ToSlice()
	assertSlice(t, []int{1, 2}, result)
}

func TestBuilderRemoveAt(t *testing.T) {
	result := NewBuilder[int]().AddAll(1, 2, 3).RemoveAt(1).ToSlice()
	assertSlice(t, []int{1, 3}, result)
}

func TestBuilderRemoveIf(t *testing.T) {
	result := NewBuilder[int]().AddAll(1, 2, 3, 4, 5).
		RemoveIf(func(n int) bool { return n%2 == 0 }).
		ToSlice()
	assertSlice(t, []int{1, 3, 5}, result)
}

func TestBuilderClear(t *testing.T) {
	b := NewBuilder[int]().AddAll(1, 2, 3).Clear()
	if !b.IsEmpty() {
		t.Error("expected empty after clear")
	}
}

func TestBuilderBuild(t *testing.T) {
	s := NewBuilder[int]().AddAll(1, 2, 3).Build()
	result := s.ToList()
	assertSlice(t, []int{1, 2, 3}, result)
}

func TestBuilderBuildAndClear(t *testing.T) {
	b := NewBuilder[int]().AddAll(1, 2, 3)
	s := b.BuildAndClear()
	result := s.ToList()
	assertSlice(t, []int{1, 2, 3}, result)
	if !b.IsEmpty() {
		t.Error("expected empty after BuildAndClear")
	}
}

func TestNewBuilderFrom(t *testing.T) {
	result := NewBuilderFrom(1, 2, 3).ToSlice()
	assertSlice(t, []int{1, 2, 3}, result)
}

func TestNewBuilderFromSlice(t *testing.T) {
	result := NewBuilderFromSlice([]int{4, 5, 6}).ToSlice()
	assertSlice(t, []int{4, 5, 6}, result)
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
