package redis_test

import (
	"testing"
	"time"

	"github.com/heaven-chp/common-library-go/database/redis"
)

func TestInitialize(t *testing.T) {
	var redis redis.Redis

	err := redis.Initialize("", "127.0.0.1:6378", 3, 240)
	if err.Error() != "dial tcp 127.0.0.1:6378: connect: connection refused" {
		t.Error(err)
	}

	err = redis.Initialize("", "127.0.0.1:6379", 3, 240)
	if err != nil {
		t.Error(err)
	}

	err = redis.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestFinalize(t *testing.T) {
	var redis redis.Redis

	err := redis.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestPing(t *testing.T) {
	var redis redis.Redis

	err := redis.Ping()
	if err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	err = redis.Initialize("", "127.0.0.1:6379", 3, 240)
	if err != nil {
		t.Error(err)
	}

	err = redis.Ping()
	if err != nil {
		t.Error(err)
	}

	err = redis.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestSelect(t *testing.T) {
	var redis redis.Redis

	err := redis.Select(0)
	if err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	err = redis.Initialize("", "127.0.0.1:6379", 3, 240)
	if err != nil {
		t.Error(err)
	}

	err = redis.Select(0)
	if err != nil {
		t.Error(err)
	}

	err = redis.Select(1024)
	if err.Error() != "ERR DB index is out of range" {
		t.Error(err)
	}

	err = redis.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestGet(t *testing.T) {
	var redis redis.Redis

	_, err := redis.Get("key")
	if err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	err = redis.Initialize("", "127.0.0.1:6379", 3, 240)
	if err != nil {
		t.Error(err)
	}

	err = redis.FlushDB()
	if err != nil {
		t.Error(err)
	}

	_, err = redis.Get("key")
	if err.Error() != "redigo: nil returned" {
		t.Error(err)
	}

	err = redis.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestSet(t *testing.T) {
	var redis redis.Redis

	err := redis.Set("key", "value")
	if err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	err = redis.Initialize("", "127.0.0.1:6379", 3, 240)
	if err != nil {
		t.Error(err)
	}

	err = redis.Set("key", "value")
	if err != nil {
		t.Error(err)
	}

	data, err := redis.Get("key")
	if err != nil {
		t.Error(err)
	}
	if data != "value" {
		t.Errorf("invalid data - data : (%s)", data)
	}

	err = redis.Del("key")
	if err != nil {
		t.Error(err)
	}

	err = redis.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestSetex(t *testing.T) {
	var redis redis.Redis

	err := redis.Setex("key", 2, "value")
	if err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	err = redis.Initialize("", "127.0.0.1:6379", 3, 240)
	if err != nil {
		t.Error(err)
	}

	existsKey, err := redis.Exists("key")
	if err != nil {
		t.Error(err)
	}
	if existsKey != 0 {
		t.Errorf("invalid data - existsKey : (%d)", existsKey)
	}

	err = redis.Setex("key", 2, "value")
	if err != nil {
		t.Error(err)
	}

	existsKey, err = redis.Exists("key")
	if err != nil {
		t.Error(err)
	}
	if existsKey != 1 {
		t.Errorf("invalid data - existsKey : (%d)", existsKey)
	}

	time.Sleep(3 * time.Second)

	existsKey, err = redis.Exists("key")
	if err != nil {
		t.Error(err)
	}
	if existsKey != 0 {
		t.Errorf("invalid data - existsKey : (%d)", existsKey)
	}

	err = redis.Finalize()
	if err != nil {
		t.Error(err)
	}
}
func TestMGet(t *testing.T) {
	var redis redis.Redis

	_, err := redis.MGet("key1", "key2")
	if err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	err = redis.Initialize("", "127.0.0.1:6379", 3, 240)
	if err != nil {
		t.Error(err)
	}

	_, err = redis.MGet("key1", "key2")
	if err != nil {
		t.Error(err)
	}
	err = redis.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestMSet(t *testing.T) {
	var redis redis.Redis

	err := redis.MSet("key1", "value1", "key2", "value2")
	if err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	err = redis.Initialize("", "127.0.0.1:6379", 3, 240)
	if err != nil {
		t.Error(err)
	}

	err = redis.MSet("key1", "value1", "key2", "value2")
	if err != nil {
		t.Error(err)
	}

	data, err := redis.MGet("key1", "key2")
	if err != nil {
		t.Error(err)
	}
	if data[0] != "value1" {
		t.Errorf("invalid data - data : (%s)", data[0])
	}
	if data[1] != "value2" {
		t.Errorf("invalid data - data : (%s)", data[0])
	}

	err = redis.Del("key1")
	if err != nil {
		t.Error(err)
	}

	err = redis.Del("key2")
	if err != nil {
		t.Error(err)
	}

	err = redis.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestDel(t *testing.T) {
	var redis redis.Redis

	err := redis.Del("key")
	if err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	err = redis.Initialize("", "127.0.0.1:6379", 3, 240)
	if err != nil {
		t.Error(err)
	}

	err = redis.Del("key")
	if err != nil {
		t.Error(err)
	}

	err = redis.Del("key1")
	if err != nil {
		t.Error(err)
	}

	err = redis.Del("key2")
	if err != nil {
		t.Error(err)
	}

	err = redis.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestFlushDB(t *testing.T) {
	var redis redis.Redis

	err := redis.FlushDB()
	if err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	err = redis.Initialize("", "127.0.0.1:6379", 3, 240)
	if err != nil {
		t.Error(err)
	}

	err = redis.Set("key", "value")
	if err != nil {
		t.Error(err)
	}

	keyCount, err := redis.DBsize()
	if err != nil {
		t.Error(err)
	}
	if keyCount == 0 {
		t.Errorf("invalid data - keyCount : (%d)", keyCount)
	}

	err = redis.FlushDB()
	if err != nil {
		t.Error(err)
	}

	keyCount, err = redis.DBsize()
	if err != nil {
		t.Error(err)
	}
	if keyCount != 0 {
		t.Errorf("invalid data - keyCount : (%d)", keyCount)
	}

	err = redis.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestFlushAll(t *testing.T) {
	var redis redis.Redis

	err := redis.FlushAll()
	if err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	err = redis.Initialize("", "127.0.0.1:6379", 3, 240)
	if err != nil {
		t.Error(err)
	}

	err = redis.Set("key", "value")
	if err != nil {
		t.Error(err)
	}

	result, err := redis.Info("Keyspace")
	if err != nil {
		t.Error(err)
	}
	if len(result) == 0 {
		t.Error(err)
	}

	err = redis.FlushAll()
	if err != nil {
		t.Error(err)
	}

	result, err = redis.Info("Keyspace")
	if err != nil {
		t.Error(err)
	}
	if result != "# Keyspace\r\n" {
		t.Errorf("invalid data : %s", result)
	}

	err = redis.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestTtl(t *testing.T) {
	var redis redis.Redis

	ttl, err := redis.Ttl("key")
	if err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	err = redis.Initialize("", "127.0.0.1:6379", 3, 240)
	if err != nil {
		t.Error(err)
	}

	err = redis.Set("key", "value")
	if err != nil {
		t.Error(err)
	}

	ttl, err = redis.Ttl("key")
	if err != nil {
		t.Error(err)
	}
	if ttl != -1 {
		t.Errorf("invalid data - ttl : (%d)", ttl)
	}

	err = redis.Del("key")
	if err != nil {
		t.Error(err)
	}

	ttl, err = redis.Ttl("key")
	if err != nil {
		t.Error(err)
	}
	if ttl != -2 {
		t.Errorf("invalid data - ttl : (%d)", ttl)
	}

	err = redis.Setex("keyex", 2, "value")
	if err != nil {
		t.Error(err)
	}

	ttl, err = redis.Ttl("keyex")
	if err != nil {
		t.Error(err)
	}
	if ttl == -1 || ttl == -2 {
		t.Errorf("invalid data - ttl : (%d)", ttl)
	}

	err = redis.Del("keyex")
	if err != nil {
		t.Error(err)
	}

	err = redis.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestInfo(t *testing.T) {
	var redis redis.Redis

	_, err := redis.Info("ALL")
	if err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	err = redis.Initialize("", "127.0.0.1:6379", 3, 240)
	if err != nil {
		t.Error(err)
	}

	result, err := redis.Info("ALL")
	if err != nil {
		t.Error(err)
	}

	if len(result) == 0 {
		t.Errorf("invalid data - result : (%s)", result)
	}

	err = redis.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestDBsize(t *testing.T) {
	var redis redis.Redis

	keyCount, err := redis.DBsize()
	if keyCount != -1 || err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	err = redis.Initialize("", "127.0.0.1:6379", 3, 240)
	if err != nil {
		t.Error(err)
	}

	keyCount1, err := redis.DBsize()
	if err != nil {
		t.Error(err)
	}

	err = redis.Set("key", "value")
	if err != nil {
		t.Error(err)
	}

	keyCount2, err := redis.DBsize()
	if err != nil {
		t.Error(err)
	}
	if keyCount2 != keyCount1+1 {
		t.Errorf("invalid data - keyCount1 : (%d), keyCount2 : (%d)", keyCount1, keyCount2)
	}

	err = redis.Del("key")
	if err != nil {
		t.Error(err)
	}

	err = redis.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestExists(t *testing.T) {
	var redis redis.Redis

	existsKey, err := redis.Exists("key")
	if err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	err = redis.Initialize("", "127.0.0.1:6379", 3, 240)
	if err != nil {
		t.Error(err)
	}

	existsKey, err = redis.Exists("key")
	if err != nil {
		t.Error(err)
	}
	if existsKey != 0 {
		t.Errorf("invalid data - existsKey : (%d)", existsKey)
	}

	err = redis.Set("key", "value")
	if err != nil {
		t.Error(err)
	}

	existsKey, err = redis.Exists("key")
	if err != nil {
		t.Error(err)
	}
	if existsKey != 1 {
		t.Errorf("invalid data - existsKey : (%d)", existsKey)
	}

	existsKey, err = redis.Exists("key", 1, 2, "3")
	if err != nil {
		t.Error(err)
	}
	if existsKey != 1 {
		t.Errorf("invalid data - existsKey : (%d)", existsKey)
	}

	err = redis.Del("key")
	if err != nil {
		t.Error(err)
	}

	existsKey, err = redis.Exists("key")
	if err != nil {
		t.Error(err)
	}
	if existsKey != 0 {
		t.Errorf("invalid data - existsKey : (%d)", existsKey)
	}

	err = redis.Finalize()
	if err != nil {
		t.Error(err)
	}
}
func TestRename(t *testing.T) {
	var redis redis.Redis

	err := redis.Rename("key", "key_rename")
	if err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	err = redis.Initialize("", "127.0.0.1:6379", 3, 240)
	if err != nil {
		t.Error(err)
	}

	err = redis.Set("key", "value")
	if err != nil {
		t.Error(err)
	}

	err = redis.Rename("key", "key_rename")
	if err != nil {
		t.Error(err)
	}

	data, err := redis.Get("key_rename")
	if err != nil {
		t.Error(err)
	}
	if data != "value" {
		t.Errorf("invalid data - data : (%s)", data)
	}

	err = redis.Del("key_rename")
	if err != nil {
		t.Error(err)
	}

	err = redis.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestRandomKey(t *testing.T) {
	var redis redis.Redis

	key, err := redis.RandomKey()
	if err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	err = redis.Initialize("", "127.0.0.1:6379", 3, 240)
	if err != nil {
		t.Error(err)
	}

	err = redis.MSet("key1", "value1", "key2", "value2")
	if err != nil {
		t.Error(err)
	}

	key, err = redis.RandomKey()
	if err != nil {
		t.Error(err)
	}

	if key != "key1" && key != "key2" {
		t.Errorf("invalid data - key : (%s)", key)
	}

	err = redis.FlushDB()
	if err != nil {
		t.Error(err)
	}

	err = redis.Finalize()
	if err != nil {
		t.Error(err)
	}
}
