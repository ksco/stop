package frame

import (
	"errors"
	"qtop/remote"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type processItem struct {
	User       string  `json:"user"`
	CPUPercent float64 `json:"cpuPercent"`
	RSS        uint64  `json:"rss"`
	Cmd        string  `json:"cmd"`
}

type diskItem struct {
	FileSystem string `json:"fileSystem"`
	Total      uint64 `json:"total"`
	Usage      uint64 `json:"usage"`
}

type load struct {
	One     float64 `json:"one"`
	Five    float64 `json:"five"`
	Fifteen float64 `json:"fifteen"`
}

type Frame struct {
	Name             string         `json:"name"`
	Addr             string         `json:"addr"`
	Uptime           time.Duration  `json:"uptime"`
	SessionsCount    uint64         `json:"sessionsCount"`
	ProcessesCount   uint64         `json:"processesCount"`
	Processes        []*processItem `json:"processes"`
	FileHandlesCount uint64         `json:"fileHandlesCount"`
	FileHandlesLimit uint64         `json:"fileHandlesLimit"`
	OsKernel         string         `json:"osKernel"`
	OsName           string         `json:"osName"`
	OsArch           string         `json:"osArch"`
	CPUName          string         `json:"cpuName"`
	CPUCores         uint64         `json:"cpuCores"`
	CPUFreq          float64        `json:"cpuFreq"`
	RAMTotal         uint64         `json:"ramTotal"`
	RAMUsage         uint64         `json:"ramUsage"`
	SwapTotal        uint64         `json:"swapTotal"`
	SwapUsage        uint64         `json:"swapUsage"`
	Disks            []*diskItem    `json:"disks"`
	DiskTotal        uint64         `json:"diskTotal"`
	DiskUsage        uint64         `json:"diskUsage"`
	ConnectionsCount uint64         `json:"connectionsCount"`
	Load             *load          `json:"load"`
}

func NewFrame(c remote.Client) (*Frame, error) {
	f := &Frame{Name: c.Alias, Addr: c.Addr}
	frameType := reflect.TypeOf(f)
	for i := 0; i < frameType.NumMethod(); i++ {
		method := frameType.Method(i)
		if strings.HasPrefix(method.Name, "Fill") {
			errValue := method.Func.Call([]reflect.Value{reflect.ValueOf(f), reflect.ValueOf(c)})
			if !errValue[0].IsNil() {
				return nil, errors.New("get frame failed")
			}
		}
	}
	return f, nil
}

func (f *Frame) FillUptime(c remote.Client) (err error) {
	r, err := c.FloatValue("cat /proc/uptime | awk '{ print $1 }'")
	if err != nil {
		return err
	}
	f.Uptime = time.Duration(r * 1e9)
	return
}

func (f *Frame) FillSessionsCount(c remote.Client) (err error) {
	f.SessionsCount, err = c.UintValue("who | wc -l")
	return
}

func (f *Frame) FillProcesses(c remote.Client) error {
	output, err := c.RunCommand("ps axc -o uname:12,pcpu,rss,cmd --sort=-pcpu,-rss --noheaders --width 120 | grep -v ' ps$' | sed 's/ \\+ / /g' | sed '/^$/d' | tr '\n' ';'")
	if err != nil {
		return err
	}
	var ps []*processItem
	processItemStrs := strings.Split(output, ";")
	for _, processItemStr := range processItemStrs[:10] {
		parts := strings.SplitN(processItemStr, " ", 4)
		if len(parts) != 4 {
			continue
		}
		cpuPercent, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return err
		}
		rss, err := strconv.ParseUint(parts[2], 10, 64)
		if err != nil {
			return err
		}

		ps = append(ps, &processItem{
			User:       parts[0],
			CPUPercent: cpuPercent,
			RSS:        rss,
			Cmd:        parts[3],
		})
	}

	f.Processes = ps
	f.ProcessesCount = uint64(len(processItemStrs))
	return nil
}

func (f *Frame) FillFileHandlesCount(c remote.Client) (err error) {
	f.FileHandlesCount, err = c.UintValue("cat /proc/sys/fs/file-nr | awk '{ print $1 }'")
	return
}

func (f *Frame) FillFileHandlesLimit(c remote.Client) (err error) {
	f.FileHandlesLimit, err = c.UintValue("cat /proc/sys/fs/file-nr | awk '{ print $3 }'")
	return
}

func (f *Frame) FillOsKernel(c remote.Client) (err error) {
	f.OsKernel, err = c.StringValue("uname -r")
	return
}

