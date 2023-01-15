package main

import (
	"fmt"
	"log"
	"net/url"
)

func main() {
	makeMap()
}

func makeMap() {
	// 1. make関数を使う
	// url.Valuesはmap[string][]stringのエイリアスなのでmake関数で初期化
	vs := make(url.Values)
	vs.Add("key1", "value1")
	vs.Add("key2", "value2")
	for k, v := range vs {
		fmt.Printf("%s: %v\n", k, v)
	}

	// 2. mapリテラルを使う
	// mapリテラルで初期化するとmap[string][]stringの型が推論される
	vs2 := url.Values{
		"key1": {"value1"},
		"key2": {"value2"},
	}

	log.Println(vs, vs2)
}
