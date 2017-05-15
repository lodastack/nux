package nux

import (
	"bufio"
	"bytes"
	"github.com/toolkits/file"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
)

func CurrentProcStat() (*ProcStat, error) {
	f := "/proc/stat"
	bs, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}

	ps := &ProcStat{Cpus: make([]*CpuUsage, NumCpu())}
	reader := bufio.NewReader(bytes.NewBuffer(bs))

	for {
		line, err := file.ReadLine(reader)
		if err == io.EOF {
			err = nil
			break
		} else if err != nil {
			return ps, err
		}
		parseLine(line, ps)
	}

	return ps, nil
}

func parseLine(line []byte, ps *ProcStat) {
	fields := strings.Fields(string(line))
	if len(fields) < 3 {
		return
	}

	fieldName := fields[0]
	if fieldName == "cpu" {
		ps.Cpu = parseCpuFields(fields)
		return
	}

	if strings.HasPrefix(fieldName, "cpu") {
		idx, err := strconv.Atoi(fieldName[3:])
		if err != nil || idx >= len(ps.Cpus) {
			return
		}

		ps.Cpus[idx] = parseCpuFields(fields)
		return
	}

	if fieldName == "ctxt" {
		ps.Ctxt, _ = strconv.ParseUint(fields[1], 10, 64)
		return
	}

	if fieldName == "processes" {
		ps.Processes, _ = strconv.ParseUint(fields[1], 10, 64)
		return
	}

	if fieldName == "procs_running" {
		ps.ProcsRunning, _ = strconv.ParseUint(fields[1], 10, 64)
		return
	}

	if fieldName == "procs_blocked" {
		ps.ProcsBlocked, _ = strconv.ParseUint(fields[1], 10, 64)
		return
	}
}

func parseCpuFields(fields []string) *CpuUsage {
	cu := new(CpuUsage)
	sz := len(fields)
	for i := 1; i < sz; i++ {
		val, err := strconv.ParseUint(fields[i], 10, 64)
		if err != nil {
			continue
		}

		cu.Total += val
		switch i {
		case 1:
			cu.User = val
		case 2:
			cu.Nice = val
		case 3:
			cu.System = val
		case 4:
			cu.Idle = val
		case 5:
			cu.Iowait = val
		case 6:
			cu.Irq = val
		case 7:
			cu.SoftIrq = val
		case 8:
			cu.Steal = val
		case 9:
			cu.Guest = val
		}
	}
	return cu
}
