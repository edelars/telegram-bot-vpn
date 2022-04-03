package freekassa

import (
	"backend-vpn/internal/dto"
	"backend-vpn/pkg/config"
	"crypto/md5"
	"errors"
	"fmt"
)

type FreeKassa struct {
	env config.Environment
}

const (
	urlApi   = "https://pay.freekassa.ru/"
	currency = "RUB"
)

func NewFreeKassa(env config.Environment) *FreeKassa {
	return &FreeKassa{env: env}
}

// GenerateUrlToPay  md5($merchant_id.':'.$order_amount.':'.$secret_word.':'.$currency.':'.$order_id);строки "ID Вашего магазина:Сумма платежа:Секретное слово:Валюта платежа:Номер
func (k FreeKassa) GenerateUrlToPay(pd dto.PayData) (err error, url string) {

	if pd.Order == "" {
		return errors.New("order not set"), url
	}
	if pd.Value <= 0 {
		return errors.New("value too small"), url
	}

	signMd5 := k.md5(pd.Value, k.env.FKId, k.env.FKSecKey1, currency, pd.Order)
	url = fmt.Sprintf("%s?m=%s&oa=%d&i=&currency=%s&o=%s&pay=PAY&s=%s", urlApi, k.env.FKId, pd.Value, currency, pd.Order, signMd5)
	return err, url
}

func (k FreeKassa) md5(Value int, Id, FKSecKey1, Currency, Order string) string {
	sign := fmt.Sprintf("%s:%d:%s:%s:%s", Id, Value, FKSecKey1, Currency, Order)
	return fmt.Sprintf("%x", md5.Sum([]byte(sign)))
}

func (k FreeKassa) SignVerify(pd dto.PayData, sign string) bool {
	return sign == k.md5(pd.Value, k.env.FKId, k.env.FKSecKey1, currency, pd.Order)
}