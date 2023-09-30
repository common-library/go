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

	if err := client.Initialize("", "127.0.0.1:6379", 3, 240); err != nil {
		t.Error(err)
	}

	if err := client.Finalize(); err != nil {
		t.Error(err)
	}
}

func TestFinalize(t *testing.T) {
	client := redis.Client{}

	if err := client.Finalize(); err != nil {
		t.Error(err)
	}
}

func TestPing(t *testing.T) {
	client := redis.Client{}

	if err := client.Ping(); err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	if err := client.Initialize("", "127.0.0.1:6379", 3, 240); err != nil {
		t.Error(err)
	}

	if err := client.Ping(); err != nil {
		t.Error(err)
	}

	if err := client.Finalize(); err != nil {
		t.Error(err)
	}
}

func TestSelect(t *testing.T) {
	client := redis.Client{}

	if err := client.Select(0); err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	if err := client.Initialize("", "127.0.0.1:6379", 3, 240); err != nil {
		t.Error(err)
	}

	if err := client.Select(0); err != nil {
		t.Error(err)
	}

	if err := client.Select(1024); err.Error() != "ERR DB index is out of range" {
		t.Error(err)
	}

	if err := client.Finalize(); err != nil {
		t.Error(err)
	}
}

func TestGet(t *testing.T) {
	client := redis.Client{}

	if _, err := client.Get("key"); err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	if err := client.Initialize("", "127.0.0.1:6379", 3, 240); err != nil {
		t.Error(err)
	}

	if err := client.FlushDB(); err != nil {
		t.Error(err)
	}

	if _, err := client.Get("key"); err.Error() != "redigo: nil returned" {
		t.Error(err)
	}

	if err := client.Finalize(); err != nil {
		t.Error(err)
	}
}

func TestSet(t *testing.T) {
	client := redis.Client{}

	if err := client.Set("key", "value"); err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	if err := client.Initialize("", "127.0.0.1:6379", 3, 240); err != nil {
		t.Error(err)
	}

	if err := client.Set("key", "value"); err != nil {
		t.Error(err)
	}

	if data, err := client.Get("key"); err != nil {
		t.Error(err)
	} else if data != "value" {
		t.Errorf("invalid data - data : (%s)", data)
	}

	if err := client.Del("key"); err != nil {
		t.Error(err)
	}

	if err := client.Finalize(); err != nil {
		t.Error(err)
	}
}

func TestSetex(t *testing.T) {
	client := redis.Client{}

	if err := client.Setex("key", 2, "value"); err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	if err := client.Initialize("", "127.0.0.1:6379", 3, 240); err != nil {
		t.Error(err)
	}

	if existsKey, err := client.Exists("key"); err != nil {
		t.Error(err)
	} else if existsKey != 0 {
		t.Errorf("invalid data - existsKey : (%d)", existsKey)
	}

	if err := client.Setex("key", 2, "value"); err != nil {
		t.Error(err)
	}

	if existsKey, err := client.Exists("key"); err != nil {
		t.Error(err)
	} else if existsKey != 1 {
		t.Errorf("invalid data - existsKey : (%d)", existsKey)
	}

	time.Sleep(2 * time.Second)

	if existsKey, err := client.Exists("key"); err != nil {
		t.Error(err)
	} else if existsKey != 0 {
		t.Errorf("invalid data - existsKey : (%d)", existsKey)
	}

	if err := client.Finalize(); err != nil {
		t.Error(err)
	}
}
func TestMGet(t *testing.T) {
	client := redis.Client{}

	if _, err := client.MGet("key1", "key2"); err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	if err := client.Initialize("", "127.0.0.1:6379", 3, 240); err != nil {
		t.Error(err)
	}

	if _, err := client.MGet("key1", "key2"); err != nil {
		t.Error(err)
	}

	if err := client.Finalize(); err != nil {
		t.Error(err)
	}
}

