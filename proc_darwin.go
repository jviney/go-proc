// +build darwin

package proc

// Based on // https://github.com/cloudfoundry/gosigar/blob/master/sigar_darwin.go

// #include <sys/sysctl.h>
import "C"

import (
  "bytes"
  "encoding/binary"
  "io"
  "strings"
  "syscall"
  "unsafe"
)

func GetProcess(pid int) (*Process, error) {
  process := Process{Pid: pid}

  err := kern_procargs(process.Pid,
    func(command string) {
      parts := strings.Split(command, "/")
      process.Command = parts[len(parts) - 1]
    },
    func(argv string) {
      process.CommandLine = process.Command + " " + argv
    },
    nil,
  )

  if err != nil {
    return nil, err
  } else {
    return &process, nil
  }
}

// wrapper around sysctl KERN_PROCARGS2
// callbacks params are optional,
// up to the caller as to which pieces of data they want
func kern_procargs(pid int,
  exe func(string),
  argv func(string),
  env func(string, string)) error {

  mib := []C.int{C.CTL_KERN, C.KERN_PROCARGS2, C.int(pid)}
  argmax := uintptr(C.ARG_MAX)
  buf := make([]byte, argmax)
  err := sysctl(mib, &buf[0], &argmax, nil, 0)
  if err != nil {
    return nil
  }

  bbuf := bytes.NewBuffer(buf)
  bbuf.Truncate(int(argmax))

  var argc int32
  binary.Read(bbuf, binary.LittleEndian, &argc)

  path, err := bbuf.ReadBytes(0)
  if exe != nil {
    exe(string(chop(path)))
  }

  // skip trailing \0's
  for {
    c, _ := bbuf.ReadByte()
    if c != 0 {
      bbuf.UnreadByte()
      break // start of argv[0]
    }
  }

  for i := 0; i < int(argc); i++ {
    arg, err := bbuf.ReadBytes(0)
    if err == io.EOF {
      break
    }
    if argv != nil {
      argv(string(chop(arg)))
    }
  }

  if env == nil {
    return nil
  }

  delim := []byte{61} // "="

  for {
    line, err := bbuf.ReadBytes(0)
    if err == io.EOF || line[0] == 0 {
      break
    }
    pair := bytes.SplitN(chop(line), delim, 2)
    env(string(pair[0]), string(pair[1]))
  }

  return nil
}

func sysctl(mib []C.int, old *byte, oldlen *uintptr,
  new *byte, newlen uintptr) (err error) {
  var p0 unsafe.Pointer
  p0 = unsafe.Pointer(&mib[0])
  _, _, e1 := syscall.Syscall6(syscall.SYS___SYSCTL, uintptr(p0),
    uintptr(len(mib)),
    uintptr(unsafe.Pointer(old)), uintptr(unsafe.Pointer(oldlen)),
    uintptr(unsafe.Pointer(new)), uintptr(newlen))
  if e1 != 0 {
    err = e1
  }
  return
}

func chop(buf []byte) []byte {
  return buf[0 : len(buf)-1]
}
