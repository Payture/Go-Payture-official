package payture

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
)

const (
	SESSION_PAY   string = "Pay"
	SESSION_ADD   string = "Add"
	SESSION_BLOCK string = "Block"
)

type APIType interface {
	getAPI() string
}

type Merchant struct {
	Key, Password, Host string
}

type Payment struct {
	OrderId string
	Amount  string
}

func (order Payment) plain() string {
	return fmt.Sprintf("OrderId=%s;Amount=%s;", order.OrderId, order.Amount)
}

type Card struct {
	CardHolder string
	EMonth     string
	EYear      string
	SecureCode string
}

type CustParams struct {
	CustomerFields map[string]string
}

type PaytureEncoder interface {
	plain() string
}

type ParamsFormer interface {
	content() map[string][]string
}

func (dict CustParams) plain() string {
	resStr := ""
	for key := range dict.CustomerFields {
		resStr += key + "=" + dict.CustomerFields[key] + ";"
	}
	return resStr
}

func sendRequestFormer(url string, params ParamsFormer) (*http.Response, error) {
	return sendRequest(url, params.content())
}

func sendRequest(url string, params map[string][]string) (*http.Response, error) {
	resp, err := http.PostForm(url, params)
	return resp, err
}

func (order Payment) Charge(apiType string, merch Merchant) (OrderResponse, error) {
	url := merch.Host + "/" + apiType + "/Charge"
	key := "Key"
	if apiType == "vwapi" {
		key = "VWID"
	}
	params := make(map[string][]string)
	params[key] = []string{merch.Key}
	params["OrderId"] = []string{order.OrderId}
	if apiType == "vwapi" || apiType == "apim" {
		params["Password"] = []string{merch.Password}
		params["Amount"] = []string{order.Amount}
	}

	httpResp, err := sendRequest(url, params)
	if err != nil {
		return OrderResponse{}, err
	}
	orderResp := OrderResponse{}
	orderResp.MapHttpRespToResp(httpResp)
	return orderResp, nil
}

func (order Payment) Unblock(apiType string, merch Merchant) (OrderResponse, error) {
	url := merch.Host + "/" + apiType + "/Unblock"
	key := "Key"
	if apiType == "vwapi" {
		key = "VWID"
	}
	params := make(map[string][]string)
	params[key] = []string{merch.Key}
	params["OrderId"] = []string{order.OrderId}
	params["Amount"] = []string{order.Amount}

	if apiType == "vwapi" || apiType == "apim" {
		params["Password"] = []string{merch.Password}
	}
	httpResp, err := sendRequest(url, params)
	if err != nil {
		return OrderResponse{}, err
	}
	orderResp := OrderResponse{}
	orderResp.MapHttpRespToResp(httpResp)
	return orderResp, nil
}

func (order Payment) Refund(apiType string, merch Merchant) (OrderResponse, error) {
	url := merch.Host + "/" + apiType + "/Refund"
	params := make(map[string][]string)
	if apiType == "vwapi" {
		params = make(map[string][]string)
		params["VWID"] = []string{merch.Key}
		params["DATA"] = []string{fmt.Sprintf("OrderId=%s;Password=%s;Amount=%s", order.OrderId, merch.Password, order.Amount)}
	} else {
		params = make(map[string][]string)
		params["Key"] = []string{merch.Key}
		params["OrderId"] = []string{order.OrderId}
		params["Amount"] = []string{order.Amount}
		params["Password"] = []string{merch.Password}
	}
	httpResp, err := sendRequest(url, params)
	if err != nil {
		return OrderResponse{}, err
	}
	orderResp := OrderResponse{}
	orderResp.MapHttpRespToResp(httpResp)
	return orderResp, nil
}

func (order Payment) GetState(merch Merchant) (OrderResponse, error) {
	url := merch.Host + "/api/GetState"
	params := make(map[string][]string)
	params["Key"] = []string{merch.Key}
	params["OrderId"] = []string{order.OrderId}
	httpResp, err := sendRequest(url, params)
	if err != nil {
		return OrderResponse{}, err
	}
	orderResp := OrderResponse{}
	orderResp.MapHttpRespToResp(httpResp)
	return orderResp, nil
}

func (order Payment) PayStatus(apiType string, merch Merchant) (OrderResponse, error) {
	url := merch.Host + "/" + apiType + "/PayStatus"
	params := make(map[string][]string)
	params["Key"] = []string{merch.Key}
	params["OrderId"] = []string{order.OrderId}

	if apiType == "vwapi" {
		params = make(map[string][]string)
		params["VWID"] = []string{merch.Key}
		params["DATA"] = []string{fmt.Sprintf("OrderId=%s", order.OrderId)}
	}

	httpResp, err := sendRequest(url, params)
	if err != nil {
		return OrderResponse{}, err
	}
	orderResp := OrderResponse{}
	orderResp.MapHttpRespToResp(httpResp)
	return orderResp, nil
}

func (order Payment) GenerateId(fixedPart string) string {
	return fixedPart + "_" + strconv.FormatInt(rand.Int63(), 10)
}

/*
Parse response
*/

func BodyByte(resp *http.Response) (body []byte, err error) {
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("response code:%d", resp.StatusCode)
		return
	}
	return body, err
}

func Parse(resp *http.Response) (responseText string, err error) {
	defer resp.Body.Close()

	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("response code:%d", resp.StatusCode)
		return
	}

	responseText = fmt.Sprintf("%s", body)
	return responseText, err
}

func ParseByteBody(body []byte) (responseText string) {
	responseText = fmt.Sprintf("%s", body)
	return responseText
}

type OrderResponse struct {
	Success   bool   `xml:"Success,attr"`
	OrderId   string `xml:"OrderId,attr"`
	Amount    int64  `xml:"Amount,attr"`
	ErrCode   string `xml:"ErrCode,attr"`
	NewAmount int64  `xml:"NewAmount,attr"`
}

func (resp *OrderResponse) MapHttpRespToResp(httpResp *http.Response) error {
	byteBody, err := BodyByte(httpResp)
	if err != nil {
		return err
	}
	xml.Unmarshal(byteBody, &resp)
	return nil
}
