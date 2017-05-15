package nux

import (
	"fmt"
)

type Proc struct {
	Pid      int
	Name     string
	Exe      string
	Mem      uint64
	TotalCpu float64
	Cpu      float64

	RBytes   uint64
	WBytes   uint64
	TcpEstab int
	FdCount  int
}

func (this *Proc) String() string {
	return fmt.Sprintf("<Pid:%d, Name:%s, Exe:%s Mem:%d Cpu:%.3f>",
		this.Pid, this.Name, this.Exe, this.Mem, this.Cpu)
}
