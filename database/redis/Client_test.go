package redis_test

import (
	"os"
	"testing"
	"time"

	"github.com/common-library/go/database/redis"
)

func getClient(t *testing.T) (*redis.Client, bool) {
	client := &redis.Client{}
	address := os.Getenv("REDIS_ADDRESS")

	if len(address) == 0 {
		return client, true
	} else if err := client.Initialize(address, "", 10, 60); err != nil {
		t.Fatal(err)
	}

	return client, false
}

func finalize(t *testing.T, client *redis.Client) {
	if err := client.Finalize(); err != nil {
		t.Fatal(err)
	}
}

func TestInitialize(t *testing.T) {
	client, stop := getClient(t)
	if stop {
		return
	}
	defer finalize(t, client)
}

func TestFinalize(t *testing.T) {
	TestInitialize(t)
}

func TestPing(t *testing.T) {
	if err := (&redis.Client{}).Ping(); err.Error() != "please call Initialize first" {
		t.Fatal(err)
	}

	client, stop := getClient(t)
	if stop {
		return
	}
	defer finalize(t, client)

	if err := client.Ping(); err != nil {
		t.Fatal(err)
	}
}

func TestSelect(t *testing.T) {
	if err := (&redis.Client{}).Select(0); err.Error() != "please call Initialize first" {
		t.Fatal(err)
	}

	client, stop := getClient(t)
	if stop {
		return
	}
	defer finalize(t, client)

	if err := client.Select(0); err != nil {
		t.Fatal(err)
	}

	if err := client.Select(1024); err.Error() != "ERR DB index is out of range" {
		t.Fatal(err)
	}
}

func TestGet(t *testing.T) {
	if _, err := (&redis.Client{}).Get("key"); err.Error() != "please call Initialize first" {
		t.Fatal(err)
	}

	client, stop := getClient(t)
	if stop {
		return
	}
	defer finalize(t, client)

	if err := client.FlushDB(); err != nil {
		t.Fatal(err)
	}

	if _, err := client.Get("key"); err.Error() != "redigo: nil returned" {
		t.Fatal(err)
	}
}

func TestSet(t *testing.T) {
	if err := (&redis.Client{}).Set("key", "value"); err.Error() != "please call Initialize first" {
		t.Fatal(err)
	}

	client, stop := getClient(t)
	if stop {
		return
	}
	defer finalize(t, client)

	if err := client.Set("key", "value"); err != nil {
		t.Fatal(err)
	}

	if data, err := client.Get("key"); err != nil {
		t.Fatal(err)
	} else if data != "value" {
		t.Fatal(data)
	}

	if err := client.Del("key"); err != nil {
		t.Fatal(err)
	}
}

func TestSetex(t *testing.T) {
	if err := (&redis.Client{}).Setex("key", 2, "value"); err.Error() != "please call Initialize first" {
		t.Fatal(err)
	}

	client, stop := getClient(t)
	if stop {
		return
	}
	defer finalize(t, client)

	if existsKey, err := client.Exists("key"); err != nil {
		t.Fatal(err)
	} else if existsKey != false {
		t.Fatal(existsKey)
	}

	const second = 1

	if err := client.Setex("key", second, "value"); err != nil {
		t.Fatal(err)
	}

	if existsKey, err := client.Exists("key"); err != nil {
		t.Fatal(err)
	} else if existsKey != true {
		t.Fatal(existsKey)
	}

	time.Sleep(1100 * time.Millisecond)

	if existsKey, err := client.Exists("key"); err != nil {
		t.Fatal(err)
	} else if existsKey != false {
		t.Fatal(existsKey)
	}
}

func TestMGet(t *testing.T) {
	if _, err := (&redis.Client{}).MGet("key1", "key2"); err.Error() != "please call Initialize first" {
		t.Fatal(err)
	}

	client, stop := getClient(t)
	if stop {
		return
	}
	defer finalize(t, client)

	if _, err := client.MGet("key1", "key2"); err != nil {
		t.Fatal(err)
	}
}

func TestMSet(t *testing.T) {
	if err := (&redis.Client{}).MSet("key1", "value1", "key2", "value2"); err.Error() != "please call Initialize first" {
		t.Fatal(err)
	}

	client, stop := getClient(t)
	if stop {
		return
	}
	defer finalize(t, client)

	if err := client.MSet("key1", "value1", "key2", "value2"); err != nil {
		t.Fatal(err)
	}

	if data, err := client.MGet("key1", "key2"); err != nil {
		t.Fatal(err)
	} else if data[0] != "value1" {
		t.Fatal(data[0])
	} else if data[1] != "value2" {
		t.Fatal(data[0])
	}

	if err := client.Del("key1"); err != nil {
		t.Fatal(err)
	}

	if err := client.Del("key2"); err != nil {
		t.Fatal(err)
	}
}

func TestDel(t *testing.T) {
	if err := (&redis.Client{}).Del("key"); err.Error() != "please call Initialize first" {
		t.Fatal(err)
	}

	client, stop := getClient(t)
	if stop {
		return
	}
	defer finalize(t, client)

	if err := client.Del("key"); err != nil {
		t.Fatal(err)
	}

	if err := client.Del("key1"); err != nil {
		t.Fatal(err)
	}

	if err := client.Del("key2"); err != nil {
		t.Fatal(err)
	}
}

