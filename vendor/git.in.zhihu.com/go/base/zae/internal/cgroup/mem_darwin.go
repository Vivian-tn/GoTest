package cgroup

import "log"

// fork from github.com/VictoriaMetrics/VictoriaMetrics/lib/memory

// This has been adapted from github.com/pbnjay/memory.
func TotalMemory() int {
	s, err := sysctlUint64("hw.memsize")
	if err != nil {
		log.Fatalf("cannot determine system memory: %s", err)
	}
	return int(s)
}
