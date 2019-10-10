package payture

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
)

type Merchant struct {
	Key, Password, Host string
}

type Payment struct {
	OrderId string
	Amount  string
}

type UrlProvider interface {
	getReqUrl(cmd string) string
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

func sendRequestFormer(reqUrl UrlProvider, cmd string, params ParamsFormer) (*http.Response, error) {
	return sendRequest(reqUrl, cmd, params.content())
}

func sendRequestFormerMap(ret Unwrapper, reqUrl UrlProvider, cmd string, params ParamsFormer) error {
	return sendRequestAndMap(ret, params.content(), reqUrl, cmd)
}

func sendRequest(reqUrl UrlProvider, cmd string, params map[string][]string) (*http.Response, error) {
	resp, err := http.PostForm(reqUrl.getReqUrl(cmd), params)
	return resp, err
}

func sendRequestAndMap(ret Unwrapper, params map[string][]string, reqUrl UrlProvider, cmd string) (err error) {
	httpResp, err := http.PostForm(reqUrl.getReqUrl(cmd), params)
	if err != nil {
		return
	}
	byteBody, err := BodyByte(httpResp)
	if err != nil {
		return
	}
	return ret.Unwrap(byteBody)
}

func MapHttpRespToResp(ret Unwrapper, httpResp *http.Response) error {
	byteBody, err := BodyByte(httpResp)
	if err != nil {
		return err
	}
	ret.Unwrap(byteBody)
	return nil
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

type Unwrapper interface {
	Unwrap(byteBody []byte) error
}

type reqPrms struct {
	requestParams map[string][]string
	created       bool
}

func (req reqPrms) createSet() reqPrms {
	req.requestParams = make(map[string][]string)
	req.created = true
	return req
}

func (req reqPrms) set(key string, data string) reqPrms {
	if !req.created {
		req = req.createSet()
	}
	req.requestParams[key] = []string{data}
	return req
}

func (req reqPrms) get() map[string][]string {
	return req.requestParams
}

func (req reqPrms) setKey(merch Merchant) reqPrms {
	return req.set(KEY, merch.Key)
}

func (req reqPrms) setPass(merch Merchant) reqPrms {
	return req.set(PASSWORD, merch.Password)
}
