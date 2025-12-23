// Package redis provides a Redis client implementation with connection pooling.
//
// This package offers a convenient wrapper around the Redigo Redis client library,
// providing connection pooling, automatic reconnection, and simplified method signatures
// for common Redis operations.
//
// Features:
//   - Connection pooling with configurable size and timeout
//   - Automatic health checks via PING
//   - Support for all common Redis data types and operations
//   - Database selection
//   - Key expiration (TTL) management
//   - Batch operations (MGET, MSET)
//
// Example:
//
//	var client redis.Client
//	client.Initialize("localhost:6379", "", 10, 60*time.Second)
//	defer client.Finalize()
//	client.Set("key", "value")
//	value, _ := client.Get("key")
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

// Initialize initializes the Redis client with connection pool settings.
//
// Parameters:
//   - address: Redis server address in the format "host:port" (e.g., "localhost:6379")
//   - password: Authentication password. Use empty string "" for no authentication
//   - maxConnection: Maximum number of idle connections in the pool
//   - timeout: Idle connection timeout duration. Connections idle longer than this will be closed
//
// Returns:
//   - error: Error if connection test fails, nil on success
//
// The function creates a connection pool and validates connectivity by sending a PING command.
// All connections in the pool are tested on borrow using PING to ensure they are still alive.
//
// Example:
//
//	var client redis.Client
//	err := client.Initialize("localhost:6379", "", 10, 60*time.Second)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer client.Finalize()
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

// Finalize closes all connections and cleans up resources.
//
// Returns:
//   - error: Always returns nil
//
// This function should be called when the client is no longer needed, typically using defer
// after Initialize. It closes both the active connection (if any) and the connection pool.
//
// Example:
//
//	var client redis.Client
//	client.Initialize("localhost:6379", "", 10, 60*time.Second)
//	defer client.Finalize() // Ensure cleanup
func (c *Client) Finalize() error {
	if c.connection != nil {
		c.connection.Close()
	}

	if c.pool != nil {
		c.pool.Close()
	}

	return nil
}

// Ping tests the connection to the Redis server.
//
// Returns:
//   - error: Error if client not initialized or PING command fails, nil if connection is healthy
//
// This function sends a PING command to the Redis server to verify connectivity.
// It returns an error if Initialize has not been called first.
//
// Example:
//
//	err := client.Ping()
//	if err != nil {
//	    log.Println("Redis connection failed:", err)
//	}
func (c *Client) Ping() error {
	if c.pool == nil {
		return errors.New("please call Initialize first")
	}

	_, err := c.do("PING")

	return err
}

// Select switches to the specified database by index.
//
// Parameters:
//   - index: Database index to select (default Redis configuration supports 0-15)
//
// Returns:
//   - error: Error if client not initialized or SELECT command fails, nil on success
//
// Redis supports multiple databases identified by a numeric index. The default database is 0.
// This function changes the currently active database for subsequent operations.
//
// Example:
//
//	err := client.Select(1) // Switch to database 1
//	if err != nil {
//	    log.Fatal(err)
//	}
func (c *Client) Select(index int) error {
	if c.pool == nil {
		return errors.New("please call Initialize first")
	}

	_, err := c.do("SELECT", index)

	return err
}

// Get retrieves the value of the specified key.
//
// Parameters:
//   - key: The key to retrieve (can be string, int, or any type convertible to string)
//
// Returns:
//   - string: The value stored at the key, or empty string if key doesn't exist
//   - error: Error if client not initialized or GET command fails, nil on success
//
// If the key does not exist, Redis returns a nil reply which is converted to an empty string.
// If the key holds a value that is not a string, an error is returned.
//
// Example:
//
//	value, err := client.Get("user:1:name")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(value) // "Alice"
func (c *Client) Get(key any) (string, error) {
	if c.pool == nil {
		return "", errors.New("please call Initialize first")
	}

	return redigo_redis.String(c.do("GET", key))
}

// Set sets the key to hold the specified string value.
//
// Parameters:
//   - key: The key to set (can be string, int, or any type convertible to string)
//   - value: The value to store (can be string, int, or any type convertible to string)
//
// Returns:
//   - error: Error if client not initialized or SET command fails, nil on success
//
// If the key already holds a value, it is overwritten. Any existing time-to-live (TTL)
// associated with the key is discarded.
//
// Example:
//
//	err := client.Set("user:1:name", "Alice")
//	if err != nil {
//	    log.Fatal(err)
//	}
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

