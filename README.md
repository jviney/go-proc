# Go process inspector

## Overview

Go process inspector exposes process data on Linux and Mac OS with a common API.

Inspired by [sys-proctable](http://github.com/djberg96/sys-proctable).

## Example

```
  cmd := exec.Command("sleep", "5")
  cmd.Start()

  if process := proc.GetProcess(cmd.Process.Pid); process != nil {
    process.Pid # <pid>
    process.Command # "sleep"
    process.ComandLine # "sleep 5"
  }
```

```
  for _, p := range proc.GetAllProcesses() {
    p.Pid
    p.Command
    p.CommandLine
  }
```

## Supported platforms

Modern flavours of darwin and linux.

## License

Apache 2.0