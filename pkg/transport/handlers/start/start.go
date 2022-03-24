package start

import "backend-vpn/pkg/transport"

type Start struct {
}

func NewStart() *Start {
	return &Start{}
}

func (p *Start) Endpoint() interface{} {
	return `/start`
}

func (p *Start) Handler() func() interface{} {
	return func() interface{} { return "blah" }
}

func (p *Start) Menu() []transport.MenuI {
	return nil
}
