package proc

import(
  "os"
  "syscall"
)

type Process struct {
  Pid int
  Command string
  CommandLine string
}

func GetProcess(pid int) *Process {
  processes := ps(pid)

  if len(processes) == 1 {
    return processes[0]
  }
  return nil
}

func GetAllProcesses() []*Process {
  return ps(-1)
}

func (p *Process) Active() bool {
  process, err := os.FindProcess(p.Pid)

  if err != nil {
    return false
  }

  err = process.Signal(syscall.Signal(0))
  return err == nil
}
