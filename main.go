package main

import (
	"fmt"
	"github.com/j-0-t/murmamau/config"
	"github.com/j-0-t/murmamau/execute"
	"github.com/j-0-t/murmamau/find"
	"github.com/j-0-t/murmamau/logging"
	"github.com/j-0-t/murmamau/system"
	"github.com/j-0-t/murmamau/user"
	"github.com/k0kubun/pp"
	"os"
)

/*
 currenlty just a test configuration
*/
var conf = config.Testing

func save(path string, a ...interface{}) {
	find.TarStringAdd(path, pp.Sprint(a))
}

func main() {
	fmt.Println("Start..............")
	var findPath string
	defaultPath := "/"
	if len(os.Args) == 2 {
		findPath = os.Args[1]
		_, err := os.Stat(findPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
	} else {
		findPath = defaultPath
	}

	/* configuration */
	configuration, err := config.ReadConfig(conf)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	if configuration.Debug == true {
		logging.SetDebug(true)
		pp.Println(configuration)
	}
	if configuration.Colors != true {
		pp.ColoringEnabled = false
	}
	/* logging */
	go logging.LogErrors()

	/* user infos */
	logging.Status("Getting system users")
	currentUser := user.CurrentUser()
	userEnv := user.Environment()
	allUsersOut, uidList, homeDirectories := user.AllUsers()
	allGroupsOut, gidList := user.AllGroups()

	/* system infos */
	logging.Status("Getting host stats")
	hostStats := system.HostStat()
	logging.Status("Getting disk stats")
	diskStats := system.DiskStat()
	logging.Status("Getting memory stats")
	memStats := system.MemStat()
	logging.Status("Getting cpu stats")
	cpuStats := system.CpuStat()
	logging.Status("Getting network stats")
	netStats := system.NetStat()
	logging.Status("Getting information on running processes")
	psStats, openFiles := system.AllProcesses()
	/* adding open files|libraries of running processes */
	for f := range openFiles {
		configuration.SearchList = append(configuration.SearchList, f)
	}

	/* execute */
	logging.Status("Execute commands")
	osCommands := execute.RunAllCommands(configuration)

	/* find */
	/* prepare find */
	find.SetUidList(uidList)
	find.SetGidList(gidList)
	find.SetHomes(homeDirectories)
	/* now set configuration */
	find.SetConfig(configuration)
	/* now start channels with new configuration */
	find.StartConfigured()
	/* saving data */
	save("/murmamau/current_user.json", currentUser)
	save("/murmamau/environment.json", userEnv)
	save("/murmamau/all_users.json", allUsersOut)
	save("/murmamau/all_groups.json", allGroupsOut)
	save("/murmamau/host_stats.json", hostStats)
	save("/murmamau/disk_stats.json", diskStats)
	save("/murmamau/mem_stats.json", memStats)
	save("/murmamau/host_stats.json", hostStats)
	save("/murmamau/network_stats.json", netStats)
	save("/murmamau/cpu_stats.json", cpuStats)
	save("/murmamau/proccesses.json", psStats)
	save("/murmamau/commands.json", osCommands)
	save("/murmamau/ps_open_files.json", openFiles)

	/* tar known path */
	logging.Status("Tar files with known path")
	find.TarKnownFiles()

	/* runninrg find and tar */
	logging.Status("Find and tar files")
	find.Files(findPath)
	logging.Status("Saving data")

	/* saving data */
	errorlog := logging.LogFile()
	find.TarFileAdd(errorlog, "/murmamau/errors.txt")
	find.SaveData()

	fmt.Println("Murmamau done")
}
