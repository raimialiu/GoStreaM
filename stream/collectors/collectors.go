package collectors

import (
	"fmt"
	"strings"
)

// ========== ToList Collector ==========

type toListCollector[T any] struct{}

func ToList[T any]() Collector[T, *[]T, []T] {
	return &toListCollector[T]{}
}

func (c *toListCollector[T]) Supplier() *[]T {
	list := make([]T, 0)
	return &list
}

func (c *toListCollector[T]) Accumulator(acc *[]T, element T) *[]T {
	*acc = append(*acc, element)
	return acc
}

func (c *toListCollector[T]) Finisher(acc *[]T) []T {
	return *acc
}

// ========== ToMap Collector ==========

type toMapCollector[T any, K comparable, V any] struct {
	keySelector   func(T) K
	valueSelector func(T) V
}

func ToMap[T any, K comparable, V any](keySelector func(T) K, valueSelector func(T) V) Collector[T, map[K]V, map[K]V] {
	return &toMapCollector[T, K, V]{
		keySelector:   keySelector,
		valueSelector: valueSelector,
	}
}

func (c *toMapCollector[T, K, V]) Supplier() map[K]V {
	return make(map[K]V)
}

func (c *toMapCollector[T, K, V]) Accumulator(acc map[K]V, element T) map[K]V {
	key := c.keySelector(element)
	value := c.valueSelector(element)
	acc[key] = value
	return acc
}

func (c *toMapCollector[T, K, V]) Finisher(acc map[K]V) map[K]V {
	return acc
}

// ========== GroupingBy Collector ==========

type groupingByCollector[T any, K comparable] struct {
	keySelector func(T) K
}

func GroupingBy[T any, K comparable](keySelector func(T) K) Collector[T, map[K][]T, map[K][]T] {
	return &groupingByCollector[T, K]{keySelector: keySelector}
}

func (c *groupingByCollector[T, K]) Supplier() map[K][]T {
	return make(map[K][]T)
}

func (c *groupingByCollector[T, K]) Accumulator(acc map[K][]T, element T) map[K][]T {
	key := c.keySelector(element)
	acc[key] = append(acc[key], element)
	return acc
}

func (c *groupingByCollector[T, K]) Finisher(acc map[K][]T) map[K][]T {
	return acc
}

// ========== PartitioningBy Collector ==========

type partitioningByCollector[T any] struct {
	predicate func(T) bool
}

type Partition[T any] struct {
	Matching    []T
	NonMatching []T
}

func PartitioningBy[T any](predicate func(T) bool) Collector[T, *Partition[T], *Partition[T]] {
	return &partitioningByCollector[T]{predicate: predicate}
}

func (c *partitioningByCollector[T]) Supplier() *Partition[T] {
	return &Partition[T]{}
}

func (c *partitioningByCollector[T]) Accumulator(acc *Partition[T], element T) *Partition[T] {
	if c.predicate(element) {
		acc.Matching = append(acc.Matching, element)
	} else {
		acc.NonMatching = append(acc.NonMatching, element)
	}
	return acc
}

func (c *partitioningByCollector[T]) Finisher(acc *Partition[T]) *Partition[T] {
	return acc
}

// ========== Joining Collector ==========

type joiningCollector struct {
	separator string
	prefix    string
	suffix    string
}

func Joining(separator string) Collector[string, *[]string, string] {
	return &joiningCollector{separator: separator}
}

func JoiningWithPrefixSuffix(separator, prefix, suffix string) Collector[string, *[]string, string] {
	return &joiningCollector{separator: separator, prefix: prefix, suffix: suffix}
}

func (c *joiningCollector) Supplier() *[]string {
	parts := make([]string, 0)
	return &parts
}

func (c *joiningCollector) Accumulator(acc *[]string, element string) *[]string {
	*acc = append(*acc, element)
	return acc
}

func (c *joiningCollector) Finisher(acc *[]string) string {
	return c.prefix + strings.Join(*acc, c.separator) + c.suffix
}

// ========== Counting Collector ==========

type countingCollector[T any] struct{}

func Counting[T any]() Collector[T, *int, int] {
	return &countingCollector[T]{}
}

func (c *countingCollector[T]) Supplier() *int {
	count := 0
	return &count
}

func (c *countingCollector[T]) Accumulator(acc *int, _ T) *int {
	*acc++
	return acc
}

func (c *countingCollector[T]) Finisher(acc *int) int {
	return *acc
}

// ========== Summing Collector ==========

type summingCollector[T any] struct {
	toFloat func(T) float64
}

func Summing[T any](toFloat func(T) float64) Collector[T, *float64, float64] {
	return &summingCollector[T]{toFloat: toFloat}
}

func (c *summingCollector[T]) Supplier() *float64 {
	sum := 0.0
	return &sum
}

