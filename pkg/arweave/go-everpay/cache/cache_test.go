package cache

import (
	"encoding/json"
	"testing"

	"github.com/Hubmakerlabs/hoover/pkg/arweave/go-everpay/cache/schema"
	paySchema "github.com/Hubmakerlabs/hoover/pkg/arweave/go-everpay/pay/schema"
	"github.com/Hubmakerlabs/hoover/pkg/eth/crypto"
	"github.com/stretchr/testify/assert"
)

func TestInternalStatus_Marshal(t *testing.T) {
	is := schema.InternalStatus{
		Status: schema.InternalStatusSuccess,
	}
	assert.Equal(t, `{"status":"success"}`, is.Marshal())

	internalErr := paySchema.NewInternalErr(1, "aaaa")
	is = schema.InternalStatus{
		Status:      schema.InternalStatusFailed,
		InternalErr: internalErr,
	}
	assert.Equal(t, `{"status":"failed","index":1,"msg":"aaaa"}`, is.Marshal())
}

func TestInternalErrToStatus(t *testing.T) {
	assert.Equal(t, `{"status":"success"}`, InternalErrToStatus(nil).Marshal())

	internalErr := &paySchema.InternalErr{
		Index: 1,
		Msg:   "aaa",
	}
	assert.Equal(t, "{\"status\":\"failed\",\"index\":1,\"msg\":\"aaa\"}",
		InternalErrToStatus(internalErr).Marshal())
}

func TestRootHash(t *testing.T) {
	hash := RootHash([]byte("1"), []byte("2"), []byte("3"))
	assert.Equal(t, crypto.Keccak256([]byte("123")), hash)

	hash = RootHash([]byte("1"), []byte("2"))
	assert.Equal(t, crypto.Keccak256([]byte("12")), hash)
}

func TestCache_txsByAcc(t *testing.T) {
	c := New()
	bundleTx := `{
            "id": "qjMCmP5B4cU9q1fSVsGkEghpHve4skGrjpd2Hx0vivA",
            "tokenSymbol": "eth",
            "action": "bundle",
            "from": "0xc8618ab07FFeBb5bb743b341570A1140d43C865b",
            "to": "0xc8618ab07FFeBb5bb743b341570A1140d43C865b",
            "data": "{\"bundle\":{\"items\":[{\"tag\":\"ethereum-usdt-0xd85476c906b5301e8e9eb58d174a6f96b9dfc5ee\",\"chainID\":\"42\",\"from\":\"0xc8618ab07FFeBb5bb743b341570A1140d43C865b\",\"to\":\"0x61EbF673c200646236B2c53465bcA0699455d5FA\",\"amount\":\"29227124\"},{\"tag\":\"ethereum-usdt-0xd85476c906b5301e8e9eb58d174a6f96b9dfc5ee\",\"chainID\":\"42\",\"from\":\"0xc8618ab07FFeBb5bb743b341570A1140d43C865b\",\"to\":\"0x258D9aeF9184d9e21f9b882E243c42fAC1466A59\",\"amount\":\"57028\"}],\"expiration\":1640702479,\"salt\":\"f2e342cc-4acf-4d3b-b202-51311c9d6a01\",\"version\":\"v1\",\"sigs\":{\"0xc8618ab07FFeBb5bb743b341570A1140d43C865b\":\"0xd2a2b38814bd9d54fd329e7f3c24b5ca7d8013260f999dc16b492ea4f0410342186296c96cb8cc50327b03eb8c059a72c0914e270c1ddfff9719d16a247e53971b\"}}}"
        }`

	tx := paySchema.Transaction{}
	err := json.Unmarshal([]byte(bundleTx), &tx)
	assert.NoError(t, err)
	c.AddTx(tx, nil)
	assert.Equal(t, len(c.txsByAcc), 3)

	burnTx := `{
            "id": "3FaV46Mhf4zLvmCRJEw2iwH7ZC8J8otGQWmQ43tMyi4",
            "tokenSymbol": "AR",
            "action": "burn",
            "from": "0xc8618ab07FFeBb5bb743b341570A1140d43C865b",
            "to": "Es8qDZi6_ExfUGI_068XK2Y2Xw58bpizDSV3B3xC6h8",
            "amount": "78633570656310"
           }`
	tx = paySchema.Transaction{}
	err = json.Unmarshal([]byte(burnTx), &tx)
	assert.NoError(t, err)
	c.AddTx(tx, nil)
	assert.Equal(t, len(c.txsByAcc), 3)

	transferTx := `{
            "id": "3FaV46Mhf4zLvmCRJEw2iwH7ZC8J8otGQWmQ43tMyi4",
            "tokenSymbol": "AR",
            "action": "transfer",
            "from": "-Vnama5Ngl1LeBEkaGJWj1Zn2hA6j1HezWzwMvY3iLA",
            "to": "0xB6eB83b3e77209f6a7c375923Fb31b2427A3Aa18",
            "amount": "78633570656310"
           }`
	tx = paySchema.Transaction{}
	err = json.Unmarshal([]byte(transferTx), &tx)
	assert.NoError(t, err)
	c.AddTx(tx, nil)
	assert.Equal(t, len(c.txsByAcc), 5)
}

