package ringbuffer

import (
	"runtime"
	"sync/atomic"
	"unsafe"
)

type RingBuffer struct {
	size uint64
	buff []*value
	tail uint64
	pad_ [3]uint64

	head uint64
}

type value struct {
	seqid    uint64
	pending  bool
	enqueued bool
	data     interface{}
}

func newValue(data interface{}) *value {
	v := &value{
		data: data,
	}
	return v
}

func New(n int) *RingBuffer {
	size := roundup(n)
	buff := make([]*value, size)
	for i := uint64(0); i < size; i++ {
		buff[i] = &value{
			seqid: i + size,
		}
	}
	return &RingBuffer{
		size: size,
		buff: buff,
		head: size,
		tail: size,
	}
}

func roundup(n int) uint64 {
	n += 1
	if (n >> 3 << 3) > 0 {
		return uint64(n &^ 7)
	}
	return uint64(8)
}

// enqueue returns false if the queue is full
func (r *RingBuffer) enqueue(v interface{}) bool {
	ticket := atomic.AddUint64(&r.tail, 1)
	pos := ticket % r.size
	target := (*unsafe.Pointer)(unsafe.Pointer(&r.buff[pos]))
load:
	vp := atomic.LoadPointer(target)
	val := (*value)(vp)

	if ticket > val.seqid { // wait other's enqueue
		print("dbg", ticket, " ", val.seqid, "\n")
		runtime.Gosched()
		goto load
	}
	if val.data == nil {
		if !val.pending { // prepare enqueue
			newv := &value{
				seqid:   ticket,
				pending: true,
			}
			atomic.StorePointer(target, unsafe.Pointer(newv))
			print("prepare ", ticket, "\n")
			goto load
		}
		// try enqueue
		newv := &value{
			seqid:   ticket + r.size,
			pending: false,
			data:    v,
		}
		atomic.CompareAndSwapPointer(target, vp, unsafe.Pointer(newv))
		print(val.seqid, " ", ticket, "\n")
		return true
	} else {
		if val.pending { // wait for dequeue
			goto load
		}
		print("full ", ticket, "\n")
		return false
	}
}