func (c *summingCollector[T]) Accumulator(acc *float64, element T) *float64 {
	*acc += c.toFloat(element)
	return acc
}

func (c *summingCollector[T]) Finisher(acc *float64) float64 {
	return *acc
}

// ========== Averaging Collector ==========

type averagingCollector[T any] struct {
	toFloat func(T) float64
}

type avgAccumulator struct {
	sum   float64
	count int
}

func Averaging[T any](toFloat func(T) float64) Collector[T, *avgAccumulator, float64] {
	return &averagingCollector[T]{toFloat: toFloat}
}

func (c *averagingCollector[T]) Supplier() *avgAccumulator {
	return &avgAccumulator{}
}

func (c *averagingCollector[T]) Accumulator(acc *avgAccumulator, element T) *avgAccumulator {
	acc.sum += c.toFloat(element)
	acc.count++
	return acc
}

func (c *averagingCollector[T]) Finisher(acc *avgAccumulator) float64 {
	if acc.count == 0 {
		return 0
	}
	return acc.sum / float64(acc.count)
}

// ========== MinBy / MaxBy Collector ==========

type minByCollector[T any] struct {
	less func(a, b T) bool
}

type optionalAccumulator[T any] struct {
	value T
	found bool
}

func MinBy[T any](less func(a, b T) bool) Collector[T, *optionalAccumulator[T], (T)] {
	return &minByCollector[T]{less: less}
}

func (c *minByCollector[T]) Supplier() *optionalAccumulator[T] {
	return &optionalAccumulator[T]{}
}

func (c *minByCollector[T]) Accumulator(acc *optionalAccumulator[T], element T) *optionalAccumulator[T] {
	if !acc.found || c.less(element, acc.value) {
		acc.value = element
		acc.found = true
	}
	return acc
}

func (c *minByCollector[T]) Finisher(acc *optionalAccumulator[T]) T {
	return acc.value
}

type maxByCollector[T any] struct {
	less func(a, b T) bool
}

func MaxBy[T any](less func(a, b T) bool) Collector[T, *optionalAccumulator[T], T] {
	return &maxByCollector[T]{less: less}
}

func (c *maxByCollector[T]) Supplier() *optionalAccumulator[T] {
	return &optionalAccumulator[T]{}
}

func (c *maxByCollector[T]) Accumulator(acc *optionalAccumulator[T], element T) *optionalAccumulator[T] {
	if !acc.found || c.less(acc.value, element) {
		acc.value = element
		acc.found = true
	}
	return acc
}

func (c *maxByCollector[T]) Finisher(acc *optionalAccumulator[T]) T {
	return acc.value
}

// ========== Reducing Collector ==========

type reducingCollector[T any] struct {
	identity    T
	accumulator func(T, T) T
}

type reduceAccumulator[T any] struct {
	value T
}

func Reducing[T any](identity T, accumulator func(T, T) T) Collector[T, *reduceAccumulator[T], T] {
	return &reducingCollector[T]{identity: identity, accumulator: accumulator}
}

func (c *reducingCollector[T]) Supplier() *reduceAccumulator[T] {
	return &reduceAccumulator[T]{value: c.identity}
}

func (c *reducingCollector[T]) Accumulator(acc *reduceAccumulator[T], element T) *reduceAccumulator[T] {
	acc.value = c.accumulator(acc.value, element)
	return acc
}

func (c *reducingCollector[T]) Finisher(acc *reduceAccumulator[T]) T {
	return acc.value
}

// ========== Mapping Collector ==========

type mappingCollector[T any, U any, A any, R any] struct {
	mapper     func(T) U
	downstream Collector[U, A, R]
}

func Mapping[T any, U any, A any, R any](mapper func(T) U, downstream Collector[U, A, R]) Collector[T, A, R] {
	return &mappingCollector[T, U, A, R]{mapper: mapper, downstream: downstream}
}

func (c *mappingCollector[T, U, A, R]) Supplier() A {
	return c.downstream.Supplier()
}

func (c *mappingCollector[T, U, A, R]) Accumulator(acc A, element T) A {
	return c.downstream.Accumulator(acc, c.mapper(element))
}

func (c *mappingCollector[T, U, A, R]) Finisher(acc A) R {
	return c.downstream.Finisher(acc)
}

// ========== Filtering Collector ==========

type filteringCollector[T any, A any, R any] struct {
	predicate  func(T) bool
	downstream Collector[T, A, R]
}

func Filtering[T any, A any, R any](predicate func(T) bool, downstream Collector[T, A, R]) Collector[T, A, R] {
	return &filteringCollector[T, A, R]{predicate: predicate, downstream: downstream}
}

func (c *filteringCollector[T, A, R]) Supplier() A {
	return c.downstream.Supplier()
}

func (c *filteringCollector[T, A, R]) Accumulator(acc A, element T) A {
	if c.predicate(element) {
		return c.downstream.Accumulator(acc, element)
	}
	return acc
}

