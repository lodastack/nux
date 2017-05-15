package nux

import (
	"github.com/toolkits/file"
	"strconv"
	"strings"
)

func LoadAvg() (*Loadavg, error) {

	loadAvg := Loadavg{}

	data, err := file.ToTrimString("/proc/loadavg")
	if err != nil {
		return nil, err
	}

	L := strings.Fields(data)
	if loadAvg.Avg1min, err = strconv.ParseFloat(L[0], 64); err != nil {
		return nil, err
	}
	if loadAvg.Avg5min, err = strconv.ParseFloat(L[1], 64); err != nil {
		return nil, err
	}
	if loadAvg.Avg15min, err = strconv.ParseFloat(L[2], 64); err != nil {
		return nil, err
	}

	return &loadAvg, nil
}
