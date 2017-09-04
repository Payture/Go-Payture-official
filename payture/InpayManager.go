package payture

import (
	"fmt"
	"net/http"
)

type InpayManager struct {
	Merchant Merchant
}

func (this InpayManager) getAPI() string {
	return "apim"
}

func (this InpayManager) Init(order Payment, sessionType string, tag string, lang string, ip string, urlReturn string, additionalFields CustParams) (*http.Response, error) {
	url := this.Merchant.Host + "/" + this.getAPI() + "/Init"
	params := make(map[string][]string)
	params["Key"] = []string{this.Merchant.Key}
	params["Data"] = []string{order.plain() + fmt.Sprintf("SessionType=%s;Language=%s;IP=%s;TemplateTag=%s;Url=%s", sessionType, lang, ip, tag, urlReturn) + additionalFields.plain()}
	return sendRequest(url, params)
}

/*//////////////
Payments command
*/
func (this InpayManager) Unblock(order Payment) (*http.Response, error) {
	return order.Unblock(this.getAPI(), this.Merchant)
}

func (this InpayManager) Refund(order Payment) (*http.Response, error) {
	return order.Refund(this.getAPI(), this.Merchant)
}

func (this InpayManager) Charge(order Payment) (*http.Response, error) {
	return order.Charge(this.getAPI(), this.Merchant)
}

func (this InpayManager) PayStatus(order Payment) (*http.Response, error) {
	return order.PayStatus(this.getAPI(), this.Merchant)
}

/*
Payments command
*/ //////////////
