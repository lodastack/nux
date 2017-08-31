package nux

import (
	"fmt"
	"time"
)

type DiskStats struct {
	Major             int
	Minor             int
	Device            string
	ReadRequests      uint64 // Total number of reads completed successfully.
	ReadMerged        uint64 // Adjacent read requests merged in a single req.
	ReadSectors       uint64 // Total number of sectors read successfully.
	MsecRead          uint64 // Total number of ms spent by all reads.
	WriteRequests     uint64 // total number of writes completed successfully.
	WriteMerged       uint64 // Adjacent write requests merged in a single req.
	WriteSectors      uint64 // total number of sectors written successfully.
	MsecWrite         uint64 // Total number of ms spent by all writes.
	IosInProgress     uint64 // Number of actual I/O requests currently in flight.
	MsecTotal         uint64 // Amount of time during which ios_in_progress >= 1.
	MsecWeightedTotal uint64 // Measure of recent I/O completion time and backlog.
	TS                time.Time
}

func (this *DiskStats) String() string {
	return fmt.Sprintf("<Device:%s, Major:%d, Minor:%d, ReadRequests:%d...>", this.Device, this.Major, this.Minor, this.ReadRequests)
}

type DiskHealthStats struct {
	Disks []PDStats
	Raids []RaidStats
}

type PDStats struct {
	DeviceID      string
	MediaError    int64
	OtherError    int64
	Temperature   int64
	FirmwareState string
}

type RaidStats struct {
	Adapter  string
	Critical int64
	Failed   int64
}
