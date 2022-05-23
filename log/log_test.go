package log

import (
	"errors"
	"testing"
)

func Test_log(t *testing.T) {
	// t.Fatal("not implemented")
	log := NewLogger(Log_Path_Stdout)
	log.Infof("This is Logger:%s", "hello")
	log.Info("This is Logger", errors.New("Hello"))
	log.Debugf("This is Logger:%s", "hello")
	log.Debug("This is Logger", errors.New("Hello"))
	log.Warnf("This is Logger:%s", "hello")
	log.Warn("This is Logger", errors.New("Hello"))
	log.Fatalf("This is Logger:%s", "hello")
	log.Fatal("This is Logger", errors.New("Hello"))

	log = NewLogger(Log_Path_Stderr).DisableColor()
	log.Infof("This is Logger:%s", "hello")
	log.Info("This is Logger", errors.New("Hello"))
	log.Debugf("This is Logger:%s", "hello")
	log.Debug("This is Logger", errors.New("Hello"))
	log.Warnf("This is Logger:%s", "hello")
	log.Warn("This is Logger", errors.New("Hello"))
	log.Fatalf("This is Logger:%s", "hello")
	log.Fatal("This is Logger", errors.New("Hello"))
}

func Benchmark_log(b *testing.B) {
	log := NewLogger("./log/test.go-*-*-*").EnableColor()
	for n := 0; n < b.N; n++ {
		log.Fatalf("This is Logger:%s", "hello")
	}
}