func TestMSet(t *testing.T) {
	client := redis.Client{}

	if err := client.MSet("key1", "value1", "key2", "value2"); err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	if err := client.Initialize("", "127.0.0.1:6379", 3, 240); err != nil {
		t.Error(err)
	}

	if err := client.MSet("key1", "value1", "key2", "value2"); err != nil {
		t.Error(err)
	}

	if data, err := client.MGet("key1", "key2"); err != nil {
		t.Error(err)
	} else if data[0] != "value1" {
		t.Errorf("invalid data - data : (%s)", data[0])
	} else if data[1] != "value2" {
		t.Errorf("invalid data - data : (%s)", data[0])
	}

	if err := client.Del("key1"); err != nil {
		t.Error(err)
	}

	if err := client.Del("key2"); err != nil {
		t.Error(err)
	}

	if err := client.Finalize(); err != nil {
		t.Error(err)
	}
}

func TestDel(t *testing.T) {
	client := redis.Client{}

	if err := client.Del("key"); err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	if err := client.Initialize("", "127.0.0.1:6379", 3, 240); err != nil {
		t.Error(err)
	}

	if err := client.Del("key"); err != nil {
		t.Error(err)
	}

	if err := client.Del("key1"); err != nil {
		t.Error(err)
	}

	if err := client.Del("key2"); err != nil {
		t.Error(err)
	}

	if err := client.Finalize(); err != nil {
		t.Error(err)
	}
}

func TestFlushDB(t *testing.T) {
	client := redis.Client{}

	if err := client.FlushDB(); err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	if err := client.Initialize("", "127.0.0.1:6379", 3, 240); err != nil {
		t.Error(err)
	}

	if err := client.Set("key", "value"); err != nil {
		t.Error(err)
	}

	if keyCount, err := client.DBsize(); err != nil {
		t.Error(err)
	} else if keyCount == 0 {
		t.Errorf("invalid data - keyCount : (%d)", keyCount)
	}

	if err := client.FlushDB(); err != nil {
		t.Error(err)
	}

	if keyCount, err := client.DBsize(); err != nil {
		t.Error(err)
	} else if keyCount != 0 {
		t.Errorf("invalid data - keyCount : (%d)", keyCount)
	}

	if err := client.Finalize(); err != nil {
		t.Error(err)
	}
}

func TestFlushAll(t *testing.T) {
	client := redis.Client{}

	if err := client.FlushAll(); err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	if err := client.Initialize("", "127.0.0.1:6379", 3, 240); err != nil {
		t.Error(err)
	}

	if err := client.Set("key", "value"); err != nil {
		t.Error(err)
	}

	if result, err := client.Info("Keyspace"); err != nil {
		t.Error(err)
	} else if len(result) == 0 {
		t.Error(err)
	}

	if err := client.FlushAll(); err != nil {
		t.Error(err)
	}

	if result, err := client.Info("Keyspace"); err != nil {
		t.Error(err)
	} else if result != "# Keyspace\r\n" {
		t.Errorf("invalid data : %s", result)
	}

	if err := client.Finalize(); err != nil {
		t.Error(err)
	}
}

func TestTtl(t *testing.T) {
	client := redis.Client{}

	if _, err := client.Ttl("key"); err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	if err := client.Initialize("", "127.0.0.1:6379", 3, 240); err != nil {
		t.Error(err)
	}

	if err := client.Set("key", "value"); err != nil {
		t.Error(err)
	}

	if ttl, err := client.Ttl("key"); err != nil {
		t.Error(err)
	} else if ttl != -1 {
		t.Errorf("invalid data - ttl : (%d)", ttl)
	}

	if err := client.Del("key"); err != nil {
		t.Error(err)
	}

	if ttl, err := client.Ttl("key"); err != nil {
		t.Error(err)
	} else if ttl != -2 {
		t.Errorf("invalid data - ttl : (%d)", ttl)
	}

	if err := client.Setex("keyex", 2, "value"); err != nil {
		t.Error(err)
	}

	if ttl, err := client.Ttl("keyex"); err != nil {
		t.Error(err)
	} else if ttl == -1 || ttl == -2 {
		t.Errorf("invalid data - ttl : (%d)", ttl)
	}

	if err := client.Del("keyex"); err != nil {
		t.Error(err)
	}

	if err := client.Finalize(); err != nil {
		t.Error(err)
	}
}

