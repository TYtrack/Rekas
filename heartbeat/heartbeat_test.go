/*
 * @Author: TYtrack
 * @Description: ...
 * @Date: 2021-12-23 21:22:34
 * @LastEditors: TYtrack
 * @LastEditTime: 2021-12-26 20:05:49
 * @FilePath: /Rekas/heartbeat/heartbeat_test.go
 */
package heartbeat

import (
	"fmt"
	"testing"
)

func TestIsOpenPort(t *testing.T) {
	for i := 8000; i < 8004; i++ {
		isOpen, err := IsOpenPort("0.0.0.0", i)
		if isOpen == false {
			fmt.Println(i, " 连接失败 ", err)
		} else {
			fmt.Println(i, " 连接成功 ")
		}
	}

}
