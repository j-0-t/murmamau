package logging

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"time"
)

var debugging bool
var logChannel chan string
var logPath string

func init() {
	debugging = false
	logPath = "/tmp/errors-2.log"
	logChannel = make(chan string, 4)

}

/*
  get the calling function, filename and line number for debugging
*/
func myCaller() string {

	// we get the callers as uintptrs - but we just need 1
	fpcs := make([]uintptr, 1)

	// skip 3 levels to get to the caller of whoever called Caller()
	n := runtime.Callers(3, fpcs)
	if n == 0 {
		return "n/a" // proper error her would be better
	}

	// get the info of the actual function that's in the pointer
	fun := runtime.FuncForPC(fpcs[0] - 1)
	if fun == nil {
		return "n/a"
	}
	filename, line := fun.FileLine(fpcs[0] - 1)
	out := fmt.Sprintf("[%s:%d] %s", filename, line, fun.Name())
	return out
}

/*
  Error sends error messages to logChannel
*/
func Error(s string) {
	currentTime := time.Now().UTC()
	caller := ""
	if debugging == true {
		caller = myCaller()
	}
	logChannel <- fmt.Sprintf("Error:\t%v()\t|%v|\t%v", caller, currentTime.Format(time.ANSIC), s)
}

/*
  Print and save error messages from logChannel
*/
func LogErrors() {
	path := logPath
	f, err := os.Create(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()

	for s := range logChannel {
		fmt.Fprintln(os.Stderr, s)
		fmt.Fprintln(w, s)
		w.Flush()
	}
}

/*
  path of the logfile
*/
func LogFile() string {
	return logPath
}

/*
  Print a status message
*/
func Status(s string) {
	//fmt.Println(s)
  logChannel <- s
}

/*
  Turn debugging on/off
*/
func SetDebug(b bool) {
	debugging = b
}

/*
  if debugging is enabled print out a debugging message
*/
func Debug(s string) {
	if debugging == true {
    logChannel <- fmt.Sprint("[DEBUG]\t" + s)
	}
}

/*
  like Debug() + error defails
*/
func DebugError(s string) {
	if debugging == true {
		currentTime := time.Now().UTC()
    caller := ""
		caller = myCaller()
    logChannel <- fmt.Sprintf("DebugError:\t%v()\t|%v|\t%v", caller, currentTime.Format(time.ANSIC), s)
	}
}
