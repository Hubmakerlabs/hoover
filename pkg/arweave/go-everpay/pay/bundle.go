package pay

import (
	"github.com/Hubmakerlabs/hoover/pkg/arweave/go-everpay/account"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/go-everpay/pay/schema"
)

func VerifyBundleSigs(stx schema.BundleWithSigs) *schema.InternalErr {
	// ensure every items through sig validation
	newSigns := make(map[string]string)
	for id, sig := range stx.Sigs {
		_, accid, err := account.IDCheck(id)
		if err != nil {
			log.Error("account.IDCheck(id)", "err", err, "id", id)
			return schema.NewInternalErr(-1, err.Error())
		}
		newSigns[accid] = sig
	}

	for idx, item := range stx.Items {
		// get acc
		acc, err := account.New(item.From)
		if err != nil {
			log.Error("invalid bundle account", "acc", item.From, "err", err)
			return schema.NewInternalErr(idx, err.Error())
		}

		// get sig
		sig, ok := newSigns[acc.ID]
		if !ok {
			log.Error("not found sig", "acc", acc.ID)
			return schema.NewInternalErr(idx, ERR_NOT_FOUND_BUNDLE_SIG.Error())
		}

		// verify sig
		err = acc.VerifySig(account.Transaction{
			Nonce: "99999999999999",
			Hash:  stx.Hash(),
			Sig:   sig,
		})
		if err != nil {
			log.Error("invalid bundle sig", "err", err)
			return schema.NewInternalErr(idx,
				account.ERR_INVALID_SIGNATURE.Error())
		}
	}

	return nil
}
