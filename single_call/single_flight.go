/*
 * @Author: TYtrack
 * @Date: 2021-12-23 16:27:07
 * @LastEditors: TYtrack
 * @LastEditTime: 2021-12-23 16:46:47
 * @FilePath: /Rekas/single_call/single_flight.go
 */

package singlecall

import "sync"

//call 代表正在进行中，或已经结束的请求。使用 sync.WaitGroup 锁避免重入
type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

//管理不同 key 的请求(call)。
type GroupCall struct {
	mu sync.Mutex // protects m
	m  map[string]*call
}

func (g *GroupCall) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err
}
