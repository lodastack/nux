package nux

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/toolkits/file"
)

const MegaCliBin = "/opt/MegaRAID/MegaCli/MegaCli"
const MegaCli64Bin = "/opt/MegaRAID/MegaCli/MegaCli64"

func ListDiskStats() ([]*DiskStats, error) {
	proc_diskstats := "/proc/diskstats"
	if !file.IsExist(proc_diskstats) {
		return nil, fmt.Errorf("%s not exists", proc_diskstats)
	}

	contents, err := ioutil.ReadFile(proc_diskstats)
	if err != nil {
		return nil, err
	}

	ret := make([]*DiskStats, 0)

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
		// shortcut the deduper and just skip disks that
		// haven't done a single read.  This elimiates a bunch
		// of loopback, ramdisk, and cdrom devices but still
		// lets us report on the rare case that we actually use
		// a ramdisk.
		if fields[3] == "0" {
			continue
		}

		size := len(fields)
		// kernel version too low
		if size != 14 {
			continue
		}

		item := &DiskStats{}
		if item.Major, err = strconv.Atoi(fields[0]); err != nil {
			return nil, err
		}

		if item.Minor, err = strconv.Atoi(fields[1]); err != nil {
			return nil, err
		}

		item.Device = fields[2]

		if item.ReadRequests, err = strconv.ParseUint(fields[3], 10, 64); err != nil {
			return nil, err
		}

		if item.ReadMerged, err = strconv.ParseUint(fields[4], 10, 64); err != nil {
			return nil, err
		}

		if item.ReadSectors, err = strconv.ParseUint(fields[5], 10, 64); err != nil {
			return nil, err
		}

		if item.MsecRead, err = strconv.ParseUint(fields[6], 10, 64); err != nil {
			return nil, err
		}

		if item.WriteRequests, err = strconv.ParseUint(fields[7], 10, 64); err != nil {
			return nil, err
		}

		if item.WriteMerged, err = strconv.ParseUint(fields[8], 10, 64); err != nil {
			return nil, err
		}

		if item.WriteSectors, err = strconv.ParseUint(fields[9], 10, 64); err != nil {
			return nil, err
		}

		if item.MsecWrite, err = strconv.ParseUint(fields[10], 10, 64); err != nil {
			return nil, err
		}

		if item.IosInProgress, err = strconv.ParseUint(fields[11], 10, 64); err != nil {
			return nil, err
		}

		if item.MsecTotal, err = strconv.ParseUint(fields[12], 10, 64); err != nil {
			return nil, err
		}

		if item.MsecWeightedTotal, err = strconv.ParseUint(fields[13], 10, 64); err != nil {
			return nil, err
		}

		item.TS = time.Now()
		ret = append(ret, item)
	}
	return ret, nil
}

