package config

const PROCESS_NAME = "./server_8080"
const CPU_AVG_FREQ = 10 // sampling frequency in miliseconds
const CPU_AVG_SAMPLES = 3600/CPU_AVG_FREQ*1000
const RAM_AVG_FREQ = 10 // sampling frequency in miliseconds
const RAM_AVG_SAMPLES = uint(float64(3600)/RAM_AVG_FREQ*1000)
