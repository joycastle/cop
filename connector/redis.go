package connector

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
)

type RedisConf struct {
	Addr           string        `yaml:"Addr"`
	Password       string        `yaml:"Password"`
	MaxActive      int           `yaml:"MaxActive"`
	MaxIdle        int           `yaml:"MaxIdle"`
	IdleTimeout    time.Duration `yaml:"IdleTimeout"`
	ConnectTimeout time.Duration `yaml:"ConnectTimeout"`
	ReadTimeout    time.Duration `yaml:"ReadTimeout"`
	WriteTimeout   time.Duration `yaml:"WriteTimeout"`
	TestInterval   time.Duration `yaml:"TestInterval"`
}

var (
	ErrRedisTemplete string                 = "[REDIS] %s"
	redisConnMapping map[string]*redis.Pool = make(map[string]*redis.Pool)
)

func InitRedisConn(configs map[string]RedisConf) {
	for sn, config := range configs {
		rp := &redis.Pool{
			MaxActive:   config.MaxActive,
			MaxIdle:     config.MaxIdle,
			IdleTimeout: config.IdleTimeout,
			Wait:        true,
			Dial: func() (redis.Conn, error) {
				var (
					addrs []string
					addr  string
					conn  redis.Conn
					err   error
				)
				addrs = strings.Split(config.Addr, ",")
				addr = addrs[rand.Intn(len(addrs))]
				conn, err = redis.DialTimeout("tcp", addr, config.ConnectTimeout, config.ReadTimeout, config.WriteTimeout)

				if err != nil {
					LogErr.Warnf(ErrRedisTemplete, fmt.Sprintf("connect failed: %s, %s", addr, err))
					return nil, err
				}

				if config.Password != "" {
					if _, err := conn.Do("AUTH", config.Password); err != nil {
						LogErr.Warnf(ErrRedisTemplete, fmt.Sprintf("connect failed: %s, %s", addr, err))
						conn.Close()
						return nil, err
					}
				}

				_, err = conn.Do("PING")
				if err != nil {
					LogErr.Warnf("connect to redis[%s] failed: %s", addr, err)
					return nil, err
				}

				return conn, nil
			},
			TestOnBorrow: func(conn redis.Conn, t time.Time) error {
				if time.Since(t) < config.TestInterval {
					return nil
				}
				_, err := conn.Do("PING")
				return err
			},
		}

		go func() {
			// PoolStats contains pool statistics.
			//type PoolStats struct {
			// ActiveCount is the number of connections in the pool. The count includes
			// idle connections and connections in use.
			// ActiveCount int
			// IdleCount is the number of idle connections in the pool.
			//IdleCount int
			//}
			for {
				time.Sleep(time.Second * 10)
				stat := rp.Stats()
				infs := fmt.Sprintf("Redis Pool ActiveCount:%d, IdleCount:%d node:%s", stat.ActiveCount, stat.IdleCount, sn)

				LogMonitor.Info(infs)
			}
		}()

		redisConnMapping[sn] = rp
	}
}

func GetRedisConn(sn string) redis.Conn {
	if v, ok := redisConnMapping[sn]; ok {
		return v.Get()
	}
	return nil
}
