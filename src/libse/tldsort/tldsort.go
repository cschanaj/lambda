package Tldsort

import (
	"strings"
)

type Order []string

func (s Order) Len() int {
	return len(s)
}

func (s Order) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Order) Less(i, j int) bool {
	h := strings.Split(s[i], ".")
	k := strings.Split(s[j], ".")

	for len(h) > 0 && len(k) > 0 {
		if m, n := h[len(h)-1], k[len(k)-1]; m != n {
			return m < n
		}

		h, k = h[:len(h)-1], k[:len(k)-1]
	}
	return len(h) < len(k)
}
