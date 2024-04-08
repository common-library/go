package redis_test

import (
	"testing"
	"time"

	"github.com/common-library/go/database/redis"
)

func TestInitialize(t *testing.T) {
	client := redis.Client{}

	err := client.Initialize("127.0.0.1:6378", "", 10, 60)
	if err.Error() != "dial tcp 127.0.0.1:6378: connect: connection refused" {
		t.Fatal(err)
	}

	if err := client.Initialize("127.0.0.1:6379", "", 10, 60); err != nil {
		t.Fatal(err)
	}

	if err := client.Finalize(); err != nil {
		t.Fatal(err)
	}
}

func TestFinalize(t *testing.T) {
	client := redis.Client{}

	if err := client.Finalize(); err != nil {
		t.Fatal(err)
	}
}

func TestPing(t *testing.T) {
	client := redis.Client{}

	if err := client.Ping(); err.Error() != "please call Initialize first" {
		t.Fatal(err)
	}

	if err := client.Initialize("127.0.0.1:6379", "", 10, 60); err != nil {
		t.Fatal(err)
	}

	if err := client.Ping(); err != nil {
		t.Fatal(err)
	}

	if err := client.Finalize(); err != nil {
		t.Fatal(err)
	}
}

func TestSelect(t *testing.T) {
	client := redis.Client{}

	if err := client.Select(0); err.Error() != "please call Initialize first" {
		t.Fatal(err)
	}

	if err := client.Initialize("127.0.0.1:6379", "", 10, 60); err != nil {
		t.Fatal(err)
	}

	if err := client.Select(0); err != nil {
		t.Fatal(err)
	}

	if err := client.Select(1024); err.Error() != "ERR DB index is out of range" {
		t.Fatal(err)
	}

	if err := client.Finalize(); err != nil {
		t.Fatal(err)
	}
}

func TestGet(t *testing.T) {
	client := redis.Client{}

	if _, err := client.Get("key"); err.Error() != "please call Initialize first" {
		t.Fatal(err)
	}

	if err := client.Initialize("127.0.0.1:6379", "", 10, 60); err != nil {
		t.Fatal(err)
	}

	if err := client.FlushDB(); err != nil {
		t.Fatal(err)
	}

	if _, err := client.Get("key"); err.Error() != "redigo: nil returned" {
		t.Fatal(err)
	}

	if err := client.Finalize(); err != nil {
		t.Fatal(err)
	}
}

func TestSet(t *testing.T) {
	client := redis.Client{}

	if err := client.Set("key", "value"); err.Error() != "please call Initialize first" {
		t.Fatal(err)
	}

	if err := client.Initialize("127.0.0.1:6379", "", 10, 60); err != nil {
		t.Fatal(err)
	}

	if err := client.Set("key", "value"); err != nil {
		t.Fatal(err)
	}

	if data, err := client.Get("key"); err != nil {
		t.Fatal(err)
	} else if data != "value" {
		t.Fatal("invalid -", data)
	}

	if err := client.Del("key"); err != nil {
		t.Fatal(err)
	}

	if err := client.Finalize(); err != nil {
		t.Fatal(err)
	}
}

func TestSetex(t *testing.T) {
	client := redis.Client{}

	if err := client.Setex("key", 2, "value"); err.Error() != "please call Initialize first" {
		t.Fatal(err)
	}

	if err := client.Initialize("127.0.0.1:6379", "", 10, 60); err != nil {
		t.Fatal(err)
	}

	if existsKey, err := client.Exists("key"); err != nil {
		t.Fatal(err)
	} else if existsKey != false {
		t.Fatal("invalid")
	}

	if err := client.Setex("key", 2, "value"); err != nil {
		t.Fatal(err)
	}

	if existsKey, err := client.Exists("key"); err != nil {
		t.Fatal(err)
	} else if existsKey != true {
		t.Fatal("invalid")
	}

	time.Sleep(3 * time.Second)

	if existsKey, err := client.Exists("key"); err != nil {
		t.Fatal(err)
	} else if existsKey != false {
		t.Fatal("invalid")
	}

	if err := client.Finalize(); err != nil {
		t.Fatal(err)
	}
}
func TestMGet(t *testing.T) {
	client := redis.Client{}

	if _, err := client.MGet("key1", "key2"); err.Error() != "please call Initialize first" {
		t.Fatal(err)
	}

	if err := client.Initialize("127.0.0.1:6379", "", 10, 60); err != nil {
		t.Fatal(err)
	}

	if _, err := client.MGet("key1", "key2"); err != nil {
		t.Fatal(err)
	}

	if err := client.Finalize(); err != nil {
		t.Fatal(err)
	}
}

