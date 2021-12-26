/*
 * @Author: TYtrack
 * @Description: 连接检测算法
 * @Date: 2021-12-23 21:12:47
 * @LastEditors: TYtrack
 * @LastEditTime: 2021-12-26 14:40:02
 * @FilePath: /Rekas/heartbeat/heartbeat.go
 */

package heartbeat

import (
	"fmt"
	"net"
	"time"
)

const KEEP_INTERVAL = 30

func IsOpenPort(address string) (isOPen bool, err error) {

	conn, err := net.DialTimeout("tcp", address, 1*time.Second)

	if err != nil {
		return false, err
	} else {
		if conn != nil {
			defer conn.Close()
			return true, nil
		} else {
			return false, nil
		}
	}

}

// 向urls发送心跳包
func FixingSendHeartBeat(urls []string) (res []string) {
	res = make([]string, 0)
	for _, url := range urls {
		url = url[7:]
		isOpen, err := IsOpenPort(url[7:])
		if isOpen == false {
			fmt.Println(url, " 连接失败 ", err)
			res = append(res, url)
		} else {
			fmt.Println(url, " 连接成功 ")
		}
	}
	return

}
