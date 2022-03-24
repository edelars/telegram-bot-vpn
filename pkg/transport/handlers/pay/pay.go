package pay

type Pay struct {
}

func NewPay() *Pay {
	return &Pay{}
}
func (p *Pay) Data() (text, unique string) {
	return "Оплатить VPN", "pay"
}
func (p *Pay) Handler() func() interface{} {
	return nil
}
