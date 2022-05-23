package log

// default log
var Default *Logger = NewLogger(LogConf{Log_Path_Stdout, 0})

func SetDefault(l *Logger) {
	Default = l
}
