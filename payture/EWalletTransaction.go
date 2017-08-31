package payture

import (
	"fmt"
	"net/http"
)

type EWalletTransaction struct {
	Cust         Customer
	Card         CardEwallet
	ConfirmCode  string
	CustomFields string
	IP           string
	Order        Payment
	SessionType  string
}

type EwalletInit struct {
	SessionType string
	Cust        Customer
	IP          string
	CardId      string
	Order       Payment
	TemplateTag string
	Language    string
}

type CardEwallet struct {
	CardHolder string
	EMonth     string
	EYear      string
	SecureCode string
	CardNumber string
	CardId     string
}

func (card CardEwallet) plain() string {
	return fmt.Sprintf("CardHolder=%s;EMonth=%s;EYear=%s;SecureCode=%s;CardNumber=%s;", card.CardHolder, card.EMonth, card.EYear, card.SecureCode, card.CardNumber)
}

type Customer struct {
	VWUserLgn   string
	VWUserPsw   string
	PhoneNumber string
	Email       string
}

func (customer Customer) plain() string {
	return fmt.Sprintf("VWUserLgn=%s;VWUserPsw=%s;PhoneNumber=%s;Email=%s;", customer.VWUserLgn, customer.VWUserPsw, customer.PhoneNumber, customer.Email)
}

func (merch Merchant) AddCardToCustomer(cust Customer, card CardEwallet) (*http.Response, error) {
	return addCard(card, cust, merch)
}

func (merch Merchant) AddCardBySession(session string) (*http.Response, error) {
	url := merch.Host + "/vwapi/Add"
	params := make(map[string][]string)
	params["SessionId"] = []string{session}
	return sendRequest(url, params)
}

func (card CardEwallet) Add(cust Customer, merch Merchant) (*http.Response, error) {
	return addCard(card, cust, merch)
}

func (cust Customer) Add(card CardEwallet, merch Merchant) (*http.Response, error) {
	return addCard(card, cust, merch)
}

func addCard(card CardEwallet, cust Customer, merch Merchant) (*http.Response, error) {
	url := merch.Host + "/vwapi/Add"
	params := make(map[string][]string)
	params["VWID"] = []string{merch.Key}
	params["DATA"] = []string{cust.plain() + card.plain()}
	return sendRequest(url, params)
}

func (card CardEwallet) Activate(cust Customer, merch Merchant, amount string) (*http.Response, error) {
	url := merch.Host + "/vwapi/Activate"
	params := make(map[string][]string)
	params["VWID"] = []string{merch.Key}
	params["DATA"] = []string{cust.plain() + fmt.Sprintf("CardId=%s;Amount=%s;", card.CardId, amount)}
	return sendRequest(url, params)
}

func (card CardEwallet) Remove(cust Customer, merch Merchant) (*http.Response, error) {
	url := merch.Host + "/vwapi/Remove"
	params := make(map[string][]string)
	params["VWID"] = []string{merch.Key}
	params["DATA"] = []string{cust.plain() + fmt.Sprintf("CardId=%s;", card.CardId)}
	return sendRequest(url, params)
}

/*
Customer
*/

func (cust Customer) GetCardList(merch Merchant) (*http.Response, error) {
	url := merch.Host + "/vwapi/GetList"
	params := make(map[string][]string)
	params["VWID"] = []string{merch.Key}
	params["DATA"] = []string{cust.plain()}
	return sendRequest(url, params)
}

func (customer Customer) RegisterCustomer(merch Merchant) (*http.Response, error) {
	return custRequest("Register", customer, merch)
}
func (customer Customer) DeleteCustomer(merch Merchant) (*http.Response, error) {
	return custRequest("Delete", customer, merch)
}

func (customer Customer) UpdateCustomer(merch Merchant) (*http.Response, error) {
	return custRequest("Update", customer, merch)
}
func (customer Customer) CheckCustomer(merch Merchant) (*http.Response, error) {
	return custRequest("Check", customer, merch)
}

func custRequest(cmd string, cust Customer, merch Merchant) (*http.Response, error) {
	url := merch.Host + "/vwapi/" + cmd
	params := make(map[string][]string)
	params["VWID"] = []string{merch.Key}
	params["DATA"] = []string{cust.plain()}
	return sendRequest(url, params)
}

func (pay EWalletTransaction) PayNoRegCard(merch Merchant) (*http.Response, error) {
	url := merch.Host + "/vwapi/Pay"
	params := make(map[string][]string)
	params["VWID"] = []string{merch.Key}
	params["DATA"] = []string{pay.Cust.plain() + pay.Card.plain() + pay.Order.plain() + fmt.Sprintf("ConfirmCode=%s;IP=%s;CustomFields=%s;SesstionType=%s;CardId=%s;", pay.ConfirmCode, pay.IP, pay.CustomFields, pay.SessionType, pay.Card.CardId)}
	return sendRequest(url, params)
}

func (pay EWalletTransaction) PayRegCard(merch Merchant) (*http.Response, error) {
	url := merch.Host + "/vwapi/Pay"
	params := make(map[string][]string)
	params["VWID"] = []string{merch.Key}
	params["DATA"] = []string{pay.Cust.plain() + pay.Order.plain() + fmt.Sprintf("SecureCode=%s;CardId=%s;SessionType=%s;ConfirmCode=%s;IP=%s;CustomFields=%s;", pay.Card.SecureCode, pay.Card.CardId,
		pay.SessionType, pay.ConfirmCode, pay.IP, pay.CustomFields)}
	return sendRequest(url, params)
}

func (merch Merchant) PayBySession(session string) (*http.Response, error) {
	url := merch.Host + "/vwapi/Pay"
	params := make(map[string][]string)
	params["SessionId"] = []string{session}
	return sendRequest(url, params)
}

func (pay EWalletTransaction) SendCode(merch Merchant) (*http.Response, error) {
	url := merch.Host + "/vwapi/SendCode"
	params := make(map[string][]string)
	params["VWID"] = []string{merch.Key}
	params["DATA"] = []string{pay.Cust.plain() + pay.Order.plain() + fmt.Sprintf("CardId=%s;", pay.Card.CardId)}
	return sendRequest(url, params)
}

func (init EwalletInit) Init(merch Merchant) (*http.Response, error) {
	url := merch.Host + "/vwapi/Init"
	params := make(map[string][]string)
	params["VWID"] = []string{merch.Key}
	params["DATA"] = []string{init.Cust.plain() + init.Order.plain() + fmt.Sprintf("CardId=%s;SessionType=%s;Language=%s;IP=%s;TemplateTag=%s;", init.CardId, init.SessionType,
		init.Language, init.IP, init.TemplateTag)}
	return sendRequest(url, params)
}

func (merch Merchant) InitEWallet(cust Customer, sessionType string, ip string, order Payment, cardId string, template string, lang string) (*http.Response, error) {
	url := merch.Host + "/vwapi/Init"
	params := make(map[string][]string)
	params["VWID"] = []string{merch.Key}
	params["DATA"] = []string{cust.plain() + order.plain() + fmt.Sprintf("CardId=%s;SessionType=%s;Language=%s;IP=%s;TemplateTag=%s;", cardId, sessionType,
		lang, ip, template)}
	return sendRequest(url, params)
}
