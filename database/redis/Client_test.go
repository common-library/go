package redis_test

import (
	"testing"
	"time"

	"github.com/heaven-chp/common-library-go/database/redis"
)

func TestInitialize(t *testing.T) {
	client := redis.Client{}

	err := client.Initialize("", "127.0.0.1:6378", 3, 240)
	if err.Error() != "dial tcp 127.0.0.1:6378: connect: connection refused" {
		t.Error(err)
	}

	err = client.Initialize("", "127.0.0.1:6379", 3, 240)
	if err != nil {
		t.Error(err)
	}

	err = client.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestFinalize(t *testing.T) {
	client := redis.Client{}

	err := client.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestPing(t *testing.T) {
	client := redis.Client{}

	err := client.Ping()
	if err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	err = client.Initialize("", "127.0.0.1:6379", 3, 240)
	if err != nil {
		t.Error(err)
	}

	err = client.Ping()
	if err != nil {
		t.Error(err)
	}

	err = client.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestSelect(t *testing.T) {
	client := redis.Client{}

	err := client.Select(0)
	if err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	err = client.Initialize("", "127.0.0.1:6379", 3, 240)
	if err != nil {
		t.Error(err)
	}

	err = client.Select(0)
	if err != nil {
		t.Error(err)
	}

	err = client.Select(1024)
	if err.Error() != "ERR DB index is out of range" {
		t.Error(err)
	}

	err = client.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestGet(t *testing.T) {
	client := redis.Client{}

	_, err := client.Get("key")
	if err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	err = client.Initialize("", "127.0.0.1:6379", 3, 240)
	if err != nil {
		t.Error(err)
	}

	err = client.FlushDB()
	if err != nil {
		t.Error(err)
	}

	_, err = client.Get("key")
	if err.Error() != "redigo: nil returned" {
		t.Error(err)
	}

	err = client.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestSet(t *testing.T) {
	client := redis.Client{}

	err := client.Set("key", "value")
	if err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	err = client.Initialize("", "127.0.0.1:6379", 3, 240)
	if err != nil {
		t.Error(err)
	}

	err = client.Set("key", "value")
	if err != nil {
		t.Error(err)
	}

	data, err := client.Get("key")
	if err != nil {
		t.Error(err)
	}
	if data != "value" {
		t.Errorf("invalid data - data : (%s)", data)
	}

	err = client.Del("key")
	if err != nil {
		t.Error(err)
	}

	err = client.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestSetex(t *testing.T) {
	client := redis.Client{}

	err := client.Setex("key", 2, "value")
	if err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	err = client.Initialize("", "127.0.0.1:6379", 3, 240)
	if err != nil {
		t.Error(err)
	}

	existsKey, err := client.Exists("key")
	if err != nil {
		t.Error(err)
	}
	if existsKey != 0 {
		t.Errorf("invalid data - existsKey : (%d)", existsKey)
	}

	err = client.Setex("key", 2, "value")
	if err != nil {
		t.Error(err)
	}

	existsKey, err = client.Exists("key")
	if err != nil {
		t.Error(err)
	}
	if existsKey != 1 {
		t.Errorf("invalid data - existsKey : (%d)", existsKey)
	}

	time.Sleep(3 * time.Second)

	existsKey, err = client.Exists("key")
	if err != nil {
		t.Error(err)
	}
	if existsKey != 0 {
		t.Errorf("invalid data - existsKey : (%d)", existsKey)
	}

	err = client.Finalize()
	if err != nil {
		t.Error(err)
	}
}
func TestMGet(t *testing.T) {
	client := redis.Client{}

	_, err := client.MGet("key1", "key2")
	if err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	err = client.Initialize("", "127.0.0.1:6379", 3, 240)
	if err != nil {
		t.Error(err)
	}

	_, err = client.MGet("key1", "key2")
	if err != nil {
		t.Error(err)
	}
	err = client.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestMSet(t *testing.T) {
	client := redis.Client{}

	err := client.MSet("key1", "value1", "key2", "value2")
	if err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	err = client.Initialize("", "127.0.0.1:6379", 3, 240)
	if err != nil {
		t.Error(err)
	}

	err = client.MSet("key1", "value1", "key2", "value2")
	if err != nil {
		t.Error(err)
	}

	data, err := client.MGet("key1", "key2")
	if err != nil {
		t.Error(err)
	}
	if data[0] != "value1" {
		t.Errorf("invalid data - data : (%s)", data[0])
	}
	if data[1] != "value2" {
		t.Errorf("invalid data - data : (%s)", data[0])
	}

	err = client.Del("key1")
	if err != nil {
		t.Error(err)
	}

	err = client.Del("key2")
	if err != nil {
		t.Error(err)
	}

	err = client.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestDel(t *testing.T) {
	client := redis.Client{}

	err := client.Del("key")
	if err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	err = client.Initialize("", "127.0.0.1:6379", 3, 240)
	if err != nil {
		t.Error(err)
	}

	err = client.Del("key")
	if err != nil {
		t.Error(err)
	}

	err = client.Del("key1")
	if err != nil {
		t.Error(err)
	}

	err = client.Del("key2")
	if err != nil {
		t.Error(err)
	}

	err = client.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestFlushDB(t *testing.T) {
	client := redis.Client{}

	err := client.FlushDB()
	if err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	err = client.Initialize("", "127.0.0.1:6379", 3, 240)
	if err != nil {
		t.Error(err)
	}

	err = client.Set("key", "value")
	if err != nil {
		t.Error(err)
	}

	keyCount, err := client.DBsize()
	if err != nil {
		t.Error(err)
	}
	if keyCount == 0 {
		t.Errorf("invalid data - keyCount : (%d)", keyCount)
	}

	err = client.FlushDB()
	if err != nil {
		t.Error(err)
	}

	keyCount, err = client.DBsize()
	if err != nil {
		t.Error(err)
	}
	if keyCount != 0 {
		t.Errorf("invalid data - keyCount : (%d)", keyCount)
	}

	err = client.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestFlushAll(t *testing.T) {
	client := redis.Client{}

	err := client.FlushAll()
	if err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	err = client.Initialize("", "127.0.0.1:6379", 3, 240)
	if err != nil {
		t.Error(err)
	}

	err = client.Set("key", "value")
	if err != nil {
		t.Error(err)
	}

	result, err := client.Info("Keyspace")
	if err != nil {
		t.Error(err)
	}
	if len(result) == 0 {
		t.Error(err)
	}

	err = client.FlushAll()
	if err != nil {
		t.Error(err)
	}

	result, err = client.Info("Keyspace")
	if err != nil {
		t.Error(err)
	}
	if result != "# Keyspace\r\n" {
		t.Errorf("invalid data : %s", result)
	}

	err = client.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestTtl(t *testing.T) {
	client := redis.Client{}

	ttl, err := client.Ttl("key")
	if err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	err = client.Initialize("", "127.0.0.1:6379", 3, 240)
	if err != nil {
		t.Error(err)
	}

	err = client.Set("key", "value")
	if err != nil {
		t.Error(err)
	}

	ttl, err = client.Ttl("key")
	if err != nil {
		t.Error(err)
	}
	if ttl != -1 {
		t.Errorf("invalid data - ttl : (%d)", ttl)
	}

	err = client.Del("key")
	if err != nil {
		t.Error(err)
	}

	ttl, err = client.Ttl("key")
	if err != nil {
		t.Error(err)
	}
	if ttl != -2 {
		t.Errorf("invalid data - ttl : (%d)", ttl)
	}

	err = client.Setex("keyex", 2, "value")
	if err != nil {
		t.Error(err)
	}

	ttl, err = client.Ttl("keyex")
	if err != nil {
		t.Error(err)
	}
	if ttl == -1 || ttl == -2 {
		t.Errorf("invalid data - ttl : (%d)", ttl)
	}

	err = client.Del("keyex")
	if err != nil {
		t.Error(err)
	}

	err = client.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestInfo(t *testing.T) {
	client := redis.Client{}

	_, err := client.Info("ALL")
	if err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	err = client.Initialize("", "127.0.0.1:6379", 3, 240)
	if err != nil {
		t.Error(err)
	}

	result, err := client.Info("ALL")
	if err != nil {
		t.Error(err)
	}

	if len(result) == 0 {
		t.Errorf("invalid data - result : (%s)", result)
	}

	err = client.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestDBsize(t *testing.T) {
	client := redis.Client{}

	keyCount, err := client.DBsize()
	if keyCount != -1 || err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	err = client.Initialize("", "127.0.0.1:6379", 3, 240)
	if err != nil {
		t.Error(err)
	}

	keyCount1, err := client.DBsize()
	if err != nil {
		t.Error(err)
	}

	err = client.Set("key", "value")
	if err != nil {
		t.Error(err)
	}

	keyCount2, err := client.DBsize()
	if err != nil {
		t.Error(err)
	}
	if keyCount2 != keyCount1+1 {
		t.Errorf("invalid data - keyCount1 : (%d), keyCount2 : (%d)", keyCount1, keyCount2)
	}

	err = client.Del("key")
	if err != nil {
		t.Error(err)
	}

	err = client.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestExists(t *testing.T) {
	client := redis.Client{}

	existsKey, err := client.Exists("key")
	if err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	err = client.Initialize("", "127.0.0.1:6379", 3, 240)
	if err != nil {
		t.Error(err)
	}

	existsKey, err = client.Exists("key")
	if err != nil {
		t.Error(err)
	}
	if existsKey != 0 {
		t.Errorf("invalid data - existsKey : (%d)", existsKey)
	}

	err = client.Set("key", "value")
	if err != nil {
		t.Error(err)
	}

	existsKey, err = client.Exists("key")
	if err != nil {
		t.Error(err)
	}
	if existsKey != 1 {
		t.Errorf("invalid data - existsKey : (%d)", existsKey)
	}

	existsKey, err = client.Exists("key", 1, 2, "3")
	if err != nil {
		t.Error(err)
	}
	if existsKey != 1 {
		t.Errorf("invalid data - existsKey : (%d)", existsKey)
	}

	err = client.Del("key")
	if err != nil {
		t.Error(err)
	}

	existsKey, err = client.Exists("key")
	if err != nil {
		t.Error(err)
	}
	if existsKey != 0 {
		t.Errorf("invalid data - existsKey : (%d)", existsKey)
	}

	err = client.Finalize()
	if err != nil {
		t.Error(err)
	}
}
func TestRename(t *testing.T) {
	client := redis.Client{}

	err := client.Rename("key", "key_rename")
	if err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	err = client.Initialize("", "127.0.0.1:6379", 3, 240)
	if err != nil {
		t.Error(err)
	}

	err = client.Set("key", "value")
	if err != nil {
		t.Error(err)
	}

	err = client.Rename("key", "key_rename")
	if err != nil {
		t.Error(err)
	}

	data, err := client.Get("key_rename")
	if err != nil {
		t.Error(err)
	}
	if data != "value" {
		t.Errorf("invalid data - data : (%s)", data)
	}

	err = client.Del("key_rename")
	if err != nil {
		t.Error(err)
	}

	err = client.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestRandomKey(t *testing.T) {
	client := redis.Client{}

	key, err := client.RandomKey()
	if err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	err = client.Initialize("", "127.0.0.1:6379", 3, 240)
	if err != nil {
		t.Error(err)
	}

	err = client.MSet("key1", "value1", "key2", "value2")
	if err != nil {
		t.Error(err)
	}

	key, err = client.RandomKey()
	if err != nil {
		t.Error(err)
	}

	if key != "key1" && key != "key2" {
		t.Errorf("invalid data - key : (%s)", key)
	}

	err = client.FlushDB()
	if err != nil {
		t.Error(err)
	}

	err = client.Finalize()
	if err != nil {
		t.Error(err)
	}
}
