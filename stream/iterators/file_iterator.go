package iterators

import (
	"bufio"
	"os"
)

type FileIterator struct {
	scanner *bufio.Scanner
	file    *os.File
	hasNext bool
}

func AsFileIterator(filePath string) *FileIterator {
	file, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(file)
	it := &FileIterator{
		scanner: scanner,
		file:    file,
		hasNext: scanner.Scan(),
	}
	return it
}

func (it *FileIterator) HasNext() bool {
	return it.hasNext
}

func (it *FileIterator) Next() string {
	if !it.HasNext() {
		return ""
	}
	text := it.scanner.Text()
	it.hasNext = it.scanner.Scan()
	return text
}

func (it *FileIterator) Close() error {
	it.hasNext = false
	return it.file.Close()
}
