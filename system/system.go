package system

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	"os"
	"strings"
	"time"
)

/*
 gopsutil struct for CPU info
*/
type CpuInfo struct {
	Info []cpu.InfoStat  `json:"info"`
	Time []cpu.TimesStat `json:"time"`
}

/*
 gopsutil struct for Memory info
*/

type MemoryInfo struct {
	VirtualMemory string `json:"virtualMemory"`
	SwapMemory    string `json:"swapMemory"`
}

/*
 gopsutil struct for disk info
*/
type DiskInfo struct {
	Usage      string                         `json:"usage"`
	Partition  []disk.PartitionStat           `json:"partition"`
	IOCounters map[string]disk.IOCountersStat `json:"iOCounters"`
}

/*
 gopsutil struct for host info
*/
type HostInfo struct {
	Info *host.InfoStat  `json:"info"`
	User []host.UserStat `json:"user"`
}

/*
 gopsutil struct for networking info
*/
type NetInfo struct {
	IOCounters    []net.IOCountersStat    `json:"iOCounters"`
	Connection    []net.ConnectionStat    `json:"connection"`
	ProtoCounters []net.ProtoCountersStat `json:"protoCounters"`
	Interface     []net.InterfaceStat     `json:"interface"`
	Filter        []net.FilterStat        `json:"filter"`
}

/*
 gopsutil struct for process info
*/
type ProcessInfo struct {
	Name           string                      `json:"name"`
	Pid            int32                       `json:"pid"`
	Ppid           int32                       `json:"ppid"`
	Exe            string                      `json:"exe"`
	Cmdline        string                      `json:"cmdline"`
	CmdlineSlice   []string                    `json:"cmdlineSlice"`
	CreateTime     int64                       `json:"createTime"`
	Cwd            string                      `json:"cwd"`
	Parent         *process.Process            `json:"parent"`
	Status         string                      `json:"status"`
	Uids           []int32                     `json:"uids"`
	Gids           []int32                     `json:"gids"`
	Terminal       string                      `json:"terminal"`
	Nice           int32                       `json:"nice"`
	IOnice         int32                       `json:"iOnice"`
	Rlimit         []process.RlimitStat        `json:"rlimit"`
	IOCounters     *process.IOCountersStat     `json:"iOCounters"`
	NumCtxSwitches *process.NumCtxSwitchesStat `json:"numCtxSwitches"`
	NumFDs         int32                       `json:"numFDs"`
	NumThreads     int32                       `json:"numThreads"`
	Threads        map[string]string           `json:"threads"`
	Times          *cpu.TimesStat              `json:"times"`
	CPUAffinity    []int32                     `json:"cpuAffinity"`
	MemoryInfo     *process.MemoryInfoStat     `json:"memoryInfo"`
	MemoryInfoEx   *process.MemoryInfoExStat   `json:"memoryInfoEx"`
	Children       []*process.Process          `json:"children"`
	OpenFiles      []process.OpenFilesStat     `json:"openFiles"`
	Connections    []net.ConnectionStat        `json:"connections"`
	NetIOCounters  []net.IOCountersStat        `json:"netIOCounters"`
	IsRunning      bool                        `json:"isRunning"`
	MemoryMaps     *[]process.MemoryMapsStat   `json:"memoryMaps"`
	ReadableStatus string
	ReadableSince  string
}

/*
 gopsutil struct for sys info
*/
type SysInfo struct {
	Host    HostInfo    `json:"host"`
	Cpu     CpuInfo     `json:"cpu"`
	Mem     MemoryInfo  `json:"mem"`
	Disk    DiskInfo    `json:"disk"`
	Net     NetInfo     `json:"net"`
	Process ProcessInfo `json:"process"`
}

/*
  get CPU information using gopsutil
*/
func CpuStat() CpuInfo {
	info, _ := cpu.Info()
	times, _ := cpu.Times(true)

	cpustat := CpuInfo{
		Info: info,
		Time: times,
	}
	return cpustat
}

/*
  get memory information using gopsutil
*/
func MemStat() MemoryInfo {
	virt, _ := mem.VirtualMemory()
	virtmem := virt.String()
	swap, _ := mem.SwapMemory()
	swapmem := swap.String()

	memstat := MemoryInfo{
		VirtualMemory: virtmem,
		SwapMemory:    swapmem,
	}

	return memstat
}

/*
  get disk information using gopsutil
*/
func DiskStat() DiskInfo {
	usage, _ := disk.Usage("./")
	partitions, _ := disk.Partitions(true)
	iOCounters, _ := disk.IOCounters()

	ue := usage.String()

	diskstat := DiskInfo{
		Usage:      ue,
		Partition:  partitions,
		IOCounters: iOCounters,
	}

	return diskstat
}

/*
  get network information using gopsutil
*/
func NetStat() NetInfo {
	iOCounters, _ := net.IOCounters(true)
	protoCounters, _ := net.ProtoCounters([]string{"tcp", "http", "udp", "snmp", "ftp"})
	filterCounters, _ := net.FilterCounters()
	connections, _ := net.Connections("tcp")

	interfaces, _ := net.Interfaces()

	netstat := NetInfo{
		IOCounters:    iOCounters,
		Connection:    connections,
		ProtoCounters: protoCounters,
		Interface:     interfaces,
		Filter:        filterCounters,
	}

	return netstat

}

