package backend

import (
	"fmt"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

func CpuMetrics() string {
	totalPercent, _ := cpu.Percent(1*time.Second, false)
	onlyPercent := strconv.FormatFloat(totalPercent[0], 'f', 2, 64)

	return onlyPercent
}

func MemoryMetrics() string {
	v, _ := mem.VirtualMemory()

	return strconv.FormatFloat(v.UsedPercent, 'f', 2, 64)
}

func BackendInit() {
	fmt.Println("CPU: " + CpuMetrics() + "%")

	fmt.Println()

	fmt.Println("VM: " + MemoryMetrics() + "%")
}
