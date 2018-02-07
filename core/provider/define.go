package provider

type IProvider interface {
	Read() (msg interface{}, err error, hasNext bool)
}
