package redis_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/common-library/go/database/redis"
	"github.com/common-library/go/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	redismodule "github.com/testcontainers/testcontainers-go/modules/redis"
)

var (
	redisContainer *redismodule.RedisContainer
	redisAddress   string
	client         *redis.Client
)

func TestMain(m *testing.M) {
	if err := setupRedisContainer(); err != nil {
		log.Fatalf("Failed to setup Redis container: %v", err)
	}

	code := m.Run()

	if err := teardownRedisContainer(); err != nil {
		log.Printf("Failed to teardown Redis container: %v", err)
	}

	os.Exit(code)
}

func setupRedisContainer() error {
	ctx := context.Background()

	container, err := redismodule.Run(ctx, testutil.RedisImage)
	if err != nil {
		return fmt.Errorf("failed to start Redis container: %w", err)
	}

	redisContainer = container

	host, err := container.Host(ctx)
	if err != nil {
		return fmt.Errorf("failed to get container host: %w", err)
	}

	port, err := container.MappedPort(ctx, "6379/tcp")
	if err != nil {
		return fmt.Errorf("failed to get container port: %w", err)
	}

	redisAddress = fmt.Sprintf("%s:%s", host, port.Port())

	client = &redis.Client{}
	if err := client.Initialize(redisAddress, "", 10, 60); err != nil {
		return fmt.Errorf("failed to initialize Redis client: %w", err)
	}

	return nil
}

func teardownRedisContainer() error {
	if client != nil {
		if err := client.Finalize(); err != nil {
			log.Printf("Failed to finalize Redis client: %v", err)
		}
	}

	if redisContainer != nil {
		ctx := context.Background()
		if err := redisContainer.Terminate(ctx); err != nil {
			return fmt.Errorf("failed to terminate Redis container: %w", err)
		}
	}

	return nil
}

func setupTest(t *testing.T) {
	t.Helper()
	require.NoError(t, client.FlushDB())
}

func TestClient_Ping(t *testing.T) {
	setupTest(t)

	err := client.Ping()
	assert.NoError(t, err)
}

func TestClient_SetAndGet(t *testing.T) {
	setupTest(t)

	key := "test_key"
	value := "test_value"

	err := client.Set(key, value)
	assert.NoError(t, err)

	result, err := client.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, value, result)
}

func TestClient_SetexAndTtl(t *testing.T) {
	setupTest(t)

	key := "test_key_ttl"
	value := "test_value_ttl"
	seconds := 1

	err := client.Setex(key, seconds, value)
	assert.NoError(t, err)

	result, err := client.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, value, result)

	ttl, err := client.Ttl(key)
	assert.NoError(t, err)
	assert.Greater(t, ttl, 0)
	assert.LessOrEqual(t, ttl, seconds)

	maxWait := 2 * time.Second
	checkInterval := 50 * time.Millisecond
	deadline := time.Now().Add(maxWait)

	expired := false
	for time.Now().Before(deadline) {
		_, err = client.Get(key)
		if err != nil {
			expired = true
			break
		}
		time.Sleep(checkInterval)
	}

	assert.True(t, expired, "Key should expire within timeout")
}

func TestClient_MSetAndMGet(t *testing.T) {
	setupTest(t)

	err := client.MSet("key1", "value1", "key2", "value2", "key3", "value3")
	assert.NoError(t, err)

	values, err := client.MGet("key1", "key2", "key3")
	assert.NoError(t, err)
	assert.Equal(t, []string{"value1", "value2", "value3"}, values)
}

