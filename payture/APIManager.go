package payture

import (
	"encoding/xml"
	"fmt"
)

type APIManager struct {
	OrderManager
}

func APIService(merch Merchant) (api APIManager) {
	api.apiType = "api"
	api.merchant = merch
	return
}

type PayAPITransaction struct {
	Key, CustomerKey string
	CustomerFields   CustParams
	PaytureId        string
	PaymentInfo      PayInfo
	Order            Payment
}

type PayInfo struct {
	Card  Card
	Order Payment
	PAN   string
}

func (info PayInfo) plain() string {

	return fmt.Sprintf("PAN=%s;EMonth=%s;EYear=%s;CardHolder=%s;SecureCode=%s;OrderId=%s;Amount=%s", info.PAN, info.Card.EMonth,
		info.Card.EYear, info.Card.CardHolder, info.Card.SecureCode, info.Order.OrderId, info.Order.Amount)
}

func (payTr PayAPITransaction) content() map[string][]string {
	return map[string][]string{
		"Key":          {payTr.Key},
		"OrderId":      {payTr.Order.OrderId},
		"Amount":       {payTr.Order.Amount},
		"CustomerKey":  {payTr.CustomerKey},
		"PayInfo":      {payTr.PaymentInfo.plain()},
		"CustomFields": {payTr.CustomerFields.plain()}}
}

func (this APIManager) Pay(order Payment, info PayInfo, additional CustParams, custKey string, paytureId string) (APIResponses, error) {
	return this.sendApi("Pay", order, info, additional, custKey, paytureId)
}

func (this APIManager) Block(order Payment, info PayInfo, additional CustParams, custKey string, paytureId string) (APIResponses, error) {
	return this.sendApi("Block", order, info, additional, custKey, paytureId)
}

func (this APIManager) sendApi(cmd string, order Payment, info PayInfo, additional CustParams, custKey string, paytureId string) (apiResponse APIResponses, err error) {
	var pay = PayAPITransaction{Order: order, PaymentInfo: info, CustomerFields: additional, CustomerKey: custKey, Key: this.merchant.Key, PaytureId: paytureId}
	err = sendRequestFormerMap(&apiResponse, this, cmd, pay)
	return
}

/*////////////
API Responses
*/
type APIResponses struct {
	AddInfo   []AdditioanalInfo `xml:"AddInfo"`
	Success   bool              `xml:"Success,attr"`
	ErrorCode string            `xml:"ErrCode,attr"`
	Amount    int64             `xml:"Amount,attr"`
	OrderId   string            `xml:"OrderId,attr"`
	State     string            `xml:"State,attr"`
	NewAmount string            `xml:"NewAmount,attr"`
	Key       string            `xml:"Key,attr"`
}

type AdditioanalInfo struct {
	Key   string `xml:"Key,attr"`
	Value string `xml:"Value,attr"`
}

/*
API Responses
*/ ///////////

func (resp *APIResponses) Unwrap(byteBody []byte) error {
	xml.Unmarshal(byteBody, &resp)
	return nil
}
