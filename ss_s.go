package nux

import (
	"bufio"
	"bytes"
	"github.com/toolkits/file"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
)

func SocketStatSummary() (m map[string]uint64, err error) {
	m = make(map[string]uint64)
	bs, err := ioutil.ReadFile("/proc/net/sockstat")
	if err != nil {
		return
	}
	reader := bufio.NewReader(bytes.NewBuffer(bs))

	for {
		var lineBytes []byte
		lineBytes, err = file.ReadLine(reader)
		if err == io.EOF {
			return
		}
		line := string(lineBytes)
		s := strings.Split(line, " ")
		if strings.HasPrefix(line, "sockets: used") {
			m["sockets.used"], _ = strconv.ParseUint(s[2], 10, 64)
		} else {
			m["sockets.tcp.inuse"], _ = strconv.ParseUint(s[2], 10, 64)
			m["tcp.timewait"], _ = strconv.ParseUint(s[6], 10, 64)
			break
		}
	}

	return
}
