package payture

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

//EwalletManager type provide the access for calling EWallet Service API in Payture system.
type EwalletManager struct {
	Customers map[string]Customer
	OrderManager
}

//EwalletService method that returns the EwalletManager.
func EwalletService(merch Merchant) (ew EwalletManager) {
	ew.apiType = "vwapi"
	ew.merchant = merch
	ew.Customers = make(map[string]Customer)
	return
}

//NotRegisteredCard type that represent information about card that you must provide for registration in Payture system.
type NotRegisteredCard struct {
	CardHolder string
	EMonth     string
	EYear      string
	SecureCode string
	CardNumber string
}

//EWalletCard type that represent information about card that you recieved then call GetCardList method.
type EWalletCard struct {
	CardName   string `xml:"CardName,attr"`
	CardId     string `xml:"CardId,attr"`
	CardHolder string `xml:"CardHolder,attr"`
	Status     string `xml:"Status,attr"`
	Expired    bool   `xml:"Expired,attr"`
	NoCVV      bool   `xml:"NoCVV,attr"`
}

//Customer type that represent information about customer.
type Customer struct {
	VWUserLgn   string
	VWUserPsw   string
	PhoneNumber string
	Email       string
	Cards       []EWalletCard `xml:"Item"`
}

func (customer Customer) plain() string {
	str := fmt.Sprintf("VWUserLgn=%s;VWUserPsw=%s;PhoneNumber=%s;Email=%s;", customer.VWUserLgn, customer.VWUserPsw, customer.PhoneNumber, customer.Email)
	fmt.Println(str)
	return str
}

func (card NotRegisteredCard) plain() string {
	return fmt.Sprintf("CardHolder=%s;EMonth=%s;EYear=%s;SecureCode=%s;CardNumber=%s;", card.CardHolder, card.EMonth, card.EYear, card.SecureCode, card.CardNumber)
}

//AddCustomerToCollection added customer to the EWalletManager's customer collection (but not in the Payture system).
func (this EwalletManager) AddCustomerToCollection(cust Customer) (EwalletManager, bool) {
	if this.Customers == nil {
		this.Customers = make(map[string]Customer)
		this.Customers[cust.VWUserLgn] = cust
		return this, true
	}
	_, exist := this.Customers[cust.VWUserLgn]
	if !exist {
		this.Customers[cust.VWUserLgn] = cust
		return this, true
	}
	return this, false
}

//RemoveCustomerFromCollection removed customer from the EWalletManager's customer collection (but not from the Payture system).
func (this EwalletManager) RemoveCustomerFromCollection(custLogin string) (EwalletManager, bool) {
	if this.Customers == nil {
		return this, false
	}
	_, exist := this.Customers[custLogin]
	if exist {
		delete(this.Customers, custLogin)
		return this, true
	}

	return this, false
}

//InitTransaction create and returned Session in Payture system by which the customer can make a payment or can registered card.
func (this EwalletManager) InitTransaction(custLogin string, sessionType string, ip string, order Payment, tag string, lang string, cardId string) (initResp InitResponse, err error) {
	cust := this.Customers[custLogin]
	err = this.sendEWRequestMap(&initResp, INIT, cust.plain()+order.plain()+fmt.Sprintf("IP=%s;TemplateTag=%s;Language=%s;CardId=%s;SessionType=%s;CustomField=Field;", ip,
		tag, lang, cardId, sessionType /*add mooooooreee*/))
	return
}

//PayRegCard makes a payment or blocks funds on customer's card. Required information about card that was attached to specified customer.
func (this EwalletManager) PayRegCard(custLogin string, cardId string, secureCode string, order Payment, ip string, confirmCode string, custFields map[string]string) (payResp EWalletPayResponse, err error) {
	cust := this.Customers[custLogin]
	err = this.sendEWRequestMap(&payResp, PAY, cust.plain()+order.plain()+fmt.Sprintf("IP=%s;ConfirmCode=%s;CardId=%s;SecureCode=%s;CustomField=Field;", ip, confirmCode, cardId, secureCode /*add mooooooreee*/))
	return
}

//PayNoRegCard makes a payment or blocks funds on customer's card. Required full information about card that wasn't attached to customer.
func (this EwalletManager) PayNoRegCard(custLogin string, card NotRegisteredCard, order Payment, ip string, confirmCode string, custFields map[string]string, addCardToCust bool) (payResp EWalletPayResponse, err error) {
	cust := this.Customers[custLogin]
	err = this.sendEWRequestMap(&payResp, PAY, cust.plain()+card.plain()+order.plain()+fmt.Sprintf("IP=%s;ConfirmCode=%s;CardId=FreePay;CustomField=Field;", ip, confirmCode /*add mooooooreee*/))
	return
}

