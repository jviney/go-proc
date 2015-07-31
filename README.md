# Go process inspector

## Overview

Go process inspector exposes process data on Linux and Mac OS with a common API.

## Example

```
  cmd := exec.Command("sleep", "5")
  cmd.Start()

  if process, err := proc.GetProcess(cmd.Process.Pid); err != nil {
    process.Command # "sleep"
    process.ComandLine # "sleep 5"
  }
```

## Supported platforms

Modern flavors of darwin and linux.

## License

Apache 2.0