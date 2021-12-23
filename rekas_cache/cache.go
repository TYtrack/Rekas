/*
 * @Author: your name
 * @Date: 2021-12-23 02:12:40
 * @LastEditTime: 2021-12-23 17:13:00
 * @LastEditors: TYtrack
 * @Description: 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 * @FilePath: /Rekas/rekas_cache/cache.go
 */

package cache

import (
	base "go_code/Rekas/base"
	"sync"
)

type cache struct {
	mu         sync.Mutex
	lru        *base.Cache
	cacheBytes int64
}

func (c *cache) add(key string, value base.ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lru == nil {
		c.lru = base.NewLru(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}

func (c *cache) get(key string) (value base.ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(key); ok {
		return v.(base.ByteView), ok
	}

	return
}
