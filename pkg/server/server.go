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
// done * 5. Available RAM over last hour
// * 6. Download url into specific folder
// * 7. Make computer "say" something
// * 8. Capture and send a screenshot
// * 9. Trigger webhook at specific time

// Request counter
type Server struct {
	counter map[int]int
	chn     chan int
	ram_usage   *ram_h
}

// Grab counter from channel when is available
func (hc *Server) grab_counter(w http.ResponseWriter) {
	cnt := <-hc.chn
	cnt++
	hc.counter[0] = cnt
	fmt.Fprint(w, "Counter = "+strconv.Itoa(cnt))
	time_now := time.Now().String()
	fmt.Fprintln(w, "	time: "+time_now)
}

// Send counter to channel for future grabs
func (hc *Server) release_counter() {
	hc.chn <- hc.counter[0]
}

// 1. Current UTC time
func (hc *Server) time_utc(w http.ResponseWriter, req *http.Request) {
	hc.grab_counter(w)
	time_now := time.Now().UTC().String()
	fmt.Fprintln(w, "Request UTC time: "+time_now+"\n")
	hc.release_counter()
	// time.Sleep(6 * time.Second)
}

// 3. Available RAM
func (hc *Server) get_ram(w http.ResponseWriter, req *http.Request) {
	hc.grab_counter(w)
	total_ram := memory.TotalMemory()
	ram_str := strconv.FormatUint(total_ram, 10)
	fmt.Fprintln(w, "Total system RAM: "+ram_str+"\n")
	hc.release_counter()
}

// func
// Server handler requests
func HandleRequests() {
	hc := Server{make(map[int]int), make(chan int), new(ram_h)}
	hc.counter[0] = 1
	// Start sampling CPU and RAM usage
	go Start_cpu_avg()
	go hc.ram_usage.start_ram_avg()

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
	go http.HandleFunc("/ram_h", hc.ram_hour_usage)
}

func StartServer() {
	const PORT_NO = "8080"
	fmt.Printf("Starting server at port %s\n", PORT_NO)


	if err := http.ListenAndServe(":"+PORT_NO, nil); err != nil {
		log.Fatal(err)
	}
}
