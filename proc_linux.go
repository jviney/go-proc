// +build linux

package proc

import (
	"bytes"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

func ps(pid int) []*ProcessInfo {
	processes := []*ProcessInfo{}
	files, _ := ioutil.ReadDir("/proc")

	for _, file := range files {
		procPid, err := strconv.Atoi(file.Name())

		// Ignore non-numeric entries
		if err != nil {
			continue
		}

		// Ignore a non-matching pid
		if pid >= 0 && procPid != pid {
			continue
		}

		process := ProcessInfo{Pid: procPid, CommandLine: []string{}}

		if commandLine, err := ioutil.ReadFile("/proc/" + file.Name() + "/cmdline"); err != nil {
			continue // Process terminated
		} else {
			args := bytes.Split(commandLine, []byte{'\x00'})
			for _, arg := range args {
				strArg := strings.TrimSpace(string(arg))
				if len(strArg) == 0 {
					continue
				}

				process.CommandLine = append(process.CommandLine, strArg)
			}
		}

		if stat, err := ioutil.ReadFile("/proc/" + file.Name() + "/stat"); err != nil {
			continue // Process terminated
		} else {
			statRegex := regexp.MustCompile("\\(([^\\)]+)\\)")
			parts := statRegex.FindStringSubmatch(string(stat))

			if len(parts) == 2 {
				process.Command = strings.TrimSpace(parts[1])
			}
		}

		processes = append(processes, &process)

		// Break if this is the one process we were looking for
		if process.Pid == pid {
			break
		}
	}

	return processes
}
