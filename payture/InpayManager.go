package payture

import (
	"fmt"
)

type InpayManager struct {
	OrderManager
}

func InPayService(merch Merchant) (inpay InpayManager) {
	inpay.merchant = merch
	inpay.apiType = "apim"
	return
}

func (this InpayManager) Init(order Payment, sessionType string, tag string, lang string, ip string, urlReturn string, additionalFields CustParams) (initResp InitResponse, err error) {
	params := map[string][]string{
		"Key":  []string{this.merchant.Key},
		"Data": []string{order.plain() + fmt.Sprintf("SessionType=%s;Language=%s;IP=%s;TemplateTag=%s;Url=%s", sessionType, lang, ip, tag, urlReturn) + additionalFields.plain()}}
	err = sendRequestAndMap(&initResp, params, this, "Init")
	return
}
