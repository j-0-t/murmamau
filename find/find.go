package find

import (
	"bufio"
	"fmt"
	"os"
	"path"
	//"path/filepath"
	"archive/tar"
	"github.com/dustin/go-humanize"
	"github.com/gobwas/glob"
	"github.com/j-0-t/murmamau/config"
	"github.com/j-0-t/murmamau/logging"
	"github.com/k0kubun/pp"
	"github.com/stretchr/powerwalk"
	"io"
	"io/ioutil"
	"strings"
	"sync"
	"syscall"
	"time"
)

type fileStat struct {
	path string
	stat *syscall.Stat_t
}

type fileLs struct {
	perms     string
	links     uint64
	username  string
	groupname string
	size      string
	timestamp string
	file      string
}

var maxTarSize int64
var everyXstatus int
var tarFullPath map[string]bool
var tarSearchListString map[string]bool
var tarSearchIgnore map[string]bool
var tarBlacklist map[string]bool
var gotFilesList map[string]bool
var uidList map[int]string
var gidList map[int]string
var homes map[string]bool
var tarball *tar.Writer
var tarSearchListGlob []glob.Glob
var filesLogfile *os.File
var filesBuffer bufio.Writer
var wg sync.WaitGroup
var tarLock sync.RWMutex

var findChannel chan string
var findPrintChannel chan fileLs
var logChannel chan string
var tarChannel chan fileStat
var gotFilesChannel chan string

func init() {
	tarFullPath = make(map[string]bool)
	tarSearchListString = make(map[string]bool)
	tarSearchIgnore = make(map[string]bool)
	tarBlacklist = make(map[string]bool)
	gotFilesList = make(map[string]bool)
	uidList = make(map[int]string)
	gidList = make(map[int]string)
	homes = make(map[string]bool)

	findChannel = make(chan string, 4)
	findPrintChannel = make(chan fileLs, 4)
	tarChannel = make(chan fileStat, 4*2)
	gotFilesChannel = make(chan string)

	now := time.Now()
	hostname, nerr := os.Hostname()
	if nerr != nil {
		logging.Error(nerr.Error())
		hostname = "default"
	}

	tarFileName := "log-murmamau-" + hostname + "-" + now.Format("20060102150x405") + ".tar"
	tarfile, err := os.Create(tarFileName) // TODO
	if err != nil {
		logging.Error(fmt.Sprintf("%s ->\t%v", tarfile, err))
	}
	logging.Status("Output logged to " + tarFileName)
	tarball = tar.NewWriter(tarfile)
}

/*
  Start go routines with configuration loaded
*/
func StartConfigured() {
	tmpfile, err := ioutil.TempFile("", "files")

	if err != nil {
		logging.Error(err.Error())
		return
	}
	filesLogfile = tmpfile
	filesBuffer := bufio.NewWriter(filesLogfile)
	go handleFind()
	go findPrint(filesBuffer)
	go tarFiles()
	go gottenFiles()
}

/*
  Save logfile as new filename into TAR archive; overwrite file and remove it
*/
func SaveLogFile(logfile *os.File, newPath string) {
	path := logfile.Name()
	err := logfile.Sync()
	if err != nil {
		logging.Error(err.Error())
	}
	logfile.Close()
	/* add to TAR archive */
	TarFileAdd(path, newPath)
	/* overwrite and delete tmp file */
	info, err := os.Stat(path)
	if err != nil {
		logging.Error(string(err.Error()))
	}
	size := info.Size()
	f, err := os.OpenFile(path, os.O_RDWR, 0700)
	overWriteBuffer := make([]byte, size)
	w := bufio.NewWriter(f)
	defer w.Flush()
	io.WriteString(w, string(overWriteBuffer))
	w.Flush()
	f.Close()
	rmerror := os.Remove(path)
	if rmerror != nil {
		logging.Error(rmerror.Error())
	}

}

/*
  save data and exit go routines
*/
func SaveData() {
	filesBuffer.Flush()
	SaveLogFile(filesLogfile, "/murmamau/files.json")
	/*
	   ending go routines for printing fliles and tarring files
	*/
	findPrintChannel <- fileLs{"", 0, "", "", "", "", "Stop findPrint()"}
	tarChannel <- fileStat{"Stop tarFiles()", nil}
	wg.Wait()

	/*
	   now close TAR file
	*/
	tarerr := tarball.Close()
	if tarerr != nil {
		logging.Error(tarerr.Error())
	}

}

