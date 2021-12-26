/*
 * @Author: your name
 * @Date: 2021-12-23 13:33:07
 * @LastEditTime: 2021-12-26 19:21:45
 * @LastEditors: TYtrack
 * @Description: 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 * @FilePath: /Rekas/net_conn/net_conn.go
 */
package netconn

import (
	"fmt"

	cache "go_code/Rekas/rekas_cache"
	peer_hash "go_code/Rekas/rekas_hash"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const (
	defaultBasePath = "/_geecache/"
	defaultReplicas = 3
)

// HTTPPool implements PeerPicker for a pool of HTTP peers.
type HTTPPool struct {
	// this peer's base URL, e.g. "https://example.net:8000"
	self     string
	basePath string
	mu       sync.Mutex // guards peers and httpGetters
	//peer.Map是一致性哈希算法的 Map，用来根据具体的 key 选择节点
	peers *peer_hash.V2RMap
	// 映射远程节点与对应的 httpGetter。每一个远程节点对应一个 httpGetter
	httpGetters map[string]*httpGetter // keyed by e.g. "http://10.0.0.2:8008"
}

// NewHTTPPool initializes an HTTP pool of peers.
func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}
func (p *HTTPPool) PrintServer() {
	log.Printf("[NetConn ] :%v\n", p.peers.GetAll())
}

// Log info with server name
func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

// ServeHTTP handle all http requests
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, p.basePath) {
		p.Log("%s %s", r.Method, r.URL.Path)
		// /<basepath>/<groupname>/<key> required
		parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)

		if len(parts) != 2 {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		groupName := parts[0]
		key := parts[1]

		group := cache.GetGroup(groupName)
		if group == nil {
			http.Error(w, "no such group: "+groupName, http.StatusNotFound)
			return
		}

		view, err := group.Get(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(view.ByteSlice())
	} else if strings.HasPrefix(r.URL.Path, "/add_server/") {

		parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)

		if len(parts) != 2 {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		serverUrl := parts[0]

		p.Add(serverUrl)

	} else if strings.HasPrefix(r.URL.Path, "/delete_server/") {
		parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)

		if len(parts) != 2 {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		serverUrl := parts[0]

		p.Delete(serverUrl)

	}
}

// Set updates the pool's list of peers.
func (p *HTTPPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers = peer_hash.NewHash(defaultReplicas, nil)
	p.peers.Add(peers...)
	p.httpGetters = make(map[string]*httpGetter, len(peers))
	for _, peer := range peers {
		p.httpGetters[peer] = &httpGetter{baseURL: peer + p.basePath}
	}
}

// Add updates the pool's list of peers.
func (p *HTTPPool) Add(peer string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers.Add(peer)
	p.httpGetters[peer] = &httpGetter{baseURL: peer + p.basePath}
}

// Delete updates the pool's list of peers.
func (p *HTTPPool) Delete(peer string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers.Delete(peer)
	delete(p.httpGetters, peer)
}

// Unset updates the pool's list of peers.
func (p *HTTPPool) Unset(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers = peer_hash.NewHash(defaultReplicas, nil)
	p.peers.Add(peers...)
	p.httpGetters = make(map[string]*httpGetter, len(peers))
	for _, peer := range peers {
		p.httpGetters[peer] = &httpGetter{baseURL: peer + p.basePath}
	}
}

// PickPeer picks a peer according to key
func (p *HTTPPool) PickPeer(key string) (cache.PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if peer := p.peers.Get(key); peer != "" && peer != p.self {
		p.Log("Pick peer %s", peer)
		return p.httpGetters[peer], true
	}
	return nil, false
}

var _ cache.PeerGetter = (*httpGetter)(nil)

type httpGetter struct {
	baseURL string
}

//使用 net_conn.Get()访问的远程节点并获取返回值，并转换为 []bytes 类型。
func (h *httpGetter) Get(group string, key string) ([]byte, error) {
	u := fmt.Sprintf(
		"%v%v/%v",
		h.baseURL,
		url.QueryEscape(group),
		url.QueryEscape(key),
	)
	res, err := http.Get(u)
	fmt.Printf("[Server] %v %v %v\n", u, "xxx", err)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned: %v", res.Status)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}

	return bytes, nil
}
