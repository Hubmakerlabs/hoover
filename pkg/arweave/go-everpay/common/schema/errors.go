package schema

import "errors"

var (
	ErrTxNotExist   = errors.New("err_not_found")
	ErrArClientNil  = errors.New("err_ar_client_nil")
	ErrEthClientNil = errors.New("err_eth_client_nil")
)
