package qiwi

import (
	"backend-vpn/internal/dto"
	"backend-vpn/pkg/config"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	urlApi   = "https://api.qiwi.com/partner/bill/v1/bills/"
	currency = "RUB"
)

type qiwi struct {
	env config.Environment
}

func NewQiwi(env config.Environment) *qiwi {
	return &qiwi{env: env}
}

func (k qiwi) GenerateUrlToPay(pd dto.PayData) (err error, url string) {

	if pd.Order == "" {
		return errors.New("order not set"), url
	}
	if pd.Value <= 0 {
		return errors.New("value too small"), url
	}

	urls := fmt.Sprintf(urlApi+"%s", pd.Order)

	client := &http.Client{}

	// marshal User to json
	jsonData, err := json.Marshal(k.NewQiwiRequest(pd))
	if err != nil {
		return
	}

	// set the HTTP method, url, and request body
	req, err := http.NewRequest(http.MethodPut, urls, bytes.NewBuffer(jsonData))
	if err != nil {
		return
	}

	// set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", k.env.QiwiSKey))
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	respJson := &qiwiResponse{}
	err = json.NewDecoder(resp.Body).Decode(respJson)
	if err != nil {
		return
	}

	if respJson.ErrorCode != "" {
		return errors.New(respJson.ErrorCode), ""
	}

	return err, respJson.PayUrl
}

func (k qiwi) SignVerify(pd dto.PayData, sign, status, qiwiSiteId string) bool {
	//invoice_parameters = {amount.currency}|{amount.value}|{billId}|{siteId}|{status.value}

	str := fmt.Sprintf("RUB|%d|%s|%s|%s", pd.Value, pd.Order, qiwiSiteId, status)

	h := hmac.New(sha256.New, []byte(k.env.QiwiSKey))
	h.Write([]byte(str))
	sha := hex.EncodeToString(h.Sum(nil))

	return sha == sign
	//hash = HMAС(SHA256, invoice_parameters, secret_key)
}

type QiwiRequest struct {
	QiwiRequestAmount  `json:"amount"`
	Comment            string `json:"comment"`
	ExpirationDateTime string `json:"expirationDateTime"`
	QiwiCustomFileds   `json:"customFields"`
}

type QiwiRequestAmount struct {
	Currency string  `json:"currency"`
	Value    float32 `json:"value"`
}

type QiwiCustomFileds struct {
	ThemeCode string `json:"currency"`
}

func (k qiwi) NewQiwiRequest(pd dto.PayData) *QiwiRequest {
	return &QiwiRequest{
		QiwiRequestAmount: QiwiRequestAmount{
			Currency: currency,
			Value:    float32(pd.Value),
		},
		ExpirationDateTime: time.Now().Add(1 * time.Hour).Format("2006-01-02T15:04:05Z07:00"),
		Comment:            pd.Order,
		QiwiCustomFileds: QiwiCustomFileds{
			ThemeCode: k.env.QiwiThemeCode,
		},
	}
}

type qiwiResponse struct {
	PayUrl      string `json:"payUrl"`
	ErrorCode   string `json:"errorCode"`
	Description string `json:"description"`
	BillId      string `json:"billId"`
}
