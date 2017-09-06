package payture

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

/*///////////
EWallet types
*/
type EwalletManager struct {
	Customers map[string]Customer
	OrderManager
}

func EwalletService(merch Merchant) (ew EwalletManager) {
	ew.apiType = "vwapi"
	ew.merchant = merch
	ew.Customers = make(map[string]Customer)
	return
}

type NotRegisteredCard struct {
	CardHolder string
	EMonth     string
	EYear      string
	SecureCode string
	CardNumber string
}

type EWalletCard struct {
	CardName   string `xml:"CardName,attr"`
	CardId     string `xml:"CardId,attr"`
	CardHolder string `xml:"CardHolder,attr"`
	Status     string `xml:"Status,attr"`
	Expired    bool   `xml:"Expired,attr"`
	NoCVV      bool   `xml:"NoCVV,attr"`
}

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

/*
EWallet types
*/ ///////////

/////////////////////////
/////////////////////////

/*////////////////////////
Manage Customer collection
*/
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

/*
Manage Customer collection
*/ ////////////////////////

/////////////////////////
/////////////////////////

/*///////////
Transactions
*/
func (this EwalletManager) InitTransaction(custLogin string, sessionType string, ip string, order Payment, tag string, lang string, cardId string) (initResp InitResponse, err error) {
	cust := this.Customers[custLogin]
	err = this.sendEWRequestMap(&initResp, "Init", cust.plain()+order.plain()+fmt.Sprintf("IP=%s;TemplateTag=%s;Language=%s;CardId=%s;SessionType=%s;CustomField=Field;", ip,
		tag, lang, cardId, sessionType /*add mooooooreee*/))
	return
}

func (this EwalletManager) PayRegCard(custLogin string, cardId string, secureCode string, order Payment, ip string, confirmCode string, custFields map[string]string) (payResp EWalletPayResponse, err error) {
	cust := this.Customers[custLogin]
	err = this.sendEWRequestMap(&payResp, "Pay", cust.plain()+order.plain()+fmt.Sprintf("IP=%s;ConfirmCode=%s;CardId=%s;SecureCode=%s;CustomField=Field;", ip, confirmCode, cardId, secureCode /*add mooooooreee*/))
	return
}

func (this EwalletManager) PayNoRegCard(custLogin string, card NotRegisteredCard, order Payment, ip string, confirmCode string, custFields map[string]string, addCardToCust bool) (payResp EWalletPayResponse, err error) {
	cust := this.Customers[custLogin]
	err = this.sendEWRequestMap(&payResp, "Pay", cust.plain()+card.plain()+order.plain()+fmt.Sprintf("IP=%s;ConfirmCode=%s;CardId=FreePay;CustomField=Field;", ip, confirmCode /*add mooooooreee*/))
	return
}

func (this EwalletManager) OperationBySession(cmd string, sessionId string) (*http.Response, error) { //For Add or Pay operation on Payture side
	params := map[string][]string{
		"SessionId": []string{sessionId}}
	return sendRequest(this, cmd, params)
}

func (this EwalletManager) SendCode(custLogin string) (*http.Response, error) {
	cust := this.Customers[custLogin]
	return this.sendEWRequest("SendCode", cust.plain())
}

/*
Transactions
*/ //////////

/////////////////////////
/////////////////////////

/*//////////
Manage Card
*/

func (this EwalletManager) CardAdd(custLogin string, card NotRegisteredCard) (cardResp CardResponse, err error) {
	cust := this.Customers[custLogin]
	err = this.sendEWRequestMap(&cardResp, "Add", cust.plain()+card.plain())
	return
}

func (this EwalletManager) CardActivate(custLogin string, cardId string, amount string) (cardResp CardResponse, err error) {
	cust := this.Customers[custLogin]
	err = this.sendEWRequestMap(&cardResp, "Activate", cust.plain()+fmt.Sprintf("CardId=%s;Amount=%s;", cardId, amount))
	return
}

func (this EwalletManager) CardRemove(custLogin string, cardId string) (cardResp CardResponse, err error) {
	cust := this.Customers[custLogin]
	err = this.sendEWRequestMap(&cardResp, "Remove", cust.plain()+fmt.Sprintf("CardId=%s;", cardId))
	return
}

/*
Manage Card
*/ /////////

/////////////////////////
/////////////////////////

/*/////////////
Manage Customer
*/

func (this EwalletManager) CustomerRegister(cust Customer) (custResp CustomerResponse, err error) {
	return this.customRequest(cust, "Register")
}

func (this EwalletManager) CustomerDelete(custLogin string) (custResp CustomerResponse, err error) {
	return this.customRequest(this.Customers[custLogin], "Delete")
}

func (this EwalletManager) CustomerUpdate(cust Customer) (custResp CustomerResponse, err error) {
	return this.customRequest(cust, "Update")
}
func (this EwalletManager) CustomerCheckBylogin(custLogin string) (custResp CustomerResponse, err error) {
	return this.customRequest(this.Customers[custLogin], "Check")
}

func (this EwalletManager) CustomerCheck(cust Customer, addToCollection bool) (custResp CustomerResponse, err error) {
	return this.customRequest(cust, "Check")
}

func (this EwalletManager) GetCardList(custLogin string) (*http.Response, error) {
	cust := this.Customers[custLogin]
	return this.sendEWRequest("GetList", cust.plain())
}

func (cust *Customer) FillCardList(resp *http.Response) {
	//defer resp.Body.Close()
	//var card GetList
	body, err := ioutil.ReadAll(resp.Body)
	xml.Unmarshal(body, &cust)
	fmt.Println(cust)
	fmt.Println(err)
}

/*
Manage Customer
*/ ////////////

/*//////////////////
Send Ewallet Request
*/

func (this EwalletManager) sendEWRequest(cmd string, data string) (*http.Response, error) {
	params := map[string][]string{
		"VWID": []string{this.merchant.Key},
		"DATA": []string{data}}
	return sendRequest(this, cmd, params)
}

func (this EwalletManager) customRequest(cust Customer, cmd string) (custResp CustomerResponse, err error) {
	err = this.sendEWRequestMap(&custResp, cmd, cust.plain())
	return
}

func (this EwalletManager) sendEWRequestMap(ret Unwrapper, cmd string, data string) (err error) {
	params := map[string][]string{
		"VWID": []string{this.merchant.Key},
		"DATA": []string{data}}
	err = sendRequestAndMap(ret, params, this, cmd)
	return
}

/*
Send Ewallet Request
*/ //////////////////

/*///////
Responses
*/

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

/*
Responses
*/ ///////
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
