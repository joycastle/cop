package connector

import "github.com/joycastle/cop/log"

var LogErr *log.Logger
var LogMonitor *log.Logger

func init() {
	LogErr = log.Default
	LogMonitor = log.Default
}

func SetErrorLogger(l *log.Logger) {
	LogErr = l
}

func SetMonitorLogger(l *log.Logger) {
	LogMonitor = l
}
