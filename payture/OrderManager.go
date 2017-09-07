package payture

import "fmt"

//OrderManager type for management of payment's state.
type OrderManager struct {
	apiType  string
	merchant Merchant
}

//GetOrderManager returned the OrderManager for specified api type and merchant.
func GetOrderManager(api string, merch Merchant) (ordM OrderManager) {
	ordM.apiType = api
	ordM.merchant = merch
	return
}

func (this OrderManager) getReqUrl(cmd string) string {
	return this.merchant.Host + "/" + this.apiType + "/" + cmd
}

//Unblock unlocking the funds on the card.
func (this OrderManager) Unblock(order Payment) (orderResp OrderResponse, err error) {
	key := KEY
	if this.apiType == EWALLET {
		key = VWID
	}
	prms := reqPrms{}.set(key, this.merchant.Key).set(ORDERID, order.OrderId).set(AMOUNT, order.Amount)

	if this.apiType == EWALLET || this.apiType == INPAY {
		prms.setPass(this.merchant)
	}
	err = sendRequestAndMap(&orderResp, prms.get(), this, UNBLOCK)
	return
}

//Refund the funds to customer.
func (this OrderManager) Refund(order Payment) (orderResp OrderResponse, err error) {
	prms := reqPrms{}
	if this.apiType == EWALLET {
		prms.set(VWID, this.merchant.Key).set(DATAUP, fmt.Sprintf("OrderId=%s;Password=%s;Amount=%s", order.OrderId, this.merchant.Password, order.Amount))
	} else {
		prms.setKey(this.merchant).set(ORDERID, order.OrderId).set(AMOUNT, order.Amount).setPass(this.merchant)
	}
	err = sendRequestAndMap(&orderResp, prms.get(), this, REFUND)
	return
}

//Refund the funds that was block from the customer's card.
func (this OrderManager) Charge(order Payment) (orderResp OrderResponse, err error) {
	key := KEY
	if this.apiType == EWALLET {
		key = VWID
	}
	prms := reqPrms{}.set(key, this.merchant.Key).set(ORDERID, order.OrderId)
	if this.apiType == EWALLET || this.apiType == INPAY {
		prms.setPass(this.merchant).set(AMOUNT, order.Amount)
	}
	err = sendRequestAndMap(&orderResp, prms.get(), this, CHARGE)
	return
}

//PayStatus returns the transaction's status for Ewallet and InPay services.
func (this OrderManager) PayStatus(order Payment) (orderResp OrderResponse, err error) {
	prms := reqPrms{}.setKey(this.merchant).set(ORDERID, order.OrderId)

	if this.apiType == EWALLET {
		prms = reqPrms{}.set(VWID, this.merchant.Key).set(DATAUP, fmt.Sprintf("OrderId=%s", order.OrderId))
	}
	err = sendRequestAndMap(&orderResp, prms.get(), this, PAYSTATUS)
	return
}

//GetState returns the transaction's status for API service.
func (this OrderManager) GetState(order Payment) (orderResp OrderResponse, err error) {
	prms := reqPrms{}.setKey(this.merchant).set(ORDERID, order.OrderId).get()
	err = sendRequestAndMap(&orderResp, prms, this, GETSTATE)
	return
}
