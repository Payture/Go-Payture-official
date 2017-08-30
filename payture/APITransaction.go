package main

import (
	"fmt"
	"net/http"
)

type PayAPITransaction struct {
	Key, CustomerKey string
	CustomerFields   CustParams
	PaytureId        string
	PayInfo          PayInfo
	order            Payment
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
		"OrderId":      {payTr.order.OrderId},
		"Amount":       {payTr.order.Amount},
		"CustomerKey":  {payTr.CustomerKey},
		"PayInfo":      {payTr.PayInfo.plain()},
		"CustomFields": {payTr.CustomerFields.plain()}}
}

func (pay PayAPITransaction) pay(merch Merchant) (*http.Response, error) {
	var url = merch.Host + "/api/Pay"
	return sendRequestFormer(url, pay)
}

func (pay PayAPITransaction) block(merch Merchant) (*http.Response, error) {
	var url = merch.Host + "/api/Block"
	return sendRequestFormer(url, pay)
}
