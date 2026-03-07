package iterators

import "io"

type ReaderIterator[T any] struct {
	rd         io.Reader
	hasNext    bool
	nextData   string
	closer     io.Closer
	bufferSize int
	computed   bool
}

func AsReaderIterator[T any](r io.Reader, bufferSize int) *ReaderIterator[T] {
	var closerInstance io.Closer
	if closer, ok := r.(io.Closer); ok {
		closerInstance = closer
	}

	return &ReaderIterator[T]{
		rd:         r,
		bufferSize: bufferSize,
		closer:     closerInstance,
	}
}

func (r *ReaderIterator[T]) HasNext() bool {
	if !r.computed {
		r.readNext()
	}
	return r.hasNext
}

func (r *ReaderIterator[T]) Next() string {
	if !r.HasNext() {
		return ""
	}

	result := r.nextData
	r.computed = false
	r.hasNext = false
	return result
}

func (r *ReaderIterator[T]) readNext() {
	buf := make([]byte, r.bufferSize)
	n, err := r.rd.Read(buf)

	r.hasNext = err == nil && n > 0
	if r.hasNext {
		r.nextData = string(buf[:n])
	}
	r.computed = true
}

func (r *ReaderIterator[T]) Close() error {
	if r.closer != nil {
		return r.closer.Close()
	}
	return nil
}
