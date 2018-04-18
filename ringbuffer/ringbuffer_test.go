package ringbuffer

import (
	"github.com/davecgh/go-spew/spew"
	"testing"
)

func TestEnqueue(t *testing.T) {
	r := New(24)
	t.Logf("%v", spew.Sdump(r))
	for i := 0; i < 30; i++ {
		go r.enqueue(i)
	}
	select {}
}
