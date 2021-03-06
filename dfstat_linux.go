package nux

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"math"
	"strings"
	"syscall"

	"github.com/toolkits/file"
)

// return: [][$fs_spec, $fs_file, $fs_vfstype, $fs_rw]
func ListMountPoint() ([][4]string, error) {
	contents, err := ioutil.ReadFile("/proc/mounts")
	if err != nil {
		return nil, err
	}

	ret := make([][4]string, 0)

	reader := bufio.NewReader(bytes.NewBuffer(contents))
	for {
		line, err := file.ReadLine(reader)
		if err == io.EOF {
			err = nil
			break
		} else if err != nil {
			return nil, err
		}

		fields := strings.Fields(string(line))
		// Docs come from the fstab(5)
		// fs_spec     # Mounted block special device or remote filesystem e.g. /dev/sda1
		// fs_file     # Mount point e.g. /data
		// fs_vfstype  # File system type e.g. ext4
		// fs_mntops   # Mount options
		// fs_freq     # Dump(8) utility flags
		// fs_passno   # Order in which filesystem checks are done at reboot time

		fs_spec := fields[0]
		fs_file := fields[1]
		fs_vfstype := fields[2]
		fs_mntops := strings.Split(fields[3], ",")
		// default value
		fs_rw := "rw"
		if len(fs_mntops) > 0 {
			fs_rw = fs_mntops[0]
		}

		if _, exist := FSSPEC_IGNORE[fs_spec]; exist {
			continue
		}

		if _, exist := FSTYPE_IGNORE[fs_vfstype]; exist {
			continue
		}

		if strings.HasPrefix(fs_vfstype, "fuse") {
			continue
		}

		if IgnoreFsFile(fs_file) {
			continue
		}

		// keep /dev/xxx device with shorter fs_file (remove mount binds)
		if strings.HasPrefix(fs_spec, "/dev") {
			deviceFound := false
			for idx := range ret {
				if ret[idx][0] == fs_spec {
					deviceFound = true
					if len(fs_file) < len(ret[idx][1]) {
						ret[idx][1] = fs_file
					}
					break
				}
			}
			if !deviceFound {
				ret = append(ret, [4]string{fs_spec, fs_file, fs_vfstype, fs_rw})
			}
		} else {
			ret = append(ret, [4]string{fs_spec, fs_file, fs_vfstype, fs_rw})
		}
	}
	return ret, nil
}

func BuildDeviceUsage(_fsSpec, _fsFile, _fsVfstype string) (*DeviceUsage, error) {
	ret := &DeviceUsage{FsSpec: _fsSpec, FsFile: _fsFile, FsVfstype: _fsVfstype}

	fs := syscall.Statfs_t{}
	err := syscall.Statfs(_fsFile, &fs)
	if err != nil {
		return nil, err
	}

	// blocks
	used := fs.Blocks - fs.Bfree
	ret.BlocksAll = uint64(fs.Frsize) * fs.Blocks / Multi / Multi
	ret.BlocksUsed = uint64(fs.Frsize) * used / Multi / Multi
	ret.BlocksFree = uint64(fs.Frsize) * fs.Bavail / Multi / Multi
	if fs.Blocks == 0 {
		ret.BlocksUsedPercent = 0.0
	} else {
		ret.BlocksUsedPercent = math.Ceil(float64(used) * 100.0 / float64(used+fs.Bavail))
	}
	ret.BlocksFreePercent = 100.0 - ret.BlocksUsedPercent

	// inodes
	ret.InodesAll = fs.Files
	ret.InodesFree = fs.Ffree
	ret.InodesUsed = fs.Files - fs.Ffree
	if fs.Files == 0 {
		ret.InodesUsedPercent = 0.0
	} else {
		ret.InodesUsedPercent = math.Ceil(float64(ret.InodesUsed) * 100.0 / float64(ret.InodesAll))
	}
	ret.InodesFreePercent = 100.0 - ret.InodesUsedPercent

	return ret, nil
}
