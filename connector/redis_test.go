package connector

import (
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"
)

func TestCase_Redis(t *testing.T) {
	configs := make(map[string]RedisConf)

	configs["default"] = RedisConf{
		Addr:           "127.0.0.1:6379,127.0.0.1:6379",
		Password:       "123456",
		MaxActive:      64,
		MaxIdle:        16,
		IdleTimeout:    time.Second * 240,
		ConnectTimeout: time.Second * 10,
		ReadTimeout:    time.Second,
		WriteTimeout:   time.Second,
		TestInterval:   time.Second * 60,
	}

	InitRedisConn(configs)

	conn := GetRedisConn("default")
	defer conn.Close()

	ret, err := redis.String(conn.Do("SET", "lifuxing", "123455"))
	if err != nil || ret != "OK" {
		t.Fatal(err)
	}
}
