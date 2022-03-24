package transport

type Transport interface {
	//	Listen()
	Send() error
}
