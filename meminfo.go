package nux

import (
	"fmt"
)

type Mem struct {
	Buffers   uint64
	Cached    uint64
	MemTotal  uint64
	MemFree   uint64
	SwapTotal uint64
	SwapUsed  uint64
	SwapFree  uint64
}

func (this *Mem) String() string {
	return fmt.Sprintf("<MemTotal:%d, MemFree:%d, Buffers:%d, Cached:%d...>", this.MemTotal, this.MemFree, this.Buffers, this.Cached)
}

var Multi uint64 = 1024

var WANT = map[string]struct{}{
	"Buffers:":   {},
	"Cached:":    {},
	"MemTotal:":  {},
	"MemFree:":   {},
	"SwapTotal:": {},
	"SwapFree:":  {},
}
