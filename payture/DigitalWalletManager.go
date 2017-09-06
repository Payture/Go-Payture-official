package payture

const (
	APPLE      string = "apple"
	MASTERPASS string = "mspass"
	ANDROID    string = "android"
)

type DigitalWalletManager struct {
	OrderManager
	digitalType string
}

func (this DigitalWalletManager) getReqUrl(cmd string) string {
	return this.merchant.Host + "/" + this.apiType + "/" + cmd
}

func AppleService(merch Merchant) (digital DigitalWalletManager) {
	digital.apiType = "api"
	digital.merchant = merch
	digital.digitalType = APPLE
	return
}

func AndroidService(merch Merchant) (digital DigitalWalletManager) {
	digital.apiType = "api"
	digital.merchant = merch
	digital.digitalType = ANDROID
	return
}

func MasterPassService(merch Merchant) (digital DigitalWalletManager) {
	digital.apiType = "api"
	digital.merchant = merch
	digital.digitalType = MASTERPASS
	return
}

func (this DigitalWalletManager) Pay(order Payment, token string, secureCode string) (dRes DigitalResponse, err error) {
	dType := this.digitalType
	switch {
	case dType == ANDROID:
		dRes, err = this.sendAndroidApple("PAY", order, token)
	case dType == APPLE:
		dRes, err = this.sendAndroidApple("PAY", order, token)
	case dType == MASTERPASS:
		dRes, err = this.sendMP("MPPay", order, secureCode, token)
	}
	return
}

func (this DigitalWalletManager) Block(order Payment, token string, secureCode string) (dRes DigitalResponse, err error) {
	dType := this.digitalType
	switch {
	case dType == ANDROID:
		dRes, err = this.sendAndroidApple("BLOCK", order, token)
	case dType == APPLE:
		dRes, err = this.sendAndroidApple("BLOCK", order, token)
	case dType == MASTERPASS:
		dRes, err = this.sendMP("MPBlock", order, secureCode, token)
	}
	return
}

func (this DigitalWalletManager) sendAndroidApple(method string, order Payment, token string) (dResp DigitalResponse, err error) {
	params := map[string][]string{
		"Key":      []string{this.merchant.Key},
		"OrderId":  []string{order.OrderId},
		"Method":   []string{method},
		"PayToken": []string{token}}
	cmd := "ApplePay"
	if this.digitalType == ANDROID {
		params["Amount"] = []string{order.Amount}
		cmd = "AndroidPay"
	}
	err = sendRequestAndMap(&dResp, params, this, cmd)
	return
}

func (this DigitalWalletManager) sendMP(cmd string, order Payment, secureCode string, token string) (dResp DigitalResponse, err error) {
	params := map[string][]string{
		"Key":     []string{this.merchant.Key},
		"OrderId": []string{order.OrderId},
		"Amount":  []string{order.Amount},
		"CVC2":    []string{secureCode},
		"Token":   []string{token}}
	err = sendRequestAndMap(&dResp, params, this, cmd)
	return
}
