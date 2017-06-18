package user

import (
	"github.com/j-0-t/murmamau/logging"
	"os"
	"os/user"
	"strconv"
	//"github.com/k0kubun/pp"
)

func Environment() []string {
	out := os.Environ()
	return out
}

func CurrentUser() *user.User {
	u, lookupError := user.Current()
	if lookupError != nil {
		logging.Error(lookupError.Error())
	}
	return u
}

func AllUsers() ([]user.User, map[int]string, map[string]bool) {
	max := 65534
	var out []user.User
	uidList := make(map[int]string)
	homes := make(map[string]bool)

	for i := 0; i < max; i++ {
		username := ""
		u, lookupError := user.LookupId(strconv.Itoa(i))
		if lookupError == nil {
			out = append(out, *u)
			username = u.Username
			homes[u.HomeDir] = true
		} else {
			username = strconv.Itoa(i)
		}
		uidList[i] = username
	}
	return out, uidList, homes
}

func AllGroups() ([]user.Group, map[int]string) {
	max := 65534
	var out []user.Group
	gidList := make(map[int]string)

	for i := 0; i < max; i++ {
		groupname := ""
		g, lookupError := user.LookupGroupId(strconv.Itoa(i))
		if lookupError == nil {
			out = append(out, *g)
			groupname = g.Name
		} else {
			groupname = strconv.Itoa(i)
		}
		gidList[i] = groupname
	}
	return out, gidList
}
