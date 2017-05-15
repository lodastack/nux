package nux

import (
	"fmt"
	"os"
)

func KernelMaxFiles() (uint64, error) {
	return 0, fmt.Errorf("windows not support KernelMaxFiles")
}

func KernelAllocateFiles() (ret uint64, err error) {
	return 0, fmt.Errorf("windows not support KernelAllocateFiles")
}

func KernelMaxProc() (uint64, error) {
	return 0, fmt.Errorf("windows not support KernelMaxProc")
}

func KernelHostname() (string, error) {
	return os.Hostname()
}
