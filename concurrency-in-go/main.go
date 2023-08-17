package main

import (
	"fmt"
	"math"
	"os"
	"sync"
	"text/tabwriter"
	"time"
)

func main() {
	producer := func(wg *sync.WaitGroup, m sync.Locker) { // ❶
		defer wg.Done()
		for i := 5; i > 0; i-- {
			m.Lock()
			m.Unlock()
			time.Sleep(1) // ❷
		}
	}

	observer := func(wg *sync.WaitGroup, rwm sync.Locker) {
		defer wg.Done()
		rwm.Lock()
		defer rwm.Unlock()
	}

	test := func(count int, mutex, rwMutex sync.Locker) time.Duration {
		var wg sync.WaitGroup
		wg.Add(count + 1)
		beginTestTime := time.Now()
		go producer(&wg, mutex)
		for i := count; i > 0; i-- {
			go observer(&wg, rwMutex)
		}

		wg.Wait()
		return time.Since(beginTestTime)
	}

	tw := tabwriter.NewWriter(os.Stdout, 0, 1, 2, ' ', 0)
	defer tw.Flush()

	var m sync.RWMutex
	fmt.Fprintf(tw, "Readers\tRWMutex\tMutex\n")
	for i := 0; i < 20; i++ {
		count := int(math.Pow(2, float64(i)))
		fmt.Fprintf(
			tw,
			"%d\t%v\t%v\n",
			count,
			// observerとしてcountつ、producerとして1つのgougineを実行する
			// 前者はobserverのロック形式がRLockになっているが、後者はLockになってる
			// 理論的には前者のほうがクリティカルセクションが少ない。
			test(count, &m, m.RLocker()),
			test(count, &m, &m),
		)
	}

	i := (uint16)(1)
}
