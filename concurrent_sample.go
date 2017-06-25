package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type PingPong struct {
	Counter int
}

func main() {
	//for i := 0; i < 100; i++ {
	//	go test(i)
	//}

	//ch := make(chan int)
	//ch <- 1
	//println(<-ch)

	var p PingPong
	chA := make(chan *PingPong)
	chB := make(chan *PingPong)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		for {
			p, ok := <-chA
			if !ok {
				break
			}

			fmt.Printf("chA: p.Counter = %d\n", p.Counter)
			p.Counter++
			if p.Counter > 6 {
				break
			}

			chB <- p
		}
		close(chB)
	}()

	go func() {
		defer wg.Done()

		for {
			p, ok := <-chB
			if !ok {
				break
			}

			fmt.Printf("chB: p.Counter = %d\n", p.Counter)
			p.Counter++
			if p.Counter > 6 {
				break
			}

			chA <- p
		}
		close(chA)
	}()

	chA <- &p
	wg.Wait()

	ch := make(chan struct{})
	go func() {
		println("w")
		ch <- struct{}{}
	}()
	<-ch
	println("r") // 読み込めるまで待つ

	chch := make(chan int, 5)
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println("Writing", i)
			chch <- i
		}
		close(chch)
	}()

	time.Sleep(1 * time.Second)

	for {
		v, ok := <-chch
		if !ok {
			break
		}
		fmt.Println("Read", v)
	}

	chchchc := make(chan int, 10)
	for i := 0; i < 10; i++ {
		chchchc <- i
	}
	close(chchchc)
	for i := 0; i < 10; i++ {
		fmt.Println("Read", <-chchchc)
	}
	fmt.Println("Read", <-chchchc) // 読み込む値だないので0を返す(変数の初期値)

	// okがtrueの場合、channelから戻ってきた
	// 値はまだ有効
	// falseの場合はすでに閉じられている
	_, ok := <-chchchc
	if !ok {
		fmt.Println("chchchc is closed.")
	}

	var mu sync.RWMutex
	var wgwg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wgwg.Add(1)
		id := i
		go func() {
			defer wgwg.Done()
			for i := 0; i < 5; i++ {
				mu.RLock() // 共有ロックを得る
				fmt.Printf("Reader %d: Acquired lock\n", id)
				time.Sleep(time.Second)
				mu.RUnlock()
			}
		}()
	}

	time.AfterFunc(3*time.Second, func() {
		fmt.Println("Writer: Acquired lock")
		mu.Lock()
	})
	time.AfterFunc(6*time.Second, func() {
		fmt.Println("Writer: Releasing lock")
		mu.Unlock()
	})

	wgwg.Wait()
	// sync.WaitGroup は厳密にどの g o r o u t i n e が終了したかを管 理しているわけではなく、
	// 単純にカウンタを保持し ていて、その数が 0 になるまで待つことができる、と いうオブジェクト

	ExampleCond()

	ExampleOnce()

	// バッファ付き chanelで同時実行数を制御する
	var wg2 sync.WaitGroup
	sem := make(chan struct{}, 5)
	for i := 0; i < 10; i++ {
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			FetchURL(sem, "https://google.com")
		}()
	}
	wg2.Wait()

	// 連番生成
	c := GenerateSeries(10)

	var wg3 sync.WaitGroup

	wg3.Add(2)
	go ReadSeries(&wg3, c)
	go ReadSeries(&wg3, c)
	wg3.Wait()
}

type Object struct {
	mu sync.Mutex
}

func (o *Object) NeedsLock() {
	// 排他ロック
	o.mu.Lock()
	// メソッドが終了したら、自動的にMutexを解放
	defer o.mu.Unlock()
}

func ExampleCond() {
	var wg sync.WaitGroup
	var mu sync.Mutex
	c := sync.NewCond(&mu)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			defer mu.Unlock()

			fmt.Printf("waiting %d\n", i)
			mu.Lock()
			c.Wait()
			fmt.Printf("go %d\n", i)

		}(i)
	}

	//for i := 0; i < 10; i++ {
	//	time.Sleep(100 * time.Millisecond)
	//	fmt.Printf("signaling!\n")
	//	c.Signal() // 一つずつ通知
	//}
	time.AfterFunc(time.Second, c.Broadcast) // 一度にすべてに通知
	wg.Wait()
}

var Global int
var initGlobal sync.Once

func ExampleOnce() {
	cb := func(wg *sync.WaitGroup) {
		defer wg.Done()
		// 一度だけしか呼び出されない
		initGlobal.Do(func() {
			Global = 1
			fmt.Println("Write Global", Global)
		})
		fmt.Println("Read Global", Global)
	}
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go cb(&wg)
	}
	wg.Wait()
}

func FetchURL(sem chan struct{}, url string) {
	sem <- struct{}{}
	defer func() { <-sem }()

	fmt.Println("Get", url)
	http.Get(url)
	time.Sleep(time.Second)
	fmt.Println("End", url)
}

func test(i int) {
	println(i)
}

func ReadSeries(wg *sync.WaitGroup, c <-chan int) {
	defer wg.Done()
	for i := range c {
		fmt.Printf("Read %d", i)
	}
}

func GenerateSeries(max int) <-chan int {
	c := make(chan int)
	go func() {
		defer close(c)
		for i := 0; i < max; i++ {
			c <- i
		}
	}()
	return c
}
