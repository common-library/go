// Package redis provides Redis client implementations.
package redis

import (
	"errors"
	"strings"
	"time"

	redigo_redis "github.com/gomodule/redigo/redis"
)

// Client is a struct that provides client related methods.
type Client struct {
	address  string
	password string

	maxConnection int
	timeout       int

	pool *redigo_redis.Pool

	redisConnection redigo_redis.Conn
}

// Initialize is initialize.
//
// ex) client.Initialize("", "127.0.0.1:6379", 3, 240)
func (this *Client) Initialize(password string, address string, maxConnection int, timeout int) error {
	this.address = address
	this.password = password
	this.maxConnection = maxConnection
	this.timeout = timeout

	this.pool = &redigo_redis.Pool{
		MaxIdle:     maxConnection,
		IdleTimeout: time.Duration(timeout) * time.Second,

		Dial: func() (redigo_redis.Conn, error) {
			return redigo_redis.Dial("tcp", address)
		},
		TestOnBorrow: func(connection redigo_redis.Conn, t time.Time) error {
			_, err := connection.Do("PING")
			return err
		},
	}

	var err error
	passwordOption := redigo_redis.DialPassword(password)
	this.redisConnection, err = redigo_redis.Dial("tcp", address, passwordOption)
	if err != nil {
		return err
	}

	_, err = this.redisConnection.Do("PING")
	return err
}

// Finalize is finalize.
//
// ex) client.Finalize()
func (this *Client) Finalize() error {
	if this.redisConnection != nil {
		this.redisConnection.Close()
	}

	if this.pool != nil {
		this.pool.Close()
	}

	return nil
}

// Ping is send ping.
//
// ex) client.Ping()
func (this *Client) Ping() error {
	if this.pool == nil {
		return errors.New("please call Initialize first")
	}

	connection := this.pool.Get()
	defer connection.Close()

	_, err := connection.Do("PING")

	return err
}

// Select is select database.
//
// ex) err := client.Select(0)
func (this *Client) Select(index int) error {
	if this.pool == nil {
		return errors.New("please call Initialize first")
	}

	connection := this.pool.Get()
	defer connection.Close()

	_, err := connection.Do("SELECT", index)

	return err
}

// Get is get data.
//
// ex) data, err := client.Get(key)
func (this *Client) Get(key interface{}) (string, error) {
	if this.pool == nil {
		return "", errors.New("please call Initialize first")
	}
	connection := this.pool.Get()
	defer connection.Close()

	return redigo_redis.String(connection.Do("GET", key))
}

// Set is set data.
//
// ex) err := client.Set(key, value)
func (this *Client) Set(key interface{}, value interface{}) error {
	if this.pool == nil {
		return errors.New("please call Initialize first")
	}

	connection := this.pool.Get()
	defer connection.Close()

	args := make([]interface{}, 2)
	args[0] = key
	args[1] = value

	_, err := connection.Do("SET", args...)

	return err
}

// Set is set data, but it is delete after specified time.
//
// ex) err := client.Setex("key", 2, "value")
func (this *Client) Setex(key interface{}, second int, value interface{}) error {
	if this.pool == nil {
		return errors.New("please call Initialize first")
	}

	connection := this.pool.Get()
	defer connection.Close()

	_, err := connection.Do("SETEX", key, second, value)

	return err
}

// MGet is multiple get data.
//
// ex) data, err := client.MGet(key1, key2)
func (this *Client) MGet(keys ...interface{}) ([]string, error) {
	if this.pool == nil {
		return nil, errors.New("please call Initialize first")
	}

	connection := this.pool.Get()
	defer connection.Close()

	return redigo_redis.Strings(connection.Do("MGET", keys...))
}

// MSet is multiple set data.
//
// ex) err := client.MSet(key1, value1, key2, value2)
func (this *Client) MSet(args ...interface{}) error {
	if this.pool == nil {
		return errors.New("please call Initialize first")
	}

	connection := this.pool.Get()
	defer connection.Close()

	_, err := connection.Do("MSET", args...)

	return err
}

// Del is delete data.
//
// ex) err := client.Del(key)
func (this *Client) Del(key interface{}) error {
	if this.pool == nil {
		return errors.New("please call Initialize first")
	}

	connection := this.pool.Get()
	defer connection.Close()

	_, err := connection.Do("DEL", key)

	return err
}

// FlushDB is delete all data in current database.
//
// ex) err := client.FlushDB()
func (this *Client) FlushDB() error {
	if this.pool == nil {
		return errors.New("please call Initialize first")
	}

	connection := this.pool.Get()
	defer connection.Close()

	_, err := connection.Do("FLUSHDB")

	return err
}

// FlushAll is delete all data in all database.
//
// ex) err := client.FlushAll()
func (this *Client) FlushAll() error {
	if this.pool == nil {
		return errors.New("please call Initialize first")
	}

	connection := this.pool.Get()
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
// ex) ttl, err := client.Ttl("key")
func (this *Client) Ttl(key interface{}) (int, error) {
	if this.pool == nil {
		return -2, errors.New("please call Initialize first")
	}

	connection := this.pool.Get()
	defer connection.Close()

	return redigo_redis.Int(connection.Do("TTL", key))
}

// Info is get redis information.
//
// kind : All, Server, Clients, Memory, Persistence, Stats, Replication, CPU, Cluster, Keyspace
//
// ex) result, err := client.Info("ALL")
func (this *Client) Info(info string) (string, error) {
	if this.pool == nil {
		return "", errors.New("please call Initialize first")
	}

	connection := this.pool.Get()
	defer connection.Close()

	info = strings.ToLower(info)

	if info == "all" {
		return redigo_redis.String(connection.Do("INFO"))
	}

	return redigo_redis.String(connection.Do("INFO", info))
}

// DBsize is key count in current database.
//
// ex) keyCount, err := client.DBsize()
func (this *Client) DBsize() (int, error) {
	if this.pool == nil {
		return -1, errors.New("please call Initialize first")
	}

	connection := this.pool.Get()
	defer connection.Close()

	return redigo_redis.Int(connection.Do("DBSIZE"))
}

// Exists is returns whether the keys exists.
//
// return value : exists - 1, not exists - 0
//
// ex 1) existsKey, err := client.Exists("key")
//
// ex 2) existsKey, err := client.Exists("key", 1, 2, "3")
func (this *Client) Exists(keys ...interface{}) (int, error) {
	if this.pool == nil {
		return -1, errors.New("please call Initialize first")
	}

	connection := this.pool.Get()
	defer connection.Close()

	return redigo_redis.Int(connection.Do("EXISTS", keys...))
}

// Rename is rename key.
//
// ex) err := client.Rename("key", "key_rename")
func (this *Client) Rename(currentKey interface{}, newKey interface{}) error {
	if this.pool == nil {
		return errors.New("please call Initialize first")
	}

	connection := this.pool.Get()
	defer connection.Close()

	_, err := connection.Do("RENAME", currentKey, newKey)

	return err
}

// RandomKey is returns one key at random.
//
// ex) key, err := client.RandomKey()
func (this *Client) RandomKey() (string, error) {
	if this.pool == nil {
		return "", errors.New("please call Initialize first")
	}

	connection := this.pool.Get()
	defer connection.Close()

	return redigo_redis.String(connection.Do("RANDOMKEY"))
}