func (f *Frame) FillOsName(c remote.Client) (err error) {
	f.OsName, err = c.StringValue("uname -s")
	return
}

func (f *Frame) FillOsArch(c remote.Client) (err error) {
	f.OsArch, err = c.StringValue("uname -m")
	return
}

func (f *Frame) FillCPUName(c remote.Client) (err error) {
	f.CPUName, err = c.StringValue("cat /proc/cpuinfo | grep 'model name' | awk -F\\: '{ print $2 }'")
	return
}

func (f *Frame) FillCPUCores(c remote.Client) error {
	r, err := c.UintValue("cat /proc/cpuinfo | grep 'model name' | awk -F\\: '{ print $2 }' | sed -e :a -e '$!N;s/\\n/\\|/;ta' | tr -cd \\| | wc -c")
	if err != nil {
		return err
	}
	f.CPUCores = r + 1
	return nil
}

func (f *Frame) FillCPUFreq(c remote.Client) (err error) {
	f.CPUFreq, err = c.FloatValue("cat /proc/cpuinfo | grep 'cpu MHz' | awk -F\\: '{ print $2 }'")
	return
}

func (f *Frame) FillRAMInfo(c remote.Client) error {
	ramTotal, err := c.FloatValue("cat /proc/meminfo | grep ^MemTotal: | awk '{ print $2 }'")
	if err != nil {
		return err
	}
	ramFree, err := c.FloatValue("cat /proc/meminfo | grep ^MemFree: | awk '{ print $2 }'")
	if err != nil {
		return err
	}
	ramCached, err := c.FloatValue("cat /proc/meminfo | grep ^Cached: | awk '{ print $2 }'")
	if err != nil {
		return err
	}
	ramBuffers, err := c.FloatValue("cat /proc/meminfo | grep ^Buffers: | awk '{ print $2 }'")
	if err != nil {
		return err
	}
	f.RAMUsage = uint64((ramTotal - (ramFree + ramCached + ramBuffers)) * 1024)
	f.RAMTotal = uint64(ramTotal * 1024)
	return nil
}

func (f *Frame) FillSwapInfo(c remote.Client) error {
	swapTotal, err := c.FloatValue("cat /proc/meminfo | grep ^SwapTotal: | awk '{ print $2 }'")
	if err != nil {
		return err
	}
	swapFree, err := c.FloatValue("cat /proc/meminfo | grep ^SwapFree: | awk '{ print $2 }'")
	if err != nil {
		return err
	}
	f.SwapUsage = uint64((swapTotal - swapFree) * 1024)
	f.SwapTotal = uint64(swapTotal * 1024)
	return nil
}

func (f *Frame) FillDiskInfo(c remote.Client) error {
	output, err := c.RunCommand("df -P -B 1 | grep '^/' | awk '{ print $1\" \"$2\" \"$3\";\" }' | sed -e :a -e '$!N;s/\\n/ /;ta' | awk '{ print $0 } END { if (!NR) print 'N/A' }'")
	if err != nil {
		return err
	}
	var ds []*diskItem
	diskItemStrs := strings.Split(output, ";")
	var diskTotal, diskUsage uint64
	for _, diskItemStr := range diskItemStrs {
		parts := strings.SplitN(diskItemStr, " ", 3)
		if len(parts) != 3 {
			continue
		}
		total, err := strconv.ParseUint(parts[1], 10, 64)
		if err != nil {
			return err
		}
		diskTotal += total
		usage, err := strconv.ParseUint(parts[2], 10, 64)
		if err != nil {
			return err
		}
		diskUsage += usage

		ds = append(ds, &diskItem{
			FileSystem: parts[0],
			Total:      total,
			Usage:      usage,
		})
	}
	f.Disks = ds
	f.DiskTotal = diskTotal
	f.DiskUsage = diskUsage
	return nil
}

func (f *Frame) FillConnectionsCount(c remote.Client) (err error) {
	f.ConnectionsCount, err = c.UintValue("netstat -tun | tail -n +3 | wc -l")
	return
}

func (f *Frame) FillLoad(c remote.Client) error {
	output, err := c.RunCommand("cat /proc/loadavg | awk '{ printf $1\" \"$2\" \"$3 }'")
	if err != nil {
		return err
	}
	parts := strings.SplitN(output, " ", 3)
	if len(parts) != 3 {
		return errors.New("parse load info failed")
	}

	l := &load{}
	l.One, err = strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return err
	}
	l.Five, err = strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return err
	}
	l.Fifteen, err = strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return err
	}
	f.Load = l
	return nil
}