/*
  get process information for a process using gopsutil
*/
func ProcessStat(pid int32) ProcessInfo {
	//pro := getSelfProcess()
	pro, _ := process.NewProcess(pid)
	processinfo := new(ProcessInfo)

	processinfo.Name, _ = pro.Name()
	processinfo.Pid = int32(os.Getpid())
	processinfo.Ppid, _ = pro.Ppid()
	processinfo.Exe, _ = pro.Exe()
	processinfo.Cmdline, _ = pro.Cmdline()
	processinfo.CmdlineSlice, _ = pro.CmdlineSlice()
	processinfo.CreateTime, _ = pro.CreateTime()
	processinfo.Cwd, _ = pro.Cwd()
	processinfo.Parent, _ = pro.Parent()
	processinfo.Status, _ = pro.Status()
	processinfo.Uids, _ = pro.Uids()
	processinfo.Gids, _ = pro.Gids()
	processinfo.Terminal, _ = pro.Terminal()
	processinfo.Nice, _ = pro.Nice()
	processinfo.IOnice, _ = pro.IOnice()
	processinfo.Rlimit, _ = pro.Rlimit()
	processinfo.IOCounters, _ = pro.IOCounters()
	processinfo.NumCtxSwitches, _ = pro.NumCtxSwitches()
	processinfo.NumFDs, _ = pro.NumFDs()
	processinfo.NumThreads, _ = pro.NumThreads()
	processinfo.Threads, _ = pro.Threads()
	processinfo.Times, _ = pro.Times()
	processinfo.CPUAffinity, _ = pro.CPUAffinity()
	processinfo.MemoryInfo, _ = pro.MemoryInfo()
	processinfo.MemoryInfoEx, _ = pro.MemoryInfoEx()
	processinfo.Children, _ = pro.Children()
	processinfo.OpenFiles, _ = pro.OpenFiles()
	processinfo.Connections, _ = pro.Connections()
	processinfo.NetIOCounters, _ = pro.NetIOCounters(true)
	processinfo.IsRunning, _ = pro.IsRunning()
	processinfo.MemoryMaps, _ = pro.MemoryMaps(true)

	pStatus := processinfo.Status
	switch {
	case pStatus == "S":
		pStatus = "Status: sleeping"
	case pStatus == "R":
		pStatus = "Status: running"
	case pStatus == "T":
		pStatus = "Status: stopped"
	case pStatus == "I":
		pStatus = "Status: idle"
	case pStatus == "Z":
		pStatus = "Status: Zombie"
	case pStatus == "W":
		pStatus = "Status: waiting"
	case pStatus == "L":
		pStatus = "Status: locked"
	}
	processinfo.ReadableStatus = pStatus
	running_since_seconds := processinfo.CreateTime
	running_since := humanize.RelTime(time.Unix(running_since_seconds/1000, 0), time.Now(), "ago", "error")
	processinfo.ReadableSince = fmt.Sprintf("Started %v", running_since)

	p := ProcessInfo(*processinfo)

	return p
}

/*
  get host information using gopsutil
*/
func HostStat() HostInfo {
	info, _ := host.Info()
	users, _ := host.Users()

	hoststat := HostInfo{
		Info: info,
		User: users,
	}
	return hoststat
}

/*
  get list of all proccess and return information about every process
*/
func AllProcesses() ([]ProcessInfo, map[string]bool) {
	LibFiles := make(map[string]bool)
	var runningProcesses []ProcessInfo

	me := int32(os.Getpid())

	all_pids, _ := process.Pids()
	for _, pid := range all_pids {
		/* ignore myself */
		if pid == me {
			continue
		}
		proc := ProcessStat(pid)
		runningProcesses = append(runningProcesses, proc)
		exe := proc.Exe
		if len(exe) < 0 {
			_, exists := LibFiles[exe]
			if !exists {
				LibFiles[exe] = true
			}
		}
		for i := 0; i < len(proc.OpenFiles); i++ {
			path := proc.OpenFiles[i].Path
			if len(path) > 0 {
				if (!strings.HasPrefix(path, "socket:[")) &&
					(!strings.HasPrefix(path, "pipe:[")) &&
					(!strings.HasPrefix(path, "anon_inode:")) &&
					(!strings.HasSuffix(path, "(deleted)")) {
					_, exists := LibFiles[path]
					if !exists {
						LibFiles[path] = true
					}
				}
			}
		}

		if proc.MemoryMaps != nil {
			//fmt.Printf("proc.MemoryMaps Lenth: %v\n", len(*proc.MemoryMaps))

			for i := 0; i < len(*proc.MemoryMaps); i++ {
				path := (*proc.MemoryMaps)[i].Path
				if len(path) > 0 {
					if (path[0] == '/') || (path[0] == '.') {
						_, exists := LibFiles[path]
						if !exists {
							LibFiles[path] = true
						}
					}
				}
			}
		}

	}

	return runningProcesses, LibFiles
}
