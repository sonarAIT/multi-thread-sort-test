package main

import (
	"fmt"
	"sync"
	"time"
)

type Elem struct {
	sync.Mutex
	Text             string
	HeavyProcessTime time.Duration
}

func RoutineReceive(ReceiveC <-chan *Elem, HeavyProcessC chan<- *Elem, PrintC chan<- *Elem) {
	for elem := range ReceiveC {
		elem.Lock()
		PrintC <- elem
		HeavyProcessC <- elem
	}
}

func RoutineHeavyProcess(HeavyProcessC <-chan *Elem) {
	for elem := range HeavyProcessC {
		// heavy process
		time.Sleep(time.Second * elem.HeavyProcessTime)

		elem.Unlock()
	}
}

func RoutinePrint(PrintC <-chan *Elem, wg *sync.WaitGroup) {
	for elem := range PrintC {
		elem.Lock()
		fmt.Println(elem.Text)
		wg.Done()
	}
}

func main() {
	// init
	fmt.Println("init")
	elems := []*Elem{
		{
			Text:             "hello",
			HeavyProcessTime: 1,
		},
		{
			Text:             "hello2",
			HeavyProcessTime: 1,
		},
	}

	var wg sync.WaitGroup
	wg.Add(len(elems))

	ReceiveC := make(chan *Elem, 1024)
	HeavyProcessC := make(chan *Elem, 1024)
	PrintC := make(chan *Elem, 1024)
	go RoutineReceive(ReceiveC, HeavyProcessC, PrintC)
	for i := 0; i < 5; i++ {
		go RoutineHeavyProcess(HeavyProcessC)
	}
	go RoutinePrint(PrintC, &wg)

	// run
	fmt.Println("run")
	for _, elem := range elems {
		ReceiveC <- elem
	}

	wg.Wait()
}
