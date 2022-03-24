package transport

type HandlerI interface {
	Endpoint() interface{}
	Handler() func() interface{}
	Menu() []MenuI
}

type MenuI interface {
	Data() (text, unique string)
	Handler() func() interface{}
}
