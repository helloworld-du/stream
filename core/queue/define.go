package queue



type Queue interface {
	Pop () (v interface{}, err error)
	TryPop() (interface{}, error)
	Push(v interface{}) error
	TryPush(v interface{}) error
	Len() int
	Cap() int
	IsClose()	bool
	Close()	error
}


