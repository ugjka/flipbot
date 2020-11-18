package main

import "strings"

func roundDuration(dur string) string {
	arr := strings.Split(dur, " ")
	if len(arr) > 2 {
		return strings.Join(arr[:4], " ")
	}
	return strings.Join(arr[:2], " ")
}
