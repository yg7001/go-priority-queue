# go-priority-queue
a simple implementation of blocking priority-queue, constructed using two channels

Using PriorityQueue.EPushBack() to insert an emergent event; using Priority.NPushBack()
to insert a normal event; using PriorityQueue.Pop() to get an event

In this implementation, there are only two kind of events: emergent events and normal events
when pop an event from the PriorityQueue, it will first pop an emergent events(if any)
and then pop a normal events(due to different setting, PriorityQueue will pop all emergnent events
and then pop the normal events or pop serveral emergent events and then pop a normal event)

One can modify the PriorityQueue to provide several different prorities by inserting more channels
when define and construct the PriorityQueue
