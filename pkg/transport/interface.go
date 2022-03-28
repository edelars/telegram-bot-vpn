package transport

import "strings"

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
	Message  string
	Id       int64
}

func (d HandlerData) GetMessageArray() []string {
	return strings.Fields(d.Message)
}