/*
  Set list of uids (for not looking the username up for each file again
*/
func SetUidList(externlUidList map[int]string) {
	uidList = externlUidList
}

/*
  Set list of gids (for not looking the groupname up for each file again
*/
func SetGidList(externlGidList map[int]string) {
	gidList = externlGidList
}

/*
  Set list of home directories
*/
func SetHomes(homeDirectories map[string]bool) {
	homes = homeDirectories
}

/*
  Set configuration
*/
func SetConfig(c config.Configuration) {
	var g glob.Glob
	maxTarSize = int64(c.MaxSize)
	everyXstatus = c.EveryXstatus
	for _, searchString := range c.SearchList {
		if strings.HasPrefix(searchString, "~/") {
			for h := range homes {
				path := h + strings.TrimPrefix(searchString, "~")
				logging.Debug(fmt.Sprint("Adding to SearchList " + path))

				_, exists := tarFullPath[path]
				if !exists {
					tarFullPath[path] = true
				}
			}
		} else {
			logging.Debug(fmt.Sprint("Adding to PathList " + searchString))
			if searchString[0] == '/' {
				_, exists := tarFullPath[searchString]
				if !exists {
					tarFullPath[searchString] = true
				}
			} else {
				_, exists := tarSearchListString[searchString]
				if !exists {
					tarSearchListString[searchString] = true
				}

			}
		}

	}
	for _, searchString := range c.GlobList {
		if strings.HasPrefix(searchString, "~/") {
			for h := range homes {
				path := h + strings.TrimPrefix(searchString, "~")
				g = glob.MustCompile(path)
				tarSearchListGlob = append(tarSearchListGlob, g)
			}
		} else {
			g = glob.MustCompile(searchString)
			tarSearchListGlob = append(tarSearchListGlob, g)
		}
	}
	for _, ignoreString := range c.IgnoreList {
		_, exists := tarSearchIgnore[ignoreString]
		if !exists {
			tarSearchIgnore[ignoreString] = true
		}
	}
	for _, blacklistString := range c.BlackList {
		_, exists := tarBlacklist[blacklistString]
		if !exists {
			tarBlacklist[blacklistString] = true
		}
	}

}

/*
  log files and give status updates
*/
func findPrint(w *bufio.Writer) {
	wg.Add(1)
	entries := 0
	/* print a status message every i files */
	i := everyXstatus
	for s := range findPrintChannel {
		ls := pp.Sprintln(s)
		fmt.Fprint(w, ls)
		entries++
		if entries%i == 0 {
			tmp := fmt.Sprintf("%s\t%v\t%+v\t%+v\t%v\t%+v\t%+v", s.perms, s.links, s.username, s.groupname, s.size, s.timestamp, s.file)
			logging.Status(fmt.Sprintf("[Status of find]: %d files\n\t%s", entries, tmp))
			w.Flush()
		}
		if s.file == "Stop findPrint()" {
			logging.Debug("Stopping findPrint()")
			wg.Done()
			return
		}

	}
	w.Flush() // Don't forget to flush!

}

/*
  handle files from findChannel
  * log them
  * check for tar and send to tarChannel
*/
func handleFind() {
	for file := range findChannel {
		/* could crash here, so adding this check even if is an extra Stat() call
		   (lower performance)
		*/
		info, err := os.Stat(file)
		if err != nil {
			logging.DebugError("Stat() " + string(err.Error()))
			continue
		}

		sys, ok := info.Sys().(*syscall.Stat_t)
		if !ok {
			logging.DebugError("info.Sys() is false")
			continue
		}
		perms := strings.ToLower(info.Mode().String())
		links := sys.Nlink
		username := uidList[int(sys.Uid)]
		groupname := gidList[int(sys.Gid)]
		size := strings.Replace(humanize.Bytes(uint64(sys.Size)), " ", "", -1)
		timestamp := info.ModTime().Format(time.UnixDate)
		findPrintChannel <- fileLs{perms, links, username, groupname, size, timestamp, file}

		if info.IsDir() == false {
			if matchForTar(file) == true {
				if ignoreForTar(file) == false {
					logging.Debug(fmt.Sprint("Match " + file))
					if sizeForTar(sys.Size) == true {
						tarChannel <- fileStat{file, sys}
					} else {
						logging.Error(fmt.Sprintf("%s :\tFile size to big (maximum=%+v) size=%s", file, maxTarSize, size))
					}
				} else {
					logging.Debug(fmt.Sprint("Ignore " + file))
				}
			} else {
				if perms[0] == 'u' {
					logging.Debug(fmt.Sprint("SUID " + file))
					tarChannel <- fileStat{file, sys}
				}
			}
		}

	}
}

