/*
 * @Author: your name
 * @Date: 2021-12-23 12:41:15
 * @LastEditTime: 2021-12-26 17:45:38
 * @LastEditors: TYtrack
 * @Description: 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 * @FilePath: /Rekas/main.go
 */
package main

import (
	"flag"
	"fmt"
	"go_code/Rekas/heartbeat"
	master "go_code/Rekas/master_server"
	netconn "go_code/Rekas/net_conn"
	cache "go_code/Rekas/rekas_cache"
	"log"
	"net/http"
	"time"
)

var db = map[string]string{
	"Tom":   "630",
	"Jack":  "589",
	"Sam":   "567",
	"Tom2":  "6302",
	"Jack2": "5892",
	"Sam2":  "5672",
	"Kevin": "2222",
	"Billy": "891",
	"Bob":   "232",
	"Zoo":   "343",
}

func createGroup() *cache.Group {
	return cache.NewGroup("scores", 2<<10, cache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
}

/**
 * @description:
 * @param {string} addr：该节点（缓存服务器）的地址以及IP
 * @param {[]string} addrs：所有节点的的地址以及IP
 * @param {*cache.Group} gee：该节点的缓存数据库结构体
 * @return {*}
 */
func startMasterServer(addr string, g1 *cache.Group) {
	p1 := netconn.NewHTTPPool(addr)
	m_master := master.NewMasterServer(g1, p1)
	m_master.AutoLookServer(nil)
	fmt.Println("[ Main ] all server :   ", m_master.ServerMap)

	log.Println("geecache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], p1))
}

func startCacheServer(addr string, addrs []string, gee *cache.Group) {
	peers := netconn.NewHTTPPool(addr)
	peers.Set(addrs...)
	gee.RegisterPeers(peers)
	peers.PrintServer()

	peers.Add("http://localhost:8005")
	gee.AgainRegisterPeers(peers)
	peers.PrintServer()

	peers.Delete("http://localhost:8002")
	gee.AgainRegisterPeers(peers)
	peers.PrintServer()

	go func() {
		tt := true
		for tt {
			addrs := []string{"http://0.0.0.0:8000", "http://0.0.0.0:8001", "http://0.0.0.0:8002", "http://0.0.0.0:8003", "http://0.0.0.0:8004", "http://0.0.0.0:8005"}
			res := heartbeat.FixingSendHeartBeat(addrs)
			fmt.Println("########", res)
			time.Sleep(time.Second * 10)
		}
	}()

	log.Println("geecache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

func startAPIServer(apiAddr string, gee *cache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := gee.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())

		}))
	log.Println("fontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))

}

func main() {
	var port int
	var api bool
	var master bool
	flag.IntVar(&port, "port", 8001, "Geecache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.BoolVar(&master, "master", false, "Is this server a Master?")
	flag.Parse()

	// 服务器API地址
	apiAddr := "http://0.0.0.0:9999"

	// 加入连接
	addrMap := map[int]string{
		8001: "http://0.0.0.0:8001",
		8002: "http://0.0.0.0:8002",
		8003: "http://0.0.0.0:8003",
	}

	fmt.Println(addrMap)

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	gee := createGroup()

	if api {
		go startAPIServer(apiAddr, gee)
	}
	if master {
		go startMasterServer(addrMap[port], gee)
	}
	startMasterServer(addrMap[port], gee)

	//startCacheServer(addrMap[port], []string(addrs), gee)
}