func DiskHealth() (*DiskHealthStats, error) {
	res := &DiskHealthStats{}
	var ds []PDStats
	var rs []RaidStats
	var MegaCli string
	if Exists(MegaCli64Bin) {
		MegaCli = MegaCli64Bin
	} else if Exists(MegaCliBin) {
		MegaCli = MegaCliBin
	} else {
		return nil, nil
	}

	idpara := ` -PDList -aALL -NoLog | grep "Device Id:" | awk -F ": " '{print $2}'`
	idcmd := fmt.Sprintf("%s%s", MegaCli, idpara)
	ids, err := cmdCall(idcmd)
	if err != nil {
		return res, err
	}
	ds = make([]PDStats, len(ids))
	for index, id := range ids {
		ds[index].DeviceID = id
	}

	mepara := ` -PDList -aALL -NoLog | grep "Media Error Count:" | awk -F ": " '{print $2}'`
	mecmd := fmt.Sprintf("%s%s", MegaCli, mepara)
	mediaErrors, err := cmdCall(mecmd)
	if err != nil {
		return res, err
	}
	if len(ds) < len(mediaErrors) {
		return res, fmt.Errorf("MediaErrors num > devices")
	}
	for index, me := range mediaErrors {
		meInt64, err := strconv.ParseInt(me, 10, 64)
		if err != nil {
			continue
		}
		ds[index].MediaError = meInt64
	}

	oepara := ` -PDList -aALL -NoLog | grep "Other Error Count:" | awk -F ": " '{print $2}'`
	oecmd := fmt.Sprintf("%s%s", MegaCli, oepara)
	otherErrors, err := cmdCall(oecmd)
	if err != nil {
		return res, err
	}
	if len(ds) < len(otherErrors) {
		return res, fmt.Errorf("otherErrors num > devices")
	}
	for index, oe := range otherErrors {
		oeInt64, err := strconv.ParseInt(oe, 10, 64)
		if err != nil {
			continue
		}
		ds[index].OtherError = oeInt64
	}

	fspara := ` -PDList -aALL -NoLog | grep "Firmware state:" | awk -F ": " '{print $2}'`
	fscmd := fmt.Sprintf("%s%s", MegaCli, fspara)
	firmwareStats, err := cmdCall(fscmd)
	if err != nil {
		return res, err
	}
	if len(ds) < len(firmwareStats) {
		return res, fmt.Errorf("firmware Stats num > devices")
	}
	for index, fs := range firmwareStats {
		ds[index].FirmwareState = fs
	}

	tmpara := ` -PDList -aALL -NoLog | grep "Drive Temperature" | awk -F ":" '{print $2}' | awk -F C '{print $1}'`
	tmcmd := fmt.Sprintf("%s%s", MegaCli, tmpara)
	tmps, err := cmdCall(tmcmd)
	if err != nil {
		return res, err
	}
	if len(ds) < len(tmps) {
		return res, fmt.Errorf("firmware Stats num > devices")
	}
	for index, tmp := range tmps {
		tmpInt64, err := strconv.ParseInt(tmp, 10, 64)
		if err != nil {
			continue
		}
		ds[index].Temperature = tmpInt64
	}

	rcpara := ` -AdpAllInfo -aAll -NoLog | grep "Adapter #" | awk -F "#" '{print $2}'`
	rccmd := fmt.Sprintf("%s%s", MegaCli, rcpara)
	rc, err := cmdCall(rccmd)
	if err != nil {
		return res, err
	}
	rs = make([]RaidStats, len(rc))
	for index, num := range rc {
		rs[index].Adapter = num
	}

	cnpara := ` -AdpAllInfo -aAll -NoLog | grep "Critical Disks" | awk -F ": " '{print $2}'`
	cncmd := fmt.Sprintf("%s%s", MegaCli, cnpara)
	cnums, err := cmdCall(cncmd)
	if err != nil {
		return res, err
	}
	if len(rs) < len(cnums) {
		return res, fmt.Errorf("raid array num < Critical nums")
	}
	for index, cn := range cnums {
		cnInt64, err := strconv.ParseInt(cn, 10, 64)
		if err != nil {
			continue
		}
		rs[index].Critical = cnInt64
	}

	fnpara := ` -AdpAllInfo -aAll -NoLog | grep "Failed Disks" | awk -F ": " '{print $2}'`
	fncmd := fmt.Sprintf("%s%s", MegaCli, fnpara)
	fnums, err := cmdCall(fncmd)
	if err != nil {
		return res, err
	}
	if len(rs) < len(fnums) {
		return res, fmt.Errorf("raid array num < Failed Disks nums")
	}
	for index, fn := range fnums {
		fnInt64, err := strconv.ParseInt(fn, 10, 64)
		if err != nil {
			continue
		}
		rs[index].Critical = fnInt64
	}

	res.Disks = ds
	res.Raids = rs
	return res, nil
}

func cmdCall(c string) ([]string, error) {
	var res []string
	cmd := exec.Command("/bin/sh", "-c", c)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	data := stdout.Bytes()
	if len(data) == 0 {
		return nil, nil
	}
	scanner := bufio.NewScanner(bytes.NewBuffer(data))
	for scanner.Scan() {
		res = append(res, scanner.Text())
	}
	return res, nil
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
