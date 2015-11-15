package mlog

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sync/atomic"
)

type LogLevel int32

const (
	// LevelTrace logs everything
	LevelTrace LogLevel = (1<<iota)

	// LevelInfo logs Info, Warnings and Errors
	LevelInfo

	// LevelWarn logs Warning and Errors
	LevelWarn

	// LevelError logs just Errors
	LevelError
)

type mlog struct {
	LogLevel int32

	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
	Fatal   *log.Logger

	LogFile *RotatingFileHandler
}

var logger mlog

//RotatingFileHandler writes log a file, if file size exceeds maxBytes,
//it will backup current file and open a new one.
//
//max backup file number is set by backupCount, it will delete oldest if backups too many.
type RotatingFileHandler struct {
	fd *os.File

	fileName    string
	maxBytes    int
	backupCount int
}

func NewRotatingFileHandler(fileName string, maxBytes int, backupCount int) (*RotatingFileHandler, error) {
	dir := path.Dir(fileName)
	os.Mkdir(dir, 0777)

	h := new(RotatingFileHandler)

	if maxBytes <= 0 {
		return nil, fmt.Errorf("invalid max bytes")
	}

	h.fileName = fileName
	h.maxBytes = maxBytes
	h.backupCount = backupCount

	var err error
	h.fd, err = os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return h, nil
}

func (h *RotatingFileHandler) Write(p []byte) (n int, err error) {
	h.doRollover()
	return h.fd.Write(p)
}

func (h *RotatingFileHandler) Close() error {
	if h.fd != nil {
		return h.fd.Close()
	}
	return nil
}

func (h *RotatingFileHandler) doRollover() {
	f, err := h.fd.Stat()
	if err != nil {
		return
	}

	// log.Println("size: ", f.Size())

	if h.maxBytes <= 0 {
		return
	} else if f.Size() < int64(h.maxBytes) {
		return
	}

	if h.backupCount > 0 {
		h.fd.Close()

		for i := h.backupCount - 1; i > 0; i-- {
			sfn := fmt.Sprintf("%s.%d", h.fileName, i)
			dfn := fmt.Sprintf("%s.%d", h.fileName, i+1)

			os.Rename(sfn, dfn)
		}

		dfn := fmt.Sprintf("%s.1", h.fileName)
		os.Rename(h.fileName, dfn)

		h.fd, _ = os.OpenFile(h.fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	}
}

func Start(level LogLevel, path string) {
	doLogging(level, path)
}

func Stop() error {
	if logger.LogFile != nil {
		return logger.LogFile.Close()
	}

	return nil
}

func doLogging(logLevel LogLevel, fileName string) {
	traceHandle := ioutil.Discard
	infoHandle := ioutil.Discard
	warnHandle := ioutil.Discard
	errorHandle := ioutil.Discard
	fatalHandle := ioutil.Discard

	var fileHandle *RotatingFileHandler
	
	switch logLevel {
	case LevelTrace :
		traceHandle = os.Stdout
		fallthrough
	case LevelInfo:
		infoHandle = os.Stdout
		fallthrough
	case LevelWarn:
		warnHandle = os.Stdout
		fallthrough
	case LevelError:
		errorHandle = os.Stderr
		fatalHandle = os.Stderr
	}

	if fileName != "" {
		var err error
		fileHandle, err = NewRotatingFileHandler(fileName, 1024*1024*1024, 10)
		if err != nil {
			log.Fatal("mlog: unable to create RotatingFileHandler: ", err)
		}

		if traceHandle == os.Stdout {
			traceHandle = io.MultiWriter(fileHandle, traceHandle)
		}

		if infoHandle == os.Stdout {
			infoHandle = io.MultiWriter(fileHandle, infoHandle)
		}

		if warnHandle == os.Stdout {
			warnHandle = io.MultiWriter(fileHandle, warnHandle)
		}

		if errorHandle == os.Stderr {
			errorHandle = io.MultiWriter(fileHandle, errorHandle)
		}

		if fatalHandle == os.Stderr {
			fatalHandle = io.MultiWriter(fileHandle, fatalHandle)
		}
	}

	logger = mlog{
		Trace:   log.New(traceHandle, "T: ", log.Ldate|log.Ltime|log.Lshortfile),
		Info:    log.New(infoHandle, "I: ", log.Ldate|log.Ltime|log.Lshortfile),
		Warning: log.New(warnHandle, "W: ", log.Ldate|log.Ltime|log.Lshortfile),
		Error:   log.New(errorHandle, "E: ", log.Ldate|log.Ltime|log.Lshortfile),
		Fatal:   log.New(errorHandle, "F: ", log.Ldate|log.Ltime|log.Lshortfile),
		LogFile: fileHandle,
	}

	atomic.StoreInt32(&logger.LogLevel, int32(logLevel))
}

//** TRACE

// Trace writes to the Trace destination
func Trace(format string, a ...interface{}) {
	logger.Trace.Output(2, fmt.Sprintf(format, a...))
}

//** INFO

// Info writes to the Info destination
func Info(format string, a ...interface{}) {
	logger.Info.Output(2, fmt.Sprintf(format, a...))
}

//** WARNING

// Warning writes to the Warning destination
func Warning(format string, a ...interface{}) {
	logger.Warning.Output(2, fmt.Sprintf(format, a...))
}

//** ERROR

// Error writes to the Error destination and accepts an err
func Error(err error) {
	logger.Error.Output(2, fmt.Sprintf("%s\n", err))
}

// Fatalf writes to the Fatal destination and exits with an error code
func Fatalf(format string, a ...interface{}) {
	logger.Fatal.Output(2, fmt.Sprintf(format, a...))
	if logger.LogFile != nil {
		logger.LogFile.fd.Sync()
	}
	os.Exit(255)
}
