package main

import "fmt"

type Func func() string

func (f Func) String() string { return f() }

type Hoge struct {
	N int
}

// 埋め込み
type Fuga struct {
	Hoge
}

type Stringer interface {
	String() string
}

type Hex int

func (h Hex) String() string {
	return fmt.Sprintf("%x", int(h))
}

// Hex2もStringerを実装
type Hex2 struct{ Hex }

type HogeHoge interface {
	M()
	N()
}

type fuga struct {
	HogeHoge
}

// Mの振る舞いを変える
func (f fuga) M() {
	fmt.Println("HI")
	f.HogeHoge.M() // 元のメソッドを呼び出す
}

func HiHoge(h HogeHoge) HogeHoge {
	return fuga{h}
}

// 関数の外では、キーワードではじまる宣言( var, func, など)が必要で、 := での暗黙的な宣言は利用できない
// aa := 2

func main() {
	f := Func(func() string { return "sample" })
	fmt.Println(f.String())

	var a interface {
		String() string
	}
	a = Func(func() string { return "test" })
	fmt.Println(a.String())

	v, ok := a.(Func)
	fmt.Println(v, ok)

	var vv interface{}
	vv = "aa"
	s, ok := vv.(int)
	fmt.Println(s, ok)

	switch v := vv.(type) {
	case int:
		fmt.Println(v * 2)
	case string:
		fmt.Println(v + "string")
	}

	fuga := Fuga{Hoge{N: 100}}
	// Hoge型のフィールドにアクセスできる
	fmt.Println(fuga.N)
	// 型名を指定してアクセスできる
	fmt.Println(fuga.Hoge.N)

	var ss Stringer
	h := Hex(100)
	ss = h
	fmt.Println(ss.String())

	h2 := Hex2{h}
	ss = h2
	fmt.Println(ss.String())

	//q := []int{2, 3, 5, 7, 11, 13}
	//fmt.Printf("%T", q)

	aa := make([]int, 5)
	printSlice("a", aa)
	for i := 0; i < 10; i++ {
		aa = append(aa, i)
	}
	printSlice("a", aa)
}

func printSlice(s string, x []int) {
	fmt.Printf("%s len=%d cap=%d %v\n",
		s, len(x), cap(x), x)
}
