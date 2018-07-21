package cache

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	c := New()
	if c.data == nil {
		t.Errorf("creation the new cache failed")
	}
}

func TestKeys(t *testing.T) {
	c := New()
	c.Set("abc", "yes")

	keys := c.Keys()

	flag := false
	for _, v := range keys {
		if v == "abc" {
			flag = true
		}
	}
	if !flag {
		t.Errorf("key 'abc' not found in %s", c.Keys())
	}
}

func TestSetTTL(t *testing.T) {
	c := New()

	c.Set("test1", "yes")
	c.SetTTL("test1", 100*time.Millisecond)

	c.Set("test2", "yes")
	c.SetTTL("test2", 200*time.Millisecond)

	time.Sleep(150 * time.Millisecond)

	_, err := c.Get("test1")
	if err != ErrorKeyNotFound {
		t.Error("expected ErrorKeyNotFound, but key found")
	}

	value, err := c.Get("test2")
	if value != "yes" {
		t.Errorf("Expected 'yes', got %s", value)
	}
	if err != nil {
		t.Errorf("expected nil error, got %s", err)
	}
}

func TestRemove(t *testing.T) {
	c := New()

	c.Set("abc", "yes")

	_, err := c.Get("abc")
	if err != nil {
		t.Errorf("expected nil error, got %s", err)
	}

	c.Remove("abc")

	_, err = c.Get("abc")
	if err != ErrorKeyNotFound {
		t.Error("expected ErrorKeyNotFound, but key found")
	}
}

func TestSetGet(t *testing.T) {
	c := New()

	c.Set("abc", "yes")

	v, err := c.Get("abc")
	if err != nil {
		t.Errorf("expected nil error, got %s", err)
	}

	if v != "yes" {
		t.Errorf("Expected 'yes', got %s", v)
	}
}

func TestPushPop(t *testing.T) {
	c := New()
	c.Push("list", "a")

	v, err := c.Pop("list")
	if err != nil {
		t.Errorf("expected nil error, got %s", err)
	}
	if v != "a" {
		t.Errorf("Expected 'a', got %s", v)
	}

	_, err = c.Pop("list")
	if err != ErrorNoItems {
		t.Errorf("expected ErrorNoItems error, got '%s'", err)
	}
}

func TestHsetHget(t *testing.T) {
	c := New()
	c.Hset("dict", "v1", "a")

	v, err := c.Hget("dict", "v1")
	if err != nil {
		t.Errorf("expected nil error, got %s", err)
	}
	if v != "a" {
		t.Errorf("Expected 'a', got %s", v)
	}
}
