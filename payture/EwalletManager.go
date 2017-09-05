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
	Merchant  Merchant
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

func (this EwalletManager) getAPI() string {
	return "vwapi"
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
func (this EwalletManager) InitTransaction(custLogin string, sessionType string, ip string, order Payment, tag string, lang string, cardId string) (*http.Response, error) {
	cust := this.Customers[custLogin]
	return this.sendEWRequest("Init", cust.plain()+order.plain()+fmt.Sprintf("IP=%s;TemplateTag=%s;Language=%s;CardId=%s;SessionType=%s;CustomField=Field;", ip,
		tag, lang, cardId, sessionType /*add mooooooreee*/))
}

func (this EwalletManager) PayRegCard(custLogin string, cardId string, secureCode string, order Payment, ip string, confirmCode string, custFields map[string]string) (*http.Response, error) {
	cust := this.Customers[custLogin]
	return this.sendEWRequest("Pay", cust.plain()+order.plain()+fmt.Sprintf("IP=%s;ConfirmCode=%s;CardId=%s;SecureCode=%s;CustomField=Field;", ip, confirmCode, cardId, secureCode /*add mooooooreee*/))
}

func (this EwalletManager) PayNoRegCard(custLogin string, card NotRegisteredCard, order Payment, ip string, confirmCode string, custFields map[string]string, addCardToCust bool) (*http.Response, error) {
	cust := this.Customers[custLogin]
	return this.sendEWRequest("Pay", cust.plain()+card.plain()+order.plain()+fmt.Sprintf("IP=%s;ConfirmCode=%s;CardId=FreePay;CustomField=Field;", ip, confirmCode /*add mooooooreee*/))
}

func (this EwalletManager) OperationBySession(cmd string, sessionId string) (*http.Response, error) { //For Add or Pay operation on Payture side
	url := this.Merchant.Host + "/" + this.getAPI() + "/" + cmd
	params := make(map[string][]string)
	params["SessionId"] = []string{sessionId}
	return sendRequest(url, params)
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

func (this EwalletManager) CardAdd(custLogin string, card NotRegisteredCard) (CardResponse, error) {
	cust := this.Customers[custLogin]
	cardResp := CardResponse{}
	httpResp, err := this.sendEWRequest("Add", cust.plain()+card.plain())
	if err != nil {
		return cardResp, err
	}
	cardResp.MapHttpRespToResp(httpResp)
	return cardResp, nil
}

func (this EwalletManager) CardActivate(custLogin string, cardId string, amount string) (CardResponse, error) {
	cust := this.Customers[custLogin]
	cardResp := CardResponse{}
	httpResp, err := this.sendEWRequest("Activate", cust.plain()+fmt.Sprintf("CardId=%s;Amount=%s;", cardId, amount))
	if err != nil {
		return cardResp, err
	}
	cardResp.MapHttpRespToResp(httpResp)
	return cardResp, nil
}

func (this EwalletManager) CardRemove(custLogin string, cardId string) (CardResponse, error) {
	cust := this.Customers[custLogin]
	cardResp := CardResponse{}
	httpResp, err := this.sendEWRequest("Remove", cust.plain()+fmt.Sprintf("CardId=%s;", cardId))
	if err != nil {
		return cardResp, err
	}
	cardResp.MapHttpRespToResp(httpResp)
	return cardResp, nil
}

/*
Manage Card
*/ /////////

/////////////////////////
/////////////////////////

/*/////////////
Manage Customer
*/

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

func (this EwalletManager) CustomerRegister(cust Customer) (CustomerResponse, error) {

	custResp := CustomerResponse{}
	httpResp, err := this.sendEWRequest("Register", cust.plain())
	if err != nil {
		return custResp, err
	}
	custResp.MapHttpRespToResp(httpResp)
	return custResp, nil
}

func (this EwalletManager) CustomerDelete(custLogin string) (CustomerResponse, error) {
	cust := this.Customers[custLogin]
	custResp := CustomerResponse{}
	httpResp, err := this.sendEWRequest("Delete", cust.plain())
	if err != nil {
		return custResp, err
	}
	custResp.MapHttpRespToResp(httpResp)
	return custResp, nil
}

func (this EwalletManager) CustomerUpdate(cust Customer) (CustomerResponse, error) {
	custResp := CustomerResponse{}
	httpResp, err := this.sendEWRequest("Update", cust.plain())
	if err != nil {
		return custResp, err
	}
	custResp.MapHttpRespToResp(httpResp)
	return custResp, nil
}
func (this EwalletManager) CustomerCheckBylogin(custLogin string) (CustomerResponse, error) {
	cust := this.Customers[custLogin]
	custResp := CustomerResponse{}
	httpResp, err := this.sendEWRequest("Check", cust.plain())
	if err != nil {
		return custResp, err
	}
	custResp.MapHttpRespToResp(httpResp)
	return custResp, nil
}

func (this EwalletManager) CustomerCheck(cust Customer, addToCollection bool) (CustomerResponse, error) {
	custResp := CustomerResponse{}
	httpResp, err := this.sendEWRequest("Check", cust.plain())
	if err != nil {
		return custResp, err
	}
	custResp.MapHttpRespToResp(httpResp)
	return custResp, nil
}

/*
Manage Customer
*/ ////////////

/*//////////////////
Send Ewallet Request
*/

func (this EwalletManager) sendEWRequest(cmd string, data string) (*http.Response, error) {
	url := this.Merchant.Host + "/" + this.getAPI() + "/" + cmd
	params := make(map[string][]string)
	params["VWID"] = []string{this.Merchant.Key}
	params["DATA"] = []string{data}
	return sendRequest(url, params)
}

/*
Send Ewallet Request
*/ //////////////////

/*//////////////
Payments command
*/
func (this EwalletManager) Unblock(order Payment) (OrderResponse, error) {
	return order.Unblock(this.getAPI(), this.Merchant)
}

func (this EwalletManager) Refund(order Payment) (OrderResponse, error) {
	return order.Refund(this.getAPI(), this.Merchant)
}

func (this EwalletManager) Charge(order Payment) (OrderResponse, error) {
	return order.Charge(this.getAPI(), this.Merchant)
}

func (this EwalletManager) PayStatus(order Payment) (OrderResponse, error) {
	return order.PayStatus(this.getAPI(), this.Merchant)
}

/*
Payments command
*/ //////////////

/*///////
Responses
*/

type CustomerResponse struct {
	Success   bool   `xml:"Success,attr"`
	ErrorCode string `xml:"ErrrCode,attr"`
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
func (resp *CustomerResponse) MapHttpRespToResp(httpResp *http.Response) error {
	byteBody, err := BodyByte(httpResp)
	if err != nil {
		return err
	}
	xml.Unmarshal(byteBody, &resp)
	return nil
}

func (resp *CardResponse) MapHttpRespToResp(httpResp *http.Response) error {
	byteBody, err := BodyByte(httpResp)
	if err != nil {
		return err
	}
	xml.Unmarshal(byteBody, &resp)
	return nil
}
