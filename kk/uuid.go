package kk

import (
	"time"
)

const twepoch = int64(1424016000000000)

func milliseconds() int64 {
	return time.Now().UnixNano() / 1e3
}

var id int64 = 0

func UUID() int64 {

	var uuid int64 = 0

	GetDispatchMain().Sync(func() {
		if id == 0 {
			id = milliseconds()
		} else {
			id = id + 1
		}

		uuid = id
	})

	return uuid
}
