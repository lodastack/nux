package nux

import (
	"github.com/shirou/gopsutil/disk"
)

// return: [][$fs_spec, $fs_file, $fs_vfstype]
func ListMountPoint() ([][4]string, error) {
	partitions, err := disk.Partitions(true)
	if err != nil {
		return nil, err
	}
	ret := make([][4]string, len(partitions))
	for i, p := range partitions {
		// CD driver in windows
		if p.Fstype == "CDFS" || p.Fstype == "UDF" {
			continue
		}
		ret[i][0] = p.Device
		ret[i][1] = p.Mountpoint
		ret[i][2] = p.Fstype
		ret[i][3] = p.Opts
	}
	return ret, nil
}

func BuildDeviceUsage(_fsSpec, _fsFile, _fsVfstype string) (*DeviceUsage, error) {
	usageStat, err := disk.Usage(_fsSpec)
	if err != nil {
		return nil, err
	}
	ret := &DeviceUsage{FsSpec: _fsSpec, FsFile: _fsFile, FsVfstype: _fsVfstype}
	ret.BlocksAll = usageStat.Total
	ret.BlocksUsed = usageStat.Used
	ret.BlocksFree = usageStat.Free
	ret.BlocksUsedPercent = usageStat.UsedPercent
	ret.BlocksFreePercent = 100 - usageStat.UsedPercent
	return ret, nil
}
