package proc

import(
  "log"
  "os/exec"
  "testing"
)

func TestBasicProcess(t *testing.T) {
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
}

func TestGetAllProcesses(t *testing.T) {
  for _, p := range GetAllProcesses() {
    log.Printf("%d %s %s", p.Pid, p.Command, p.CommandLine)
  }
}
