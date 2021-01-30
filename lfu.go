package lfu

type dbNode struct {
	key, value interface{}
	prev, next *dbNode
	freqNode   *freqNode
}

type dbList struct {
	head, tail *dbNode
	total      int
}

type freqNode struct {
	freq       int
	dl         *dbList
	prev, next *freqNode
}

type freqList struct {
	head *freqNode
	tail *freqNode
}

func (fl *freqList) removeNode(node *freqNode) {
	node.next.prev = node.prev
	node.prev.next = node.next
}

func (fl *freqList) lastFreq() *freqNode {
	return fl.tail.prev
}

func (fl *freqList) addNode(node *dbNode) {
	if fqNode := fl.lastFreq(); fqNode.freq == 1 {
		node.freqNode = fqNode
		fqNode.dl.addToHead(node)
	} else {
		newNode := &freqNode{
			freq: 1,
			dl:   initdbList(),
		}

		node.freqNode = newNode
		newNode.dl.addToHead(node)

		fqNode.next = newNode
		newNode.prev = fqNode
		newNode.next = fl.tail
		fl.tail.prev = newNode
	}
}

func (dbl *dbList) isEmpty() bool {
	return dbl.total == 0
}

func (dbl *dbList) GetTotal() int {
	return dbl.total
}

func (dbl *dbList) addToHead(node *dbNode) {
	node.next = dbl.head.next
	node.prev = dbl.head
	dbl.head.next.prev = node
	dbl.head.next = node
	dbl.total++
}

func (dbl *dbList) removeNode(node *dbNode) {
	node.next.prev = node.prev
	node.prev.next = node.next
	dbl.total--
}

func (dbl *dbList) moveToHead(node *dbNode) {
	dbl.removeNode(node)
	dbl.addToHead(node)
}

func (dbl *dbList) removeTail() *dbNode {
	node := dbl.tail.prev
	dbl.removeNode(node)
	return node
}

func initNode(k, v interface{}) *dbNode {
	return &dbNode{
		key:   k,
		value: v,
	}
}

func initdbList() *dbList {
	l := dbList{
		head: initNode(0, 0),
		tail: initNode(0, 0),
	}
	l.head.next = l.tail
	l.tail.prev = l.head

	return &l
}

type LFUCache struct {
	cache          map[interface{}]*dbNode
	size, capacity int
	freqList       *freqList
}

func NewLFU(capacity int) LFUCache {
	ca := LFUCache{
		capacity: capacity,
		cache:    make(map[interface{}]*dbNode),
	}
	ca.freqList = &freqList{
		head: &freqNode{},
		tail: &freqNode{},
	}
	ca.freqList.head.next = ca.freqList.tail
	ca.freqList.tail.prev = ca.freqList.head
	return ca
}

func (lfu *LFUCache) incrFreq(node *dbNode) {
	curfreqNode := node.freqNode
	curdbNode := curfreqNode.dl

	if curfreqNode.prev.freq == curfreqNode.freq+1 {
		curdbNode.removeNode(node)
		curfreqNode.prev.dl.addToHead(node)
		node.freqNode = curfreqNode.prev
	} else if curdbNode.GetTotal() == 1 {
		curfreqNode.freq++
	} else {
		curdbNode.removeNode(node)
		newFreqNode := &freqNode{
			freq: curfreqNode.freq + 1,
			dl:   initdbList(),
		}
		newFreqNode.dl.addToHead(node)
		node.freqNode = newFreqNode
		newFreqNode.next = curfreqNode
		newFreqNode.prev = curfreqNode.prev
		curfreqNode.prev.next = newFreqNode
		curfreqNode.prev = newFreqNode
	}

	if curdbNode.isEmpty() {
		lfu.freqList.removeNode(curfreqNode)
	}
}

func (lfu *LFUCache) Get(key interface{}) interface{} {
	if n, ok := lfu.cache[key]; ok {
		lfu.incrFreq(n)
		return n.value
	}
	return -1
}

func (lfu *LFUCache) Put(key interface{}, value interface{}) {
	if lfu.capacity == 0 {
		return
	}
	if n, ok := lfu.cache[key]; ok {
		n.value = value
		lfu.incrFreq(n)
	} else {
		if lfu.size >= lfu.capacity {
			fqNode := lfu.freqList.lastFreq()
			node := fqNode.dl.removeTail()
			lfu.size--
			delete(lfu.cache, node.key)
		}

		newNode := initNode(key, value)
		lfu.cache[key] = newNode
		lfu.freqList.addNode(newNode)
		lfu.size++
	}
}

func (lfu *LFUCache) GetIterator() func() *dbNode {
	curFreqNode := lfu.freqList.head.next
	var dump *dbNode
	return func() *dbNode {
		for {
			if dump == nil {
				dump = curFreqNode.dl.head.next
			}

			for {
				if dump == curFreqNode.dl.tail {
					break
				}
				ret := dump
				dump = dump.next
				return ret
			}

			curFreqNode = curFreqNode.next
			if curFreqNode == lfu.freqList.tail {
				return nil
			}
			dump = curFreqNode.dl.head.next
		}
	}
}

func (lfu *LFUCache) GetAll() []interface{} {
	var ret []interface{}

	curFreqNode := lfu.freqList.head.next
	for {
		if curFreqNode == lfu.freqList.tail {
			return ret
		}

		dump := curFreqNode.dl.head.next
		for {
			if dump == curFreqNode.dl.tail {
				break
			}
			ret = append(ret, dump.value)
			dump = dump.next

		}
		curFreqNode = curFreqNode.next
	}
}

