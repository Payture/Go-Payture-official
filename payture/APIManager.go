package payture

import (
	"fmt"
	"net/http"
)

type APIManager struct {
	Merchant Merchant
}

func (this APIManager) getAPI() string {
	return "api"
}

type PayAPITransaction struct {
	Key, CustomerKey string
	CustomerFields   CustParams
	PaytureId        string
	PaymentInfo      PayInfo
	Order            Payment
}

type PayInfo struct {
	card  Card
	order Payment
	PAN   string
}

func (info PayInfo) plain() string {
	return fmt.Sprintf("PAN=%s;EMonth=%s;EYear=%s;CardHolder=%s;SecureCode=%s;OrderId=%s;Amount=%s", info.PAN, info.card.EMonth,
		info.card.EYear, info.card.CardHolder, info.card.SecureCode, info.order.OrderId, info.order.Amount)
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

/*//////////////
Payments command
*/
func (this APIManager) Unblock(order Payment) (*http.Response, error) {
	return order.Unblock(this.getAPI(), this.Merchant)
}

func (this APIManager) Refund(order Payment) (*http.Response, error) {
	return order.Refund(this.getAPI(), this.Merchant)
}

func (this APIManager) Charge(order Payment) (*http.Response, error) {
	return order.Charge(this.getAPI(), this.Merchant)
}

func (this APIManager) GetState(order Payment) (*http.Response, error) {
	return order.GetState(this.Merchant)
}

/*
Payments command
*/ //////////////
func (this APIManager) Pay(order Payment, info PayInfo, additional CustParams, custKey string, paytureId string) (*http.Response, error) {
	return this.sendApi("Pay", order, info, additional, custKey, paytureId)
}

func (this APIManager) Block(order Payment, info PayInfo, additional CustParams, custKey string, paytureId string) (*http.Response, error) {
	return this.sendApi("Block", order, info, additional, custKey, paytureId)
}

func (this APIManager) sendApi(cmd string, order Payment, info PayInfo, additional CustParams, custKey string, paytureId string) (*http.Response, error) {
	var url = this.Merchant.Host + "/" + this.getAPI() + "/" + cmd
	var pay = PayAPITransaction{Order: order, PaymentInfo: info, CustomerFields: additional, CustomerKey: custKey, Key: this.Merchant.Key, PaytureId: paytureId}
	return sendRequestFormer(url, pay)
}
