package main

import "sync"

type Elem struct {
	Text             string
	HeavyProcessTime int

	state struct {
		sync.Mutex
	}
}

func (elem *Elem) HeavyProcess() {

}

func (elem *Elem) GetText() {

}