func TestGetMintTargetTxHash(t *testing.T) {
	mintArTx := paySchema.Transaction{}
	err := json.Unmarshal([]byte(`{
            "chainType": "arweave,ethereum",
			"data": "{\"format\": 2, \"id\": \"3aVtIp0afWPn_HDOp-EXGzhpFn2HEgTonDl753qVxTM\", \"last_tx\": \"TvWQO5Z9-Yqdpu1eiH2FrNkN0FNRrBPGaEXOv0J8t6-EOFRI6hyW_ibsuZzopdhM\", \"owner\": \"zOrp5BKdstzcr0-UCdzwPDDl0ywf4L1zFs7E9tnLaXcseQERKLN5_Zcj8uMbyl43eiRhHfh1s6ZxNErw_8ENXyTukNQ9dWOKszHpPLCt_m8inzFEZ53oqeRiBOQVZuR2nWqI_HYMbpRdo2eb9j8_gDWHvAHp85Ru0Sl6Fb2Byzx6j7olYBHcC4WbdTq1SI_gNqhtxXZkBOYDxkOJWOpt6ts3-aZUlNTaB3F411Ho5CXp6xUaUyDtMqWG-MM23LunJwBkVQZf_Yh99G9mSwz8WPSnVsYa3uySoAfxVJ3opxIBDv2Prl3FAfYdrtbN1RytXmVX40IvchMpYZWC9P2Z_-Oaq6ZwGYLf7unYQ6zhZwo1M69_YqdvzhVVI7eN9ym-DzgsE6ecHY-5YFoah6B0s3lw8cYXLk2zj0jab6GnKfih5r3hatrkwpa9C3910Yd31UGi74FyfDi1iczMnHJu8NACra--Vv2-iFBuyQ3Hee0XTmCVfpUenlrkPmBKd_3m4SItwvQBbf1-kD1deVc7oy2eiisnDz_VxeHqakLt-rsRTWmG3amKoFryFkODLtBZK9aCe7UsjUklzvDDLyG0RGf_sTGo7poWdhJC046vOgJuawdp0zNqsJMxJfTJFL-BQWJvPu-aaOY10mnWE84fWj3Xm2YRTCURhEW7qtkwD8c\", \"tags\": [{\"name\": \"U2lnbmluZy1DbGllbnQ\", \"value\": \"QXJDb25uZWN0\"}, {\"name\": \"U2lnbmluZy1DbGllbnQtVmVyc2lvbg\", \"value\": \"MC40LjA\"}], \"target\": \"dH-_dwLlN86fitrFZzi86IVEEQFyYpTzWcqnFh460ys\", \"quantity\": \"10000000000\", \"data\": \"\", \"data_size\": \"0\", \"data_tree\": [], \"data_root\": \"\", \"reward\": \"562392\", \"signature\": \"lntmZaMBxmKcC-BaRq2EuONVN1wFLIoDlOlEuc8IR2TI26WFVZJKoGMy9U0YXBfwcqxo4Vx8_wbRCfETIUycGBXqs4hmmhXgsNqxpApi72R2DNufAD09CNlAZuM1uvveUaAQ1xaxw5zFXLbOpX9kHQNCvcd5iC_PsD-9It2AKGyX8NSFEVTgglD2yuMtzRDmE_O9JBluyryE-sO5RQExRa-iDRkwmotUlApZci-Ms8uIXpFRs3x4Xt1N1v2GiA_CJO8Yve-FsUVLxlpkg66Dwz6ntkbZZZ8j-T420wsPjTc4Tu3VCIYw0YLZN5A47kSTg9SpvNwj4b2SguHN6yRcn2JFsC-SgPYJDdybnCZ1v5UXq-FfscfnaXtoto-CdccGE_7J30fhAnqeJ9a9f9pqLLnPYE8Ry5yBkpPWkuslOtHw7dvPKWa-Dfz4RVKC7EYYqdX6IfYW9wGSJuZ2--FTT5onFiXhhJl7VP4zXtpy2ELjGMVLzZj6Rw1KxTI6vCczIk_TDc_8fQ9Vs7C-XjjVGlHB_K5eDyPwU56YEKZfSsfAXc_5GXDbIDhaxMn1Q2aFSqYIB4857iOB9uStOCpBcDFYu1DZ1wSSHKdvItHkNVA0Z1dTcPVZoiVmDdpM0pZxghEs1a7QdQ2uIB97BNXFD4bwhRm2sxe0VhRUHv6J9jY\", \"targetChainType\": \"arweave\"}"
           }`), &mintArTx)
	assert.NoError(t, err)
	targetChainTxHash, err := GetMintTargetTxHash(mintArTx.ChainType,
		mintArTx.Data)
	assert.NoError(t, err)
	assert.Equal(t, "3aVtIp0afWPn_HDOp-EXGzhpFn2HEgTonDl753qVxTM",
		targetChainTxHash)

	mintEthArTx := paySchema.Transaction{}
	err = json.Unmarshal([]byte(`{
            "chainType": "arweave,ethereum",
			"data": "{\"hash\": \"0xfab7713ef2ff85ca205bac3a0bb536414454b5c72d5c90e21cef572e13663993\", \"nonce\": \"0x6d\", \"blockHash\": \"0x9e66558c355c049a8839efd7ab3465dd830a4671632c45be2007dbcea694d47a\", \"blockNumber\": \"0xd4d630\", \"transactionIndex\": \"0x62\", \"chainId\": \"0x1\", \"condition\": null, \"creates\": null, \"from\": \"0x0048E848C17F29Daf066DC0Bf4770e3f94CA7602\", \"to\": \"0x4FaDC7A98f2Dc96510e42dD1A74141eEae0C1543\", \"value\": \"0x0\", \"gas\": \"0x186a0\", \"gasPrice\": \"0x1c55c18cbd\", \"input\": \"0xa9059cbb00000000000000000000000038741a69785e84399fcf7c5ad61d572f7ecb1dab000000000000000000000000000000000000000000000000000003a5a6a02400\", \"r\": \"0x9381e69707e9e6fb1993923c8a00cb7921b25903ab3de3365e0993513f83f729\", \"s\": \"0x3d4f8545e0ee6974562c70a8ca2281ca586d4c809c434316cc5e23421276f649\", \"v\": \"0x26\", \"targetChainType\": \"ethereum\"}"
			}`), &mintEthArTx)
	assert.NoError(t, err)
	targetChainTxHash, err = GetMintTargetTxHash(mintEthArTx.ChainType,
		mintEthArTx.Data)
	assert.NoError(t, err)
	assert.Equal(t,
		"0xfab7713ef2ff85ca205bac3a0bb536414454b5c72d5c90e21cef572e13663993",
		targetChainTxHash)

	mintEthTx := paySchema.Transaction{}
	err = json.Unmarshal([]byte(`{
            "chainType": "ethereum",
			"data": "{\"hash\": \"0x860ca4d9d2332f0be07b25dde0628335636c42dc2e6d5f1d579d60ee2643c915\", \"nonce\": \"0x1\", \"blockHash\": \"0x6d6fcbe95f55383b32615f98cf538d993de9479e0dc888c27358dbb3456906ad\", \"blockNumber\": \"0xd4d167\", \"transactionIndex\": \"0x73\", \"chainId\": \"0x1\", \"condition\": null, \"creates\": null, \"from\": \"0xc62438F6421d89DEbB0DEe5775fE1694F8D38b27\", \"to\": \"0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48\", \"value\": \"0x0\", \"gas\": \"0x186a0\", \"gasPrice\": \"0x23166b8300\", \"input\": \"0xa9059cbb00000000000000000000000038741a69785e84399fcf7c5ad61d572f7ecb1dab000000000000000000000000000000000000000000000000000000001dcd6500\", \"r\": \"0x7660634dc0263761da357bf730cbe3e6e2c20f315513245f305b3523f03ef76d\", \"s\": \"0x11df741a0ce2a5c3f4a8bb6b974a57c82f19bd872cab5372da8c46c39ca29408\", \"v\": \"0x25\"}"
		}`), &mintEthTx)
	assert.NoError(t, err)
	targetChainTxHash, err = GetMintTargetTxHash(mintEthTx.ChainType,
		mintEthTx.Data)
	assert.NoError(t, err)
	assert.Equal(t,
		"0x860ca4d9d2332f0be07b25dde0628335636c42dc2e6d5f1d579d60ee2643c915",
		targetChainTxHash)
}
