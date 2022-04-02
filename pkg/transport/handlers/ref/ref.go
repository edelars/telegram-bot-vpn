package ref

import (
	"backend-vpn/pkg/transport"
	"fmt"
)

type Ref struct {
}

func NewRef() *Ref {
	return &Ref{}
}

func (p *Ref) Endpoint() interface{} {
	return `/ref`
}

func (p *Ref) Handler() func(data transport.HandlerData) interface{} {
	return func(data transport.HandlerData) interface{} {
		return fmt.Sprint("")
	} //TODO referral
}

func (p *Ref) Menu() []transport.MenuI {
	return nil
}
