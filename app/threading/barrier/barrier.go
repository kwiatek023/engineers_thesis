package barrier

import (
	"sync"
)

type Barrier struct {
	workerCounter        int
	nofWorkers           int
	m                    sync.Mutex
	firstBarrierChannel  chan int
	secondBarrierChannel chan int
}

func New(nofWorkers int) *Barrier {
	b := Barrier{
		nofWorkers:           nofWorkers,
		firstBarrierChannel:  make(chan int, 1),
		secondBarrierChannel: make(chan int, 1),
	}
	// close 2nd barrier
	b.secondBarrierChannel <- 1
	return &b
}

func (b *Barrier) WaitAtFirstBarrier() {
	b.m.Lock()
	b.workerCounter += 1
	if b.workerCounter == b.nofWorkers {
		// close 2nd barrier
		<-b.secondBarrierChannel
		// open 1st barrier
		b.firstBarrierChannel <- 1
	}
	b.m.Unlock()
	<-b.firstBarrierChannel
	b.firstBarrierChannel <- 1
}

func (b *Barrier) WaitAtSecondBarrier() {
	b.m.Lock()
	b.workerCounter -= 1
	if b.workerCounter == 0 {
		// close 1st barrier
		<-b.firstBarrierChannel
		// open 2st barrier
		b.secondBarrierChannel <- 1
	}
	b.m.Unlock()
	<-b.secondBarrierChannel
	b.secondBarrierChannel <- 1
}

func (b *Barrier) SetNofWorkers(nofActiveWorkers int) {
	b.nofWorkers = nofActiveWorkers
}

func (b *Barrier) Close() {
	close(b.firstBarrierChannel)
	close(b.secondBarrierChannel)
}
