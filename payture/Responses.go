package payture

import (
	"encoding/xml"
)

type EWalletPayResponse struct {
	OrderId string `xml:"OrderId,attr"`
	Amount  int64  `xml:"Amount,attr"`
	CustomerResponse
}

type CustomerResponse struct {
	Success   bool   `xml:"Success,attr"`
	ErrorCode string `xml:"ErrCode,attr"`
	VWUserLgn string `xml:"VWUserLgn,attr"`
}

type CardResponse struct {
	CustomerResponse
	CardId   string `xml:"CardId,attr"`
	CardName string `xml:"CardName,attr"`
}

type CardListResponse struct {
	CustomerResponse
	Cards []Item `xml:"Item"`
}

type Item struct {
	CardName   string `xml:"CardName,attr"`
	CardId     string `xml:"CardId,attr"`
	CardHolder string `xml:"CardHolder,attr"`
	Status     string `xml:"Status,attr"`
	NoCVV      bool   `xml:"NoCVV,attr"`
	Expired    bool   `xml:"Expired,attr"`
}

func (resp *CustomerResponse) Unwrap(byteBody []byte) error {
	xml.Unmarshal(byteBody, &resp)
	return nil
}

func (resp *CardResponse) Unwrap(byteBody []byte) error {
	xml.Unmarshal(byteBody, &resp)
	return nil
}

func (resp *EWalletPayResponse) Unwrap(byteBody []byte) error {
	xml.Unmarshal(byteBody, &resp)
	return nil
}

/*
Digital Responses
*/

type DigitalResponse struct {
	Success   string `xml:"Success,attr"`
	ErrorCode string `xml:"ErrCode,attr"`
	OrderId   string `xml:"OrderId,attr"`
	Amount    int64  `xml:"Amount,attr"`
	Key       string `xml:"Key,attr"`
}

func (resp *DigitalResponse) Unwrap(byteBody []byte) error {
	xml.Unmarshal(byteBody, &resp)
	return nil
}

type OrderResponse struct {
	Success   bool   `xml:"Success,attr"`
	OrderId   string `xml:"OrderId,attr"`
	Amount    int64  `xml:"Amount,attr"`
	ErrCode   string `xml:"ErrCode,attr"`
	NewAmount int64  `xml:"NewAmount,attr"`
}

func (resp *OrderResponse) Unwrap(byteBody []byte) error {
	xml.Unmarshal(byteBody, &resp)
	return nil
}

/*
Init  Response
*/

type InitResponse struct {
	SessionId string `xml:"SessionId,attr"`
	Success   string `xml:"Success,attr"`
	Amount    int64  `xml:"Amount,attr"`
	ErrorCode string `xml:"ErrCode,attr"`
	OrderId   string `xml:"OrderId,attr"`
}

func (resp *InitResponse) Unwrap(byteBody []byte) error {
	xml.Unmarshal(byteBody, &resp)
	return nil
}

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

func (resp *APIResponses) Unwrap(byteBody []byte) error {
	xml.Unmarshal(byteBody, &resp)
	return nil
}