func TestClient_Del(t *testing.T) {
	setupTest(t)

	key := "test_key_del"
	value := "test_value_del"

	err := client.Set(key, value)
	assert.NoError(t, err)

	exists, err := client.Exists(key)
	assert.NoError(t, err)
	assert.True(t, exists)

	err = client.Del(key)
	assert.NoError(t, err)

	exists, err = client.Exists(key)
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestClient_Select(t *testing.T) {
	setupTest(t)

	err := client.Select(0)
	assert.NoError(t, err)

	key := "test_key_select"
	value := "test_value_select"
	err = client.Set(key, value)
	assert.NoError(t, err)

	err = client.Select(1)
	assert.NoError(t, err)

	exists, err := client.Exists(key)
	assert.NoError(t, err)
	assert.False(t, exists)

	err = client.Select(0)
	assert.NoError(t, err)

	exists, err = client.Exists(key)
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestClient_Exists(t *testing.T) {
	setupTest(t)

	exists, err := client.Exists("non_existent_key")
	assert.NoError(t, err)
	assert.False(t, exists)

	key := "test_key_exists"
	value := "test_value_exists"
	err = client.Set(key, value)
	assert.NoError(t, err)

	exists, err = client.Exists(key)
	assert.NoError(t, err)
	assert.True(t, exists)

	err = client.Set("key1", "value1")
	assert.NoError(t, err)
	err = client.Set("key2", "value2")
	assert.NoError(t, err)

	exists, err = client.Exists("key1", "key2")
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestClient_Rename(t *testing.T) {
	setupTest(t)

	oldKey := "old_key"
	newKey := "new_key"
	value := "test_value_rename"

	err := client.Set(oldKey, value)
	assert.NoError(t, err)

	err = client.Rename(oldKey, newKey)
	assert.NoError(t, err)

	result, err := client.Get(newKey)
	assert.NoError(t, err)
	assert.Equal(t, value, result)

	exists, err := client.Exists(oldKey)
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestClient_RandomKey(t *testing.T) {
	setupTest(t)

	_, err := client.RandomKey()
	assert.Error(t, err)

	testKeys := []string{"key1", "key2", "key3"}
	for _, key := range testKeys {
		err := client.Set(key, "value")
		assert.NoError(t, err)
	}

	randomKey, err := client.RandomKey()
	assert.NoError(t, err)
	assert.Contains(t, testKeys, randomKey)
}

func TestClient_Info(t *testing.T) {
	setupTest(t)

	info, err := client.Info("ALL")
	assert.NoError(t, err)
	assert.NotEmpty(t, info)
	assert.Contains(t, info, "redis_version")

	serverInfo, err := client.Info("server")
	assert.NoError(t, err)
	assert.NotEmpty(t, serverInfo)
	assert.Contains(t, serverInfo, "redis_version")
}

func TestClient_DBsize(t *testing.T) {
	setupTest(t)

	size, err := client.DBsize()
	assert.NoError(t, err)
	assert.Equal(t, 0, size)

	testKeys := []string{"key1", "key2", "key3"}
	for _, key := range testKeys {
		err := client.Set(key, "value")
		assert.NoError(t, err)
	}

	size, err = client.DBsize()
	assert.NoError(t, err)
	assert.Equal(t, len(testKeys), size)
}

func TestClient_FlushDB(t *testing.T) {
	setupTest(t)

	for i := 0; i < 5; i++ {
		err := client.Set("key"+strconv.Itoa(i), "value"+strconv.Itoa(i))
		assert.NoError(t, err)
	}

	size, err := client.DBsize()
	assert.NoError(t, err)
	assert.Equal(t, 5, size)

	err = client.FlushDB()
	assert.NoError(t, err)

	size, err = client.DBsize()
	assert.NoError(t, err)
	assert.Equal(t, 0, size)
}

func TestClient_FlushAll(t *testing.T) {
	for db := 0; db < 3; db++ {
		err := client.Select(db)
		assert.NoError(t, err)

		err = client.Set("key", fmt.Sprintf("value_db_%d", db))
		assert.NoError(t, err)
	}

	err := client.Select(0)
	assert.NoError(t, err)

	err = client.FlushAll()
	assert.NoError(t, err)

	for db := 0; db < 3; db++ {
		err := client.Select(db)
		assert.NoError(t, err)

		size, err := client.DBsize()
		assert.NoError(t, err)
		assert.Equal(t, 0, size)
	}
}
