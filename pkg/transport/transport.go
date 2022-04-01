package transport

type Transport interface {
	Send(tgUserId int64, message string) error
}
