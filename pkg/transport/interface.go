package transport

type HandlerI interface {
	Endpoint() interface{}
	Handler() func(data HandlerData) interface{}
	Menu() []MenuI
}

type MenuI interface {
	Data() (text, unique string)
	Handler() func(data HandlerData) interface{}
}

type HandlerData struct {
	Username string
}
