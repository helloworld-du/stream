package queue

import (
	"sync"
	"errors"
)

type ChannelQueue struct {
	ch chan interface{}
	closeFlagLock sync.Mutex
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
			return nil, errors.New("[ChannelQueue]Queue has been closed!")
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
			return nil, errors.New("[ChannelQueue]Queue has been closed!")
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
		return errors.New("[ChannelQueue]Queue has been closed!")
	}
	q.ch <- v
	return nil
}

func (q *ChannelQueue) TryPush(v interface{}) (err error) {
	if q.IsClose() {
		return errors.New("[ChannelQueue]Queue has been closed!")
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
