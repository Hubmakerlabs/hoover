package uploader

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// GetBody reads out the body from a http.Response
func GetBody(res *http.Response) (s string, err error) {
	if res == nil {
		return
	}
	if res.Body == nil {
		return
	}
	defer res.Body.Close()
	var body []byte
	body, err = io.ReadAll(res.Body)
	if err != nil {
		return
	}
	s = string(body)
	return
}


func GetBalance(endpoint, account string) (bal int64, err error) {
	var res *http.Response
	address := fmt.Sprintf("%s/wallet/%s/balance", endpoint, account)
	res, err = http.Get(address)
	if err != nil {
		return
	}
	var body string
	if body, err = GetBody(res); err != nil {
		return
	}
	return strconv.ParseInt(body, 10, 64)
}


func Upload() (code int, err error) {
	return
}
