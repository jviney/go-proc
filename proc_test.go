package proc

import(
  "os/exec"
  "sync"
  "testing"
)

func TestGetProcessNegativePid(t *testing.T) {
  if GetProcess(-1) != nil {
    t.Errorf("expected nil")
  }
}

func TestGetProcessZero(t *testing.T) {
  GetProcess(0)
}

func TestGetProcessMissingPid(t *testing.T) {
  if GetProcess(99999999) != nil {
    t.Errorf("expected nil")
  }
}

func TestProcessFields(t *testing.T) {
  cmd := exec.Command("sleep", "5")
  cmd.Start()

  if process := GetProcess(cmd.Process.Pid); process == nil {
    t.Errorf("failed to find process")
  } else {
    if process.Command != "sleep" {
      t.Errorf("expected %s got %s", "sleep", process.Command)
    }

    if process.CommandLine != "sleep 5" {
      t.Errorf("expected '%s' got '%s'", "sleep 5", process.CommandLine)
    }
  }

  cmd.Process.Kill()
  cmd.Wait()
}

func TestGetProcessGoRoutines(t *testing.T) {
  cmd := exec.Command("sleep", "5")
  cmd.Start()

  count := 100

  var wg sync.WaitGroup
  wg.Add(count)

  for i := 0; i < count; i++ {
    go func() {
      if GetProcess(cmd.Process.Pid) == nil {
        t.Errorf("failed to get process from goroutine")
      }
      wg.Done()
    }()
  }

  wg.Wait()
}

func TestGetAllProcesses(t *testing.T) {
  processes := GetAllProcesses()

  if len(processes) < 10 {
    t.Errorf("expected at least 10 processes")
  }
}
