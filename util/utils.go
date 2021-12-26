/*
 * @Author: TYtrack
 * @Description: ...
 * @Date: 2021-12-26 19:02:33
 * @LastEditors: TYtrack
 * @LastEditTime: 2021-12-26 19:09:24
 * @FilePath: /Rekas/util/utils.go
 */

package util

import "container/list"

func MapToSlice(m map[string]string) []string {
	s := make([]string, 0, len(m))
	for _, v := range m {
		s = append(s, v)
	}
	return s
}

func Contains(l *list.List, value string) (bool, *list.Element) {
	for e := l.Front(); e != nil; e = e.Next() {
		if e.Value == value {
			return true, e
		}
	}
	return false, nil
}
