// +build darwin

package proc

// Based on https://github.com/cloudfoundry/gosigar/blob/master/sigar_darwin.go

/*
#include <sys/sysctl.h>
typedef struct kinfo_proc kInfoProc;
*/
import "C"

import (
  "bytes"
  "encoding/binary"
  "io"
  "strings"
  "syscall"
  "unsafe"
)

func ps(pid int) []*ProcessInfo {
  processes := []*ProcessInfo{}

  mib := []C.int{C.CTL_KERN, C.KERN_PROC, C.KERN_PROC_ALL, 0}
  length := uintptr(0)

  if err := sysctl(mib, nil, &length, nil, 0); err != nil {
    return nil
  }

  buf := make([]byte, length)
  if err := sysctl(mib, &buf[0], &length, nil, 0); err != nil {
    return nil
  }

  kInfoProcSize := int(unsafe.Sizeof(C.kInfoProc{}))
  count := int(length) / kInfoProcSize

  for i := 0; i < count; i++ {
    proc := (*C.kInfoProc) (unsafe.Pointer(&buf[i * kInfoProcSize]))
    procPid := int(proc.kp_proc.p_pid)

    if pid >= 0 && procPid != pid {
      continue
    }

    process := ProcessInfo{Pid: procPid}

    command, argv, err := kern_procargs(process.Pid)

    if err != nil {
      continue
    }

    commandParts := strings.Split(command, "/")
    process.Command = strings.TrimSpace(commandParts[len(commandParts) - 1])
    process.CommandLine = argv

    processes = append(processes, &process)

    if process.Pid == pid {
      break
    }
  }

  return processes
}

func kern_procargs(pid int) (command string, argv []string, err error) {
  mib := []C.int{C.CTL_KERN, C.KERN_PROCARGS2, C.int(pid)}
  argmax := uintptr(C.ARG_MAX)
  buf := make([]byte, argmax)
  err = sysctl(mib, &buf[0], &argmax, nil, 0)

  if err != nil {
    return
  }

  bbuf := bytes.NewBuffer(buf)
  bbuf.Truncate(int(argmax))

  var argc int32
  binary.Read(bbuf, binary.LittleEndian, &argc)

  path, err := bbuf.ReadBytes(0)
  command = string(chop(path))

  // skip trailing \0's
  for {
    c, _ := bbuf.ReadByte()
    if c != 0 {
      bbuf.UnreadByte()
      break // start of argv[0]
    }
  }

  argv = make([]string, argc)

  for i := 0; i < int(argc); i++ {
    arg, err := bbuf.ReadBytes(0)
    if err == io.EOF {
      break
    }
    argv[i] = string(chop(arg))
  }

  return
}


func sysctl(mib []C.int, old *byte, oldlen *uintptr, new *byte, newlen uintptr) (err error) {
  _, _, e1 := syscall.Syscall6(
    syscall.SYS___SYSCTL,
    uintptr(unsafe.Pointer(&mib[0])),
    uintptr(len(mib)),
    uintptr(unsafe.Pointer(old)),
    uintptr(unsafe.Pointer(oldlen)),
    uintptr(unsafe.Pointer(new)),
    uintptr(newlen))

  if e1 != 0 {
    err = e1
  }

  return
}

func chop(buf []byte) []byte {
  return buf[0 : len(buf)-1]
}
