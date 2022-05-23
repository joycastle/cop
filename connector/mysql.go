package connector

import (
	"errors"
	"fmt"
	"time"

	"github.com/joycastle/cop/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	ErrMysqlTemplete      = "[MYSQL] %s"
	ErrMysqlNodeNotExists = errors.New("[MYSQL] node not exists")

	loggerSlow       *log.Logger
	MysqlConnMapping map[string]MysqlConn = make(map[string]MysqlConn)
)

func init() {
	loggerSlow = log.Default
}

type MysqlConf struct {
	Dsn         string        `yaml:"Dsn"`
	MaxIdle     int           `yaml:"MaxIdle"`
	MaxOpen     int           `yaml:"MaxOpen"`
	MaxLifeTime time.Duration `yaml:"MaxLifeTime"`
}

type MysqlNodeConf struct {
	Master MysqlConf `yaml:"Master"`
	Slave  MysqlConf `yaml:"Slave"`
}

type MysqlConn struct {
	Master *gorm.DB
	Slave  *gorm.DB
}

func SetMysqlSlowLogger(l *log.Logger) {
	loggerSlow = l
}

func InitMysqlConn(configs map[string]MysqlNodeConf) error {
	for node, cfg := range configs {
		var nConn MysqlConn
		if conn, err := GetMysqlConn(cfg.Master); err != nil {
			return fmt.Errorf(ErrMysqlTemplete, "connect error [master]: %s %s", err, cfg.Master.Dsn)
		} else {
			nConn.Master = conn
		}

		if conn, err := GetMysqlConn(cfg.Slave); err != nil {
			return fmt.Errorf(ErrMysqlTemplete, "connect error [slave]: %s %s", err, cfg.Master.Dsn)
		} else {
			nConn.Slave = conn
		}

		MysqlConnMapping[node] = nConn
	}

	return nil
}

func GetMysqlMaster(node string) (*gorm.DB, error) {
	if v, ok := MysqlConnMapping[node]; ok {
		return v.Master, nil
	}
	return nil, ErrMysqlNodeNotExists
}

func GetMysqlSlave(node string) (*gorm.DB, error) {
	if v, ok := MysqlConnMapping[node]; ok {
		return v.Slave, nil
	}
	return nil, ErrMysqlNodeNotExists
}

func GetMysqlConn(config MysqlConf) (*gorm.DB, error) {

	slowLogger := logger.New(loggerSlow.Logger, logger.Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  logger.Warn,
		IgnoreRecordNotFoundError: false,
		Colorful:                  true,
	})

	gdb, err := gorm.Open(mysql.Open(config.Dsn), &gorm.Config{Logger: slowLogger})
	if err != nil {
		return nil, err
	}

	if sqlDb, err := gdb.DB(); err != nil {
		return nil, err
	} else {
		sqlDb.SetMaxIdleConns(config.MaxIdle)
		sqlDb.SetMaxOpenConns(config.MaxOpen)
		sqlDb.SetConnMaxLifetime(config.MaxLifeTime)

		//stats monitor
		go func() {
			//type DBStats struct {
			//	MaxOpenConnections int // Maximum number of open connections to the database.

			// Pool Status
			//	OpenConnections int // The number of established connections both in use and idle.
			//	InUse           int // The number of connections currently in use.
			//	Idle            int // The number of idle connections.

			// Counters
			//	WaitCount         int64         // The total number of connections waited for.
			//	WaitDuration      time.Duration // The total time blocked waiting for a new connection.
			//	MaxIdleClosed     int64         // The total number of connections closed due to SetMaxIdleConns.
			//	MaxIdleTimeClosed int64         // The total number of connections closed due to SetConnMaxIdleTime.
			//	MaxLifetimeClosed int64         // The total number of connections closed due to SetConnMaxLifetime.
			//}
			for {
				time.Sleep(time.Second * 10)
				stat := sqlDb.Stats()
				infos := fmt.Sprintf("MYSQL Connection open:%d, inUse:%d, idle:%d, waitCount:%d, waitDuration:%v dsn:%s",
					stat.OpenConnections,
					stat.InUse,
					stat.Idle,
					stat.WaitCount,
					stat.WaitDuration,
					config.Dsn)

				LogMonitor.Infof(infos)
			}
		}()
	}

	return gdb, nil
}
