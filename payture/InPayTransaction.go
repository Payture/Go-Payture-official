package payture

import (
	"fmt"
	"net/http"
)

type InPayTransaction struct {
	Language         string
	TemplateTag      string
	Url              string
	IP               string
	Order            Payment
	SessionType      string
	AdditionalFields CustParams
}

func (init InPayTransaction) Init(merch Merchant) (*http.Response, error) {
	url := merch.Host + "/apim/Init"
	params := make(map[string][]string)
	params["Key"] = []string{merch.Key}
	params["Data"] = []string{init.Order.plain() + fmt.Sprintf("SessionType=%s;Language=%s;IP=%s;TemplateTag=%s;Url=%s", init.SessionType,
		init.Language, init.IP, init.TemplateTag, init.Url) + init.AdditionalFields.plain()}
	return sendRequest(url, params)
}
