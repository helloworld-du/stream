package reader


type IReader interface {
	Run() error
	Stop() error
	isStop() bool
}
