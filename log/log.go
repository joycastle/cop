package log

import (
	"fmt"
	"log"
	"os"
	"time"
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

var (
	color_debug = fmt.Sprintf("[%sDEBUG%s] ", blue, reset)
	color_info  = fmt.Sprintf("[%sINFO%s] ", green, reset)
	color_warn  = fmt.Sprintf("[%sWARN%s] ", red, reset)
	color_fatal = fmt.Sprintf("[%sFATAL%s] ", yellow, reset)

	color_close_debug = "[DEBUG] "
	color_close_info  = "[INFO] "
	color_close_warn  = "[WARN] "
	color_close_fatal = "[FATAL] "

	color_debug_using = color_debug
	color_info_using  = color_info
	color_warn_using  = color_warn
	color_fatal_using = color_fatal
)

type Logger struct {
	Logger *log.Logger
	fpath  string
	Fptr   *os.File
	Fname  string
}

func DisableColor() {
	color_debug_using = color_close_debug
	color_info_using = color_close_info
	color_warn_using = color_close_warn
	color_fatal_using = color_close_fatal
}

func EnableColor() {
	color_debug_using = color_debug
	color_info_using = color_info
	color_warn_using = color_warn
	color_fatal_using = color_fatal
}

func NewLoggerDefault() *Logger {
	l := &Logger{
		fpath: "",
		Fptr:  os.Stdout,
		Fname: "",
	}

	l.Logger = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)

	return l
}

func NewLogger(path string) *Logger {
	l := &Logger{
		fpath: path,
	}

	l.Logger = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	l.setup_file()

	return l
}

func (l *Logger) Printf(f string, v ...any) {
	l.Logger.Output(2, color_info_using+fmt.Sprintf(f, v...))
}

func (l *Logger) Println(v ...any) {
	l.Logger.Output(2, color_info_using+fmt.Sprintln(v...))
}

func (l *Logger) Debugf(f string, v ...any) {
	l.Logger.Output(2, color_debug_using+fmt.Sprintf(f, v...))
}

func (l *Logger) Debug(v ...any) {
	l.Logger.Output(2, color_debug_using+fmt.Sprintln(v...))
}

func (l *Logger) Infof(f string, v ...any) {
	l.Logger.Output(2, color_info_using+fmt.Sprintf(f, v...))
}

func (l *Logger) Info(v ...any) {
	l.Logger.Output(2, color_info_using+fmt.Sprintln(v...))
}

func (l *Logger) Warnf(f string, v ...any) {
	l.Logger.Output(2, color_warn_using+fmt.Sprintf(f, v...))
}

func (l *Logger) Warn(v ...any) {
	l.Logger.Output(2, color_warn_using+fmt.Sprintln(v...))
}

func (l *Logger) Fatalf(f string, v ...any) {
	l.Logger.Output(2, color_fatal_using+fmt.Sprintf(f, v...))
}

func (l *Logger) Fatal(v ...any) {
	l.Logger.Output(2, color_fatal_using+fmt.Sprintln(v...))
}

func (l *Logger) setup_file() {
	var (
		fname    string
		fp       *os.File
		deadline time.Time
		err      error
	)

	if len(l.fpath) <= 0 {
		return
	}

	fname, deadline = parse_log_fname(l.fpath)
	if fp, err = open_log_file(fname); err != nil {
		fp = os.Stdout
	}

	l.Fptr = fp
	l.Fname = fname
	l.Logger.SetOutput(fp)

	go func() {
		select {
		case <-time.After(deadline.Sub(time.Now())):
			l.setup_file()
		}
	}()
}
