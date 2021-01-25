package main

type Frame struct {
	Name     string `json:"name"`
	Addr     string `json:"addr"`
	RAMTotal uint64 `json:"ramTotal"`
	RAMUsage uint64 `json:"ramUsage"`
}

func (f Frame) MemoryUsage() float64 {
	return float64(f.RAMUsage) / float64(f.RAMTotal)
}
