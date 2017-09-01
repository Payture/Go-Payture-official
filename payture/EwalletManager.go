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
	api       APIType
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
	return fmt.Sprintf("VWUserLgn=%s;VWUserPsw=%s;PhoneNumber=%s;Email=%s;", customer.VWUserLgn, customer.VWUserPsw, customer.PhoneNumber, customer.Email)
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
func (this EwalletManager) AddCustomer(cust Customer) {

}

func (this EwalletManager) RemoveCustomerFromCollection(custLogin string) {

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
	url := this.Merchant.Host + "/vwapi/" + cmd
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

func (this EwalletManager) CardAdd(custLogin string, card NotRegisteredCard) (*http.Response, error) {
	cust := this.Customers[custLogin]
	return this.sendEWRequest("Add", cust.plain()+card.plain())
}

func (this EwalletManager) CardActivate(custLogin string, cardId string, amount string) (*http.Response, error) {
	cust := this.Customers[custLogin]
	return this.sendEWRequest("Activate", cust.plain()+fmt.Sprintf("CardId=%s;Amount=%s;", cardId, amount))
}

func (this EwalletManager) CardRemove(custLogin string, cardId string) (*http.Response, error) {
	cust := this.Customers[custLogin]
	return this.sendEWRequest("Remove", cust.plain()+fmt.Sprintf("CardId=%s;", cardId))
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

func (this EwalletManager) CustomerRegister(cust Customer) (*http.Response, error) {
	return this.sendEWRequest("Register", cust.plain())
}
func (this EwalletManager) CustomerDelete(custLogin string) (*http.Response, error) {
	cust := this.Customers[custLogin]
	return this.sendEWRequest("Delete", cust.plain())
}

func (this EwalletManager) CustomerUpdate(cust Customer) (*http.Response, error) {
	return this.sendEWRequest("Update", cust.plain())
}
func (this EwalletManager) CustomerCheckBylogin(custLogin string) (*http.Response, error) {
	cust := this.Customers[custLogin]
	return this.sendEWRequest("Check", cust.plain())
}

func (this EwalletManager) CustomerCheck(cust Customer, addToCollection bool) (*http.Response, error) {
	return this.sendEWRequest("Check", cust.plain())
}

/*
Manage Customer
*/ ////////////

/*
Send Ewallet Request
*/

func (this EwalletManager) sendEWRequest(cmd string, data string) (*http.Response, error) {
	url := this.Merchant.Host + "/vwapi/" + cmd
	params := make(map[string][]string)
	params["VWID"] = []string{this.Merchant.Key}
	params["DATA"] = []string{data}
	return sendRequest(url, params)
}
