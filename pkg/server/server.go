package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"github.com/pbnjay/memory"
)

var Cpu_avgh_ cpu_h

// done * 1. Current UTC time
// done * 2. Current CPU usage
// done * 3. Available RAM
// done * 4. CPU usage over last hour
// * 5. Available RAM over last hour
// * 6. Download url into specific folder
// * 7. Make computer "say" something
// * 8. Capture and send a screenshot
// * 9. Trigger webhook at specific time

// Request counter
type HttpCounter struct {
	counter map[int]int
	chn     chan int
}

// Grab counter from channel when is available
func (hc *HttpCounter) grab_counter(w http.ResponseWriter) {
	cnt := <-hc.chn
	cnt++
	hc.counter[0] = cnt
	fmt.Fprint(w, "Counter = "+strconv.Itoa(cnt))
	time_now := time.Now().String()
	fmt.Fprintln(w, "	time: "+time_now)
}

// Send counter to channel for future grabs
func (hc *HttpCounter) release_counter() {
	hc.chn <- hc.counter[0]
}

// 1. Current UTC time
func (hc *HttpCounter) time_utc(w http.ResponseWriter, req *http.Request) {
	hc.grab_counter(w)
	time_now := time.Now().UTC().String()
	fmt.Fprintln(w, "Request UTC time: "+time_now+"\n")
	hc.release_counter()
	// time.Sleep(6 * time.Second)
}



// 3. Available RAM
func (hc *HttpCounter) get_ram(w http.ResponseWriter, req *http.Request) {
	hc.grab_counter(w)
	total_ram := memory.TotalMemory()
	ram_str := strconv.FormatUint(total_ram, 10)
	fmt.Fprintln(w, "Total system RAM: "+ram_str+"\n")
	hc.release_counter()
}

// Server handler requests
func HandleRequests() {
	hc := HttpCounter{make(map[int]int), make(chan int)}
	hc.counter[0] = 1
	go func() {
		for {
			hc.chn <- hc.counter[0]
			hc.counter[0] = <-hc.chn
		}
	}()
	go http.HandleFunc("/time", hc.time_utc)
	go http.HandleFunc("/cpu", hc.cpu_usage)
	go http.HandleFunc("/ram", hc.get_ram)
    go http.HandleFunc("/cpu_h", hc.cpu_hour_usage)
}

func StartServer() {
	const PORT_NO = "8080"
	fmt.Printf("Starting server at port %s\n", PORT_NO)
	
	go Start_cpu_avg()
	
	if err := http.ListenAndServe(":"+PORT_NO, nil); err != nil {
		log.Fatal(err)
	}
}
