A simple logging module for go, with a rotating file feature and console logging.

## Installation
go get github.com/jbrodriguez/mlog

## Usage
Sample usage

Write to stdout/stderr and create a rotating logfile
```
package main

import (
	"github.com/jbrodriguez/mlog"
)

func main() {
	mlog.Start(mlog.LevelInfo, "app.log")

	mlog.Info("Hello World !")

	ipsum := "ipsum"
	mlog.Warning("Lorem %s", ipsum)
}
```

Write to stdout/stderr only
```
package main

import (
	"github.com/jbrodriguez/mlog"
)

func main() {
	mlog.Start(mlog.LevelInfo, "")

	mlog.Info("Hello World !")

	ipsum := "ipsum"
	mlog.Warning("Lorem %s", ipsum)
}
```

## Output

```
I: 2015/05/15 07:09:45 main.go:10: Hello World !
W: 2015/05/15 07:09:45 main.go:13: Lorem ipsum
```