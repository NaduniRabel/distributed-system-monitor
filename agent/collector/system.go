package collector

import (
	"fmt"

	"github.com/shirou/gopsutil/v4/mem"
)

func CollectMetrics() {
	v, err := mem.VirtualMemory()

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)

}
