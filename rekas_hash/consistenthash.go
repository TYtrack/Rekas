/*
 * @Author: TYtrack
 * @Date: 2021-12-23 14:07:43
 * @LastEditTime: 2021-12-23 22:37:12
 * @LastEditors: TYtrack
 * @Description: 一致性哈希
 * @FilePath: /Rekas/rekas_hash/consistenthash.go
 */
package peer_hash

import (
	"fmt"
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash maps bytes to uint32
type Hash func(data []byte) uint32

// Map constains all hashed keys
//虚拟节点与真实节点的映射 结构
type V2RMap struct {
	hash     Hash
	replicas int            // 虚拟节点倍数
	keys     []int          // 虚拟节点（已排好序的）
	hashMap  map[int]string // 虚拟节点与真实节点的映射表 hashMap
}

// New creates a Map instance
func NewHash(replicas int, fn Hash) *V2RMap {
	m := &V2RMap{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add adds some keys to the hash.
func (m *V2RMap) Add(keys ...string) {
	fmt.Println("*****", keys)
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

// Add deletes some keys to the hash.
func (m *V2RMap) Delete(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			idx := sort.Search(len(m.keys), func(i int) bool {
				return m.keys[i] == hash
			})
			m.keys = append(m.keys[:idx], m.keys[idx+1:]...)
			delete(m.hashMap, hash)

		}
	}

}

// Get gets the closest item in the hash to the provided key.
func (m *V2RMap) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))
	// Binary search for appropriate replica.
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	//因为sort.Search找不到满足条件时不是返回-1，而是数组的长度
	return m.hashMap[m.keys[idx%len(m.keys)]]
}

func (m *V2RMap) GetAll() string {
	return fmt.Sprintf("[consistent] :%v\n", m.hashMap)
}
