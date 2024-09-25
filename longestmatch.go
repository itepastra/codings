package main

import (
	"slices"
)

func findMatch(search []byte, lookahead []byte) (offset int, length int, next *byte) {

	if !slices.Contains(search, lookahead[0]) {
		return 0, 0, &lookahead[0]
	}

	n := 0
	for i := len(search) - 1; i >= 0; i -= 1 {
		n += 1
		for j := 0; cmp(search, lookahead, i, j); j += 1 {
			if j > length {
				length = j
				offset = n
				if len(lookahead) > n {
					next = &lookahead[n]
				}
			}
		}
	}

	return
}

func cmp(search []byte, lookahead []byte, index int, offset int) bool {
	searchlen := len(search)
	var val byte
	if index+offset >= searchlen {
		val = lookahead[index+offset-searchlen]
	} else {
		val = search[index+offset]
	}
	return val == lookahead[offset]
}
