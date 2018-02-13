package reader

import (
	"github.com/helloworld-du/stream/core/queue"
	"github.com/helloworld-du/stream/core/provider"

	"errors"
	"sync"
	"time"
	//"fmt"
)
func NewsEasyReader(p provider.IProvider, msgQ, errQ queue.IQueue, workerNum int) (*EasyReader, error) {
	if workerNum <= 0 {
		return nil, errors.New("[EasyReader]worker num in below 0")
	}
	return &EasyReader{workerNum, false, &sync.RWMutex{},&sync.WaitGroup{},
	p, msgQ, errQ}, nil
}


type EasyReader struct {
	workerNum int
	isStop bool
	lock *sync.RWMutex
	wg *sync.WaitGroup
	dataProvider provider.IProvider
	msgQueue queue.IQueue
	errQueue queue.IQueue
}

func (r *EasyReader) Stop() (err error) {
	/*
	先设stop置标记位，这样run函数在最后一次读取后就会推出for循环
	然后等所有的worker都完成了再关闭queue
	给等待worker完成一个超时时间（1秒）
	最后关闭queue，如果所有worker都在1秒内完成，以后读出来的数据就缓存在queue中，不然会丢失
	 */
	r.lock.Lock()
	if r.isStop == true {
		r.lock.Unlock()
		return
	}
	r.isStop = true
	r.lock.Unlock()

	defer func() {
		if err1 := r.msgQueue.Close(); err1 != nil && err == nil  {
			err = err1
		}

		if err1 := r.errQueue.Close(); err1 != nil && err == nil  {
			err = err1
		}
	}()

	select {
		case <- time.After(time.Second):
			err = errors.New("did not stop after 1 second")
		case <- func () (chan bool) {
			r.wg.Wait()
			ch := make(chan bool, 1)
			ch <- true
			return ch
		}():
			err = nil
	}
	return
}



func (r *EasyReader) IsStop() bool {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.isStop
}



func (r *EasyReader) Run() error {
	/**
	判断stop标记位，如果已经stop，就不要继续run了
	拉起对应的worker进程，r.wg可以用来表示是否所有的worker都完成了工作
	 */
	if r.IsStop() {
		return errors.New("[EasyReader]can not run on a closed queue")
	}
	for i := 0; i < r.workerNum; i++ {
		r.wg.Add(1)
		go func (r *EasyReader, wg *sync.WaitGroup){
			canStop := false;
			for ; !r.IsStop(); {
				msg, err, hasNext := r.dataProvider.Read()
				//fmt.Println("r.dataProvider.Read", msg, err, hasNext)
				if err != nil {
					r.errQueue.TryPush(err)
				} else {
					r.msgQueue.Push(msg)
				}
				if !hasNext {
					canStop = true
					break
				}
			}
			wg.Done()
			// data provider标记没有更多数据了，通知reader推出
			// 先标记本worker已完成，再stop
			if canStop {
				r.Stop()
			}
		}(r, r.wg)
	}

	return nil
}