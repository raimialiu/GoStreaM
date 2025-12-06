package iterators

import (
	"bufio"
	"os"
)

type FileIterator struct {
	scanner *bufio.Scanner
	closed  bool
	file    *os.File
	text    string
}

func AsFileIterator(filePath string) *FileIterator {
	file, fileReadError := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
	if fileReadError != nil {
		panic(fileReadError)
	}
	scanner := bufio.NewScanner(file)

	return &FileIterator{
		scanner: scanner,
		closed:  false,
		file:    file,
	}
}

func (it *FileIterator) Next() string {
	if !it.HasNext() {
		return ""
	}
	value := it.text
	it.readFromFile()
	return value
}

func (it *FileIterator) readFromFile() string {
	if it.closed {
		return ""
	}
	currentLine := it.scanner.Text()
	it.scanner.Scan()
	if err := it.scanner.Err(); err != nil {
		it.closed = true
	}
	it.text = currentLine
	return currentLine
}

func (it *FileIterator) HasNext() bool {
	return it.closed
}

func (it *FileIterator) Close() error {
	closeError := it.file.Close()
	if closeError != nil {
		panic(closeError)
	}

	it.closed = true
	return nil
}