// Setex sets the key to hold the value with an expiration time in seconds.
//
// Parameters:
//   - key: The key to set (can be string, int, or any type convertible to string)
//   - second: Time to live in seconds before the key is automatically deleted
//   - value: The value to store (can be string, int, or any type convertible to string)
//
// Returns:
//   - error: Error if client not initialized or SETEX command fails, nil on success
//
// This is an atomic operation that sets both the value and expiration time. The key will be
// automatically deleted after the specified number of seconds. This is useful for temporary
// data like sessions, caches, or rate limiting.
//
// Example:
//
//	// Store session data that expires in 1 hour (3600 seconds)
//	err := client.Setex("session:abc123", 3600, "user_data")
//	if err != nil {
//	    log.Fatal(err)
//	}
func (c *Client) Setex(key any, second int, value any) error {
	if c.pool == nil {
		return errors.New("please call Initialize first")
	}

	_, err := c.do("SETEX", key, second, value)

	return err
}

// MGet retrieves the values of all specified keys.
//
// Parameters:
//   - keys: Variable number of keys to retrieve
//
// Returns:
//   - []string: Slice of values corresponding to the keys. Empty string for non-existing keys
//   - error: Error if client not initialized or MGET command fails, nil on success
//
// For every key that does not hold a string value or does not exist, the slice contains
// an empty string. The returned values are in the same order as the requested keys.
//
// Example:
//
//	values, err := client.MGet("key1", "key2", "key3")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for i, val := range values {
//	    fmt.Printf("Key %d: %s\n", i+1, val)
//	}
func (c *Client) MGet(keys ...any) ([]string, error) {
	if c.pool == nil {
		return nil, errors.New("please call Initialize first")
	}

	return redigo_redis.Strings(c.do("MGET", keys...))
}

// MSet sets multiple keys to multiple values in a single atomic operation.
//
// Parameters:
//   - args: Variable number of alternating key-value pairs (key1, value1, key2, value2, ...)
//
// Returns:
//   - error: Error if client not initialized or MSET command fails, nil on success
//
// This operation is atomic, meaning all keys are set at once. Existing values are overwritten.
// The number of arguments must be even (pairs of keys and values).
//
// Example:
//
//	err := client.MSet(
//	    "user:1:name", "Alice",
//	    "user:1:email", "alice@example.com",
//	    "user:1:age", "30",
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
func (c *Client) MSet(args ...any) error {
	if c.pool == nil {
		return errors.New("please call Initialize first")
	}

	_, err := c.do("MSET", args...)

	return err
}

// Del deletes the specified key.
//
// Parameters:
//   - key: The key to delete (can be string, int, or any type convertible to string)
//
// Returns:
//   - error: Error if client not initialized or DEL command fails, nil on success
//
// If the key does not exist, the operation succeeds without error. This function removes
// the key and its associated value from the database.
//
// Example:
//
//	err := client.Del("session:abc123")
//	if err != nil {
//	    log.Fatal(err)
//	}
func (c *Client) Del(key any) error {
	if c.pool == nil {
		return errors.New("please call Initialize first")
	}

	_, err := c.do("DEL", key)

	return err
}

// FlushDB deletes all keys in the currently selected database.
//
// Returns:
//   - error: Error if client not initialized or FLUSHDB command fails, nil on success
//
// This operation removes all keys from the current database only. Keys in other databases
// are not affected. Use with caution as this operation cannot be undone.
//
// Warning: This is a destructive operation that will delete all data in the current database.
//
// Example:
//
//	client.Select(1) // Switch to database 1
//	err := client.FlushDB() // Only database 1 is cleared
//	if err != nil {
//	    log.Fatal(err)
//	}
func (c *Client) FlushDB() error {
	if c.pool == nil {
		return errors.New("please call Initialize first")
	}

	_, err := c.do("FLUSHDB")

	return err
}

// FlushAll deletes all keys in all databases.
//
// Returns:
//   - error: Error if client not initialized or FLUSHALL command fails, nil on success
//
// This operation removes all keys from all databases on the Redis server. Use with extreme
// caution as this operation cannot be undone and affects all databases.
//
// Warning: This is a destructive operation that will delete ALL data across ALL databases.
//
// Example:
//
//	err := client.FlushAll() // Clears all databases
//	if err != nil {
//	    log.Fatal(err)
//	}
func (c *Client) FlushAll() error {
	if c.pool == nil {
		return errors.New("please call Initialize first")
	}

	_, err := c.do("FLUSHALL")

	return err
}

