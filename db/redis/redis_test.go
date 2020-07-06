package redis

import (
	"testing"
)

func TestInitialize(t *testing.T) {
	var redis Redis

	err := redis.Initialize("", "127.0.0.1:6379", 3, 240)
	if err != nil {
		t.Error(err)
	}

	err = redis.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestFinalize(t *testing.T) {
	var redis Redis

	err := redis.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestPing(t *testing.T) {
	var redis Redis

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

func TestGet(t *testing.T) {
	var redis Redis

	_, err := redis.Get("key")
	if err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	err = redis.Initialize("", "127.0.0.1:6379", 3, 240)
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
	var redis Redis

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

	var data string
	data, err = redis.Get("key")
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

func TestMGet(t *testing.T) {
	var redis Redis

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
	var redis Redis

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
	var redis Redis

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

/*
func TestAccept1(t *testing.T) {
    var server Server
    defer server.Finalize()

    _, err := server.accept()
    if err.Error() != "Listen first before Accept" {
        t.Error(err)
    }
}

*/
