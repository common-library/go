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
	pool *redigo_redis.Pool

	connection redigo_redis.Conn
}

// Initialize is initialize.
//
// ex) err := client.Initialize("127.0.0.1:6379", "", 10, 60)
func (this *Client) Initialize(address, password string, maxConnection int, timeout int64) error {
	this.pool = &redigo_redis.Pool{
		Dial: func() (redigo_redis.Conn, error) {
			dialOption := redigo_redis.DialPassword(password)
			return redigo_redis.Dial("tcp", address, dialOption)
		},

		TestOnBorrow: func(connection redigo_redis.Conn, t time.Time) error {
			_, err := connection.Do("PING")
			return err
		},

		MaxIdle: maxConnection,

		IdleTimeout: time.Duration(timeout) * time.Second,
	}

	return this.Ping()
}

// Finalize is finalize.
//
// ex) err := client.Finalize()
func (this *Client) Finalize() error {
	if this.connection != nil {
		this.connection.Close()
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

	_, err := this.do("PING")

	return err
}

// Select is select database.
//
// ex) err := client.Select(0)
func (this *Client) Select(index int) error {
	if this.pool == nil {
		return errors.New("please call Initialize first")
	}

	_, err := this.do("SELECT", index)

	return err
}

// Get is get data.
//
// ex) data, err := client.Get(key)
func (this *Client) Get(key any) (string, error) {
	if this.pool == nil {
		return "", errors.New("please call Initialize first")
	}

	return redigo_redis.String(this.do("GET", key))
}

// Set is set data.
//
// ex) err := client.Set(key, value)
func (this *Client) Set(key, value any) error {
	if this.pool == nil {
		return errors.New("please call Initialize first")
	}

	args := make([]any, 2)
	args[0] = key
	args[1] = value

	_, err := this.do("SET", args...)

	return err
}

// Set is set data, but it is delete after specified time.
//
// ex) err := client.Setex("key", 2, "value")
func (this *Client) Setex(key any, second int, value any) error {
	if this.pool == nil {
		return errors.New("please call Initialize first")
	}

	_, err := this.do("SETEX", key, second, value)

	return err
}

// MGet is multiple get data.
//
// ex) data, err := client.MGet(key1, key2)
func (this *Client) MGet(keys ...any) ([]string, error) {
	if this.pool == nil {
		return nil, errors.New("please call Initialize first")
	}

	return redigo_redis.Strings(this.do("MGET", keys...))
}

// MSet is multiple set data.
//
// ex) err := client.MSet(key1, value1, key2, value2)
func (this *Client) MSet(args ...any) error {
	if this.pool == nil {
		return errors.New("please call Initialize first")
	}

	_, err := this.do("MSET", args...)

	return err
}

// Del is delete data.
//
// ex) err := client.Del(key)
func (this *Client) Del(key any) error {
	if this.pool == nil {
		return errors.New("please call Initialize first")
	}

	_, err := this.do("DEL", key)

	return err
}

// FlushDB is delete all data in current database.
//
// ex) err := client.FlushDB()
func (this *Client) FlushDB() error {
	if this.pool == nil {
		return errors.New("please call Initialize first")
	}

	_, err := this.do("FLUSHDB")

	return err
}

// FlushAll is delete all data in all database.
//
// ex) err := client.FlushAll()
func (this *Client) FlushAll() error {
	if this.pool == nil {
		return errors.New("please call Initialize first")
	}

	_, err := this.do("FLUSHALL")

	return err
}

// Ttl is returns valid time.
//
// If there is not exist key, -2 is returned.
//
// If the expire time is not set, -1 is returned.
//
// ex) ttl, err := client.Ttl("key")
func (this *Client) Ttl(key any) (int, error) {
	if this.pool == nil {
		return -2, errors.New("please call Initialize first")
	}

	return redigo_redis.Int(this.do("TTL", key))
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

	info = strings.ToLower(info)

	if info == "all" {
		return redigo_redis.String(this.do("INFO"))
	}

	return redigo_redis.String(this.do("INFO", info))
}

// DBsize is key count in current database.
//
// ex) keyCount, err := client.DBsize()
func (this *Client) DBsize() (int, error) {
	if this.pool == nil {
		return -1, errors.New("please call Initialize first")
	}

	return redigo_redis.Int(this.do("DBSIZE"))
}

// Exists is returns whether the keys exists.
//
// ex 1) existsKey, err := client.Exists("key")
// ex 2) existsKey, err := client.Exists("key", 1, 2, "3")
func (this *Client) Exists(keys ...any) (bool, error) {
	if this.pool == nil {
		return false, errors.New("please call Initialize first")
	}

	return redigo_redis.Bool(this.do("EXISTS", keys...))
}

// Rename is rename key.
//
// ex) err := client.Rename("key", "key_rename")
func (this *Client) Rename(currentKey, newKey any) error {
	if this.pool == nil {
		return errors.New("please call Initialize first")
	}

	_, err := this.do("RENAME", currentKey, newKey)
	return err
}

// RandomKey is returns one key at random.
//
// ex) key, err := client.RandomKey()
func (this *Client) RandomKey() (string, error) {
	if this.pool == nil {
		return "", errors.New("please call Initialize first")
	}

	return redigo_redis.String(this.do("RANDOMKEY"))
}

func (this *Client) do(command string, args ...any) (any, error) {
	connection := this.pool.Get()
	defer connection.Close()

	return connection.Do(command, args...)
}
