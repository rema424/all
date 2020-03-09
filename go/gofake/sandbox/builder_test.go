package sandbox

import (
	"bytes"
	"strings"
	"testing"
)

func ConcatWithBytesBuffer(ss ...string) []byte {
	var b bytes.Buffer
	for _, s := range ss {
		b.WriteString(s)
	}
	return b.Bytes()
}

func ConcatWithStringsBuilder(ss ...string) []byte {
	var b strings.Builder
	for _, s := range ss {
		b.WriteString(s)
	}
	return []byte(b.String())
}

func ConcatWithStringsBuilderGrow(ss ...string) []byte {
	var b strings.Builder
	b.Grow(1024)
	for _, s := range ss {
		b.WriteString(s)
	}
	return []byte(b.String())
}

func Benchmark_Buffer(b *testing.B) {
	ss := []string{"aaa", "bbb", "ccc", "ddd", "eee", "fff", "ggg", "hhh", "iii"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ConcatWithBytesBuffer(ss...)
	}
}

func Benchmark_Builder(b *testing.B) {
	ss := []string{"aaa", "bbb", "ccc", "ddd", "eee", "fff", "ggg", "hhh", "iii"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ConcatWithStringsBuilder(ss...)
	}
}

func Benchmark_BufferGrow(b *testing.B) {
	ss := []string{"aaa", "bbb", "ccc", "ddd", "eee", "fff", "ggg", "hhh", "iii"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ConcatWithStringsBuilderGrow(ss...)
	}
}
