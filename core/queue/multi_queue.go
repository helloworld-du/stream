package queue

import (
	"errors"
	"strings"
)

type MultilQueue struct {
	output []IQueue
}

func NewMultiQueue(qu []IQueue) *MultilQueue {
	return &MultilQueue{qu}
}


func (mq *MultilQueue) IsClose() bool {
	return mq.output[0].IsClose()
}

func (mq *MultilQueue) Close() (err error) {
	errStr := []string{}
	for _, q := range mq.output {
		err1 := q.Close()
		if err1 != nil {
			errStr = append(errStr, err1.Error())
		}

	}
	if len(errStr) > 0 {
		return errors.New(strings.Join(errStr, ";"))
	}
	return nil
}


func (mq *MultilQueue) Pop() (interface{}, error) {
	return nil, nil
}


func (mq *MultilQueue) TryPop() (interface{}, error) {
	return nil, nil

}

func (mq *MultilQueue) Push(v interface{}) (error) {
	errStr := []string{}
	for _, q := range mq.output {
		err1 := q.Push(v)
		if err1 != nil {
			errStr = append(errStr, err1.Error())
		}

	}
	if len(errStr) > 0 {
		return errors.New(strings.Join(errStr, ";"))
	}
	return nil
}

func (mq *MultilQueue) TryPush(v interface{}) (err error) {
	errStr := []string{}
	for _, q := range mq.output {
		err1 := q.TryPush(v)
		if err1 != nil {
			errStr = append(errStr, err1.Error())
		}

	}
	if len(errStr) > 0 {
		return errors.New(strings.Join(errStr, ";"))
	}
	return nil
}


func (mq *MultilQueue) Cap() int {
	return mq.output[0].Cap()
}

func (mq *MultilQueue) Len() int {
	return mq.output[0].Len()
}