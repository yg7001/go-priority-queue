/*
*Impletation of a simple blocking priority queue, which is constructed using channels. 
*When Priority Queue is empty, read from it will be blocked
*Using NPushBack() and EPushBack() to insert normal events and emergent events, (only) using Pop() to get an event
*from the queue
*the Mutex is used to protect the QueueSize: NormalQueueSize and EmergentQueueSize
*
*Since the QueueSize is not strictly synchronized with the elements number in channels, the QueueSize may be < 0 some times, 
*But finally, it will be zero(Evently, the QueueSize is synchronized to the elements number of channel)
*
*Author zhang jingqing, zhang_jingqing@126.com
*/

package priority_queue

import (
	"errors"
	"sync"

	"fmt"
	"time"
)


type Packet struct {
	Msg string
}


type PriorityQueue struct {
	NormalQueue chan *Packet
	NormalQueueLen int
	NormalQueueSize int
	EmergentQueue chan *Packet
	EmergentQueueLen int
	EmergentQueueSize int
	Level int
	ECount int
	sync.Mutex
}

func NewPriorityQueue(nQueueLen, eQueueLen, level int) (*PriorityQueue, error) {
	if nQueueLen <= 0 || eQueueLen <= 0 {
		err := errors.New("NewPriorityQueue() failed, queue length must be positive")
		fmt.Println(err)
		return nil, err
	}

	queue := &PriorityQueue {
		NormalQueueLen : nQueueLen,
		EmergentQueueLen : eQueueLen,
		Level : level,
		NormalQueue : make(chan *Packet, nQueueLen),
		EmergentQueue : make(chan *Packet, eQueueLen),
	}
	return queue, nil
}

func (this *PriorityQueue) NPushBack(packet *Packet) error {
	fmt.Println("\n", time.Now().String(), " in NPushBack ", "\n")
	if packet == nil {
		err := errors.New("Priority.NPushBack() failed, got a nil pointer")
		fmt.Println(err)
		return err
	}
	this.NormalQueue<-packet

	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	this.NormalQueueSize++
	return nil
}

//non-blocking pop
/*
since the QueueSize is not strictly sync with the elemments in chann, 
in nPop() and ePop(), the condition in if must be <= 0 if the Mutex is Locked before reading from the queue

the below implemetation is more safe: condition must be <= 0 and the Mutex is locked after reading from the queue
*/
func (this *PriorityQueue) nPop() *Packet {
	if this.NormalQueueSize <= 0 {
		return nil
	}
	packet := <-this.NormalQueue
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	this.NormalQueueSize--
	return packet
}

func (this *PriorityQueue) EPushBack(packet *Packet) error {
	if packet == nil {
		err := errors.New("PriorityQueue.EPushBack() failed, got a nil pointer")
		fmt.Println(err)
		return err
	}
	this.EmergentQueue <- packet
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	this.EmergentQueueSize++
	return nil
}

//non-blocking pop
func (this *PriorityQueue) ePop() *Packet {
	if this.EmergentQueueSize <= 0 {
		return nil
	}
	packet := <-this.EmergentQueue
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	this.EmergentQueueSize--
	return packet
}

//blocking pop
func (this *PriorityQueue) Pop() *Packet {
	packet := this.prioPop()

	if packet == nil {
		select {
		case packet = <-this.EmergentQueue:
			this.Mutex.Lock()
			this.EmergentQueueSize--
			this.Mutex.Unlock()
		case packet = <-this.NormalQueue:
			this.Mutex.Lock()
			this.NormalQueueSize--
			this.Mutex.Unlock()
		}
	}
	return packet
}

//non-blocking pop
func (this *PriorityQueue) prioPop() *Packet {
	var packet *Packet
	if this.Level <= 0 || this.Level > 3 {
		//default case, handle all emergent events then handle normal events
		packet = this.ePop()
		if packet == nil {
			packet = this.nPop()
		}
	} else {
		//only in this function, the ECount will be modified, so it is not necessary to protect it by a Lock
		if this.ECount < this.Level {
			packet = this.ePop()
			this.ECount++
			if packet == nil {
				packet = this.nPop()
			}
		} else {
			this.ECount = 0
			packet = this.nPop()
			if packet == nil {
				packet = this.ePop()
			}
		}
	}
	return packet
}
