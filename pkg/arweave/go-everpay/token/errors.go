package token

import "errors"

var (
	ERR_TX_MINTED                = errors.New("err_tx_minted")
	ERR_NIL_AMOUNT               = errors.New("err_nil_amount")
	ERR_NEGATIVE_AMOUNT          = errors.New("err_negative_amount")
	ERR_NEGATIVE_FEE             = errors.New("err_negative_fee")
	ERR_LARGER_AMOUNT            = errors.New("err_larger_amount")
	ERR_INSUFFICIENT_BALANCE     = errors.New("err_insufficient_balance")
	ERR_INSUFFICIENT_TOTALSUPPLY = errors.New("err_insufficient_totalsupply")
	ERR_INVALID_ACTION           = errors.New("err_invalid_action")
	ERR_INVALID_MINT_TX          = errors.New("err_invalid_mint_tx")
	ERR_INVALID_BURN_TX          = errors.New("err_invalid_burn_tx")
	ERR_INVALID_TOKEN            = errors.New("err_invalid_token")
	ERR_NOT_FOUND_TARGET_INFO    = errors.New("err_not_found_target_chain_info")
)
