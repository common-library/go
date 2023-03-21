// Package redis provides redis interface.
//
// used "github.com/gomodule/redigo/redis".
package redis

import (
	"errors"
	"strings"
	"time"

	redigo_redis "github.com/gomodule/redigo/redis"
)

// Redis is object that provides redis interface.
type Redis struct {
	address  string
	password string

	maxConnection int
	timeout       int

	pool *redigo_redis.Pool

	redisConnection redigo_redis.Conn
}

// Initialize is initialize.
//
// ex) redis.Initialize("", "127.0.0.1:6379", 3, 240)
func (redis *Redis) Initialize(password string, address string, maxConnection int, timeout int) error {
	redis.address = address
	redis.password = password
	redis.maxConnection = maxConnection
	redis.timeout = timeout

	redis.pool = &redigo_redis.Pool{
		MaxIdle:     maxConnection,
		IdleTimeout: time.Duration(timeout) * time.Second,

		Dial: func() (redigo_redis.Conn, error) {
			connection, err := redigo_redis.Dial("tcp", address)
			if err != nil {
				return nil, err
			}
			return connection, err
		},
		TestOnBorrow: func(connection redigo_redis.Conn, t time.Time) error {
			_, err := connection.Do("PING")
			return err
		},
	}

	var err error
	passwordOption := redigo_redis.DialPassword(password)
	redis.redisConnection, err = redigo_redis.Dial("tcp", address, passwordOption)
	if err != nil {
		return err
	}

	_, err = redis.redisConnection.Do("PING")
	if err != nil {
		return err
	}

	return nil
}

// Finalize is finalize.
//
// ex) redis.Finalize()
func (redis *Redis) Finalize() error {
	if redis.redisConnection != nil {
		redis.redisConnection.Close()
	}

	if redis.pool != nil {
		redis.pool.Close()
	}

	return nil
}

// Ping is send ping.
//
// ex) redis.Ping()
func (redis *Redis) Ping() error {
	if redis.pool == nil {
		return errors.New("please call Initialize first")
	}

	connection := redis.pool.Get()
	defer connection.Close()

	_, err := connection.Do("PING")

	return err
}

// Select is select database.
//
// ex) err := redis.Select(0)
func (redis *Redis) Select(index int) error {
	if redis.pool == nil {
		return errors.New("please call Initialize first")
	}

	connection := redis.pool.Get()
	defer connection.Close()

	_, err := connection.Do("SELECT", index)

	return err
}

// Get is get data.
//
// ex) data, err := redis.Get(key)
func (redis *Redis) Get(key interface{}) (string, error) {
	if redis.pool == nil {
		return "", errors.New("please call Initialize first")
	}
	connection := redis.pool.Get()
	defer connection.Close()

	return redigo_redis.String(connection.Do("GET", key))
}

// Set is set data.
//
// ex) err := Set(key, value)
func (redis *Redis) Set(key interface{}, value interface{}) error {
	if redis.pool == nil {
		return errors.New("please call Initialize first")
	}

	connection := redis.pool.Get()
	defer connection.Close()

	args := make([]interface{}, 2)
	args[0] = key
	args[1] = value

	_, err := connection.Do("SET", args...)

	return err
}

// Set is set data, but it is delete after specified time.
//
// ex) err := redis.Setex("key", 2, "value")
func (redis *Redis) Setex(key interface{}, second int, value interface{}) error {
	if redis.pool == nil {
		return errors.New("please call Initialize first")
	}

	connection := redis.pool.Get()
	defer connection.Close()

	_, err := connection.Do("SETEX", key, second, value)

	return err
}

// MGet is multiple get data.
//
// ex) data, err := redis.MGet(key1, key2)
func (redis *Redis) MGet(keys ...interface{}) ([]string, error) {
	if redis.pool == nil {
		return nil, errors.New("please call Initialize first")
	}

	connection := redis.pool.Get()
	defer connection.Close()

	return redigo_redis.Strings(connection.Do("MGET", keys...))
}

