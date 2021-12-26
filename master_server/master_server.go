/*
 * @Author: TYtrack
 * @Description: ...
 * @Date: 2021-12-23 21:31:46
 * @LastEditors: TYtrack
 * @LastEditTime: 2021-12-26 19:20:24
 * @FilePath: /Rekas/master_server/master_server.go
 */

package master

import (
	"fmt"
	"go_code/Rekas/heartbeat"
	netconn "go_code/Rekas/net_conn"
	cache "go_code/Rekas/rekas_cache"
	"go_code/Rekas/util"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

type MasterServer struct {
	ServerList []string
	NumServer  int

	gee   *cache.Group
	peers *netconn.HTTPPool

	mu sync.Mutex
}

func NewMasterServer(g1 *cache.Group, p1 *netconn.HTTPPool) (m_server *MasterServer) {
	return &MasterServer{
		ServerList: make([]string, 0),
		NumServer:  0,
		gee:        g1,
		peers:      p1,
	}
}

func readAllServerFromConfig() (res []string) {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("read config failed: %v", err)
	}

	map_res := viper.GetStringMapString("server")
	res = util.MapToSlice(map_res)
	log.Println("[master]: ", res)
	return
}

func (m_server *MasterServer) AutoLookServer(arg []string) {
	if arg == nil {
		arg = readAllServerFromConfig()
	}
	m_server.mu.Lock()
	defer m_server.mu.Unlock()
	m_server.ServerList = make([]string, 0)
	for _, url := range arg {
		url = strings.TrimLeft(url, "http://")
		url = strings.TrimLeft(url, "https://")
		isOpen, err := heartbeat.IsOpenPort(url)
		if err != nil {
			log.Println("[ Master ]: ", err)
		}
		if isOpen {
			m_server.AddServer(url)
		}
	}
	m_server.NumServer = len(m_server.ServerList)

	log.Println("[ Master ]: Rekas cache add a server: ", m_server.ServerList)
}

func (m_server *MasterServer) NotifyServer() {
	res := readAllServerFromConfig()
	arg := make([]string, 0)

	arg = append(arg, res...)

	for _, send_url := range res {
		new_send_url := fmt.Sprintf(
			"%v/server_list",
			send_url,
		)

		_, err := SendHttpPost(new_send_url, arg)
		if err != nil {
			log.Println("[ Master ] Send Message error : ", err)
			continue
		}
	}
}

func (m_server *MasterServer) NotifyAddServer(add_url string) {
	res := readAllServerFromConfig()

	for _, send_url := range res {
		new_send_url := fmt.Sprintf(
			"%v/add_server/%v",
			send_url,
			url.QueryEscape(add_url),
		)
		_, err := SendHttpGet(new_send_url)
		if err != nil {
			log.Println("[ Master ] Send Message error : ", err)
			continue
		}
	}

}

func (m_server *MasterServer) NotifyDeleteServer(del_url string) {
	res := readAllServerFromConfig()

	for _, send_url := range res {
		new_send_url := fmt.Sprintf(
			"%v/delete_server/%v",
			send_url,
			url.QueryEscape(del_url),
		)
		_, err := SendHttpGet(new_send_url)
		if err != nil {
			log.Println("[ Master ] Send Message error : ", err)
			continue
		}
	}

}

func SendHttpGet(send_url string) ([]byte, error) {
	res, err := http.Get(send_url)
	fmt.Printf("[Server] %v %v %v\n", send_url, "xxx", err)

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

func SendHttpPost(send_url string, serverSlice []string) ([]byte, error) {
	data := make(url.Values)
	data["ServerList"] = serverSlice

	res, err := http.PostForm(send_url, data)
	fmt.Printf("[Server] %v %v %v\n", send_url, "xxx", err)

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

func (m_server *MasterServer) AddServer(address string) {

	m_server.mu.Lock()
	defer m_server.mu.Unlock()

	m_server.NumServer++
	m_server.ServerList = append(m_server.ServerList, address)
	m_server.peers.Set(address)

	m_server.gee.AgainRegisterPeers(m_server.peers)
	log.Println("[ Master ]: Rekas cache add a server:", address)
	// log.Fatal(http.ListenAndServe(address[7:], peers))
}

func (m_server *MasterServer) RemoveServer(address string) {

	m_server.mu.Lock()
	defer m_server.mu.Unlock()

	m_server.NumServer--
	maxIdx := len(m_server.ServerList) - 1
	for i := maxIdx; i >= 0; i-- {
		if address == m_server.ServerList[i] {
			m_server.ServerList = append(m_server.ServerList[:i], m_server.ServerList[i+1:]...)
		}
	}

	m_server.peers.Unset(address)
	m_server.gee.AgainRegisterPeers(m_server.peers)
	log.Println("[ Master ]: Rekas cache remove a server:", address)
	// log.Fatal(http.ListenAndServe(address[7:], peers))
}
