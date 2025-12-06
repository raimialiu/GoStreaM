package iterators

import "io"

type ReaderIterator[T any] struct {
	rd         io.Reader
	hasNext    bool
	nextData   string
	closer     io.Closer
	bufferSize int
}

func (r ReaderIterator[T]) HasNext() bool {
	if !r.hasNext {
		r.readNext()
	}

	return r.hasNext
}

func (r ReaderIterator[T]) Next() string {
	if r.HasNext() {
		return r.nextData
	}

	return ""
}

func (r ReaderIterator[T]) readNext() {
	buf := make([]byte, r.bufferSize)
	vl, err := r.rd.Read(buf)

	r.hasNext = err == nil && vl > 0
	r.nextData = string(buf[:vl])
}

func (r ReaderIterator[T]) Close() error {
	if r.closer != nil {
		return r.closer.Close()
	}

	return nil
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
