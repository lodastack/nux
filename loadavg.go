package nux

import (
	"fmt"
)

type Loadavg struct {
	Avg1min  float64
	Avg5min  float64
	Avg15min float64
}

func (this *Loadavg) String() string {
	return fmt.Sprintf("<1min:%f, 5min:%f, 15min:%f>", this.Avg1min, this.Avg5min, this.Avg15min)
}
