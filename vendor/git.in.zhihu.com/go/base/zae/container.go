package zae

import (
	"runtime"

	"git.in.zhihu.com/go/base/zae/internal/cgroup"
)

var totalCPU, totalMemory = cgroup.TotalCPU(), cgroup.TotalMemory()

func init() {
	runtime.GOMAXPROCS(totalCPU)
}

func TotalCPU() int {
	return totalCPU
}

func TotalMemory() int {
	return totalMemory
}
