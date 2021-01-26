package main

type processItem struct {
	User       string  `json:"user"`
	CPUPercent float64 `json:"cpuPercent"`
	RSS        uint64  `json:"rss"`
	Cmd        string  `json:"cmd"`
}

type Frame struct {
	Name            string  `json:"name"`
	Addr            string  `json:"addr"`
	RAMTotal        uint64  `json:"ramTotal"`
	RAMUsage        uint64  `json:"ramUsage"`
	CPUUsagePercent float64 `json:"cpuUsagePercent"`
	DiskTotal       uint64  `json:"diskTotal"`
	DiskUsage       uint64  `json:"diskUsage"`
}

func (f *Frame) MemoryUsagePercent() float64 {
	return 100 * (float64(f.RAMUsage) / float64(f.RAMTotal))
}

func (f *Frame) DiskUsagePercent() float64 {
	return 100 * (float64(f.DiskUsage) / float64(f.DiskTotal))
}
