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

func MemInfo() (*Mem, error) {
	contents, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		return nil, err
	}

	memInfo := &Mem{}

	reader := bufio.NewReader(bytes.NewBuffer(contents))

	for {
		line, err := file.ReadLine(reader)
		if err == io.EOF {
			err = nil
			break
		} else if err != nil {
			return nil, err
		}

		fields := strings.Fields(string(line))
		fieldName := fields[0]

		_, ok := WANT[fieldName]
		if ok && len(fields) == 3 {
			val, numerr := strconv.ParseUint(fields[1], 10, 64)
			if numerr != nil {
				continue
			}
			switch fieldName {
			case "Buffers:":
				memInfo.Buffers = val / Multi
			case "Cached:":
				memInfo.Cached = val / Multi
			case "MemTotal:":
				memInfo.MemTotal = val / Multi
			case "MemFree:":
				memInfo.MemFree = val / Multi
			case "SwapTotal:":
				memInfo.SwapTotal = val / Multi
			case "SwapFree:":
				memInfo.SwapFree = val / Multi
			}
		}
	}

	memInfo.SwapUsed = memInfo.SwapTotal - memInfo.SwapFree

	return memInfo, nil
}
