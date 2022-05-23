package log

// default log
var Default *Logger = NewLogger(Log_Path_Stdout)

func SetDefault(l *Logger) {
	Default = l
}
