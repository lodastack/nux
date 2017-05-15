package nux

import (
	"fmt"
)

func SocketStatSummary() (m map[string]uint64, err error) {
	return nil, fmt.Errorf("windows not support socket stat")
}
