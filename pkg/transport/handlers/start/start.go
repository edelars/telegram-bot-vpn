package start

import (
	"backend-vpn/pkg/transport"
	"fmt"
)

type Start struct {
}

func NewStart() *Start {
	return &Start{}
}

func (p *Start) Endpoint() interface{} {
	return `/start`
}

func (p *Start) Handler() func(data transport.HandlerData) interface{} {
	return func(data transport.HandlerData) interface{} {
		return fmt.Sprintf("Доступны команды:\n\n/price - Стоимость VPN\n/info - Информация об аакаунте VPN\n/help - Инструкции по настройке")
	}
}

func (p *Start) Menu() []transport.MenuI {
	return nil
}
