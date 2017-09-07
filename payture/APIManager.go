package payture

import (
	"fmt"
)

//APIManager type provide the access for calling API Service  in Payture system.
type APIManager struct {
	OrderManager
}

//APIService method returned the APIManager.
func APIService(merch Merchant) (api APIManager) {
	api.apiType = "api"
	api.merchant = merch
	return
}

//PayAPITransaction collection of fields that you must supplied for make pay or block operation.
type PayAPITransaction struct {
	Key, CustomerKey string
	CustomerFields   CustParams
	PaytureId        string
	PaymentInfo      PayInfo
	Order            Payment
}

//PayInfo detailed information about transaction.
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
	return reqPrms{}.set(KEY, payTr.Key).
		set(ORDERID, payTr.Order.OrderId).
		set(AMOUNT, payTr.Order.Amount).
		set(CUSTKEY, payTr.CustomerKey).
		set(PAYINFO, payTr.PaymentInfo.plain()).
		set(CUSTFIELDS, payTr.CustomerFields.plain()).get()
}

//Pay makes one stage charge of funds from card.
func (this APIManager) Pay(order Payment, info PayInfo, additional CustParams, custKey string, paytureId string) (APIResponses, error) {
	return this.sendApi(PAY, order, info, additional, custKey, paytureId)
}

//Block blocking funds on card. For two-stage charge of funds from card.s
func (this APIManager) Block(order Payment, info PayInfo, additional CustParams, custKey string, paytureId string) (APIResponses, error) {
	return this.sendApi(BLOCK, order, info, additional, custKey, paytureId)
}

func (this APIManager) sendApi(cmd string, order Payment, info PayInfo, additional CustParams, custKey string, paytureId string) (apiResponse APIResponses, err error) {
	var pay = PayAPITransaction{Order: order, PaymentInfo: info, CustomerFields: additional, CustomerKey: custKey, Key: this.merchant.Key, PaytureId: paytureId}
	err = sendRequestFormerMap(&apiResponse, this, cmd, pay)
	return
}
