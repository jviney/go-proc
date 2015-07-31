package proc

import(
  "os/exec"
  "testing"
)

func TestBasicProcess(t *testing.T) {
  cmd := exec.Command("sleep", "5")
  cmd.Start()

  if process, err := GetProcess(cmd.Process.Pid); err != nil {
    t.Errorf(err.Error())
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
