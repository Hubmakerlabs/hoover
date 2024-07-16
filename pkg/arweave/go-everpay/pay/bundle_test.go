package pay

import (
	"encoding/hex"
	"encoding/json"
	"strings"
	"testing"

	"github.com/Hubmakerlabs/hoover/pkg/arweave/go-everpay/pay/schema"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/goether"
	"github.com/Hubmakerlabs/hoover/pkg/eth/common/hexutil"
	"github.com/Hubmakerlabs/hoover/pkg/eth/crypto"
	"github.com/stretchr/testify/assert"
)

func TestVerifyBundleSigs(t *testing.T) {
	signer1, err := goether.NewSigner("338f76e7463ed64f98e883aa0f522c92cc5881cbce113894559d703d515a55e1")
	assert.NoError(t, err)
	signer2, err := goether.NewSigner("ad1dcf8f1c449e7af21a7b8341eba5f053055819fff9948f1251ea94a0184cae")
	assert.NoError(t, err)

	bundle := schema.Bundle{
		Items: []schema.BundleItem{
			{
				From: signer1.Address.String(),
			},
			{
				From: signer2.Address.String(),
			},
		},
		Expiration: 1728069948,
		Salt:       "",
		Version:    schema.BundleTxVersionV1,
	}
	signData := []byte(bundle.String())

	sig1, _ := signer1.SignMsg(signData)
	sig2, _ := signer2.SignMsg(signData)

	bundleWithSigs := schema.BundleWithSigs{
		Bundle: bundle,
		Sigs: map[string]string{
			signer1.Address.String():                  hexutil.Encode(sig1),
			strings.ToLower(signer2.Address.String()): hexutil.Encode(sig2),
		},
	}

	assert.Nil(t, VerifyBundleSigs(bundleWithSigs))
}

