package util

import (
	"bufio"
)

func ReadAll(br *bufio.Reader, data []byte) (err error) {
	tl, n, t := len(data), 0, 0
	for {
		if t, err = br.Read(data[n:]); err != nil {
			return
		}
		if n += t; n == tl {
			break
		}
	}
	return
}
