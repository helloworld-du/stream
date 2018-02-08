package reader

import (
	"github.com/helloworld-du/stream/core/provider"
	"github.com/helloworld-du/stream/core/queue"

	"testing"
	"strconv"
	"time"
)

func Test_ChannelQueue_1(t *testing.T) {
	q := queue.NewChannelQueue(10)
	p := provider.NewFileProvider("reader_test.dat", '\n')
	r, err := NewsEasyReader(p, q, q, 1)
	if r == nil || err != nil{
		t.Errorf("new obj err: %v", err)
		t.FailNow()
	}
	err = r.Run()
	if err != nil{
		t.Errorf("run err: %v", err)
		t.FailNow()
	}
	for i := 0; i < 31; i++ {
		ele, err := q.Pop()
		//fmt.Println(i, ele, err)
		if err != nil {
			t.Errorf("pop err: %d, %s, %v", i, ele, err)
			t.FailNow()
		}
		if e, ok := ele.(string); !ok {
			t.Errorf("pop msg err: %d, %v", i, ele)
			t.FailNow()
		} else if ie, err := strconv.Atoi(e); err != nil || ie != i {
			t.Errorf("pop msg val err: %d, %v", i, ele)
			t.FailNow()
			//fmt.Println(ie)
		}
	}
	ele, err := q.Pop()
	if ele != "" || err != nil {
		t.Errorf("pop err:%v, %v", ele, err)
		t.FailNow()
	}
	time.Sleep(time.Second*1)
	if q.IsClose() != true {
		t.Errorf("do not stop")
		t.FailNow()
	}
	ele, err = q.Pop()
	if ele != nil || err != queue.QUEUE_CLOSED_ERROR {
			t.Errorf("pop err: %s, %v", ele, err)
			t.FailNow()
	}
}