func TestMSet(t *testing.T) {
	client := redis.Client{}

	if err := client.MSet("key1", "value1", "key2", "value2"); err.Error() != "please call Initialize first" {
		t.Fatal(err)
	}

	if err := client.Initialize("127.0.0.1:6379", "", 10, 60); err != nil {
		t.Fatal(err)
	}

	if err := client.MSet("key1", "value1", "key2", "value2"); err != nil {
		t.Fatal(err)
	}

	if data, err := client.MGet("key1", "key2"); err != nil {
		t.Fatal(err)
	} else if data[0] != "value1" {
		t.Fatal("invalid -", data[0])
	} else if data[1] != "value2" {
		t.Fatal("invalid -", data[0])
	}

	if err := client.Del("key1"); err != nil {
		t.Fatal(err)
	}

	if err := client.Del("key2"); err != nil {
		t.Fatal(err)
	}

	if err := client.Finalize(); err != nil {
		t.Fatal(err)
	}
}

func TestDel(t *testing.T) {
	client := redis.Client{}

	if err := client.Del("key"); err.Error() != "please call Initialize first" {
		t.Fatal(err)
	}

	if err := client.Initialize("127.0.0.1:6379", "", 10, 60); err != nil {
		t.Fatal(err)
	}

	if err := client.Del("key"); err != nil {
		t.Fatal(err)
	}

	if err := client.Del("key1"); err != nil {
		t.Fatal(err)
	}

	if err := client.Del("key2"); err != nil {
		t.Fatal(err)
	}

	if err := client.Finalize(); err != nil {
		t.Fatal(err)
	}
}

func TestFlushDB(t *testing.T) {
	client := redis.Client{}

	if err := client.FlushDB(); err.Error() != "please call Initialize first" {
		t.Fatal(err)
	}

	if err := client.Initialize("127.0.0.1:6379", "", 10, 60); err != nil {
		t.Fatal(err)
	}

	if err := client.Set("key", "value"); err != nil {
		t.Fatal(err)
	}

	if keyCount, err := client.DBsize(); err != nil {
		t.Fatal(err)
	} else if keyCount == 0 {
		t.Fatal("invalid -", keyCount)
	}

	if err := client.FlushDB(); err != nil {
		t.Fatal(err)
	}

	if keyCount, err := client.DBsize(); err != nil {
		t.Fatal(err)
	} else if keyCount != 0 {
		t.Fatal("invalid -", keyCount)
	}

	if err := client.Finalize(); err != nil {
		t.Fatal(err)
	}
}

func TestFlushAll(t *testing.T) {
	client := redis.Client{}

	if err := client.FlushAll(); err.Error() != "please call Initialize first" {
		t.Fatal(err)
	}

	if err := client.Initialize("127.0.0.1:6379", "", 10, 60); err != nil {
		t.Fatal(err)
	}

	if err := client.Set("key", "value"); err != nil {
		t.Fatal(err)
	}

	if result, err := client.Info("Keyspace"); err != nil {
		t.Fatal(err)
	} else if len(result) == 0 {
		t.Fatal(err)
	}

	if err := client.FlushAll(); err != nil {
		t.Fatal(err)
	}

	if result, err := client.Info("Keyspace"); err != nil {
		t.Fatal(err)
	} else if result != "# Keyspace\r\n" {
		t.Fatal("invalid -", result)
	}

	if err := client.Finalize(); err != nil {
		t.Fatal(err)
	}
}

func TestTtl(t *testing.T) {
	client := redis.Client{}

	if _, err := client.Ttl("key"); err.Error() != "please call Initialize first" {
		t.Fatal(err)
	}

	if err := client.Initialize("127.0.0.1:6379", "", 10, 60); err != nil {
		t.Fatal(err)
	}

	if err := client.Set("key", "value"); err != nil {
		t.Fatal(err)
	}

	if ttl, err := client.Ttl("key"); err != nil {
		t.Fatal(err)
	} else if ttl != -1 {
		t.Fatal("invalid -", ttl)
	}

	if err := client.Del("key"); err != nil {
		t.Fatal(err)
	}

	if ttl, err := client.Ttl("key"); err != nil {
		t.Fatal(err)
	} else if ttl != -2 {
		t.Fatal("invalid -", ttl)
	}

	if err := client.Setex("keyex", 2, "value"); err != nil {
		t.Fatal(err)
	}

	if ttl, err := client.Ttl("keyex"); err != nil {
		t.Fatal(err)
	} else if ttl == -1 || ttl == -2 {
		t.Fatal("invalid -", ttl)
	}

	if err := client.Del("keyex"); err != nil {
		t.Fatal(err)
	}

	if err := client.Finalize(); err != nil {
		t.Fatal(err)
	}
}

