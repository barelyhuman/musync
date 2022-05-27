package utils

import "fmt"

type Writer struct {
	silent bool
}

type WriterOptions func(w *Writer)

func (wr *Writer) Success(msg string, values ...interface{}) {
	if wr.silent {
		return
	}
	fmt.Printf("\033[32m✓\033[0m " + msg + "\n\n")
}

func (wr *Writer) Error(msg string, values ...interface{}) {
	fmt.Printf(msg, values...)
}

func (wr *Writer) Info(msg string) {
	if wr.silent {
		return
	}
	fmt.Printf("\033[32m>>\033[0m " + msg + "\n\n")
}

func NewWriter(
	options ...WriterOptions,
) *Writer {
	writer := &Writer{}
	for _, opt := range options {
		opt(writer)
	}
	return writer
}

func WithSilent(silent bool) WriterOptions {
	return func(w *Writer) {
		w.silent = silent
	}
}
