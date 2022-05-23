package log

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/joycastle/cop/util"
)

const (
	green   = "\033[97;42m"
	white   = "\033[90;47m"
	yellow  = "\033[90;43m"
	red     = "\033[97;41m"
	blue    = "\033[97;44m"
	magenta = "\033[97;45m"
	cyan    = "\033[97;46m"
	reset   = "\033[0m"
)

const (
	level_debug = iota
	level_info
	level_warn
	level_fatal
)

const (
	Log_Path_Stdout = "!#%WDCXFT97"
	Log_Path_Stderr = "^&$IJHVZSC8"
)

type LogConf struct {
	LogName    string `yaml:"LogName"`
	ExpireDays int64  `yaml:"ExpireDays"`
}

type Logger struct {
	Logger *log.Logger
	Fptr   *os.File
	Fname  string

	logName      string
	colorMap   map[int]string
	mu         *sync.RWMutex
	expireDays int64
	oldFiles   map[int64]string
	oMu        *sync.Mutex
}

func NewLogger(cfg LogConf) *Logger {
	l := &Logger{
		logName:      cfg.LogName,
		colorMap:   make(map[int]string),
		mu:         new(sync.RWMutex),
		expireDays: cfg.ExpireDays,
		oldFiles:   make(map[int64]string),
		oMu:        new(sync.Mutex),
	}

	l.EnableColor()

	if cfg.LogName == Log_Path_Stdout {
		l.Logger = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	} else if cfg.LogName == Log_Path_Stderr {
		l.Logger = log.New(os.Stderr, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	} else {
		l.Logger = log.New(os.Stderr, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)
		l.setup_file()
		l.delete_file()
	}

	return l
}

func (l *Logger) SetExpireDays(days int64) {
	l.expireDays = days
}

func (l *Logger) EnableColor() *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.colorMap[level_debug] = fmt.Sprintf("[%sDEBUG%s] ", blue, reset)
	l.colorMap[level_info] = fmt.Sprintf("[%sINFO%s] ", green, reset)
	l.colorMap[level_warn] = fmt.Sprintf("[%sWARN%s] ", red, reset)
	l.colorMap[level_fatal] = fmt.Sprintf("[%sFATAL%s] ", yellow, reset)
	return l
}

func (l *Logger) DisableColor() *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.colorMap[level_debug] = "[DEBUG] "
	l.colorMap[level_info] = "[INFO] "
	l.colorMap[level_warn] = "[WARN] "
	l.colorMap[level_fatal] = "[FATAL] "
	return l
}

func (l *Logger) pf(level int, f string, v ...any) {
	l.mu.RLock()
	color_prefix := l.colorMap[level]
	l.mu.RUnlock()
	l.Logger.Output(2, color_prefix+fmt.Sprintf(f, v...))
}

func (l *Logger) pln(level int, v ...any) {
	l.mu.RLock()
	color_prefix := l.colorMap[level]
	l.mu.RUnlock()
	l.Logger.Output(2, color_prefix+fmt.Sprintln(v...))
}

func (l *Logger) Debugf(f string, v ...any) {
	l.pf(level_debug, f, v...)
}

func (l *Logger) Debug(v ...any) {
	l.pln(level_debug, v...)
}

func (l *Logger) Infof(f string, v ...any) {
	l.pf(level_info, f, v...)
}

func (l *Logger) Info(v ...any) {
	l.pln(level_info, v...)
}

func (l *Logger) Warnf(f string, v ...any) {
	l.pf(level_warn, f, v...)
}

func (l *Logger) Warn(v ...any) {
	l.pln(level_warn, v...)
}

func (l *Logger) Fatalf(f string, v ...any) {
	l.pf(level_fatal, f, v...)
}

func (l *Logger) Fatal(v ...any) {
	l.pln(level_fatal, v...)
}

func (l *Logger) setup_file() {
	var (
		fname    string
		fp       *os.File
		deadline time.Time
		err      error
	)

	if len(l.logName) <= 0 {
		return
	}

	fname, deadline = parse_log_fname(l.logName)
	if fp, err = open_log_file(fname); err != nil {
		fp = os.Stdout
	}

	l.Fptr = fp
	l.Fname = fname
	l.Logger.SetOutput(fp)

	createTime := time.Now().Unix()
	l.oMu.Lock()
	l.oldFiles[createTime] = fname
	l.oMu.Unlock()

	go func() {
		select {
		case <-time.After(deadline.Sub(time.Now())):
			l.setup_file()
		}
	}()
}

func (l *Logger) delete_file() {
	go func() {
		for {
			time.Sleep(time.Second * 3600)

			if l.expireDays <= 0 {
				continue
			}

			exTime := time.Now().Unix() - l.expireDays*86400

			var ndel []string

			l.oMu.Lock()
			for createTime, fileName := range l.oldFiles {
				if createTime <= exTime {
					ndel = append(ndel, fileName)
					delete(l.oldFiles, createTime)
				}
			}
			l.oMu.Unlock()

			if len(ndel) > 0 {
				for _, fileName := range ndel {
					util.DeleteFile(fileName)
				}
			}
		}
	}()
}
