package nux

import (
	"fmt"
)

type NetIf struct {
	Iface          string
	InBytes        uint64
	InPackages     uint64
	InErrors       uint64
	InDropped      uint64
	InFifoErrs     uint64
	InFrameErrs    uint64
	InCompressed   uint64
	InMulticast    uint64
	OutBytes       uint64
	OutPackages    uint64
	OutErrors      uint64
	OutDropped     uint64
	OutFifoErrs    uint64
	OutCollisions  uint64
	OutCarrierErrs uint64
	OutCompressed  uint64
	TotalBytes     uint64
	TotalPackages  uint64
	TotalErrors    uint64
	TotalDropped   uint64
	Speed          uint64
}

func (this *NetIf) String() string {
	return fmt.Sprintf("<Iface:%s,InBytes:%d,InPackages:%d...>", this.Iface, this.InBytes, this.InPackages)
}
