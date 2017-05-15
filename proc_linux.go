package nux

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/toolkits/file"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
)

func Procs(cmdlines map[string]string) (ps []*Proc, err error) {
	var dirs []string
	dirs, err = file.DirsUnder("/proc")
	if err != nil {
		return
	}

	size := len(dirs)
	if size == 0 {
		return
	}

	for i := 0; i < size; i++ {
		pid, e := strconv.Atoi(dirs[i])
		if e != nil {
			continue
		}
		statusFile := fmt.Sprintf("/proc/%d/status", pid)
		exeFile := fmt.Sprintf("/proc/%d/exe", pid)
		if !file.IsExist(statusFile) || !file.IsExist(exeFile) {
			continue
		}

		name, memory, e := ReadNameAndMem(statusFile)
		if e != nil {
			continue
		}

		exe, e := os.Readlink(exeFile)
		if e != nil {
			continue
		}

		if _, ok := cmdlines[exe]; !ok {
			continue
		}

		p := Proc{Pid: pid, Name: name, Exe: exe, Mem: memory}
		ps = append(ps, &p)
	}

	jiffyTotal := readJiffy()

	for _, p := range ps {
		jiffy := readProcJiffy(p.Pid)
		p.Cpu = float64(jiffy)
		p.TotalCpu = float64(jiffyTotal)
	}
	return
}

func ReadNameAndMem(path string) (name string, memory uint64, err error) {
	var content []byte
	content, err = ioutil.ReadFile(path)
	if err != nil {
		return
	}

	reader := bufio.NewReader(bytes.NewBuffer(content))

	for {
		var bs []byte
		bs, err = file.ReadLine(reader)
		if err == io.EOF {
			return
		}

		line := string(bs)
		colonIndex := strings.Index(line, ":")

		if strings.TrimSpace(line[0:colonIndex]) == "Name" {
			name = strings.TrimSpace(line[colonIndex+1:])
		} else if strings.TrimSpace(line[0:colonIndex]) == "VmRSS" {
			kbIndex := strings.Index(line, "kB")
			memory, _ = strconv.ParseUint(strings.TrimSpace(line[colonIndex+1:kbIndex]), 10, 64)
			break
		}

	}
	return
}

func readJiffy() uint64 {
	f, err := os.Open("/proc/stat")
	if err != nil {
		return 0
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Scan()
	s := scanner.Text()
	if !strings.HasPrefix(s, "cpu ") {
		return 0
	}
	ss := strings.Split(s, " ")
	var ret uint64
	for _, x := range ss {
		if x == "" || x == "cpu" {
			continue
		}
		if v, e := strconv.ParseUint(x, 10, 64); e == nil {
			ret += v
		}
	}
	return ret
}

func readProcFd(pid int) int {
	var fds []string
	fds, err := file.FilesUnder(fmt.Sprintf("/proc/%d/fd", pid))
	if err != nil {
		return 0
	}
	return len(fds)
}

func readProcJiffy(pid int) uint64 {
	f, err := os.Open(fmt.Sprintf("/proc/%d/stat", pid))
	if err != nil {
		return 0
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Scan()
	s := scanner.Text()
	ss := strings.Split(s, " ")
	var ret uint64
	for i := 13; i < 15; i++ {
		v, e := strconv.ParseUint(ss[i], 10, 64)
		if e != nil {
			return 0
		}
		ret += v
	}
	return ret
}

func readTcp() map[uint64]bool {
	res := make(map[uint64]bool)
	f, err := os.Open("/proc/net/tcp")
	if err != nil {
		return res
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Scan()
	for scanner.Scan() {
		b := scanner.Bytes()
		if len(b) != 149 || b[34] != 48 || b[35] != 49 { //only established
			continue
		}
		start := 91
		end := start
		for end = start; b[end] != 32; end++ {
		}
		inode, err := strconv.ParseUint(string(b[start:end]), 10, 64)
		if err != nil {
			continue
		}
		res[inode] = true
	}
	return res
}

func tcpEstablishCount(inodes map[uint64]bool, pid int) int {
	res := 0
	dir := fmt.Sprintf("/proc/%d/fd", pid)
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		return res
	}
	for _, fi := range fis {
		link, err := os.Readlink(path.Join(dir, fi.Name()))
		if err != nil || !strings.HasPrefix(link, "socket:") {
			continue
		}
		s := link[8 : len(link)-1]
		inode, err := strconv.ParseUint(s, 10, 64)
		if err == nil && inodes[inode] {
			res++
		}
	}
	return res
}

func readIO(pid int) (r uint64, w uint64) {
	f, err := os.Open(fmt.Sprintf("/proc/%d/io", pid))
	if err != nil {
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		s := scanner.Text()
		if strings.HasPrefix(s, "read_bytes") || strings.HasPrefix(s, "write_bytes") {
			v := strings.Split(s, " ")
			if len(v) == 2 {
				value, _ := strconv.ParseUint(v[1], 10, 64)
				if s[0] == 'r' {
					r = value
				} else {
					w = value
				}
			}
		}
	}
	return
}
