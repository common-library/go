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
func (c *Client) Initialize(address, password string, maxConnection int, timeout time.Duration) error {
	c.pool = &redigo_redis.Pool{
		Dial: func() (redigo_redis.Conn, error) {
			dialOption := redigo_redis.DialPassword(password)
			return redigo_redis.Dial("tcp", address, dialOption)
		},

		TestOnBorrow: func(connection redigo_redis.Conn, t time.Time) error {
			_, err := connection.Do("PING")
			return err
		},

		MaxIdle: maxConnection,

		IdleTimeout: timeout * time.Second,
	}

	return c.Ping()
}

// Finalize is finalize.
//
// ex) err := client.Finalize()
func (c *Client) Finalize() error {
	if c.connection != nil {
		c.connection.Close()
	}

	if c.pool != nil {
		c.pool.Close()
	}

	return nil
}

// Ping is send ping.
//
// ex) client.Ping()
func (c *Client) Ping() error {
	if c.pool == nil {
		return errors.New("please call Initialize first")
	}

	_, err := c.do("PING")

	return err
}

// Select is select database.
//
// ex) err := client.Select(0)
func (c *Client) Select(index int) error {
	if c.pool == nil {
		return errors.New("please call Initialize first")
	}

	_, err := c.do("SELECT", index)

	return err
}

// Get is get data.
//
// ex) data, err := client.Get(key)
func (c *Client) Get(key any) (string, error) {
	if c.pool == nil {
		return "", errors.New("please call Initialize first")
	}

	return redigo_redis.String(c.do("GET", key))
}

// Set is set data.
//
// ex) err := client.Set(key, value)
func (c *Client) Set(key, value any) error {
	if c.pool == nil {
		return errors.New("please call Initialize first")
	}

	args := make([]any, 2)
	args[0] = key
	args[1] = value

	_, err := c.do("SET", args...)

	return err
}

// Set is set data, but it is delete after specified time.
//
// ex) err := client.Setex("key", 2, "value")
func (c *Client) Setex(key any, second int, value any) error {
	if c.pool == nil {
		return errors.New("please call Initialize first")
	}

	_, err := c.do("SETEX", key, second, value)

	return err
}

// MGet is multiple get data.
//
// ex) data, err := client.MGet(key1, key2)
func (c *Client) MGet(keys ...any) ([]string, error) {
	if c.pool == nil {
		return nil, errors.New("please call Initialize first")
	}

	return redigo_redis.Strings(c.do("MGET", keys...))
}

// MSet is multiple set data.
//
// ex) err := client.MSet(key1, value1, key2, value2)
func (c *Client) MSet(args ...any) error {
	if c.pool == nil {
		return errors.New("please call Initialize first")
	}

	_, err := c.do("MSET", args...)

	return err
}

// Del is delete data.
//
// ex) err := client.Del(key)
func (c *Client) Del(key any) error {
	if c.pool == nil {
		return errors.New("please call Initialize first")
	}

	_, err := c.do("DEL", key)

	return err
}

// FlushDB is delete all data in current database.
//
// ex) err := client.FlushDB()
func (c *Client) FlushDB() error {
	if c.pool == nil {
		return errors.New("please call Initialize first")
	}

	_, err := c.do("FLUSHDB")

	return err
}

// FlushAll is delete all data in all database.
//
// ex) err := client.FlushAll()
func (c *Client) FlushAll() error {
	if c.pool == nil {
		return errors.New("please call Initialize first")
	}

	_, err := c.do("FLUSHALL")

	return err
}

// Ttl is returns valid time.
//
// If there is not exist key, -2 is returned.
//
// If the expire time is not set, -1 is returned.
//
// ex) ttl, err := client.Ttl("key")
func (c *Client) Ttl(key any) (int, error) {
	if c.pool == nil {
		return -2, errors.New("please call Initialize first")
	}

	return redigo_redis.Int(c.do("TTL", key))
}

// Info is get redis information.
//
// kind : All, Server, Clients, Memory, Persistence, Stats, Replication, CPU, Cluster, Keyspace
//
// ex) result, err := client.Info("ALL")
func (c *Client) Info(info string) (string, error) {
	if c.pool == nil {
		return "", errors.New("please call Initialize first")
	}

	info = strings.ToLower(info)

	if info == "all" {
		return redigo_redis.String(c.do("INFO"))
	}

	return redigo_redis.String(c.do("INFO", info))
}

// DBsize is key count in current database.
//
// ex) keyCount, err := client.DBsize()
func (c *Client) DBsize() (int, error) {
	if c.pool == nil {
		return -1, errors.New("please call Initialize first")
	}

	return redigo_redis.Int(c.do("DBSIZE"))
}

// Exists is returns whether the keys exists.
//
// ex 1) existsKey, err := client.Exists("key")
// ex 2) existsKey, err := client.Exists("key", 1, 2, "3")
func (c *Client) Exists(keys ...any) (bool, error) {
	if c.pool == nil {
		return false, errors.New("please call Initialize first")
	}

	return redigo_redis.Bool(c.do("EXISTS", keys...))
}

// Rename is rename key.
//
// ex) err := client.Rename("key", "key_rename")
func (c *Client) Rename(currentKey, newKey any) error {
	if c.pool == nil {
		return errors.New("please call Initialize first")
	}

	_, err := c.do("RENAME", currentKey, newKey)
	return err
}

// RandomKey is returns one key at random.
//
// ex) key, err := client.RandomKey()
func (c *Client) RandomKey() (string, error) {
	if c.pool == nil {
		return "", errors.New("please call Initialize first")
	}

	return redigo_redis.String(c.do("RANDOMKEY"))
}

func (c *Client) do(command string, args ...any) (any, error) {
	connection := c.pool.Get()
	defer connection.Close()

	return connection.Do(command, args...)
}
