package nux

import (
	"github.com/shirou/gopsutil/mem"
)

func MemInfo() (*Mem, error) {
	memStat, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	memInfo := &Mem{}
	memInfo.MemTotal = memStat.Total / Multi / Multi
	memInfo.MemFree = memStat.Available / Multi / Multi
	memInfo.Buffers = memStat.Buffers / Multi / Multi
	memInfo.Cached = memStat.Cached / Multi / Multi

	return memInfo, nil
}
