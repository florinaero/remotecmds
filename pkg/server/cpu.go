package server

import (
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/florinaero/remotecmds/pkg/config"
)

// Get cpu usage by reading output of 'ps' sheel instruction
func get_cpu_usage() float64 {
	var sum float64
	args := "-A -o %cpu"
	cmd := exec.Command("ps", strings.Split(args, " ")...)

	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return -1
	}
	res := string(stdout)
	f := strings.Fields(res)
	// Skip first item that is string "CPU"
	count_proc := 0
	for i := 1; i < len(f); i++ {
		num, err := strconv.ParseFloat(f[i], 64)
		if err != nil {
			fmt.Println(err)
		}
		count_proc++
		sum = sum + num
	}
	return sum
}

// 2. Current CPU usage
func (hc *HttpCounter) cpu_usage(w http.ResponseWriter, req *http.Request) {
	hc.grab_counter(w)

	// cpu_num := runtime.NumCPU()
	cpu_usage := get_cpu_usage()
	fmt.Fprintln(w, "cpu usage is "+strconv.FormatFloat(cpu_usage, 'f', 2, 64)+"\n")
	hc.release_counter()
}


// 4. CPU usage over last hour
func (hc *HttpCounter) cpu_hour_usage(w http.ResponseWriter, req *http.Request) {
	hc.grab_counter(w)
	avg := Cpu_avgh_.get_cpu_avg()
	fmt.Fprintf(w, "CPU average on last hour is: %3.2f\n", avg)
	hc.release_counter()
}

type cpu_h struct {
	pid string
	process_name string
	process_time string
	cpu_usage float64
	cpu_avg float64
	cpu_counter_vect int
	cpu_counter_smpl int
	cpu_vector_hour [config.CPU_AVG_SAMPLES]float64
	cpu_sum float64
	mu sync.RWMutex
}

// Parse output data from 'ps' command for server's process
func (cp* cpu_h) get_process_data() int {
	instr := "ps -f -o %cpu | grep "+config.PROCESS_NAME 
	cmd := exec.Command("bash", "-c", instr)

	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return -1
	}
	res := string(stdout)
	// fmt.Println(res)
	words := strings.Fields(res)
	
	if len(words)>0 {
		cp.pid = words[1]
		cp.process_name = words[7]
		cp.process_time = words[6]
		cpu ,err := strconv.ParseFloat(words[8],32)
		if err!=nil {
			fmt.Println(err)
			return -1
		}
		cp.cpu_usage = cpu
	}
	return -1
}

// CPU usage of server for the last hour, sampled at defined period
// In the first hour average is until sampling moment, then on last hour
func get_cpu_hour() {
	var avg float64

	time.Sleep(config.CPU_AVG_FREQ * time.Millisecond)
	Cpu_avgh_.cpu_counter_smpl++
	Cpu_avgh_.cpu_counter_vect++
	// Get cpu usage from "ps" command
	Cpu_avgh_.get_process_data()
	// Fill vector until is completed
	if Cpu_avgh_.cpu_counter_smpl < config.CPU_AVG_SAMPLES {
		Cpu_avgh_.cpu_vector_hour[Cpu_avgh_.cpu_counter_vect] = Cpu_avgh_.cpu_usage 
		Cpu_avgh_.cpu_sum += Cpu_avgh_.cpu_usage
		avg = Cpu_avgh_.cpu_sum / float64(Cpu_avgh_.cpu_counter_vect)
		Cpu_avgh_.set_cpu_avg(avg)
	} else {
		Cpu_avgh_.cpu_counter_vect = 0
		out := Cpu_avgh_.cpu_vector_hour[Cpu_avgh_.cpu_counter_vect]
		Cpu_avgh_.cpu_sum -= out
		Cpu_avgh_.cpu_sum += Cpu_avgh_.cpu_usage
		Cpu_avgh_.cpu_vector_hour[Cpu_avgh_.cpu_counter_vect] = Cpu_avgh_.cpu_usage 		
		avg = Cpu_avgh_.cpu_sum / config.CPU_AVG_SAMPLES
		Cpu_avgh_.set_cpu_avg(avg)
	}
}

func (cp* cpu_h) set_cpu_avg(data float64) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	cp.cpu_avg = data
}

func (cp* cpu_h) get_cpu_avg() float64 {
	cp.mu.RLock()
	defer cp.mu.RUnlock()
	return cp.cpu_avg
}

// Start cpu average computation
func Start_cpu_avg() float64 {
	for {
		get_cpu_hour()		
	}
}
