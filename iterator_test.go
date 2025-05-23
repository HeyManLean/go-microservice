package main

import (
	"fmt"
	"iter"
	"testing"
)

// 生成器，go1.23+
// iter.Seq[T] 是生成器返回的类型，传入 func(yield func(T) bool)
func genMultiple[T ~int](multiple int) iter.Seq[T] {
	/*
		按照 multiple 倍数，从 1 生成倍数，1, 3, 9, 27, 91...
	*/
	var i T = 1
	return func(yield func(T) bool) {
		for {
			if !yield(i) {
				break
			}
			i *= T(multiple)
		}
	}
}

func TestIterator(t *testing.T) {
	/*
		生成器，动态生成数据，避免占用内存过多
	*/
	for num := range genMultiple[int](3) {
		if num >= 100 {
			break
		}
		fmt.Println(num)
	}
}