// MSet is multiple set data.
//
// ex) err := redis.MSet(key1, value1, key2, value2)
func (redis *Redis) MSet(args ...interface{}) error {
	if redis.pool == nil {
		return errors.New("please call Initialize first")
	}

	connection := redis.pool.Get()
	defer connection.Close()

	_, err := connection.Do("MSET", args...)

	return err
}

// Del is delete data.
//
// ex) err := redis.Del(key)
func (redis *Redis) Del(key interface{}) error {
	if redis.pool == nil {
		return errors.New("please call Initialize first")
	}

	connection := redis.pool.Get()
	defer connection.Close()

	_, err := connection.Do("DEL", key)

	return err
}

// FlushDB is delete all data in current database.
//
// ex) err := redis.FlushDB()
func (redis *Redis) FlushDB() error {
	if redis.pool == nil {
		return errors.New("please call Initialize first")
	}

	connection := redis.pool.Get()
	defer connection.Close()

	_, err := connection.Do("FLUSHDB")

	return err
}

// FlushAll is delete all data in all database.
//
// ex) err := redis.FlushAll()
func (redis *Redis) FlushAll() error {
	if redis.pool == nil {
		return errors.New("please call Initialize first")
	}

	connection := redis.pool.Get()
	defer connection.Close()

	_, err := connection.Do("FLUSHALL")

	return err
}

// Ttl is returns valid time.
//
// If there is not exist key, -2 is returned.
//
// If the expire time is not set, -1 is returned.
//
// ex) ttl, err := redis.Ttl("key")
func (redis *Redis) Ttl(key interface{}) (int, error) {
	if redis.pool == nil {
		return -2, errors.New("please call Initialize first")
	}

	connection := redis.pool.Get()
	defer connection.Close()

	return redigo_redis.Int(connection.Do("TTL", key))
}

// Info is get redis information.
//
// kind : All, Server, Clients, Memory, Persistence, Stats, Replication, CPU, Cluster, Keyspace
//
// ex) result, err := redis.Info("ALL")
func (redis *Redis) Info(info string) (string, error) {
	if redis.pool == nil {
		return "", errors.New("please call Initialize first")
	}

	connection := redis.pool.Get()
	defer connection.Close()

	info = strings.ToLower(info)

	if info == "all" {
		return redigo_redis.String(connection.Do("INFO"))
	}

	return redigo_redis.String(connection.Do("INFO", info))
}

// DBsize is key count in current database.
//
// ex) keyCount, err := redis.DBsize()
func (redis *Redis) DBsize() (int, error) {
	if redis.pool == nil {
		return -1, errors.New("please call Initialize first")
	}

	connection := redis.pool.Get()
	defer connection.Close()

	return redigo_redis.Int(connection.Do("DBSIZE"))
}

// Exists is returns whether the keys exists.
//
// return value : exists - 1, not exists - 0
//
// ex 1) existsKey, err := redis.Exists("key")
//
// ex 2) existsKey, err := redis.Exists("key", 1, 2, "3")
func (redis *Redis) Exists(keys ...interface{}) (int, error) {
	if redis.pool == nil {
		return -1, errors.New("please call Initialize first")
	}

	connection := redis.pool.Get()
	defer connection.Close()

	return redigo_redis.Int(connection.Do("EXISTS", keys...))
}

// Rename is rename key.
//
// ex) err := redis.Rename("key", "key_rename")
func (redis *Redis) Rename(currentKey interface{}, newKey interface{}) error {
	if redis.pool == nil {
		return errors.New("please call Initialize first")
	}

	connection := redis.pool.Get()
	defer connection.Close()

	_, err := connection.Do("RENAME", currentKey, newKey)

	return err
}

// RandomKey is returns one key at random.
//
// ex) key, err := redis.RandomKey()
func (redis *Redis) RandomKey() (string, error) {
	if redis.pool == nil {
		return "", errors.New("please call Initialize first")
	}

	connection := redis.pool.Get()
	defer connection.Close()

	return redigo_redis.String(connection.Do("RANDOMKEY"))
}