// Ttl returns the remaining time to live of a key in seconds.
//
// Parameters:
//   - key: The key to check (can be string, int, or any type convertible to string)
//
// Returns:
//   - int: Time to live in seconds, or special values:
//     -2 if the key does not exist
//     -1 if the key exists but has no associated expiration
//     Positive integer for the remaining seconds until expiration
//   - error: Error if client not initialized or TTL command fails, nil on success
//
// Example:
//
//	ttl, err := client.Ttl("session:abc123")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	switch ttl {
//	case -2:
//	    fmt.Println("Key does not exist")
//	case -1:
//	    fmt.Println("Key has no expiration")
//	default:
//	    fmt.Printf("Key expires in %d seconds\n", ttl)
//	}
func (c *Client) Ttl(key any) (int, error) {
	if c.pool == nil {
		return -2, errors.New("please call Initialize first")
	}

	return redigo_redis.Int(c.do("TTL", key))
}

// Info retrieves information and statistics about the Redis server.
//
// Parameters:
//   - info: Category name (case-insensitive) or "ALL" for all categories. Valid categories:
//     "Server", "Clients", "Memory", "Persistence", "Stats", "Replication",
//     "CPU", "Cluster", "Keyspace"
//
// Returns:
//   - string: Server information in text format with key-value pairs
//   - error: Error if client not initialized or INFO command fails, nil on success
//
// The returned string contains multiple lines with information about the selected category.
// Each line is in the format "field:value". Sections are separated by "# Section" lines.
//
// Example:
//
//	memoryInfo, err := client.Info("Memory")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(memoryInfo)
//	// Output: # Memory
//	// used_memory:1234567
//	// used_memory_human:1.18M
//	// ...
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

// DBsize returns the number of keys in the currently selected database.
//
// Returns:
//   - int: Number of keys in the current database, or -1 if client not initialized
//   - error: Error if client not initialized or DBSIZE command fails, nil on success
//
// This command provides a count of all keys in the current database, regardless of their type.
// The operation is very fast as Redis maintains a count internally.
//
// Example:
//
//	count, err := client.DBsize()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Database has %d keys\n", count)
func (c *Client) DBsize() (int, error) {
	if c.pool == nil {
		return -1, errors.New("please call Initialize first")
	}

	return redigo_redis.Int(c.do("DBSIZE"))
}

// Exists checks if one or more keys exist.
//
// Parameters:
//   - keys: Variable number of keys to check
//
// Returns:
//   - bool: true if at least one key exists, false if none exist
//   - error: Error if client not initialized or EXISTS command fails, nil on success
//
// When multiple keys are provided, this function returns true if at least one of them exists.
// The actual EXISTS command returns the count of existing keys, but this function converts
// it to a boolean.
//
// Example:
//
//	// Check single key
//	exists, err := client.Exists("user:1")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	if exists {
//	    fmt.Println("User exists")
//	}
//
//	// Check multiple keys
//	exists, err = client.Exists("key1", "key2", "key3")
func (c *Client) Exists(keys ...any) (bool, error) {
	if c.pool == nil {
		return false, errors.New("please call Initialize first")
	}

	return redigo_redis.Bool(c.do("EXISTS", keys...))
}

// Rename renames a key to a new name.
//
// Parameters:
//   - currentKey: The existing key to rename
//   - newKey: The new name for the key
//
// Returns:
//   - error: Error if client not initialized, key doesn't exist, or RENAME command fails
//
// If newKey already exists, it will be overwritten. If currentKey does not exist, an error
// is returned. The operation is atomic.
//
// Example:
//
//	err := client.Rename("old_session", "new_session")
//	if err != nil {
//	    log.Fatal(err)
//	}
func (c *Client) Rename(currentKey, newKey any) error {
	if c.pool == nil {
		return errors.New("please call Initialize first")
	}

	_, err := c.do("RENAME", currentKey, newKey)
	return err
}

// RandomKey returns a random key from the currently selected database.
//
// Returns:
//   - string: A random key from the database, or empty string if database is empty
//   - error: Error if client not initialized or RANDOMKEY command fails, nil on success
//
// This function selects a random key from the current database. If the database is empty,
// an empty string is returned. This is useful for sampling or debugging.
//
// Example:
//
//	key, err := client.RandomKey()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	if key != "" {
//	    fmt.Printf("Random key: %s\n", key)
//	} else {
//	    fmt.Println("Database is empty")
//	}
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
