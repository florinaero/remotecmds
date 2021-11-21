package server

import (
	"fmt"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/florinaero/remotecmds/pkg/config"
	// "strings"
)

type ram_h struct {
	ram_avg_h   float64
	ram_sum 	uint
	ram_mtx_avg sync.RWMutex
	ram_vector [config.RAM_AVG_SAMPLES]uint
	ram_cnt_reset uint
	ram_cnt_total uint
}

// Return free RAM in bytes using command vm_stat
func (rm *ram_h) get_free_ram() uint {
	cmd := exec.Command("vm_stat")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return 0
	}
	res := string(out)
	reg_ps := regexp.MustCompile("page size of ([0-9]*) bytes")
	reg_pf := regexp.MustCompile("Pages free: *([0-9]*)\\.")
	ps := reg_ps.FindStringSubmatch(res)
	pf := reg_pf.FindStringSubmatch(res)

	page_size, err1 := strconv.Atoi(pf[1])
	page_free, err2 := strconv.Atoi(ps[1])
	if err1 != nil || err2 != nil {
		return 0
	}
	free_ram := uint(page_size * page_free)

	return free_ram
}

// Return average of last hour of free RAM, sampled defined in config file
func (rm *ram_h) get_free_ram_avg() float64 {
	var avg float64
	var samples uint 

	// Sampling period 
	time.Sleep(config.RAM_AVG_FREQ * time.Millisecond)
	
	// Count samples until achieves the desired number
	if rm.ram_cnt_total < config.RAM_AVG_SAMPLES {
		samples = rm.ram_cnt_reset + 1
	} else {
		samples = config.RAM_AVG_SAMPLES
	}
	free_ram := rm.get_free_ram()
	old := rm.ram_vector[rm.ram_cnt_reset]
	rm.ram_vector[rm.ram_cnt_reset] = free_ram
	rm.ram_sum -= old
	rm.ram_sum += free_ram
	avg = float64(rm.ram_sum / samples)
	rm.ram_cnt_total++

	// Reset position counter for updating old values in vector
	if rm.ram_cnt_reset < config.RAM_AVG_SAMPLES-1 {
		rm.ram_cnt_reset++
	} else {
		rm.ram_cnt_reset = 0
	}
	rm.set_ram_avg(avg)
	return avg
}

func (rm *ram_h) start_ram_avg() {
	for {
		rm.get_free_ram_avg()
	}
}

// 5. Available RAM over last hour
func (hc *Server) ram_hour_usage(w http.ResponseWriter, req *http.Request) {
	hc.grab_counter(w)
	avg := hc.ram_usage.get_ram_avg() * 1e-9
	fmt.Fprintf(w, "Available RAM over last hour is: %3.2f GB\n", avg)
	hc.release_counter()
}

func (rm *ram_h) set_ram_avg(data float64) {
	rm.ram_mtx_avg.Lock()
	defer rm.ram_mtx_avg.Unlock()
	rm.ram_avg_h = data
}

func (rm *ram_h) get_ram_avg() float64 {
	rm.ram_mtx_avg.RLock()
	defer rm.ram_mtx_avg.RUnlock()
	return rm.ram_avg_h
}
