package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
)

func main() {
	b, err := ioutil.ReadFile("in.csv")
	if err != nil {
		panic(err)
	}
	r := csv.NewReader(bytes.NewBuffer(b))
	r.LazyQuotes = true
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		fmt.Println(record, len(record))
	}
}
