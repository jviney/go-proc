// +build linux

package proc

import (
  "io/ioutil"
  "regexp"
  "strings"
  "strconv"
)

func ps(pid int) []*Process {
  processes := []*Process{}
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

    process := Process{Pid: procPid}

    if commandLine, err := ioutil.ReadFile("/proc/" + file.Name() + "/cmdline"); err != nil {
      continue // Process terminated
    } else {
      process.CommandLine = strings.TrimSpace(strings.Replace(string(commandLine), "\000", " ", -1))
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

    if process.CommandLine == "" {
      process.CommandLine = process.Command
    }

    processes = append(processes, &process)

    if process.Pid == pid {
      break
    }
  }

  return processes
}
