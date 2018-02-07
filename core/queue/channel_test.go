package queue


import (
	"testing"
)



func Test_ChannelQueue_1(t *testing.T) {
	qu := NewChannelQueue(100)
	if qu.IsClose() != false {
		t.Errorf("queue is closed")
		t.FailNow()
	}

	/****************** pop from empty queue******************/
	msg, err := qu.TryPop()
	if err != nil || msg != nil {
		t.Errorf("pop something from an empty queue; err:%v, msg:%v", err, msg)
		t.FailNow()
	}

	/****************** basic opt******************/
	err = qu.TryPush("123")
	if err != nil {
		t.Errorf("TryPush error")
		t.FailNow()
	}
	err = qu.Push("234")
	if err != nil {
		t.Errorf("Push error")
		t.FailNow()
	}
	msg, err = qu.TryPop()
	if err != nil || msg == nil {
		t.Errorf("TryPop error")
		t.FailNow()
	}
	if strMsg, ok := msg.(string); !ok || strMsg != "123" {
		t.Errorf("pop wrong msg")
		t.FailNow()
	}
	msg, err = qu.Pop()
	if err != nil || msg == nil {
		t.Errorf("TryPop error")
		t.FailNow()
	}
	if strMsg, ok := msg.(string); !ok || strMsg != "234" {
		t.Errorf("pop wrong msg")
		t.FailNow()
	}

	qu.Push("888")
	qu.Push("999")

	/******************close******************/
	err = qu.Close()
	if err != nil {
		t.Errorf("close errerr:%v", err)
		t.FailNow()
	}
	if qu.IsClose() != true {
		t.Errorf("close flag is wrong")
		t.FailNow()
	}
	err = qu.Close()
	if err != nil {
		t.Errorf("close errerr:%v", err)
		t.FailNow()
	}

	/******************after close******************/
	err = qu.TryPush("345")
	if err == nil {
		t.Errorf("TryPush to closed queue")
		t.FailNow()
	}
	err = qu.Push("456")
	if err == nil {
		t.Errorf("Push error")
		t.FailNow()
	}

	/******************support pop after close******************/
	msg, err = qu.TryPop()
	if err != nil || msg == nil {
		t.Errorf("TryPop error")
		t.FailNow()
	}
	if strMsg, ok := msg.(string); !ok || strMsg != "888" {
		t.Errorf("pop wrong msg")
		t.FailNow()
	}

	msg, err = qu.Pop()
	if err != nil || msg == nil {
		t.Errorf("TryPop error")
		t.FailNow()
	}
	if strMsg, ok := msg.(string); !ok || strMsg != "999" {
		t.Errorf("pop wrong msg")
		t.FailNow()
	}

	/******************queue is empty******************/
	msg, err = qu.TryPop()
	if err == nil || msg != nil {
		t.Errorf("pop something from an empty queue 2; err:%v, msg:%v", err, msg)
		t.FailNow()
	}

}