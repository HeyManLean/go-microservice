package main

import (
	"fmt"
)

func findKthNumber(n int, k int) int {
	/*
	   返回cur开头的数字个数
	   - 如果个数小于k，则表示k可以包含这些数字，扣减个数后，需要将cur加1
	   - 如果大于k，则表示k在cur开头的数字里面，将cur * 10，k-=1 继续遍历
	*/
	cur := 1
	k -= 1
	for k > 0 {
		steps := getSteps(n, cur)
		fmt.Println(k, cur, steps)
		if steps <= k {
			k -= steps
			cur += 1
		} else {
			k -= 1
			cur *= 10 // 下一个数
		}
	}
	return cur
}

func getSteps(n int, cur int) (steps int) {
	/*
	   找出cur开头的数字，在n内还有多少个数字
	*/
	first, last := cur, cur
	for first <= n {
		steps += min(last, n) - first + 1
		first *= 10        // 10,100,100
		last = last*10 + 9 // 19,199,1999
	}
	return
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func main() {
	res := findKthNumber(10000, 10000)
	fmt.Println(res, 9999)
}
