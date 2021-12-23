/*
 * @Author: your name
 * @Date: 2021-12-23 12:37:45
 * @LastEditTime: 2021-12-23 17:11:13
 * @LastEditors: TYtrack
 * @Description: 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 * @FilePath: /Rekas/rekas_cache/group.go
 */
package cache

import (
	"fmt"
	base "go_code/Rekas/base"
	singlecall "go_code/Rekas/single_call"
	"log"
	"sync"
)

// A Group is a cache namespace and associated data loaded spread over
type Group struct {
	name      string
	getter    Getter
	mainCache cache
	peers     PeerPicker

	loader *singlecall.GroupCall
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

// RegisterPeers registers a PeerPicker for choosing remote peer
func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

// NewGroup create a new instance of Group
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
		loader:    &singlecall.GroupCall{},
	}
	groups[name] = g

	return g
}

// GetGroup returns the named group previously created with NewGroup, or
// nil if there's no such group.
func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

// Get value for a key from cache
func (g *Group) Get(key string) (base.ByteView, error) {
	if key == "" {
		return base.ByteView{}, fmt.Errorf("key is required")
	}

	if v, ok := g.mainCache.get(key); ok {
		log.Println("[GeeCache] hit")
		return v, nil
	}

	return g.load(key)
}

func (g *Group) load(key string) (value base.ByteView, err error) {

	viewi, err := g.loader.Do(key, func() (interface{}, error) {
		if g.peers != nil {
			if peer, ok := g.peers.PickPeer(key); ok {
				if value, err = g.getFromPeer(peer, key); err == nil {
					return value, nil
				}
				log.Println("[GeeCache] Failed to get from peer", err)
			}
		}

		return g.getLocally(key)
	})

	if err == nil {
		return viewi.(base.ByteView), nil
	}
	return
}

// 使用实现了 PeerGetter 接口的 httpGetter 从访问远程节点，获取缓存值。
func (g *Group) getFromPeer(peer PeerGetter, key string) (base.ByteView, error) {
	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return base.ByteView{}, err
	}
	return base.GetByteView(bytes), nil
}

func (g *Group) getLocally(key string) (base.ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return base.ByteView{}, err

	}
	value := base.GetByteView(bytes)
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value base.ByteView) {
	g.mainCache.add(key, value)
}