func TestInfo(t *testing.T) {
	client := redis.Client{}

	if _, err := client.Info("ALL"); err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	if err := client.Initialize("", "127.0.0.1:6379", 3, 240); err != nil {
		t.Error(err)
	}

	if result, err := client.Info("ALL"); err != nil {
		t.Error(err)
	} else if len(result) == 0 {
		t.Errorf("invalid data - result : (%s)", result)
	}

	if err := client.Finalize(); err != nil {
		t.Error(err)
	}
}

func TestDBsize(t *testing.T) {
	client := redis.Client{}

	if keyCount, err := client.DBsize(); keyCount != -1 || err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	if err := client.Initialize("", "127.0.0.1:6379", 3, 240); err != nil {
		t.Error(err)
	}

	if keyCount1, err := client.DBsize(); err != nil {
		t.Error(err)
	} else if err := client.Set("key", "value"); err != nil {
		t.Error(err)
	} else if keyCount2, err := client.DBsize(); err != nil {
		t.Error(err)
	} else if keyCount2 != keyCount1+1 {
		t.Errorf("invalid data - keyCount1 : (%d), keyCount2 : (%d)", keyCount1, keyCount2)
	}

	if err := client.Del("key"); err != nil {
		t.Error(err)
	}

	if err := client.Finalize(); err != nil {
		t.Error(err)
	}
}

func TestExists(t *testing.T) {
	client := redis.Client{}

	if _, err := client.Exists("key"); err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	if err := client.Initialize("", "127.0.0.1:6379", 3, 240); err != nil {
		t.Error(err)
	}

	if existsKey, err := client.Exists("key"); err != nil {
		t.Error(err)
	} else if existsKey != 0 {
		t.Errorf("invalid data - existsKey : (%d)", existsKey)
	}

	if err := client.Set("key", "value"); err != nil {
		t.Error(err)
	}

	if existsKey, err := client.Exists("key"); err != nil {
		t.Error(err)
	} else if existsKey != 1 {
		t.Errorf("invalid data - existsKey : (%d)", existsKey)
	}

	if existsKey, err := client.Exists("key", 1, 2, "3"); err != nil {
		t.Error(err)
	} else if existsKey != 1 {
		t.Errorf("invalid data - existsKey : (%d)", existsKey)
	}

	if err := client.Del("key"); err != nil {
		t.Error(err)
	}

	if existsKey, err := client.Exists("key"); err != nil {
		t.Error(err)
	} else if existsKey != 0 {
		t.Errorf("invalid data - existsKey : (%d)", existsKey)
	}

	if err := client.Finalize(); err != nil {
		t.Error(err)
	}
}

func TestRename(t *testing.T) {
	client := redis.Client{}

	if err := client.Rename("key", "key_rename"); err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	if err := client.Initialize("", "127.0.0.1:6379", 3, 240); err != nil {
		t.Error(err)
	}

	if err := client.Set("key", "value"); err != nil {
		t.Error(err)
	}

	if err := client.Rename("key", "key_rename"); err != nil {
		t.Error(err)
	}

	if data, err := client.Get("key_rename"); err != nil {
		t.Error(err)
	} else if data != "value" {
		t.Errorf("invalid data - data : (%s)", data)
	}

	if err := client.Del("key_rename"); err != nil {
		t.Error(err)
	}

	if err := client.Finalize(); err != nil {
		t.Error(err)
	}
}

func TestRandomKey(t *testing.T) {
	client := redis.Client{}

	if _, err := client.RandomKey(); err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	if err := client.Initialize("", "127.0.0.1:6379", 3, 240); err != nil {
		t.Error(err)
	}

	if err := client.MSet("key1", "value1", "key2", "value2"); err != nil {
		t.Error(err)
	}

	if key, err := client.RandomKey(); err != nil {
		t.Error(err)
	} else if key != "key1" && key != "key2" {
		t.Errorf("invalid data - key : (%s)", key)
	}

	if err := client.FlushDB(); err != nil {
		t.Error(err)
	}

	if err := client.Finalize(); err != nil {
		t.Error(err)
	}
}
