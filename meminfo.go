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
	// MemAvailable is in /proc/meminfo (kernel 3.14+)
	// https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git/commit/?id=34e431b0ae398fc54ea69ff85ec700722c9da773
	// https://www.kernel.org/doc/Documentation/filesystems/proc.txt
	MemAvailable  uint64
	MemAvaSupport bool
}

func (this *Mem) String() string {
	return fmt.Sprintf("<MemTotal:%d, MemFree:%d, Buffers:%d, Cached:%d...>", this.MemTotal, this.MemFree, this.Buffers, this.Cached)
}

var Multi uint64 = 1024

var WANT = map[string]struct{}{
	"Buffers:":      {},
	"Cached:":       {},
	"MemTotal:":     {},
	"MemFree:":      {},
	"SwapTotal:":    {},
	"SwapFree:":     {},
	"MemAvailable:": {},
}