/*
  Walk the given directory
*/
func Files(root string) {
	var err error
	err = powerwalk.Walk(root, func(myPath string, info os.FileInfo, err error) error {
		if info == nil {
			logging.DebugError(err.Error())
		}
		fullPath := path.Clean(myPath)
		findChannel <- fullPath
		return nil

	})
	if err != nil {
		logging.Error(err.Error())
	}
}

/*
  check if filename should be saved
*/
func matchForTar(filename string) bool {
	_, exists := tarSearchListString[filename]
	if exists {
		logging.Debug("Match: " + filename + "\tin searchlist")
		return true
	}
	_, existsBase := tarSearchListString[path.Base(filename)]
	if existsBase {
		logging.Debug("Match: " + filename + "\t basename in searchlist")
		return true
	}

	for _, g := range tarSearchListGlob {
		if g.Match(filename) == true {
			logging.Debug("Match: " + filename + "\tmatches globlist")
			return true
		}
	}
	return false
}

/*
  check if filename should not be saved
    * already saved
    * file extension on the list of ignored extensions
    * filename on blacklist
*/
func ignoreForTar(filename string) bool {
	_, existsIgnore := tarSearchIgnore[path.Ext(filename)]
	if existsIgnore {
		logging.Debug("Ignore: " + filename + "\textension in ignorelist")
		return true
	}

	_, existsLoaded := gotFilesList[filename]
	if existsLoaded {
		logging.Debug("Ignore: " + filename + "\talready saved")
		return true
	}
	_, blacklisted := tarBlacklist[filename]
	if blacklisted {
		logging.Debug("Ignore: " + filename + "\tfilename on blacklist")
		return true
	}
	return false
}

/*
  check if filesize is to big for saving
*/
func sizeForTar(size int64) bool {
	if size < maxTarSize {
		return true
	} else {
		return false
	}
}

/*
  tar a file (filename as argument)
*/
func tarThisFile(path string) {
	logging.Debug(fmt.Sprint("Tar output to " + path))
	file, err := os.Stat(path)
	if err != nil {
		logging.DebugError(fmt.Sprintf("%s ->\t%v", path, err))
		return
	}
	sys, ok := file.Sys().(*syscall.Stat_t)
	if !ok {
		logging.Error("Cannot get syscall.Stat_t on " + path)
		return // TODO: error
	}
	tarChannel <- fileStat{path, sys}

}

/*
  try to tar all files on the list
*/
func TarKnownFiles() {
	for s := range tarFullPath {
		tarThisFile(s)
	}
}

/*
  store filenames of saved files for avoiding to save them more than once
*/
func gottenFiles() {
	for path := range gotFilesChannel {
		_, exists := gotFilesList[path]
		if !exists {
			tmp := make(map[string]bool)
			for k, v := range gotFilesList {
				tmp[k] = v
			}
			tmp[path] = true
			gotFilesList = tmp
		}
	}
}

/*
  reads files from tarChannel and saves them
*/
func tarFiles() {
	wg.Add(1)
	for file := range tarChannel {
		if file.path == "Stop tarFiles()" {
			logging.Debug("Stopping tarFiles()")
			wg.Done()
			return
		}
		tarFile(file)
	}
}

/*
  save a file with a new filenname in the TAR archive
*/
func TarFileAdd(path string, tarPath string) {
	logging.Debug(fmt.Sprintf("Tar: adding %s as %s", path, tarPath))
	stat, err := os.Stat(path)
	if err != nil {
		logging.DebugError(fmt.Sprintf("%s ->\t%v", path, err))
		return
	}
	mode := stat.Mode()
	if mode.IsRegular() == false {
		logging.Debug("Not a regular file: " + path)
		return
	}

	/*
	   Opening file and return if it does nt work (no permissions, etc)
	*/

	f, err := os.Open(path)
	if err != nil {
		logging.DebugError(string(err.Error()))
		return
	}
	/*
	   ignoring empty files (and devices)
	   if not devices like /dev/urandom are reaad and it crashes "fatal error: runtime: out of memory"
	*/

	if stat.Size() == 0 {
		logging.Debug(fmt.Sprint("Tar: ignored because it seems to be empty:\t" + path))
		f.Close()
		wg.Done()
		return
	}
	// xxxxxxxxxxxxxxx
	/*
		info, err := f.Stat()
	  if err != nil {
			logging.Error(path + "\t" + err.Error())
		}
		header, err := tar.FileInfoHeader(info, path)
	*/
	header, err := tar.FileInfoHeader(stat, path)
	if err != nil {
		logging.Error(path + "\t" + err.Error())
	}
	header.Name = tarPath
	tarAdd(header, f)

}

