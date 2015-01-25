package main

import (
	"fmt"
	"time"
	"strconv"
	"math/rand"
	"os"

	"go-priority-queue/priority_queue"
)

func main() {

	args := os.Args
	if len(args) < 3 {
		fmt.Println("Usage: ./test_priority Ntot level")
		return
	}
	Ntot, err := strconv.Atoi(args[1])
	level, err := strconv.Atoi(args[2])

	normLen := 100
	emergLen := 200
	fmt.Println("\n", "Ntot: ", Ntot, " level: ", level, "\n")
	pQueue, err := priority_queue.NewPriorityQueue(normLen, emergLen, level)
	if err != nil {
		fmt.Println("Error ", err)
	}

	sig := make(chan bool, 2)
	go PrintMessage(pQueue, Ntot, sig)
	go PutMessage(pQueue, Ntot, sig)
	<-sig
	<-sig
}

func PutMessage(pQueue *priority_queue.PriorityQueue, N int, sig chan bool) {
	label := []string{"Normal", "Emergent"}
	for idx := 0; idx < N; idx++ {
		rand.Seed(time.Now().UnixNano())
		itype := rand.Intn(2)
//		itype := idx%2
		packet := &priority_queue.Packet{label[itype] + strconv.Itoa(idx)}
		switch itype {
		case 0:
			pQueue.NPushBack(packet)
		case 1:
			pQueue.EPushBack(packet)
		}
		fmt.Println("\n", time.Now().String(), "Put a packet: ", packet.Msg, "\n")
		n := rand.Intn(10)
		time.Sleep(time.Duration(n*1e7))
	}
	sig<-true
}

func PrintMessage(pQueue *priority_queue.PriorityQueue, N int, sig chan bool) {
	for idx := 0; idx < N; idx++ {
		packet := pQueue.Pop()
		fmt.Println("\n", time.Now().String(), "Got a packet: ", packet.Msg, "\n")
		time.Sleep(2*time.Millisecond)
	}
	sig<-true
}

