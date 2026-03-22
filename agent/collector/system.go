package collector

import (
	"time"

	"github.com/shirou/gopsutil/v4/cpu"

	"github.com/shirou/gopsutil/v4/mem"

	"github.com/shirou/gopsutil/v4/disk"
)

func CollectMetrics() (SystemMetrics, error) {
	/*Calling methods to get CPU, memory and disk data*/
	cpuMetrics, err := getCPUMetrics()
	if err != nil {
		return SystemMetrics{}, err
	}

	memMetrics, err := getMemoryMetrics()
	if err != nil {
		return SystemMetrics{}, err
	}

	diskMetrics, err := getDiskMetrics()
	if err != nil {
		return SystemMetrics{}, err
	}

	return SystemMetrics{
		CPU:    cpuMetrics,
		Memory: memMetrics,
		Disk:   diskMetrics,
	}, nil

}

/*Structs*/
type SystemMetrics struct {
	CPU    CPUMetrics
	Memory MemMetrics
	Disk   DiskMetrics
}
type CPUMetrics struct {
	Usage []float64
}

type MemMetrics struct {
	Total          uint64
	Free           uint64
	UsedPercentage float64
	Available      uint64
}

type DiskMetrics struct {
	Total uint64
	Used  uint64
	Free  uint64
}

func getCPUMetrics() (CPUMetrics, error) {

	/*get total CPU usage within 5 seconds*/
	totalPercent, err := cpu.Percent(5*time.Second, false)

	if err != nil {
		return CPUMetrics{}, err
	}

	return CPUMetrics{
		Usage: totalPercent,
	}, nil
}

func getMemoryMetrics() (MemMetrics, error) {

	/*Get memory stats*/
	v, err := mem.VirtualMemory()

	if err != nil {
		return MemMetrics{}, err
	}

	return MemMetrics{
		Total:          v.Total,
		Free:           v.Free,
		UsedPercentage: v.UsedPercent,
		Available:      v.Available,
	}, nil
}

func getDiskMetrics() (DiskMetrics, error) {

	/*Get total used and free root directory usage*/
	info, err := disk.Usage("/")

	if err != nil {
		return DiskMetrics{}, err
	}

	return DiskMetrics{
		Total: info.Total,
		Used:  info.Used,
		Free:  info.Free,
	}, nil
}
