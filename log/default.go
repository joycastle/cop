package log

// default log
var Default *Logger = NewLogger(Log_Path_Stdout)

// error log
var Error *Logger = NewLogger(Log_Path_Stderr)

// run log
var Run *Logger = Default

// monitor log
var Monitor *Logger = Default

func SetDefault(l *Logger) {
	Default = l
}

func SetError(l *Logger) {
	Error = l
}

func SetRun(l *Logger) {
	Run = l
}

func SetMonitor(l *Logger) {
	Run = l
}
