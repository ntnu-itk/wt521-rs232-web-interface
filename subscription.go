package main

import (
	"container/list"
	"sync"
)

type StateSubscription struct {
	ch          chan State
	persistence SubscriptionPersistence
}

type SubscriptionPersistence int

const (
	Permanent SubscriptionPersistence = -1
	OneShot   SubscriptionPersistence = 1
)

func NewStateSubscription(persistence SubscriptionPersistence) *StateSubscription {
	return &StateSubscription{ch: make(chan State, 0), persistence: persistence}
}

type StateSubscriptionLedger struct {
	subs *list.List
	sync.Mutex
}

func NewStateSubscriptionLedger() *StateSubscriptionLedger {
	return &StateSubscriptionLedger{
		subs: list.New()}
}

func (ssl *StateSubscriptionLedger) Manage(newSubscriptionChannel <-chan *StateSubscription) {
	for {
		ss := <-newSubscriptionChannel

		ssl.Lock()

		ssl.subs.PushBack(ss)

		ssl.Unlock()
	}
}

func (ssl *StateSubscriptionLedger) Broadcast(state State) (n int) {
	ssl.Lock()
	defer ssl.Unlock()

	for e := ssl.subs.Front(); e != nil; e = e.Next() {
		ss := e.Value.(*StateSubscription)

		ss.ch <- state
		n++

		if ss.persistence == OneShot {
			ssl.subs.Remove(e)
		}
	}

	return
}
