package main

import "fmt"

func main() {
	s := "sheinerintheworld"
	lens := len(s)
	lo := 0
	hi := 0
	maxx := 0
	ans := make([]string, 0)
	cnt := make(map[uint8]bool)
	for hi < lens {
		ch := s[hi]
		if _, ok := cnt[ch]; ok {
			for lo < hi && s[lo] != ch {
				delete(cnt, s[lo])
				lo++
			}
			lo++
		}
		if maxx < hi-lo+1 {
			maxx = hi - lo + 1
			ans = make([]string, 0)
			ans = append(ans, s[lo:hi+1])
		} else if maxx == hi-lo+1 {
			ans = append(ans, s[lo:hi+1])
		}
		cnt[ch] = true
		hi++
	}

	for _, key := range ans {
		fmt.Println(key)
	}
}