//OperationBySession makes a payment, blocks funds on customer's card or registers card in Payture system by SessionId.
func (this EwalletManager) OperationBySession(cmd string, sessionId string) (*http.Response, error) { //For Add or Pay operation on Payture side
	prms := reqPrms{}.set(SESSIONID, sessionId).get()
	return sendRequest(this, cmd, prms)
}

//SendCode to customer in the process of payment for secure reasons.
func (this EwalletManager) SendCode(custLogin string) (*http.Response, error) {
	cust := this.Customers[custLogin]
	return this.sendEWRequest(SENDCODE, cust.plain())
}

//CardAdd method for add supplied card to specified customer.
func (this EwalletManager) CardAdd(custLogin string, card NotRegisteredCard) (cardResp CardResponse, err error) {
	cust := this.Customers[custLogin]
	err = this.sendEWRequestMap(&cardResp, ADD, cust.plain()+card.plain())
	return
}

//CardActivate method for activate card that was attached to specified customer.
func (this EwalletManager) CardActivate(custLogin string, cardId string, amount string) (cardResp CardResponse, err error) {
	cust := this.Customers[custLogin]
	err = this.sendEWRequestMap(&cardResp, ACTIVATE, cust.plain()+fmt.Sprintf("CardId=%s;Amount=%s;", cardId, amount))
	return
}

//CardRemove method for remove card with specified CardId from cutomer cards list.
func (this EwalletManager) CardRemove(custLogin string, cardId string) (cardResp CardResponse, err error) {
	cust := this.Customers[custLogin]
	err = this.sendEWRequestMap(&cardResp, REMOVE, cust.plain()+fmt.Sprintf("CardId=%s;", cardId))
	return
}

//CustomerRegister method registers supplied customer in Payture system.
func (this EwalletManager) CustomerRegister(cust Customer) (custResp CustomerResponse, err error) {
	return this.customRequest(cust, REGISTER)
}

//CustomerDelete delete specified customer from Payture system.
func (this EwalletManager) CustomerDelete(custLogin string) (custResp CustomerResponse, err error) {
	return this.customRequest(this.Customers[custLogin], DELETE)
}

//CustomerUpdate changed some information in customer's account in Payture system.
func (this EwalletManager) CustomerUpdate(cust Customer) (custResp CustomerResponse, err error) {
	return this.customRequest(cust, UPDATE)
}

//CustomerCheckBylogin check is the customer was registered in Payture system by login.
//The Customer must be in the customer's collection at the time when the method was called.
func (this EwalletManager) CustomerCheckBylogin(custLogin string) (custResp CustomerResponse, err error) {
	return this.customRequest(this.Customers[custLogin], CHECK)
}

//CustomerCheck check is the customer was registered in Payture system.
func (this EwalletManager) CustomerCheck(cust Customer, addToCollection bool) (custResp CustomerResponse, err error) {
	return this.customRequest(cust, CHECK)
}

//GetCardList returned the card's collection that was attached to specified customer.
func (this EwalletManager) GetCardList(custLogin string) (*http.Response, error) {
	cust := this.Customers[custLogin]
	return this.sendEWRequest(GETLIST, cust.plain())
}

func (cust *Customer) FillCardList(resp *http.Response) {
	//defer resp.Body.Close()
	//var card GetList
	body, err := ioutil.ReadAll(resp.Body)
	xml.Unmarshal(body, &cust)
	fmt.Println(cust)
	fmt.Println(err)
}

func (this EwalletManager) sendEWRequest(cmd string, data string) (*http.Response, error) {
	par := reqPrms{}.set(VWID, this.merchant.Key).set(DATAUP, data)
	return sendRequest(this, cmd, par.get())
}

func (this EwalletManager) customRequest(cust Customer, cmd string) (custResp CustomerResponse, err error) {
	err = this.sendEWRequestMap(&custResp, cmd, cust.plain())
	return
}

func (this EwalletManager) sendEWRequestMap(ret Unwrapper, cmd string, data string) (err error) {
	prms := reqPrms{}.set(VWID, this.merchant.Key).set(DATAUP, data)
	err = sendRequestAndMap(ret, prms.get(), this, cmd)
	return
}
