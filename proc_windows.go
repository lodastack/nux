package nux

import (
	"fmt"
)

func Procs(cmdlines map[string]string) (ps []*Proc, err error) {
	return nil, fmt.Errorf("windows not support proc")
}
