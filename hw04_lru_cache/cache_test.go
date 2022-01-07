package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(3)

		wasInCache := c.Set("a", 1)
		require.False(t, wasInCache)
		wasInCache = c.Set("b", 2)
		require.False(t, wasInCache)
		wasInCache = c.Set("c", 3) // {c:3, b:2, a:1}
		require.False(t, wasInCache)
		wasInCache = c.Set("d", 4) // {d:4, c:3, b:2}, a -> out
		require.False(t, wasInCache)

		val, ok := c.Get("a")
		require.False(t, ok)
		require.Nil(t, val)

		wasInCache = c.Set("c", 5) // {c:5, d:4, b:2} c -> first
		require.True(t, wasInCache)
		wasInCache = c.Set("b", 6) // {b:6, c:5, d:4} b -> first
		require.True(t, wasInCache)

		val, ok = c.Get("d") // {d:4, b:6, c:5} d -> first
		require.True(t, ok)
		require.Equal(t, 4, val)

		wasInCache = c.Set("e", 7) // {e:7, d:4, b:6} c -> out
		require.False(t, wasInCache)

		val, ok = c.Get("c")
		require.False(t, ok)
		require.Nil(t, val)

		c.Clear()
		val, ok = c.Get("e")
		require.False(t, ok)
		require.Nil(t, val)
	})
}

func TestCacheMultithreading(t *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
