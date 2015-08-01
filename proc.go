package proc

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
