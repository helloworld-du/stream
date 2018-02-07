package queue


import (
	"testing"
)



func Test_multi_queue_1(t *testing.T) {
	qu1 := NewChannelQueue(100)
	qu2 := NewChannelQueue(10)
	mq := NewMultiQueue([]IQueue{qu1, qu2})


	if mq.IsClose() != false {
		t.Errorf("queue is closed")
		t.FailNow()
	}


	/****************** basic opt******************/
	err := mq.TryPush("123")
	if err != nil {
		t.Errorf("TryPush error")
		t.FailNow()
	}
	err = mq.Push("234")
	if err != nil {
		t.Errorf("Push error")
		t.FailNow()
	}
	msg, err := qu1.TryPop()
	if err != nil || msg == nil {
		t.Errorf("TryPop error")
		t.FailNow()
	}
	if strMsg, ok := msg.(string); !ok || strMsg != "123" {
		t.Errorf("pop wrong msg")
		t.FailNow()
	}
	msg, err = qu1.Pop()
	if err != nil || msg == nil {
		t.Errorf("TryPop error")
		t.FailNow()
	}
	if strMsg, ok := msg.(string); !ok || strMsg != "234" {
		t.Errorf("pop wrong msg")
		t.FailNow()
	}

	msg, err = qu2.TryPop()
	if err != nil || msg == nil {
		t.Errorf("TryPop error")
		t.FailNow()
	}
	if strMsg, ok := msg.(string); !ok || strMsg != "123" {
		t.Errorf("pop wrong msg")
		t.FailNow()
	}
	msg, err = qu2.Pop()
	if err != nil || msg == nil {
		t.Errorf("TryPop error")
		t.FailNow()
	}
	if strMsg, ok := msg.(string); !ok || strMsg != "234" {
		t.Errorf("pop wrong msg")
		t.FailNow()
	}


	mq.Push("888")
	mq.Push("999")

	/******************close******************/
	err = mq.Close()
	if err != nil {
		t.Errorf("close errerr:%v", err)
		t.FailNow()
	}
	if mq.IsClose() != true {
		t.Errorf("close flag is wrong")
		t.FailNow()
	}
	err = mq.Close()
	if err != nil {
		t.Errorf("close errerr:%v", err)
		t.FailNow()
	}

	/******************after close******************/
	err = mq.TryPush("345")
	if err == nil {
		t.Errorf("TryPush to closed queue")
		t.FailNow()
	}
	err = mq.Push("456")
	if err == nil {
		t.Errorf("Push error")
		t.FailNow()
	}

	/******************support pop after close******************/
	msg, err = qu1.TryPop()
	if err != nil || msg == nil {
		t.Errorf("TryPop error")
		t.FailNow()
	}
	if strMsg, ok := msg.(string); !ok || strMsg != "888" {
		t.Errorf("pop wrong msg")
		t.FailNow()
	}

	msg, err = qu1.Pop()
	if err != nil || msg == nil {
		t.Errorf("TryPop error")
		t.FailNow()
	}
	if strMsg, ok := msg.(string); !ok || strMsg != "999" {
		t.Errorf("pop wrong msg")
		t.FailNow()
	}

	/******************queue is empty******************/
	msg, err = qu1.TryPop()
	if err == nil || msg != nil {
		t.Errorf("pop something from an empty queue 2; err:%v, msg:%v", err, msg)
		t.FailNow()
	}

	/******************support pop after close******************/
	msg, err = qu2.TryPop()
	if err != nil || msg == nil {
		t.Errorf("TryPop error")
		t.FailNow()
	}
	if strMsg, ok := msg.(string); !ok || strMsg != "888" {
		t.Errorf("pop wrong msg")
		t.FailNow()
	}

	msg, err = qu2.Pop()
	if err != nil || msg == nil {
		t.Errorf("TryPop error")
		t.FailNow()
	}
	if strMsg, ok := msg.(string); !ok || strMsg != "999" {
		t.Errorf("pop wrong msg")
		t.FailNow()
	}

	/******************queue is empty******************/
	msg, err = qu2.TryPop()
	if err == nil || msg != nil {
		t.Errorf("pop something from an empty queue 2; err:%v, msg:%v", err, msg)
		t.FailNow()
	}

}