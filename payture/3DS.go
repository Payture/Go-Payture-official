package payture

import (
	"net/http"
)

//APIThreeDSPAY for completing the one-stage charging on a card with 3-D Secure.
func (this APIManager) APIThreeDSPAY(orderId string, paRes string) (*http.Response, error) {
	return this.send3DS(PAY3DS, orderId, paRes)
}

//APIThreeDSPAY for completing blocking stage in two-stage charging on a card with 3-D Secure.
func (this APIManager) APIThreeDSBlock(orderId string, paRes string) (*http.Response, error) {
	return this.send3DS(BLOCK3DS, orderId, paRes)
}

func (this APIManager) send3DS(method string, orderId string, paRes string) (*http.Response, error) {
	prms := reqPrms{}.set(KEY, this.merchant.Key).set(ORDERID, orderId).set(PARES, paRes).get()
	return sendRequest(this, method, prms)
}

//APIThreeDSPAY for completing payment on a card with 3-D Secure in EWallet Service.
func (this EwalletManager) PaySubmit3DS(md string, paRes string) (*http.Response, error) {
	prms := reqPrms{}.set(MD, md).set(PARES, paRes).get()
	return sendRequest(this, PAYSUBMIT3DS, prms)
}
