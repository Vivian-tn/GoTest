package cgroup

import (
	"log"
	"syscall"
)

const maxInt = int(^uint(0) >> 1)

func TotalMemory() int {
	var si syscall.Sysinfo_t
	if err := syscall.Sysinfo(&si); err != nil {
		log.Fatalf("syscall.Sysinfo failed %s", err)
	}
	totalMem := maxInt
	if uint64(maxInt)/uint64(si.Totalram) > uint64(si.Unit) {
		totalMem = int(uint64(si.Totalram) * uint64(si.Unit))
	}
	mem := getMemoryLimit()
	if mem <= 0 || int64(int(mem)) != mem || int(mem) > totalMem {
		// Try reading hierachical memory limit.
		// See https://github.com/VictoriaMetrics/VictoriaMetrics/issues/699
		mem = getHierarchicalMemoryLimit()
		if mem <= 0 || int64(int(mem)) != mem || int(mem) > totalMem {
			return totalMem
		}
	}
	return int(mem)
}

// GetMemoryLimit returns cgroup memory limit
func getMemoryLimit() int64 {
	// Try determining the amount of memory inside docker container.
	// See https://stackoverflow.com/questions/42187085/check-mem-limit-within-a-docker-container
	//
	// Read memory limit according to https://unix.stackexchange.com/questions/242718/how-to-find-out-how-much-memory-lxc-container-is-allowed-to-consume
	// This should properly determine the limit inside lxc container.
	// See https://github.com/VictoriaMetrics/VictoriaMetrics/issues/84
	n, err := readInt64("/sys/fs/cgroup/memory/memory.limit_in_bytes", "cat /sys/fs/cgroup/memory$(cat /proc/self/cgroup | grep memory | cut -d: -f3)/memory.limit_in_bytes")
	if err != nil {
		return 0
	}
	return n
}

// GetHierarchicalMemoryLimit returns hierarchical memory limit
func getHierarchicalMemoryLimit() int64 {
	// See https://github.com/VictoriaMetrics/VictoriaMetrics/issues/699
	n, err := readInt64FromCommand("cat /sys/fs/cgroup/memory/memory.stat | grep hierarchical_memory_limit | cut -d' ' -f 2")
	if err == nil {
		return n
	}
	n, err = readInt64FromCommand(
		"cat /sys/fs/cgroup/memory$(cat /proc/self/cgroup | grep memory | cut -d: -f3)/memory.stat | grep hierarchical_memory_limit | cut -d' ' -f 2")
	if err != nil {
		return 0
	}
	return n
}
