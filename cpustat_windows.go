package nux

import (
	"fmt"

	"github.com/shirou/gopsutil/cpu"
)

func CurrentProcStat() (*ProcStat, error) {
	cpuStat, err := cpu.Times(false)
	if err != nil {
		return nil, err
	}
	if len(cpuStat) < 1 {
		return nil, fmt.Errorf("can not found CPU")
	}
	ps := &ProcStat{Cpus: make([]*CpuUsage, 1)}
	c := cpuStat[0]
	ps.Cpu = &CpuUsage{}
	ps.Cpu.User = uint64(c.User)
	ps.Cpu.Idle = uint64(c.Idle)
	ps.Cpu.Nice = uint64(c.Nice)
	ps.Cpu.System = uint64(c.System)
	ps.Cpu.Irq = uint64(c.Irq)
	ps.Cpu.Iowait = uint64(c.Iowait)
	ps.Cpu.Steal = uint64(c.Steal)
	ps.Cpu.Guest = uint64(c.Guest)
	ps.Cpu.SoftIrq = uint64(c.Softirq)
	ps.Cpu.Total = uint64(c.User + c.Idle + c.Nice + c.System +
		c.Softirq + c.Irq + c.Iowait + c.Steal +
		c.Guest)
	//TODO: per core
	ps.Cpus[0] = ps.Cpu

	return ps, nil
}
