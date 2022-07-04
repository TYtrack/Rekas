/*
 * @Author: TYtrack
 * @Description: ...
 * @Date: 2021-12-26 19:42:45
 * @LastEditors: TYtrack
 * @LastEditTime: 2021-12-26 19:42:46
 * @FilePath: /Rekas/bloomfilter/bloomx_test.go
 */

package bloomx

import (
	"fmt"
	"testing"
)

func TestBloomx(t *testing.T) {
	filter := NewBloomFilter()
	fmt.Println(filter.Funcs[1].Seed)
	str1 := "hello,bloom filter!"
	filter.Add(str1)
	str2 := "A happy day"
	filter.Add(str2)
	str3 := "Greate wall"
	filter.Add(str3)

	fmt.Println(filter.Set.Count())
	fmt.Println(filter.Contains(str1))
	fmt.Println(filter.Contains(str2))
	fmt.Println(filter.Contains(str3))
	fmt.Println(filter.Contains("blockchain technology"))
}
func TestCuckoofilter(t *testing.T) {
	cf := NewCuckBloom(1000)
	fmt.Println(cf.Bloom.Count())
	cf.Insert([]byte("hello"))
	fmt.Println(cf.Bloom.Count())
	cf.Insert([]byte("world"))
	fmt.Println(cf.Bloom.Count())
	cf.Delete([]byte("hello"))
	fmt.Println(cf.Bloom.Count())
}
