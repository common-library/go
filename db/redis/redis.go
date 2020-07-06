// Package redis provides redis interface.
// used "github.com/gomodule/redigo/redis".
package redis

import (
	"errors"
	redigo_redis "github.com/gomodule/redigo/redis"
	"time"
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
//  ex) redis.Initialize("", "127.0.0.1:6379", 3, 240)
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
//  ex) redis.Finalize()
func (redis *Redis) Finalize() error {
	if redis.redisConnection != nil {
		redis.redisConnection.Close()
	}

	if redis.pool != nil {
		redis.pool.Close()
	}

	return nil
}

// Ping is send ping
// ex) redis.Ping()
func (redis *Redis) Ping() error {
	if redis.pool == nil {
		return errors.New("please call Initialize first")
	}

	connection := redis.pool.Get()
	defer connection.Close()

	_, err := redigo_redis.String(connection.Do("PING"))

	return err
}

// Get is get data
//  ex) data, err := redis.Get(key)
func (redis *Redis) Get(key interface{}) (string, error) {
	if redis.pool == nil {
		return "", errors.New("please call Initialize first")
	}
	connection := redis.pool.Get()
	defer connection.Close()

	return redigo_redis.String(connection.Do("GET", key))
}

// Set is set data
//  ex) err := Set(key, value)
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

// MGet is multiple get data
//  ex) data, err := redis.MGet(key1, key2)
func (redis *Redis) MGet(keys ...interface{}) ([]string, error) {
	if redis.pool == nil {
		return []string{}, errors.New("please call Initialize first")
	}

	connection := redis.pool.Get()
	defer connection.Close()

	return redigo_redis.Strings(connection.Do("MGET", keys...))
}

// MSet is multiple set data
//  ex) err := redis.MSet(key1, value1, key2, value2)
func (redis *Redis) MSet(args ...interface{}) error {
	if redis.pool == nil {
		return errors.New("please call Initialize first")
	}

	connection := redis.pool.Get()
	defer connection.Close()

	_, err := connection.Do("MSET", args...)

	return err
}

// Del is delete data
//  ex) err := redis.Del(key)
func (redis *Redis) Del(key interface{}) error {
	if redis.pool == nil {
		return errors.New("please call Initialize first")
	}

	connection := redis.pool.Get()
	defer connection.Close()

	_, err := connection.Do("DEL", key)

	return err
}
