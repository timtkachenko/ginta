package main

import "sync"

type counterSync struct {
	mx    *sync.RWMutex
	store int
}

func (cs *counterSync) Incr() {
	cs.mx.Lock()
	defer cs.mx.Unlock()
	cs.store++
}
func (cs *counterSync) Get() int {
	cs.mx.RLock()
	defer cs.mx.RUnlock()
	return cs.store
}
func (cs *counterSync) Reset() {
	cs.mx.Lock()
	defer cs.mx.Unlock()
	cs.store = 0
}

type queueSync struct {
	mx    *sync.RWMutex
	store [][]byte
}

func (qs *queueSync) Set(val []byte) {
	qs.mx.Lock()
	defer qs.mx.Unlock()
	qs.store = append(qs.store, val)
}

func (qs *queueSync) DrainAll() [][]byte {
	qs.mx.RLock()
	all := qs.store
	qs.store = [][]byte{}
	qs.mx.RUnlock()
	return all
}