func TestVerifyBundleSigs2(t *testing.T) {
	// sigs address verify
	data := `{
    "bundle": {
        "items": [
            {
                "tag": "ethereum-eth-0x0000000000000000000000000000000000000000",
                "chainID": "42",
                "from": "0x911f42b0229c15bbb38d648b7aa7ca480ed977d6",
                "to": "0x61ebf673c200646236b2c53465bca0699455d5fa",
                "amount": "100000000000000000"
            },
            {
                "tag": "ethereum-usdt-0xd85476c906b5301e8e9eb58d174a6f96b9dfc5ee",
                "chainID": "42",
                "from": "0x61ebf673c200646236b2c53465bca0699455d5fa",
                "to": "0x911f42b0229c15bbb38d648b7aa7ca480ed977d6",
                "amount": "296147410"
            }
        ],
        "expiration": 1645336839,
        "salt": "af2b2d0a-d979-4d15-90d2-de7d7fc0bbd9",
        "version": "v1",
        "sigs": {
            "0x911F42b0229c15bBB38D648B7Aa7CA480eD977d6": "0xb231b845a843c6e2f813739e00519dc069ae8f2240d0707ab2ae82a5ed46a1242267759f4f4cc82d302584895de0249405fe3f877d596b39fc370aad2c65a59a1c",
            "0x61ebf673c200646236b2c53465bca0699455d5fa": "0x112fe6e8981b2f2592a19b18fa99e4ebefbe6d8eee9cc27a1743600439962f1e18946f178717235173e146f037a9ea8efd3b5953d092fb63331d4024086898e41c"
        }
    }
}`

	bundleData := schema.BundleData{}
	err := json.Unmarshal([]byte(data), &bundleData)
	assert.NoError(t, err)

	assert.Nil(t, VerifyBundleSigs(bundleData.Bundle))

	// arweave address sig
	data = `{
    "bundle": {
        "items": [
            {
                "tag": "ethereum-eth-0x0000000000000000000000000000000000000000",
                "chainID": "42",
                "from": "E1YK40az7mbpAYrdvLNp9PdzacT65DaUeJAkobxskyU",
                "to": "0x053Dcc1E2DfD3D60ff46363FDC0f17D5f0667F34",
                "amount": "100000000000000000"
            },
            {
                "tag": "ethereum-usdt-0xd85476c906b5301e8e9eb58d174a6f96b9dfc5ee",
                "chainID": "42",
                "from": "0x053Dcc1E2DfD3D60ff46363FDC0f17D5f0667F34",
                "to": "E1YK40az7mbpAYrdvLNp9PdzacT65DaUeJAkobxskyU",
                "amount": "300000000"
            }
        ],
        "expiration": 1630596000677,
        "salt": "48104227-6335-47dd-a58c-b758a928f2fd",
        "version": "v1",
        "sigs": {
            "E1YK40az7mbpAYrdvLNp9PdzacT65DaUeJAkobxskyU": "e-rQ2HTKCJ48pIeEJgyMsEyYlz8uhXaHGUVSwAFPBwV1pLRKrn3qkVMhCCW_kklO6UGPh7664znRHQ8YsWHGcSulXYdZ-iFNZzKuHM1R9RPMJM0CY9hfW86kfRHbrpkPSVCttQBZKOHDZ9DzDC9u4HVW5Gvmz56fWngXdeY6SsuPRwJ2Abb8j4hqapYsYVd6ZM-Jex2FukLRNrStU4UEyKAen_KD4gIMHTStQX5UhUS9YRX5ngQnoN0W6tS7Jo1iix_m6-RO4M9V6H_Cr_jg8v8rexIJGszyy3ZFwI4QZuEP-2Mp_KE7OOdzLBeiRF3bTMhljGuHTUJ36UrYM00YIV4KUxCBP8ZFk0w1Td3jjxYr_FhacHE8xnhT0hDB-pDHjeMnS_DFmV6vHbqjpQT7Z8MsMppo6v87jFTVuEIG5gPBysKgzNnYmVB8-LKFXY1PCq6gYmJ_LA5O9MqO7t5VmaHLpXfcAYugXcuPSTmdY5HWBCCGDUwfzOmED-QDSIRoo1u_EZa1MHwLObnPiKkT2mu6m37N9C01tMaf6b2cun3WUUNal5LRa8e2oDZ6WbMjJz7M6HI3NaOdYJr0ElxYmsLjgLiwpTh5uy6lUndAkRny0dmdfeAZPCUKiCooSwl-ZBLUyCklnh5nUKBkxAsSfh42XKgbmfQX8tox5OapdlA,vF41vGo9coJK5ZmZ8jh3hOa7FY6v4f_x5tPqah_Nr9OdO_jEvz2iBbEKDjCQepTaAzwCm2qNl-6ou_4br_MfpHSv6Dwl--jmHNgBqaEWTWu7EwuVKqctdkneFKVhGrMwUCCddh4c0NRI-f0kz3PyL1Lrv_fZ2NhfLVkUgrNqTToA3Paiz98uVv7aDzDVm9a-SKBs-oP0ClgFxVpavSee7z3AtT-BIgrrDXVPEVNR_VH_FN-a-UKJSsnvyYYrfn2CaB-9lAGLJd8olKPq6ml225uxh-CTKZ1X2__09er9JvasB10FXb6PyJrM8R1HhkqBlqzj4xwuVl6d8Lhy-mAi71IEmOXL8vqeH3w4wE9HMgf--KZ2MiV5yvTWDsKVdCA-f6Y9sKLyHim5etTOdEQxQDgshe070Cn6_nf4uXNhSi1ortFfBO92LPQOESJZSFg8nPH90ft0zF16TqoOBmmthSP5w__SjlAwMlyhvuKKG6hNHgKa9QO674zQnY_Neh3GAjYuvAVI9Jnmq39xhlA9cNtWPO0OyBLS4CxYrTCi4dxiQnX83JR8TzQcnj09da_KBqpWowB2TAXlIubsuwSetPNmdNByVSPh7_sLmeHPn-ke1dCeCd4sVryenb8pXuVN8ES_l99rza7u8Xtx0xfivEVlqq1d97gNZGyRl-HQ5yE",
            "0x053Dcc1E2DfD3D60ff46363FDC0f17D5f0667F34": "0xeb547488936ec7160b97a90273aa6adea33617f473dd30fe44e4819b2699348f4bc4b3a5a7d587e26558fb96b7a793094fed3d07f2f2d7b63328b51ca4b2825c1c"
        }
    }
}`
	bundleData = schema.BundleData{}
	err = json.Unmarshal([]byte(data), &bundleData)
	assert.NoError(t, err)

	assert.Nil(t, VerifyBundleSigs(bundleData.Bundle))
}

type Pk struct {
	Private string
	Address string
}

func generateTestPrivateKey(count int) []Pk {
	pk := make([]Pk, 0)
	for i := 0; i < count; i++ {
		priv, _ := crypto.GenerateKey()
		addr := crypto.PubkeyToAddress(priv.PublicKey)
		pk = append(pk, Pk{
			Private: hex.EncodeToString(crypto.FromECDSA(priv)),
			Address: addr.String(),
		})
	}
	return pk
}

func TestAsTokenTx(t *testing.T) {
	pks := generateTestPrivateKey(2)
	t.Log(pks)
}
