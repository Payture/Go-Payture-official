package payture

import (
	"net/http"
)

func APIThreeDSPAY(merch Merchant, orderId string, paRes string) (*http.Response, error) {
	return send3DS(merch, "Pay3DS", orderId, paRes)
}

func APIThreeDSBlock(merch Merchant, orderId string, paRes string) (*http.Response, error) {
	return send3DS(merch, "Block3DS", orderId, paRes)
}

func send3DS(merch Merchant, method string, orderId string, paRes string) (*http.Response, error) {
	url := merch.Host + "/api/" + method
	params := make(map[string][]string)
	params["Key"] = []string{merch.Key}
	params["OrderId"] = []string{orderId}
	params["PaRes"] = []string{paRes}
	return sendRequest(url, params)
}

func PaySubmit3DS(merch Merchant, md string, paRes string) (*http.Response, error) {
	url := merch.Host + "/vwapi/PaySubmit3DS"
	params := make(map[string][]string)
	params["MD"] = []string{md}
	params["PaRes"] = []string{paRes}
	return sendRequest(url, params)
}
