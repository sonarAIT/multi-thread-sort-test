package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Elem struct {
	sync.Mutex
	Text             string
	HeavyProcessTime time.Duration
}

func Handler(w http.ResponseWriter, r *http.Request, ReceiveC chan<- *Elem) {
	// setting
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	r.Header.Set("Content-Type", "application/json")

	// unmarshal
	jsonBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintln(w, `{"status":"Unavailable"}`)
		fmt.Println("Can't catch Elem Data(io error)", err)
		return
	}

	var rec struct {
		Text             string `json:"Text"`
		HeavyProcessTime int    `json:"HeavyProcessTime"`
	}
	if err := json.Unmarshal(jsonBytes, &rec); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintln(w, `{"status":"Unavailable"}`)
		fmt.Println("Can't catch Elem Data(JSON Unmarshal error)", err)
		return
	}

	// send
	elem := Elem{
		Text:             rec.Text,
		HeavyProcessTime: time.Duration(rec.HeavyProcessTime),
	}

	ReceiveC <- &elem

	// res
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, `{"status":"available"}`)
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

func RoutinePrint(PrintC <-chan *Elem) {
	for elem := range PrintC {
		elem.Lock()
		fmt.Println(elem.Text)
	}
}

func main() {
	// init
	fmt.Println("init")

	ReceiveC := make(chan *Elem, 1024)
	HeavyProcessC := make(chan *Elem, 1024)
	PrintC := make(chan *Elem, 1024)
	go RoutineReceive(ReceiveC, HeavyProcessC, PrintC)
	for i := 0; i < 5; i++ {
		go RoutineHeavyProcess(HeavyProcessC)
	}
	go RoutinePrint(PrintC)

	// metrics
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()

	// run
	fmt.Println("run")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Handler(w, r, ReceiveC)
	})
	http.ListenAndServe(":8080", nil)
}
