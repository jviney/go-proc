// +build linux

package proc

import (
  "io/ioutil"
  "regexp"
  "strings"
  "strconv"
)

func GetProcess(pid int) (*Process, error) {
  process := Process{Pid: pid}
  pidStr := strconv.Itoa(process.Pid)

  if commandLine, err := ioutil.ReadFile("/proc/" + pidStr + "/cmdline"); err != nil {
    return nil, err
  } else {
    process.CommandLine = strings.TrimSpace(strings.Replace(string(commandLine), "\000", " ", -1))
  }

  if stat, err := ioutil.ReadFile("/proc/" + pidStr + "/stat"); err != nil {
    return nil, err
  } else {
    statRegex := regexp.MustCompile("\\(([^\\)]+)\\)")
    parts := statRegex.FindStringSubmatch(string(stat))

    if len(parts) == 2 {
      process.Command = parts[1]
    }
  }

  return &process, nil
}