func (c *filteringCollector[T, A, R]) Finisher(acc A) R {
	return c.downstream.Finisher(acc)
}

// ========== ToSet Collector ==========

type toSetCollector[T comparable] struct{}

func ToSet[T comparable]() Collector[T, map[T]struct{}, map[T]struct{}] {
	return &toSetCollector[T]{}
}

func (c *toSetCollector[T]) Supplier() map[T]struct{} {
	return make(map[T]struct{})
}

func (c *toSetCollector[T]) Accumulator(acc map[T]struct{}, element T) map[T]struct{} {
	acc[element] = struct{}{}
	return acc
}

func (c *toSetCollector[T]) Finisher(acc map[T]struct{}) map[T]struct{} {
	return acc
}

// ========== FlatMapping Collector ==========

type flatMappingCollector[T any, U any, A any, R any] struct {
	mapper     func(T) []U
	downstream Collector[U, A, R]
}

func FlatMapping[T any, U any, A any, R any](mapper func(T) []U, downstream Collector[U, A, R]) Collector[T, A, R] {
	return &flatMappingCollector[T, U, A, R]{mapper: mapper, downstream: downstream}
}

func (c *flatMappingCollector[T, U, A, R]) Supplier() A {
	return c.downstream.Supplier()
}

func (c *flatMappingCollector[T, U, A, R]) Accumulator(acc A, element T) A {
	for _, u := range c.mapper(element) {
		acc = c.downstream.Accumulator(acc, u)
	}
	return acc
}

func (c *flatMappingCollector[T, U, A, R]) Finisher(acc A) R {
	return c.downstream.Finisher(acc)
}

// ========== Stringify Collector ==========

type toStringCollector[T any] struct {
	separator string
	formatter func(T) string
}

func ToStringCollector[T any](separator string, formatter func(T) string) Collector[T, *[]string, string] {
	return &toStringCollector[T]{separator: separator, formatter: formatter}
}

func (c *toStringCollector[T]) Supplier() *[]string {
	parts := make([]string, 0)
	return &parts
}

func (c *toStringCollector[T]) Accumulator(acc *[]string, element T) *[]string {
	*acc = append(*acc, c.formatter(element))
	return acc
}

func (c *toStringCollector[T]) Finisher(acc *[]string) string {
	return strings.Join(*acc, c.separator)
}

// ========== Collect Function ==========

// Collect applies a Collector to a slice of elements.
// Use this with GoStream.ToSlice() or directly with slices.
func Collect[T any, A any, R any](elements []T, collector Collector[T, A, R]) R {
	acc := collector.Supplier()
	for _, element := range elements {
		acc = collector.Accumulator(acc, element)
	}
	return collector.Finisher(acc)
}

// ========== Custom Collector Helper ==========

// FuncCollector allows creating a Collector from functions without defining a struct.
type FuncCollector[T any, A any, R any] struct {
	SupplierFn    func() A
	AccumulatorFn func(A, T) A
	FinisherFn    func(A) R
}

func NewCollector[T any, A any, R any](
	supplier func() A,
	accumulator func(A, T) A,
	finisher func(A) R,
) Collector[T, A, R] {
	return &FuncCollector[T, A, R]{
		SupplierFn:    supplier,
		AccumulatorFn: accumulator,
		FinisherFn:    finisher,
	}
}

func (c *FuncCollector[T, A, R]) Supplier() A            { return c.SupplierFn() }
func (c *FuncCollector[T, A, R]) Accumulator(a A, t T) A { return c.AccumulatorFn(a, t) }
func (c *FuncCollector[T, A, R]) Finisher(a A) R         { return c.FinisherFn(a) }

// ========== ForEach Collector (side-effect only) ==========

type forEachCollector[T any] struct {
	action func(T)
}

type unit struct{}

func ForEach[T any](action func(T)) Collector[T, *unit, *unit] {
	return &forEachCollector[T]{action: action}
}

func (c *forEachCollector[T]) Supplier() *unit {
	return &unit{}
}

func (c *forEachCollector[T]) Accumulator(acc *unit, element T) *unit {
	c.action(element)
	return acc
}

func (c *forEachCollector[T]) Finisher(acc *unit) *unit {
	return acc
}

// ========== Tee Collector (fan out to multiple collectors) ==========

type teeResult[R1 any, R2 any] struct {
	First  R1
	Second R2
}

// TeeResult holds the results from two collectors applied in parallel.
type TeeResult[R1, R2 any] struct {
	First  R1
	Second R2
}

// FormatCollector creates a collector that formats elements using fmt.Sprintf.
func FormatCollector[T any](format string, separator string) Collector[T, *[]string, string] {
	return ToStringCollector[T](separator, func(t T) string {
		return fmt.Sprintf(format, t)
	})
}
