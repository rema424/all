package lib

import "testing"

func BenchmarkDictionaryWithMap(b *testing.B) {
	num := 11
	commands := [11]string{
		"insert AAA",
		"insert AAC",
		"find AAA",
		"find CCC",
		"insert CCC",
		"find CCC",
		"find AA",
		"find GG",
		"insert GG",
		"find AA",
		"find GG",
	}
	for i := 0; i < b.N; i++ {
		dictionaryWithMap(num, commands[:])
	}
}

func BenchmarkDictionaryWithArray(b *testing.B) {
	num := 11
	commands := [11]string{
		"insert AAA",
		"insert AAC",
		"find AAA",
		"find CCC",
		"insert CCC",
		"find CCC",
		"find AA",
		"find GG",
		"insert GG",
		"find AA",
		"find GG",
	}
	for i := 0; i < b.N; i++ {
		dictionaryWithArray(num, commands[:])
	}
}

func BenchmarkDictionaryWithArrayOptimized(b *testing.B) {
	num := 11
	commands := [11]string{
		"insert AAA",
		"insert AAC",
		"find AAA",
		"find CCC",
		"insert CCC",
		"find CCC",
		"find AA",
		"find GG",
		"insert GG",
		"find AA",
		"find GG",
	}
	for i := 0; i < b.N; i++ {
		dictionaryWithArrayOptimized(num, commands[:])
	}
}
