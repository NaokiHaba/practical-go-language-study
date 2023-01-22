package main

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"sync"
	"time"
)

type Author struct {
	FirstName string
	LastName  string
}

type Book struct {
	Title      string
	Author     Author
	Publisher  string
	ReleasedAt time.Time
}

// NewAuthor creates a new Author factory
// ファクトリー関数は、構造体の初期化を行う関数
// この関数を使うことで、構造体の初期化を隠蔽できる
// また、構造体の初期化を行うロジックを変更するときに、
// この関数を変更するだけで済む
func NewAuthor(firstName, lastName string) *Author {
	return &Author{
		FirstName: firstName,
		LastName:  lastName,
	}
}

type Struct struct {
	v int
}

// NewStruct creates a new Struct factory
func NewStruct(v int) *Struct {
	return &Struct{
		v: v,
	}
}

// PrintStatus prints the status of the struct
// メソッドは、構造体のメソッドレシーバーを指定する
// メソッドレシーバーは、構造体のポインターを指定する
// これにより、構造体の値を変更できる
// インスタンスのレシーバーを指定すると、値渡しになる
// そのため、構造体の値を変更できない
func (s *Struct) PrintStatus() {
	log.Println("struct", s.v)
}

// String returns the string representation of the struct
// [T any]のように、メソッドレシーバーに型を指定することで、ジェネリクスを実現できる
func String[T any](v T) string {
	return fmt.Sprintf("%v", v)
}

// GenericStruct is a generic struct
type GenericStruct struct {
	t interface{}
}

// GenericString returns the string representation of the struct
func (s GenericStruct) GenericString() string {
	return fmt.Sprintf("%v", s.t)
}

// 構造体の埋め込みで共通部分を使い回す

type UrlBook struct {
	Title string
	// ISBNとは、国際標準図書番号のこと
	ISBN string
}

func (b UrlBook) GetAmazonURL() string {
	return "https://amazon.co.jp/dp/" + b.ISBN
}

type OraillyBook struct {
	UrlBook
	ISBN13 string
}

func (o OraillyBook) GetOraillyURL() string {
	return "https://orailly.co.jp/dp/" + o.ISBN13
}

func Decode(target interface{}, src map[string]string) error {
	// reflect.ValueOfは、インターフェースの値を返す
	v := reflect.ValueOf(target)

	// .Elem()は、ポインターの値を返す
	e := v.Elem()

	return decode(e, src)
}

func decode(e reflect.Value, src map[string]string) error {
	// reflect.Value.Type()は、型を返す
	t := e.Type()

	// reflect.Value.NumField()は、フィールドの数を返す
	for i := 0; i < t.NumField(); i++ {

		// reflect.Value.Field()は、フィールドの値を返す
		f := t.Field(i)

		// .Anonymousは、埋め込みフィールドかどうかを返す
		if f.Anonymous {
			if err := decode(e.Field(i), src); err != nil {
				return err
			}
			continue
		}

		// .Kind()は、型の種類を返す
		// reflect.Structは、構造体を表す
		if f.Type.Kind() == reflect.Struct {
			if err := decode(e.Field(i), src); err != nil {
				return err
			}
			continue
		}

		// .Tag.Get()は、タグの値を返す
		key := f.Tag.Get("map")
		if key == "" {
			key = f.Name
		}

		// 元データになければスキップ
		sv, ok := src[key]
		if !ok {
			continue
		}

		// フィールドの型を取得
		var k reflect.Kind
		var isP bool
		if f.Type.Kind() != reflect.Ptr {
			k = f.Type.Kind()
		} else {
			k = f.Type.Elem().Kind()
			// ポインターのポインターは無視
			if k == reflect.Ptr {
				continue
			}
			isP = true
		}
		switch k {
		case reflect.String:
			if isP {
				// .Field.Set()は、フィールドの値を設定する
				e.Field(i).Set(reflect.ValueOf(&sv))
			} else {
				// .Field.SetString()は、フィールドの値を文字列で設定する
				e.Field(i).SetString(sv)
			}
		case reflect.Bool:
			b, err := strconv.ParseBool(sv)
			if err == nil {
				if isP {
					e.Field(i).Set(reflect.ValueOf(&b))
				} else {
					e.Field(i).SetBool(b)
				}
			}
		case reflect.Int:
			n64, err := strconv.ParseInt(sv, 10, 64)
			if err == nil {
				if isP {
					n := int(n64)
					e.Field(i).Set(reflect.ValueOf(&n))
				} else {
					e.Field(i).SetInt(n64)
				}
			}
		}
	}
	return nil
}

type MapStruct struct {
	// mapのタグを指定する
	// mapのタグを指定しない場合、フィールド名がキーになる
	// mapのタグを指定する場合、タグの値がキーになる
	Str     string  `map:"str"`
	StrPtr  *string `map:"str"`
	Bool    bool    `map:"bool"`
	BoolPtr *bool   `map:"bool"`
	Int     int     `map:"int"`
	IntPtr  *int    `map:"int"`
}

// 空の構造体を使ってゴルーチン間での通知を行う

func PrintCh() {
	// make(chan struct{})は、空の構造体のチャネルを作成する
	wait := make(chan struct{})

	// go func()でゴルーチンを作成する
	go func() {
		// 何かの処理
		fmt.Println("送信")

		// wait <- struct{}{}で、空の構造体を送信する
		wait <- struct{}{}
	}()
	fmt.Println("受信待ち")

	// <-waitで、空の構造体を受信する
	<-wait

	fmt.Println("受信")
}

type BigStruct struct {
	Member string
}

// sync.Poolを使ってメモリの再利用を行う
var pool = sync.Pool{
	// New()は、sync.Poolに値がない場合に呼ばれる
	New: func() interface{} {
		return &BigStruct{}
	},
}

func main() {
	s := NewStruct(1)
	s.PrintStatus()

	log.Println(GenericStruct{t: 1}.GenericString())

	ob := OraillyBook{
		UrlBook: UrlBook{
			Title: "Mithril",
			ISBN:  "4873119030",
		},
		ISBN13: "9784873119030",
	}

	log.Println(ob.GetAmazonURL())
	log.Println(ob.GetOraillyURL())

	log.Println(Decode(&ob, map[string]string{
		"Title": "Mithril",
		"ISBN":  "4873119030",
	}))

	PrintCh()

	// Get()でsync.Poolから値を取得する
	b := pool.Get().(*BigStruct)

	// .Put()でsync.Poolに値を返す
	pool.Put(b)
}
