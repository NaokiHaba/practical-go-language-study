package main

import (
	"encoding/json"
	"fmt"
)

type ConfidentialCustomer struct {
	CustomerID int64
	CreditCard CreditCard
}

type CreditCard string

// Stringerインターフェースを実装しているので、fmt.Println(c)で呼ばれる
// StringerインターフェースはfmtパッケージのString()メソッドを実装している
func (c CreditCard) String() string {
	return "xxxx-xxxx-xxxx-xxxx"
}

// GoString GoStringはfmtパッケージのGoString()メソッドを実装している
func (c CreditCard) GoString() string {
	return "xxxx-xxxx-xxxx-xxxx"
}

func main() {
	// 構造体の初期化
	c := ConfidentialCustomer{
		CustomerID: 1,
		CreditCard: "4111-1111-1111-1111",
	}

	fmt.Println(c)

	// %v は構造体のフィールド名を表示しない
	fmt.Printf("%v\n", c)

	// %+v は構造体のフィールド名を表示する
	fmt.Printf("%+v\n", c)

	// %#v は構造体のフィールド名と型を表示する
	fmt.Printf("%#v\n", c)

	// JSONに変換
	bytes, _ := json.Marshal(c)

	// JSONを出力
	fmt.Println("JSON: ", string(bytes)) // 元通り利用可能
}
