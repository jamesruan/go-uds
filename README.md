# go-uds
Useful data structure for Go

## RingBuffer

RingBuffer is used mostly for a multi-producer multi-consumer queue.

Plan to implement a wait-free ring buffer.

### level of progress
- Obstruction-Free: At least one thread can make progress. Dead-lock free.
- Lock-Free: At least one thread can always make progress. Dead-lock free, live-lock free.
- Wait-Free: At least one thread can always make progress in finite number of steps. Dead-lock free, live-lock free, starvation free.

