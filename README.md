A simple logging module for go, with a rotating file feature and console logging.

## Installation
go get github.com/jbrodriguez/mlog

## Usage
Sample usage

Write to stdout/stderr and create a rotating logfile
```go
package main

import "github.com/jbrodriguez/mlog"

func main() {
	mlog.Start(mlog.LevelInfo, "app.log")

	mlog.Info("Hello World !")

	ipsum := "ipsum"
	mlog.Warning("Lorem %s", ipsum)
}
```

Write to stdout/stderr only
```go
package main

import "github.com/jbrodriguez/mlog"

func main() {
	mlog.Start(mlog.LevelInfo, "")

	mlog.Info("Hello World !")

	ipsum := "ipsum"
	mlog.Warning("Lorem %s", ipsum)
}
```

By default, the log will be rolled over to a backup file when its size reaches 10Mb and 10 such files will be created (and eventually reused).

Alternatively, you can specify the max size of the log file before it gets rotated, and the number of backup files you want to create, with the StartEx function.

```go
package main

import "github.com/jbrodriguez/mlog"

func main() {
    mlog.StartEx(mlog.LevelInfo, "app.log", 5*1024*1024, 5)

    mlog.Info("Hello World !")

    ipsum := "ipsum"
    mlog.Warning("Lorem %s", ipsum)
}
```
This will rotate the file when it reaches 5Mb and 5 backup files will eventually be created.


## Output

```
I: 2015/05/15 07:09:45 main.go:10: Hello World !
W: 2015/05/15 07:09:45 main.go:13: Lorem ipsum
```
