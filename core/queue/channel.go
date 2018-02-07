package queue

import (
	"sync"
	"errors"
)

var QUEUE_CLOSED_ERROR = errors.New("[ChannelQueue]Queue has been closed!")

type ChannelQueue struct {
	ch chan interface{}
	closeFlagLock sync.RWMutex
	popLock sync.Mutex
	ifClose bool
}

func NewChannelQueue(cap int) *ChannelQueue {
	return &ChannelQueue{
		ch: make(chan interface{}, cap),
		ifClose: false,
	}
}


func (q *ChannelQueue) IsClose() bool {
	q.closeFlagLock.RLock()
	defer q.closeFlagLock.RUnlock()
	return q.ifClose
}

func (q *ChannelQueue) Close() (err error) {
	q.closeFlagLock.Lock()
	if q.ifClose == true {
		q.closeFlagLock.Unlock()
		return nil
	}
	defer func() {
		r := recover()
		if r != nil {
			if e, ok := r.(error); ok {
				err = e
			}
		}
		q.closeFlagLock.Unlock()
	}()
	q.ifClose = true
	close(q.ch)
	return nil
}


func (q *ChannelQueue) Pop() (interface{}, error) {
	//如果close了ch 但是ch中还有数据，就允许继续pop
	if q.IsClose() {
		q.popLock.Lock()
		defer q.popLock.Unlock()
		if q.Len() == 0 {
			return nil, QUEUE_CLOSED_ERROR
		}
	}
	e, ok:= <- q.ch
	if ok {
		return e, nil
	} else {
		return e, errors.New("[ChannelQueue]Queue is closing!")
	}
}


func (q *ChannelQueue) TryPop() (interface{}, error) {
	//如果close了ch 但是ch中还有数据，就允许继续pop
	if q.IsClose() {
		q.popLock.Lock()
		defer q.popLock.Unlock()
		if q.Len() == 0 {
			return nil, QUEUE_CLOSED_ERROR
		}
	}
	select{
	case e, ok := <- q.ch:
		if ok {
			return e, nil
		} else {
			return e, errors.New("[ChannelQueue]Queue is closing!")
		}
	default:
		return nil, nil
	}
}

func (q *ChannelQueue) Push(v interface{}) (error) {
	if q.IsClose() {
		return QUEUE_CLOSED_ERROR
	}
	q.ch <- v
	return nil
}

func (q *ChannelQueue) TryPush(v interface{}) (err error) {
	if q.IsClose() {
		return QUEUE_CLOSED_ERROR
	}
	defer func() {
		r := recover()
		if r != nil {
			if e, ok := r.(error); ok {
				err = e
			}
		}
	}()
	select{
	case q.ch <- v:
		return nil
	default:
		return errors.New("[ChannelQueue]Queue is full!")
	}
	return nil
}


func (q *ChannelQueue) Cap() int {
	return cap(q.ch)
}

func (q *ChannelQueue) Len() int {
	return len(q.ch)
}
