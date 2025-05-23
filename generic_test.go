package main

import (
	"fmt"
	"testing"
)

// 泛型，go1.18+
// 通过 [T any] 声明泛型类型，通过 User[T] 声明一个结构体
// 使用一个泛型结构体，需要先将泛型指定类型后，才能开始使用
type User[P any, A comparable] struct {
	name P
	age  A
}

func (u *User[P, A]) Show() (P, A) {
	return u.name, u.age
}

// ~int 波浪号表示允许底层类型是 int 的类型传入，如 Myint 底层是 int，直接传入 Myint 会报错
// User[string, int] 注入泛型，声明一个具体的结构体
func ShowUser[T ~int](u *User[string, int], rank T) {
	name, age := u.Show()
	fmt.Println("Name: ", name)
	fmt.Println("Age: ", age)
	fmt.Println("Rank: ", rank)
}

type MyInt int

func TestGeneric(t *testing.T) {
	/*
		泛型
	*/
	var rank MyInt = 1
	user := &User[string, int]{"zhangsan", int(rank)}

	// 自动识别泛型类型
	ShowUser(user, rank)
	// 显式声明泛型类型
	ShowUser[MyInt](user, rank)
}
