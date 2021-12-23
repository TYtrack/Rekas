/*
 * @Author: TYtrack
 * @Date: 2021-12-23 14:50:09
 * @LastEditors: TYtrack
 * @LastEditTime: 2021-12-23 16:57:11
 * @FilePath: /Rekas/rekas_cache/peers.go
 */

package cache

// PeerPicker is the interface that must be implemented to locate
// the peer that owns a specific key.
// PickPeer() 方法用于根据传入的 key 选择相应节点 PeerGetter。
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter is the interface that must be implemented by a peer.
// 接口 PeerGetter 的 Get() 方法用于从对应 group 查找缓存值。PeerGetter 就对应于上述流程中的 HTTP 客户端。
type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}