func TestFlushDB(t *testing.T) {
	if err := (&redis.Client{}).FlushDB(); err.Error() != "please call Initialize first" {
		t.Fatal(err)
	}

	client, stop := getClient(t)
	if stop {
		return
	}
	defer finalize(t, client)

	if err := client.Set("key", "value"); err != nil {
		t.Fatal(err)
	}

	if keyCount, err := client.DBsize(); err != nil {
		t.Fatal(err)
	} else if keyCount == 0 {
		t.Fatal(keyCount)
	}

	if err := client.FlushDB(); err != nil {
		t.Fatal(err)
	}

	if keyCount, err := client.DBsize(); err != nil {
		t.Fatal(err)
	} else if keyCount != 0 {
		t.Fatal(keyCount)
	}
}

func TestFlushAll(t *testing.T) {
	if err := (&redis.Client{}).FlushAll(); err.Error() != "please call Initialize first" {
		t.Fatal(err)
	}

	client, stop := getClient(t)
	if stop {
		return
	}
	defer finalize(t, client)

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
		t.Fatal(result)
	}
}

func TestTtl(t *testing.T) {
	if _, err := (&redis.Client{}).Ttl("key"); err.Error() != "please call Initialize first" {
		t.Fatal(err)
	}

	client, stop := getClient(t)
	if stop {
		return
	}
	defer finalize(t, client)

	if err := client.Set("key", "value"); err != nil {
		t.Fatal(err)
	}

	if ttl, err := client.Ttl("key"); err != nil {
		t.Fatal(err)
	} else if ttl != -1 {
		t.Fatal(ttl)
	}

	if err := client.Del("key"); err != nil {
		t.Fatal(err)
	}

	if ttl, err := client.Ttl("key"); err != nil {
		t.Fatal(err)
	} else if ttl != -2 {
		t.Fatal(ttl)
	}

	if err := client.Setex("keyex", 2, "value"); err != nil {
		t.Fatal(err)
	}

	if ttl, err := client.Ttl("keyex"); err != nil {
		t.Fatal(err)
	} else if ttl == -1 || ttl == -2 {
		t.Fatal(ttl)
	}

	if err := client.Del("keyex"); err != nil {
		t.Fatal(err)
	}
}

func TestInfo(t *testing.T) {
	if _, err := (&redis.Client{}).Info("ALL"); err.Error() != "please call Initialize first" {
		t.Fatal(err)
	}

	client, stop := getClient(t)
	if stop {
		return
	}
	defer finalize(t, client)

	if result, err := client.Info("ALL"); err != nil {
		t.Fatal(err)
	} else if len(result) == 0 {
		t.Fatal(result)
	}
}

func TestDBsize(t *testing.T) {
	if keyCount, err := (&redis.Client{}).DBsize(); keyCount != -1 || err.Error() != "please call Initialize first" {
		t.Fatal(err)
	}

	client, stop := getClient(t)
	if stop {
		return
	}
	defer finalize(t, client)

	if keyCount1, err := client.DBsize(); err != nil {
		t.Fatal(err)
	} else if err := client.Set("key", "value"); err != nil {
		t.Fatal(err)
	} else if keyCount2, err := client.DBsize(); err != nil {
		t.Fatal(err)
	} else if keyCount2 != keyCount1+1 {
		t.Fatal(keyCount1, ",", keyCount2)
	}

	if err := client.Del("key"); err != nil {
		t.Fatal(err)
	}
}

func TestExists(t *testing.T) {
	if _, err := (&redis.Client{}).Exists("key"); err.Error() != "please call Initialize first" {
		t.Fatal(err)
	}

	client, stop := getClient(t)
	if stop {
		return
	}
	defer finalize(t, client)

	if existsKey, err := client.Exists("key"); err != nil {
		t.Fatal(err)
	} else if existsKey != false {
		t.Fatal(existsKey)
	}

	if err := client.Set("key", "value"); err != nil {
		t.Fatal(err)
	}

	if existsKey, err := client.Exists("key"); err != nil {
		t.Fatal(err)
	} else if existsKey != true {
		t.Fatal(existsKey)
	}

	if existsKey, err := client.Exists("key", 1, 2, "3"); err != nil {
		t.Fatal(err)
	} else if existsKey != true {
		t.Fatal(existsKey)
	}

	if err := client.Del("key"); err != nil {
		t.Fatal(err)
	}

	if existsKey, err := client.Exists("key"); err != nil {
		t.Fatal(err)
	} else if existsKey != false {
		t.Fatal(existsKey)
	}
}

func TestRename(t *testing.T) {
	if err := (&redis.Client{}).Rename("key", "key_rename"); err.Error() != "please call Initialize first" {
		t.Fatal(err)
	}

	client, stop := getClient(t)
	if stop {
		return
	}
	defer finalize(t, client)

	if err := client.Set("key", "value"); err != nil {
		t.Fatal(err)
	}

	if err := client.Rename("key", "key_rename"); err != nil {
		t.Fatal(err)
	}

	if data, err := client.Get("key_rename"); err != nil {
		t.Fatal(err)
	} else if data != "value" {
		t.Fatal(data)
	}

	if err := client.Del("key_rename"); err != nil {
		t.Fatal(err)
	}
}

func TestRandomKey(t *testing.T) {
	if _, err := (&redis.Client{}).RandomKey(); err.Error() != "please call Initialize first" {
		t.Fatal(err)
	}

	client, stop := getClient(t)
	if stop {
		return
	}
	defer finalize(t, client)

	if err := client.MSet("key1", "value1", "key2", "value2"); err != nil {
		t.Fatal(err)
	}

	if key, err := client.RandomKey(); err != nil {
		t.Fatal(err)
	} else if key != "key1" && key != "key2" {
		t.Fatal(key)
	}

	if err := client.FlushDB(); err != nil {
		t.Fatal(err)
	}
}
