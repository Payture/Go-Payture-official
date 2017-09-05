package payture

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
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

/*//////////////
Payments command
*/
/*
func (this APIManager) Unblock(order Payment) (*http.Response, error) {
	return order.Unblock(this.getAPI(), this.Merchant)
}*/

func (this APIManager) Unblock(order Payment) (OrderResponse, error) {
	return order.Unblock(this.getAPI(), this.Merchant)
}

func (this APIManager) Refund(order Payment) (OrderResponse, error) {
	return order.Refund(this.getAPI(), this.Merchant)
}

func (this APIManager) Charge(order Payment) (OrderResponse, error) {
	return order.Charge(this.getAPI(), this.Merchant)
}

func (this APIManager) GetState(order Payment) (OrderResponse, error) {
	return order.GetState(this.Merchant)
}

/*
Payments command
*/ //////////////
func (this APIManager) Pay(order Payment, info PayInfo, additional CustParams, custKey string, paytureId string) (APIResponses, error) {
	apiResponse := APIResponses{}
	httpResp, err := this.sendApi("Pay", order, info, additional, custKey, paytureId)
	if err != nil {
		return apiResponse, err
	}
	body, err2 := BodyByte(httpResp)
	if err2 != nil {
		return apiResponse, err2
	}
	apiResponse.ParseAPIByte(body)
	return apiResponse, nil
}

func (this APIManager) Block(order Payment, info PayInfo, additional CustParams, custKey string, paytureId string) (APIResponses, error) {
	apiResponse := APIResponses{}
	httpResp, err := this.sendApi("Block", order, info, additional, custKey, paytureId)
	if err != nil {
		return apiResponse, err
	}
	body, err2 := BodyByte(httpResp)
	if err2 != nil {
		return apiResponse, err2
	}
	apiResponse.ParseAPIByte(body)
	return apiResponse, nil
}

func (this APIManager) sendApi(cmd string, order Payment, info PayInfo, additional CustParams, custKey string, paytureId string) (*http.Response, error) {
	var url = this.Merchant.Host + "/" + this.getAPI() + "/" + cmd
	var pay = PayAPITransaction{Order: order, PaymentInfo: info, CustomerFields: additional, CustomerKey: custKey, Key: this.Merchant.Key, PaytureId: paytureId}
	return sendRequestFormer(url, pay)
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

func (response *APIResponses) ParseAPI(resp *http.Response) {
	//defer resp.Body.Close()
	//var card GetList
	body, err := ioutil.ReadAll(resp.Body)
	xml.Unmarshal(body, &response)
	fmt.Println(response)
	fmt.Println(err)
}

func (response *APIResponses) ParseAPIByte(respBody []byte) {
	xml.Unmarshal(respBody, &response)
	fmt.Println(response)
}
