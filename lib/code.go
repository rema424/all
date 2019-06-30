package lib

import (
	"bufio"
	"os"
	"strings"
)

var sc = bufio.NewScanner(os.Stdin)

// NextLine は標準入力から次の1行を読み込みます。
func NextLine() string {
	sc.Scan()
	return sc.Text()
}

// SplitByWhiteSpace は空白区切で文字列を分解します。
func SplitByWhiteSpace(str string) []string {
	return strings.Split(str, " ")
}

// SplitByCommma はコンマ区切で文字列を分解します。
func SplitByCommma(str string) []string {
	return strings.Split(str, ",")
}
