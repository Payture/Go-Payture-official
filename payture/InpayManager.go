package payture

import (
	"fmt"
)

//InpayManager for payment management in InPay Service.
type InpayManager struct {
	OrderManager
}

//InPayService method returns the InpayManager.
func InPayService(merch Merchant) (inpay InpayManager) {
	inpay.merchant = merch
	inpay.apiType = "apim"
	return
}

//Init create and returned Session in Payture system by which the customer can make a payment.
func (this InpayManager) Init(order Payment, sessionType string, tag string, lang string, ip string, urlReturn string, additionalFields CustParams) (initResp InitResponse, err error) {
	prms := reqPrms{}.setKey(this.merchant).set(DATA, order.plain()+fmt.Sprintf("SessionType=%s;Language=%s;IP=%s;TemplateTag=%s;Url=%s", sessionType, lang, ip, tag, urlReturn)+additionalFields.plain()).get()
	err = sendRequestAndMap(&initResp, prms, this, INIT)
	return
}
