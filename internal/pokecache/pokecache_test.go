package pokecache

import (
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	cache := NewCache(5 * time.Second)

	key := "https://example.com"
	value := []byte("testdata")

	cache.Add(key, value)
	data, found := cache.Get(key)
	if !found {
		t.Fatalf("Expected to find cache entry")
	}

	if string(data) != "testdata" {
		t.Fatalf("Expected to retrieve correct data")
	}
}

func TestReapLoop(t *testing.T) {
	cache := NewCache(5 * time.Millisecond)

	key := "https://example.com"
	value := []byte("testdata")
	cache.Add(key, value)

	// Wait for the cache entry to expire
	time.Sleep(10 * time.Millisecond)

	_, found := cache.Get(key)
	if found {
		t.Fatalf("Expected cache entry to be reaped")
	}
}
