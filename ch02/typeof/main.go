package main

import (
	"fmt"
	"io"
	"os"
	"reflect"
)

func main() {
	typeOf()
}

func typeOf() {
	// 型が一致するかを判定する
	// reflect.TypeOf は型を返す

	// int
	fmt.Println(reflect.TypeOf(1))

	// string
	fmt.Println(reflect.TypeOf("string"))

	// []int
	fmt.Println(reflect.TypeOf([]int{1, 2, 3}))

	// map[string]int
	fmt.Println(reflect.TypeOf(map[string]int{"one": 1, "two": 2}))

	// func()
	fmt.Println(reflect.TypeOf(func() {}))

	writerType := reflect.TypeOf((*io.Writer)(nil)).Elem()

	fileType := reflect.TypeOf((*os.File)(nil))

	// fileType.Implements(writerType) は fileType が io.Writer を実装しているかを返す
	fmt.Println(fileType.Implements(writerType))
}
