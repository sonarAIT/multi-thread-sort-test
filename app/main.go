package main

import (
	"fmt"
	"sync"
)

func RoutineReceive(ch chan string) {
	for {
		s, isActive := <-ch
		if !isActive {
			return
		}
		fmt.Println(s)
	}
}

func RoutineHeavyProcess() {

}

func RoutinePrint(wg *sync.WaitGroup) {

}

func main() {
	ch := make(chan string)
	go RoutineReceive(ch)
	ch <- "hello"
	ch <- "hello2"
	ch <- "hello3"
}