func TestInfo(t *testing.T) {
	client := redis.Client{}

	if _, err := client.Info("ALL"); err.Error() != "please call Initialize first" {
		t.Fatal(err)
	}

	if err := client.Initialize("127.0.0.1:6379", "", 10, 60); err != nil {
		t.Fatal(err)
	}

	if result, err := client.Info("ALL"); err != nil {
		t.Fatal(err)
	} else if len(result) == 0 {
		t.Fatal("invalid -", result)
	}

	if err := client.Finalize(); err != nil {
		t.Fatal(err)
	}
}

func TestDBsize(t *testing.T) {
	client := redis.Client{}

	if keyCount, err := client.DBsize(); keyCount != -1 || err.Error() != "please call Initialize first" {
		t.Fatal(err)
	}

	if err := client.Initialize("127.0.0.1:6379", "", 10, 60); err != nil {
		t.Fatal(err)
	}

	if keyCount1, err := client.DBsize(); err != nil {
		t.Fatal(err)
	} else if err := client.Set("key", "value"); err != nil {
		t.Fatal(err)
	} else if keyCount2, err := client.DBsize(); err != nil {
		t.Fatal(err)
	} else if keyCount2 != keyCount1+1 {
		t.Fatal("invalid -", keyCount1, ",", keyCount2)
	}

	if err := client.Del("key"); err != nil {
		t.Fatal(err)
	}

	if err := client.Finalize(); err != nil {
		t.Fatal(err)
	}
}

func TestExists(t *testing.T) {
	client := redis.Client{}

	if _, err := client.Exists("key"); err.Error() != "please call Initialize first" {
		t.Fatal(err)
	}

	if err := client.Initialize("127.0.0.1:6379", "", 10, 60); err != nil {
		t.Fatal(err)
	}

	if existsKey, err := client.Exists("key"); err != nil {
		t.Fatal(err)
	} else if existsKey != false {
		t.Fatal("invalid")
	}

	if err := client.Set("key", "value"); err != nil {
		t.Fatal(err)
	}

	if existsKey, err := client.Exists("key"); err != nil {
		t.Fatal(err)
	} else if existsKey != true {
		t.Fatal("invalid")
	}

	if existsKey, err := client.Exists("key", 1, 2, "3"); err != nil {
		t.Fatal(err)
	} else if existsKey != true {
		t.Fatal("invalid")
	}

	if err := client.Del("key"); err != nil {
		t.Fatal(err)
	}

	if existsKey, err := client.Exists("key"); err != nil {
		t.Fatal(err)
	} else if existsKey != false {
		t.Fatal("invalid")
	}

	if err := client.Finalize(); err != nil {
		t.Fatal(err)
	}
}

func TestRename(t *testing.T) {
	client := redis.Client{}

	if err := client.Rename("key", "key_rename"); err.Error() != "please call Initialize first" {
		t.Fatal(err)
	}

	if err := client.Initialize("127.0.0.1:6379", "", 10, 60); err != nil {
		t.Fatal(err)
	}

	if err := client.Set("key", "value"); err != nil {
		t.Fatal(err)
	}

	if err := client.Rename("key", "key_rename"); err != nil {
		t.Fatal(err)
	}

	if data, err := client.Get("key_rename"); err != nil {
		t.Fatal(err)
	} else if data != "value" {
		t.Fatal("invalid -", data)
	}

	if err := client.Del("key_rename"); err != nil {
		t.Fatal(err)
	}

	if err := client.Finalize(); err != nil {
		t.Fatal(err)
	}
}

func TestRandomKey(t *testing.T) {
	client := redis.Client{}

	if _, err := client.RandomKey(); err.Error() != "please call Initialize first" {
		t.Fatal(err)
	}

	if err := client.Initialize("127.0.0.1:6379", "", 10, 60); err != nil {
		t.Fatal(err)
	}

	if err := client.MSet("key1", "value1", "key2", "value2"); err != nil {
		t.Fatal(err)
	}

	if key, err := client.RandomKey(); err != nil {
		t.Fatal(err)
	} else if key != "key1" && key != "key2" {
		t.Fatal("invalid -", key)
	}

	if err := client.FlushDB(); err != nil {
		t.Fatal(err)
	}

	if err := client.Finalize(); err != nil {
		t.Fatal(err)
	}
}
