package kk

import (
	"sync"
	"time"
)

const (
	twepoch = int64(1424016000000000)
	aidBits = uint(5)
	nidBits = uint(9)
	sBits   = uint(10)

	/*
	 * 1 符号位  |  39 时间戳                                     | 5 区域  |  9 节点       | 10 （毫秒内）自增ID
	 * 0        |  0000000 00000000 00000000 00000000 00000000  | 00000  | 000000 000   |  000000 0000
	 *
	 */

	nidShift = sBits //左移次数
	aidShift = sBits + nidBits
	tShift   = sBits + aidBits + nidBits
	sMask    = -1 ^ (-1 << sBits)
)

func Milliseconds() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

type IID struct {
	s    int64 //序号
	aid  int64 //区域ID
	nid  int64 //节点ID
	lock sync.Mutex
}

func NewIID(aid int64, nid int64) *IID {
	v := IID{}
	v.s = 0
	v.nid = nid
	v.aid = aid
	v.lock = sync.Mutex{}
	return &v
}

func (v *IID) NewID() int64 {
	v.lock.Lock()
	defer v.lock.Unlock()

	id := Milliseconds()
	v.s = (v.s + 1) & sMask

	return ((id - twepoch) << tShift) | (v.aid << aidShift) | (v.nid << nidShift) | v.s
}
