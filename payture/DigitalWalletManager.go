package payture

//DigitalWalletManager type provide the access for calling API Services for digital payments: Apple, Android and MasterPass.
type DigitalWalletManager struct {
	OrderManager
	digitalType string
}

func (this DigitalWalletManager) getReqUrl(cmd string) string {
	return this.merchant.Host + "/" + this.apiType + "/" + cmd
}

//AppleService returns DigitalWalletManager for Apple system.
func AppleService(merch Merchant) (digital DigitalWalletManager) {
	digital.apiType = "api"
	digital.merchant = merch
	digital.digitalType = APPLE
	return
}

//AndroidService returns DigitalWalletManager for Android system.
func AndroidService(merch Merchant) (digital DigitalWalletManager) {
	digital.apiType = "api"
	digital.merchant = merch
	digital.digitalType = ANDROID
	return
}

//MasterPassService returns DigitalWalletManager for MasterPass system.
func MasterPassService(merch Merchant) (digital DigitalWalletManager) {
	digital.apiType = "api"
	digital.merchant = merch
	digital.digitalType = MASTERPASS
	return
}

//Pay makes one-stage charging of funds.
func (this DigitalWalletManager) Pay(order Payment, token string, secureCode string) (dRes DigitalResponse, err error) {
	dType := this.digitalType
	switch {
	case dType == ANDROID:
		dRes, err = this.sendAndroidApple("PAY", order, token)
	case dType == APPLE:
		dRes, err = this.sendAndroidApple("PAY", order, token)
	case dType == MASTERPASS:
		dRes, err = this.sendMP(MPPAY, order, secureCode, token)
	}
	return
}

//Block perform bloking funds in two-stage charging.
func (this DigitalWalletManager) Block(order Payment, token string, secureCode string) (dRes DigitalResponse, err error) {
	dType := this.digitalType
	switch {
	case dType == ANDROID:
		dRes, err = this.sendAndroidApple("BLOCK", order, token)
	case dType == APPLE:
		dRes, err = this.sendAndroidApple("BLOCK", order, token)
	case dType == MASTERPASS:
		dRes, err = this.sendMP(MPBLOCK, order, secureCode, token)
	}
	return
}

func (this DigitalWalletManager) sendAndroidApple(method string, order Payment, token string) (dResp DigitalResponse, err error) {
	prms := reqPrms{}.set(KEY, this.merchant.Key).set(ORDERID, order.OrderId).set(METHOD, method).set(PAYTOKEN, token)
	cmd := "ApplePay"
	if this.digitalType == ANDROID {
		prms.set(AMOUNT, order.Amount)
		cmd = "AndroidPay"
	}
	err = sendRequestAndMap(&dResp, prms.get(), this, cmd)
	return
}

func (this DigitalWalletManager) sendMP(cmd string, order Payment, secureCode string, token string) (dResp DigitalResponse, err error) {
	prms := reqPrms{}.set(KEY, this.merchant.Key).set(ORDERID, order.OrderId).set(AMOUNT, order.Amount).set(CVC2, secureCode).set(TOKEN, token).get()
	err = sendRequestAndMap(&dResp, prms, this, cmd)
	return
}
