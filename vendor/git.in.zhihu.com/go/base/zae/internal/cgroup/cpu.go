package cgroup

import (
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"
)

// fork from github.com/VictoriaMetrics/VictoriaMetrics/lib/cgroup

func TotalCPU() int {
	var cpu float64
	if v := os.Getenv("GOMAXPROCS"); v != "" {
		cpu, _ = strconv.ParseFloat(v, 64)
	}
	if cpu <= 0 {
		cpu = getCPUQuota()
	}
	if cpu <= 0 {
		cpu = float64(runtime.NumCPU())
	}
	return int(cpu + 0.5)
}

func getCPUQuota() float64 {
	quotaUS, err := readInt64("/sys/fs/cgroup/cpu/cpu.cfs_quota_us", "cat /sys/fs/cgroup/cpu$(cat /proc/self/cgroup | grep cpu, | cut -d: -f3)/cpu.cfs_quota_us")
	if err != nil {
		return 0
	}
	if quotaUS <= 0 {
		// The quota isn't set. This may be the case in multilevel containers.
		// See https://github.com/VictoriaMetrics/VictoriaMetrics/issues/685#issuecomment-674423728
		return getOnlineCPUCount()
	}
	periodUS, err := readInt64("/sys/fs/cgroup/cpu/cpu.cfs_period_us", "cat /sys/fs/cgroup/cpu$(cat /proc/self/cgroup | grep cpu, | cut -d: -f3)/cpu.cfs_period_us")
	if err != nil {
		return 0
	}
	return float64(quotaUS) / float64(periodUS)
}

func getOnlineCPUCount() float64 {
	// See https://github.com/VictoriaMetrics/VictoriaMetrics/issues/685#issuecomment-674423728
	data, err := ioutil.ReadFile("/sys/devices/system/cpu/online")
	if err != nil {
		return -1
	}
	n := float64(countCPUs(string(data)))
	if n <= 0 {
		return -1
	}
	// Add a half of CPU core, since it looks like actual cores is usually bigger than online cores.
	// See https://github.com/VictoriaMetrics/VictoriaMetrics/issues/685#issuecomment-674423728
	return n + 0.5
}

func countCPUs(data string) int {
	data = strings.TrimSpace(data)
	n := 0
	for _, s := range strings.Split(data, ",") {
		n++
		if !strings.Contains(s, "-") {
			if _, err := strconv.Atoi(s); err != nil {
				return -1
			}
			continue
		}
		bounds := strings.Split(s, "-")
		if len(bounds) != 2 {
			return -1
		}
		start, err := strconv.Atoi(bounds[0])
		if err != nil {
			return -1
		}
		end, err := strconv.Atoi(bounds[1])
		if err != nil {
			return -1
		}
		n += end - start
	}
	return n
}
