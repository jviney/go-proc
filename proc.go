package proc

type ProcessInfo struct {
	Pid         int
	Command     string
	CommandLine []string
}

func GetProcessInfo(pid int) *ProcessInfo {
	processes := ps(pid)

	if len(processes) == 1 {
		return processes[0]
	}

	return nil
}

func GetAllProcessesInfo() []*ProcessInfo {
	return ps(-1)
}
