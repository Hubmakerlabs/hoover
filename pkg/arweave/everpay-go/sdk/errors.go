package sdk

import "errors"

var (
	ErrTokenNotExist   = errors.New("err_not_exist_token")
	ErrBurnFeeNotExist = errors.New("err_not_exist_burn_fee")
	ErrNotBundleTx     = errors.New("err_not_bundle_tx")
	ErrNotJsonData     = errors.New("err_not_json_data")
)
