package pubsub

import "sync"

type PubSub struct {
	mu   sync.RWMutex
	subs map[chan<- string]struct{}
}

func NewPubSub() *PubSub {
	ps := &PubSub{}
	ps.subs = make(map[chan<- string]struct{})
	return ps
}

func (ps *PubSub) Subscribe(ch chan<- string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.subs[ch] = struct{}{}
}

func (ps *PubSub) Unsubscribe(ch chan<- string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	delete(ps.subs, ch)
}

func (ps *PubSub) Publish(msg string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	for ch, _ := range ps.subs {
		// non-blocking push
		go func(ch chan<- string) { ch <- msg }(ch)
	}
}

func (ps *PubSub) SubscribersNum() int {
	return len(ps.subs)
}