/*
  save a string as given file inside of the TAR archive
*/
func TarStringAdd(path string, data string) {
	wg.Add(1)
	currentTime := time.Now()
	header := new(tar.Header)
	header.Name = path
	header.Mode = 0600
	header.ModTime = currentTime
	header.AccessTime = currentTime
	header.ChangeTime = currentTime
	header.Size = int64(len([]byte(data)))
	tarLock.Lock()
	err := tarball.WriteHeader(header)
	if err != nil {
		logging.Error(path + "\t" + err.Error())
	}

	_, err = tarball.Write([]byte(data))
	if err != nil {
		logging.Error(path + "\t" + err.Error())
	}
	err = tarball.Flush()
	if err != nil {
		logging.Error(path + "\t" + err.Error())
	}
	tarLock.Unlock()

	logging.Debug(fmt.Sprint("Tar : added data as " + header.Name))
	wg.Done()
}

/*
  save a file into the TAR archive
*/
func tarFile(file fileStat) {
	path := file.path
	stat := file.stat
	logging.Debug(fmt.Sprint("Tar: " + path))

	fstat, err := os.Stat(path)
	if err != nil {
		logging.DebugError(fmt.Sprintf("%s ->\t%v", path, err))
		return
	}

	mode := fstat.Mode()
	if mode.IsRegular() == false {
		logging.Debug("Not a regular file: " + path)
		return
	}

	/*
	   Opening file and return if it does nt work (no permissions, etc)
	*/
	f, err := os.Open(file.path)
	if err != nil {
		logging.DebugError(string(err.Error()))
		return
	}

	/*
	   ignoring empty files (and devices)
	   if not devices like /dev/urandom are reaad and it crashes "fatal error: runtime: out of memory"
	*/
	if stat.Size == 0 {
		logging.Debug(fmt.Sprint("Tar: ignored because it seems to be empty:\t" + path))
		f.Close()
		return
	}
	info, err := f.Stat()

	/*
	   create TAR header from Stat() of the file
	   and set the required path and size
	*/
	header, err := tar.FileInfoHeader(info, path)
	if err != nil {
		logging.Error(path + "\t" + err.Error())
		f.Close()
		return
	}
	header.Name = path
	header.Size = info.Size()

	tarAdd(header, f)
}

func tarAdd(header *tar.Header, f *os.File) {
	wg.Add(1)
	path := header.Name
	/*
	   start a lock on the tar file (to be sure that nothing else is writing now)
	*/
	tarLock.Lock()
	err := tarball.WriteHeader(header)
	if err != nil {
		logging.Error(path + ":\t" + err.Error())
		f.Close()
		wg.Done()
		return
	}
	/*
	   copy file to tar
	   io.copyN to be sure not to write more than in TAR header.size
	*/
	n, err := io.CopyN(tarball, f, header.Size)
	if err != nil {
		logging.Error(path + "\t" + err.Error())
	}
	/*
	   double check error and fix common issues
	*/
	if header.Size != n {
		logging.Debug(fmt.Sprintf("%s\tsize=%d\twritten=%d", path, header.Size, n))
		/*
		   In some cases the real size is lower than the size in stat.
		   This would write less data into the TAR file and mess up the whole rest of the TAR file.
		   Adding some 0 bytes fixes this issue.
		   Typical example is:
		     /sys/devices/virtual/tty/tty0/active
		         size: 4096  real size 5
		*/
		if n < header.Size {
			needed := header.Size - n
			fillBuffer := make([]byte, needed)
			tarball.Write(fillBuffer)
			logging.Debug("Added a fill buffer")
		}
	}
	//tarball.Flush()
	tarLock.Unlock()
	gotFilesChannel <- path
	f.Close()
	wg.Done()
}
