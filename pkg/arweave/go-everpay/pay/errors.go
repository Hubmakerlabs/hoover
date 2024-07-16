package pay

import (
	"errors"
)

var (
	ERR_LARGER_DATA               = errors.New("err_larger_data")
	ERR_INVALID_OWNER             = errors.New("err_invalid_owner")
	ERR_INVALID_TX_VERSION        = errors.New("err_invalid_tx_version")
	ERR_INVALID_AMOUNT            = errors.New("err_invalid_amount")
	ERR_INVALID_FEE               = errors.New("err_invalid_fee")
	ERR_INVALID_TARGET_CHAIN_TYPE = errors.New("err_invalid_target_chain_type")
	ERR_INVALID_FEE_RECIPIENT     = errors.New("err_invalid_fee_recipient")
	ERR_INVALID_CHAINID           = errors.New("err_invalid_chainid")
	ERR_INVALID_BUNDLE_DATA       = errors.New("err_invalid_bundle_data")
	ERR_BUNDLE_EXECUTED           = errors.New("err_bundle_executed")
	ERR_BUNDLE_EXPIRED            = errors.New("err_bundle_expired")
	ERR_BUNDLE_SALT               = errors.New("err_bundle_salt")
	ERR_BUNDLE_VERSION            = errors.New("err_bundle_version")
	ERR_NOT_FOUND_TOKEN           = errors.New("err_token_not_found")
	ERR_NOT_FOUND_BUNDLE_SIG      = errors.New("err_not_found_bundle_sig")
	ERR_NOT_FOUND_BUNDLE_ITEMS    = errors.New("err_not_found_bundle_items")
)
