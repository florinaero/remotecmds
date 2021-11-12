package server

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// done * Current UTC time
// done * Current CPU usage
// * Available RAM
// * CPU usage over last hour
// * Available RAM over last hour
// * Download url into specific folder
// * Make computer "say" something
// * Capture and send a screenshot
// * Trigger webhook at specific time

// Request counter
type httpCounter struct {
	counter map[int]int
	chn chan int
}

// Grab counter from channel when is available
func (hc *httpCounter) grab_counter(w http.ResponseWriter) {
	cnt := <-hc.chn
	cnt++
	hc.counter[0] = cnt
	fmt.Fprint(w, "Counter = " + strconv.Itoa(cnt))
	time_now := time.Now().String()
	fmt.Fprintln(w, "	time: " + time_now)
}

// Send counter to channel for future grabs
func (hc *httpCounter) release_counter() {
	hc.chn <- hc.counter[0]
}

func (hc *httpCounter) time_utc(w http.ResponseWriter, req *http.Request) {
	hc.grab_counter(w)
	time_now := time.Now().UTC().String()
	fmt.Fprintln(w, "Request UTC time: "+time_now+"\n")
	hc.release_counter()
	// time.Sleep(6 * time.Second)
}

func get_cpu_usage() float64 {
	var sum float64
	args := "-A -o %cpu"
	cmd := exec.Command("ps",strings.Split(args," ")... ) 
	
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return -1
	}
	res := string(stdout)
	f := strings.Fields(res)
	// Skip first item that is string "CPU"
	count_proc := 0
	for i:=1;i<len(f);i++ {
		num,err := strconv.ParseFloat(f[i],64)
		if err!=nil{
			fmt.Println(err)
		}
		count_proc++
		sum = sum + num
	}
	return sum
}

func (hc *httpCounter) cpu_usage(w http.ResponseWriter, req *http.Request) {
	hc.grab_counter(w)

	// cpu_num := runtime.NumCPU()
	cpu_usage := get_cpu_usage()
	fmt.Fprintln(w, "cpu usage is "+strconv.FormatFloat(cpu_usage,'f',2,64)+"\n")
	hc.release_counter()
}

func HandleRequests() {
	hc := httpCounter{make(map[int]int), make(chan int)}
	hc.counter[0] = 1
	go func (){
		for{
			hc.chn <- hc.counter[0]
			hc.counter[0] = <-hc.chn
		}
		}()
		go http.HandleFunc("/time", hc.time_utc)
		go http.HandleFunc("/cpu", hc.cpu_usage)
}

func StartServer() {
	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
