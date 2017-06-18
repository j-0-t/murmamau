package execute

import (
	"fmt"
	"github.com/contiv/executor"
	"github.com/j-0-t/murmamau/config"
	"github.com/j-0-t/murmamau/logging"
	"golang.org/x/net/context"
	"os/exec"
	"runtime"
	"time"
)

type cmdOut struct {
	cmd        string
	stdout     string
	stderr     string
	returncode int
}

/*
  execute a command
*/
func Command(cmd string) cmdOut {
	shell := "/bin/sh"
	options := "-c"
	cmdTimeout := 10 * time.Second

	e := executor.NewCapture(exec.Command(shell, options, cmd))
	e.Start()
	ctx, _ := context.WithTimeout(context.Background(), cmdTimeout)
	o, err := e.Wait(ctx) // wait for only 10 seconds
	if err != nil {
		//logging.Error(cmd + " :\t" + err.Error())
		logging.Debug(cmd + " :\t" + err.Error())
	}
	out := cmdOut{cmd, o.Stdout, o.Stderr, o.ExitStatus}
	return out
}

/*
  execute all commands configured for the operating system
*/
func RunAllCommands(conf config.Configuration) []cmdOut {
	var out []cmdOut
	commands := conf.Commands
	operatingSystem := runtime.GOOS
	osCommands := commands[operatingSystem]
	for _, c := range osCommands {
		logging.Debug(fmt.Sprintf("Execute:\t%v", c))
		o := Command(c)
		out = append(out, o)
	}
	return out
}
