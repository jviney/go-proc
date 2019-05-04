package proc

import (
	"os/exec"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetProcessNegativePid(t *testing.T) {
	assert.Nil(t, GetProcessInfo(-1))
}

func TestGetProcessZero(t *testing.T) {
	assert.Nil(t, GetProcessInfo(0))
}

func TestGetProcessMissingPid(t *testing.T) {
	assert.Nil(t, GetProcessInfo(99999999))
}

func TestProcessFields(t *testing.T) {
	cmd := exec.Command("sleep", "5")
	cmd.Start()

	process := GetProcessInfo(cmd.Process.Pid)
	assert.NotNil(t, process)

	assert.Equal(t, "sleep", process.Command)
	assert.Equal(t, []string{"sleep", "5"}, process.CommandLine)

	cmd.Process.Kill()
	cmd.Wait()
}

func TestLongCommandLine(t *testing.T) {
	cmd := exec.Command("dd", "if=/dev/zero", "of=/dev/null")
	cmd.Start()

	process := GetProcessInfo(cmd.Process.Pid)
	assert.NotNil(t, process)

	assert.Equal(t, "dd", process.Command)
	assert.Equal(t, []string{"dd", "if=/dev/zero", "of=/dev/null"}, process.CommandLine)

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
			assert.NotNil(t, GetProcessInfo(cmd.Process.Pid))
			wg.Done()
		}()
	}

	wg.Wait()
}

func TestGetAllProcesses(t *testing.T) {
	processes := GetAllProcessesInfo()

	assert.Greater(t, len(processes), 10)
}
