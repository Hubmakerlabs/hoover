package sdk

import (
	"testing"
)

func TestSDK_Transfer(t *testing.T) {
	// priv01 := `{"kty":"RSA","e":"AQAB","n":"yB4ENow16LnZuqOrxBhtfCgj53JDSxMq57q4bj49dmAbdy1RzM9Hm6jWsGy2oIySjbd-kBn0CKhPLz_49QRZ6CDexSODMrmu4W-MioipIJwWh_TFMyzCz8F1dT0Falz3tvKX0g9Lp4LU5TNZVhykHImpRSkVOKwEp2BD-vKGEZrUFuA3gbYpuStz7tccu5sRWX3IV5W7yidcEPjkxcjLOmsl2iVVlmfvJNU5olX5Uohr7rCIaeV3X2uBKR07OAGTxlZhtlIQfcEbzeFKxZavtjYub7W-xNid5yJB1zhFfxA2sYhFdm-zoOyKeH5RqHnzJvHehjYF8dFrpjIqJGasbbvT9djMKju3RUD1p2m7aCqTSgw4ewsgJxZK3xHndldAPJODddLivUff-pzf3zXXHHTSRjZdsIAn98p2r9aVqKiTXyjSSK8ZoQIUubHdSRsj3lgKmI1OY8SAeiXVGQ-dM0WfRzX2MA2C_e5uJG-9PnMx8SMlSY7WINZqvsnj6xNyT3VxnfGV-lPr6_S8_G6oqAuXL13YuEccJE9sEL-pcuWzMqdHsdMOzrcYdY7djbUrltZv4gwoWWRjNxO_Ib4SQQJ6lVZ4dF42HYSonX6f9DY7505xAxLQ4TaQ960PHWeRISSxXXl0bLsdwtqUnInacS_wCBTnMqOfRIKMuYh4Eek","d":"W331uSazgvpthDpOTrJIZCpJ-AgyogFabWXd6nFEUnmwNPWgX34pJA62USTy_kYmvuo00Bc9jjO9OEIMRaE-pbVZhFIvKqLgRKtmWvQONx2F45TSEAqX6rptRDOo5S8qBhe1t9FMaoKTaVLntkoT8gQyQq98x1NGMWv-yJ3XE_dwniDW39P0EXe_YmQw-8jyzrKGRRztJYiHPC8_EPoziDryrGDDV6efG3ztA75A5i3_ArrtM1U5FYQmIjD-ISX_dP-LXvr4FTaTyrYFJH83SJxO8L3sQT00PDdaWNaeV3CHQPWnpHV_b8SVrm86UEepFpYggcXzQ9LWX76xq8DRde76UuYyVSlMhm5l9LaESXBSdx3ldRCVLLzAffPrdiy4gta9PTq14z3UHK1bmI_ILE8COHAttES0lPvjE0movJRou5JPMCp2zZX2H5Dy7PDWn_3cbnstzOGWRtQcI-hv13B1Nk95l3KGYJx4Nncmce4eDWUP4RBTg2SFpY2tY6opGlhK0Th5M2pqCzmFhgqYc0-k4jRpvOrl7_afFWrFo88-XWTsywP7mflGvEEh0D6QK1nwosfx5gTHfZe0g-DJHwU0F_fMZhhmv0YRjEO_xNj9Ugva3CB1c68rxdRuT5hCqh_dZi5L82En0KN69KN6iuVCtn12KJ2jp2TqsoDco50","p":"60H6zepl0rZ2Byt__8r2Rr7ZMW7gGrWEf3RfuycpSggXMs_2JA2ky3GyjQqnt4E2FowO3RnACNG3ftK7oc1d4DdipGi9Y0mCbK76TpCnCQA2_Cc4RpKmFjdAHiCmi33yOZBnPOE8was8G7RLPXkX-UXyVZ95e1YFZ3uTBylVqRs0rwuIpjx2EJnFBbRR7bdWYLoy1wBdVjy3pq1DqXpR4zjfSw0jBbAfGeRP6WFY9Ana9JzKLPozQOULfQobPtPhldny3y53q6-qD8UINl8kvBx6dNV0jnvrv0dXQmq28KO7vld6hotqY3zz94GRUky6eJpcosYK7VBWtJ7zQNQO5w","q":"2cLg1YGg561hgYmsI4mpmW4-bPWu8tLpCueHIsrX41KnIaEwhA-Pz-sg9glehveVU_3R9rXrVhpADV9JDJDxYNswDbJ_BYXWsF1I4A-wy0EmDUW9HglYKvLEiA9JHbnW4w8duOJPv50ej0PbrfTjW_4p7DofH0_JEfNzZfORf2LDPVqLaS5nvP-wkhZU14FEVpOHC8VrUGcrr9G5uAwYF0N6IT-xdUTMZkCDmJD25vGDVn3Ca2ne9mB6Jd4H7gpijHGiIjkdJhC4IBElDlNPtszFf-6WESWiFL0bByXJ1VxK0LfEJGu5fxCkALdvd5I9IUgwim5ofCXbtlARc9zOrw","dp":"MwUzJFcyT_lcLX_kmY_lyz_2kH7wQTqeipmtbUQ38yNADLgHNJh05d45j9cAuo5eMZOcYZ0FBaSdu_Zt1NBaDJdWYDuK3y6BB5sHE-TRJOWUBylmTf_d3zGKST5hgB0HXC6SBST_7Sx6s7NtId7SMTBXNvSH_xSPxXqKT2JKfos06MWClDLGlCEVFoCdyAUbyYx4HCKrG1m9pnsrEmVBvxqFapxlvQUOKapXHlELXpmIj8Y1Hn4AgZFq5Wo8sGp5IOuMfZRxRZ2qLxNXDZuitt0iPZZRWdlMLkaFiRTlqdiIjeYg32762qqtqj7CSmvzgNZRQsfadM8YjsFuFeAnfw","dq":"VKSdRCBI7QTFu9ZJpN6jn9HsTeoJgLVehDCOpIV3-RJtiHLhKtPpsIXSoA_wQEIIN0eXz2_S8_rsHsaE3G-Sg3VvbkONgBYP5ym7Y-x1aev-4HXVFtHHBZqrrb9TSkysLEH56Z1-JhrqgAF-aFWh8mYO4ZWN91vJ6kJY_q34Ri2bekOxoMa66AnzFjW78LB8YbKicX7hQbV4k7TPnayFyLUfycC6N7zwPmahQDJI5mfGB16Grb3PPrEtiX6OUoaS28hnnynYHK2vBDfl0XWsrH9X3WxdxHh-UdVXpiWYGGjxY8OqaAW_apaLSQQEPdQIEuG-jHByCJ-mkWz5-7E-XQ","qi":"ng5p_Pvv0yqThHgNEheE1y_L358UFri6mQo6RfVbzaKyxE-x8zXef_IIvIObW45vyleCom3RhGgAnpRdJVG8PMb2ZGiShD3uF9VrJeWosSrRnowbMvakpGTeGD-V3wP8J3iA6WWhj7ExUvrfonBhOkJ8KOsvp1s9idwgXTscXq_TK1gOnTKdKYsEAiPbQkEBlyZX5T84pxJ1crBjBsI769V5_tclX5PWnVeiWzwOBR8A8-f_COtQ8JoBGscGc2F5wM6WtLf8Vm5JJsCXnMuOPfLPAdqyhhknmYZ3HWUOAbJP4h6NOLwn0_HyUUqmFQbshHp9sznyn2KNtCZjW2s72Q"}`
	// // addr01 := "Fkj5J8CDLC9Jif4CzgtbiXJBnwXLSrp5AaIllleH_yY"
	// signer, err := goar.NewSigner([]byte(priv01))
	// assert.NoError(t, err)
	// testSDK, err := New(signer, "https://api-dev.everpay.io")
	// assert.NoError(t, err)
	// to := "0xa2026731B31E4DFBa78314bDBfBFDC8cF5F761F8"
	// amount := big.NewInt(100000)
	// result, err := testSDK.Transfer("usdt", amount, to, `{"msg": "hello"}`)
	// assert.NoError(t, err)
	// t.Log(result.HexHash())
}

func TestClient_SubmitBundleTx(t *testing.T) {
	// addr01 := "0xa06b79E655Db7D7C3B3E7B2ccEEb068c3259d0C9"
	// priv01 := "1a7ffbdae668acf43251ed8913596f7db0ce0f90bcd27d4aa85b2bd8a3d0c550"
	// priv02 := "53e11a3eeb52f6105ce81638d8460feb5cc0a8bcdec54b49eb2cc8adffa84c27"
	// addr02 := "0x3314183F9F3CAcf8e4915dA59f754568345aF4D3"
	//
	// addr03 := "cSYOy8-p1QFenktkDBFyRM3cwZSTrQ_J4EsELLho_UE"
	//
	// items := []paySchema.BundleItem{
	// 	{
	// 		Tag:     "ethereum-eth-0x0000000000000000000000000000000000000000",
	// 		ChainID: "42",
	// 		From:    addr01,
	// 		To:      addr02,
	// 		Amount:  "99999",
	// 	},
	// 	{
	// 		Tag:     "ethereum-eth-0x0000000000000000000000000000000000000000",
	// 		ChainID: "42",
	// 		From:    addr01,
	// 		To:      addr03,
	// 		Amount:  "888888",
	// 	},
	// 	{
	// 		Tag:     "ethereum-usdt-0xd85476c906b5301e8e9eb58d174a6f96b9dfc5ee",
	// 		ChainID: "42",
	// 		From:    addr01,
	// 		To:      addr03,
	// 		Amount:  "12345",
	// 	},
	// 	{
	// 		Tag:     "ethereum-eth-0x0000000000000000000000000000000000000000",
	// 		ChainID: "42",
	// 		From:    addr02,
	// 		To:      addr03,
	// 		Amount:  "6666",
	// 	},
	// }
	//
	// txNonce := time.Now().UnixNano() / 1e6
	// expiration :=  txNonce/1000 + 1000
	// bundle := GenBundle(items, expiration)
	//
	// signer01 , _ := goether.NewSigner(priv01)
	// signer02 , _ := goether.NewSigner(priv02)
	// sdk01 ,err := New(signer01,"https://api-dev.everpay.io")
	// assert.NoError(t, err)
	// sdk02 ,err:= New(signer02,"https://api-dev.everpay.io")
	// assert.NoError(t, err)
	//
	// bundleData01, err := sdk01.SignBundleData(bundle)
	// assert.NoError(t, err)
	// bundleData02, err := sdk02.SignBundleData(bundle)
	// assert.NoError(t, err)
	//
	// bundleSigs := paySchema.BundleWithSigs{
	// 	Bundle: bundle,
	// 	Sigs: map[string]string{
	// 		sdk01.AccId: bundleData01.Sigs[sdk01.AccId],
	// 		sdk02.AccId: bundleData02.Sigs[sdk02.AccId],
	// 	},
	// }
	//
	// res, err := sdk01.Bundle("ETH",addr01,nil,bundleSigs)
	// assert.NoError(t, err)
	// t.Log(res.HexHash())
}

func TestClient_SubmitBundleTxSignByRSA(t *testing.T) {
	// priv01 := `{"kty":"RSA","e":"AQAB","n":"yB4ENow16LnZuqOrxBhtfCgj53JDSxMq57q4bj49dmAbdy1RzM9Hm6jWsGy2oIySjbd-kBn0CKhPLz_49QRZ6CDexSODMrmu4W-MioipIJwWh_TFMyzCz8F1dT0Falz3tvKX0g9Lp4LU5TNZVhykHImpRSkVOKwEp2BD-vKGEZrUFuA3gbYpuStz7tccu5sRWX3IV5W7yidcEPjkxcjLOmsl2iVVlmfvJNU5olX5Uohr7rCIaeV3X2uBKR07OAGTxlZhtlIQfcEbzeFKxZavtjYub7W-xNid5yJB1zhFfxA2sYhFdm-zoOyKeH5RqHnzJvHehjYF8dFrpjIqJGasbbvT9djMKju3RUD1p2m7aCqTSgw4ewsgJxZK3xHndldAPJODddLivUff-pzf3zXXHHTSRjZdsIAn98p2r9aVqKiTXyjSSK8ZoQIUubHdSRsj3lgKmI1OY8SAeiXVGQ-dM0WfRzX2MA2C_e5uJG-9PnMx8SMlSY7WINZqvsnj6xNyT3VxnfGV-lPr6_S8_G6oqAuXL13YuEccJE9sEL-pcuWzMqdHsdMOzrcYdY7djbUrltZv4gwoWWRjNxO_Ib4SQQJ6lVZ4dF42HYSonX6f9DY7505xAxLQ4TaQ960PHWeRISSxXXl0bLsdwtqUnInacS_wCBTnMqOfRIKMuYh4Eek","d":"W331uSazgvpthDpOTrJIZCpJ-AgyogFabWXd6nFEUnmwNPWgX34pJA62USTy_kYmvuo00Bc9jjO9OEIMRaE-pbVZhFIvKqLgRKtmWvQONx2F45TSEAqX6rptRDOo5S8qBhe1t9FMaoKTaVLntkoT8gQyQq98x1NGMWv-yJ3XE_dwniDW39P0EXe_YmQw-8jyzrKGRRztJYiHPC8_EPoziDryrGDDV6efG3ztA75A5i3_ArrtM1U5FYQmIjD-ISX_dP-LXvr4FTaTyrYFJH83SJxO8L3sQT00PDdaWNaeV3CHQPWnpHV_b8SVrm86UEepFpYggcXzQ9LWX76xq8DRde76UuYyVSlMhm5l9LaESXBSdx3ldRCVLLzAffPrdiy4gta9PTq14z3UHK1bmI_ILE8COHAttES0lPvjE0movJRou5JPMCp2zZX2H5Dy7PDWn_3cbnstzOGWRtQcI-hv13B1Nk95l3KGYJx4Nncmce4eDWUP4RBTg2SFpY2tY6opGlhK0Th5M2pqCzmFhgqYc0-k4jRpvOrl7_afFWrFo88-XWTsywP7mflGvEEh0D6QK1nwosfx5gTHfZe0g-DJHwU0F_fMZhhmv0YRjEO_xNj9Ugva3CB1c68rxdRuT5hCqh_dZi5L82En0KN69KN6iuVCtn12KJ2jp2TqsoDco50","p":"60H6zepl0rZ2Byt__8r2Rr7ZMW7gGrWEf3RfuycpSggXMs_2JA2ky3GyjQqnt4E2FowO3RnACNG3ftK7oc1d4DdipGi9Y0mCbK76TpCnCQA2_Cc4RpKmFjdAHiCmi33yOZBnPOE8was8G7RLPXkX-UXyVZ95e1YFZ3uTBylVqRs0rwuIpjx2EJnFBbRR7bdWYLoy1wBdVjy3pq1DqXpR4zjfSw0jBbAfGeRP6WFY9Ana9JzKLPozQOULfQobPtPhldny3y53q6-qD8UINl8kvBx6dNV0jnvrv0dXQmq28KO7vld6hotqY3zz94GRUky6eJpcosYK7VBWtJ7zQNQO5w","q":"2cLg1YGg561hgYmsI4mpmW4-bPWu8tLpCueHIsrX41KnIaEwhA-Pz-sg9glehveVU_3R9rXrVhpADV9JDJDxYNswDbJ_BYXWsF1I4A-wy0EmDUW9HglYKvLEiA9JHbnW4w8duOJPv50ej0PbrfTjW_4p7DofH0_JEfNzZfORf2LDPVqLaS5nvP-wkhZU14FEVpOHC8VrUGcrr9G5uAwYF0N6IT-xdUTMZkCDmJD25vGDVn3Ca2ne9mB6Jd4H7gpijHGiIjkdJhC4IBElDlNPtszFf-6WESWiFL0bByXJ1VxK0LfEJGu5fxCkALdvd5I9IUgwim5ofCXbtlARc9zOrw","dp":"MwUzJFcyT_lcLX_kmY_lyz_2kH7wQTqeipmtbUQ38yNADLgHNJh05d45j9cAuo5eMZOcYZ0FBaSdu_Zt1NBaDJdWYDuK3y6BB5sHE-TRJOWUBylmTf_d3zGKST5hgB0HXC6SBST_7Sx6s7NtId7SMTBXNvSH_xSPxXqKT2JKfos06MWClDLGlCEVFoCdyAUbyYx4HCKrG1m9pnsrEmVBvxqFapxlvQUOKapXHlELXpmIj8Y1Hn4AgZFq5Wo8sGp5IOuMfZRxRZ2qLxNXDZuitt0iPZZRWdlMLkaFiRTlqdiIjeYg32762qqtqj7CSmvzgNZRQsfadM8YjsFuFeAnfw","dq":"VKSdRCBI7QTFu9ZJpN6jn9HsTeoJgLVehDCOpIV3-RJtiHLhKtPpsIXSoA_wQEIIN0eXz2_S8_rsHsaE3G-Sg3VvbkONgBYP5ym7Y-x1aev-4HXVFtHHBZqrrb9TSkysLEH56Z1-JhrqgAF-aFWh8mYO4ZWN91vJ6kJY_q34Ri2bekOxoMa66AnzFjW78LB8YbKicX7hQbV4k7TPnayFyLUfycC6N7zwPmahQDJI5mfGB16Grb3PPrEtiX6OUoaS28hnnynYHK2vBDfl0XWsrH9X3WxdxHh-UdVXpiWYGGjxY8OqaAW_apaLSQQEPdQIEuG-jHByCJ-mkWz5-7E-XQ","qi":"ng5p_Pvv0yqThHgNEheE1y_L358UFri6mQo6RfVbzaKyxE-x8zXef_IIvIObW45vyleCom3RhGgAnpRdJVG8PMb2ZGiShD3uF9VrJeWosSrRnowbMvakpGTeGD-V3wP8J3iA6WWhj7ExUvrfonBhOkJ8KOsvp1s9idwgXTscXq_TK1gOnTKdKYsEAiPbQkEBlyZX5T84pxJ1crBjBsI769V5_tclX5PWnVeiWzwOBR8A8-f_COtQ8JoBGscGc2F5wM6WtLf8Vm5JJsCXnMuOPfLPAdqyhhknmYZ3HWUOAbJP4h6NOLwn0_HyUUqmFQbshHp9sznyn2KNtCZjW2s72Q"}`
	// addr01 := "Fkj5J8CDLC9Jif4CzgtbiXJBnwXLSrp5AaIllleH_yY"
	//
	// priv02 := "53e11a3eeb52f6105ce81638d8460feb5cc0a8bcdec54b49eb2cc8adffa84c27"
	// addr02 := "0x3314183F9F3CAcf8e4915dA59f754568345aF4D3"
	//
	// items := []paySchema.BundleItem{
	// 	{
	// 		Tag:     "ethereum-eth-0x0000000000000000000000000000000000000000",
	// 		ChainID: "42",
	// 		From:    addr02,
	// 		To:      addr01,
	// 		Amount:  "999",
	// 	},
	// 	{
	// 		Tag:     "ethereum-usdt-0xd85476c906b5301e8e9eb58d174a6f96b9dfc5ee",
	// 		ChainID: "42",
	// 		From:    addr01,
	// 		To:      addr02,
	// 		Amount:  "666",
	// 	},
	// }
	//
	// txNonce := time.Now().UnixNano() / 1e6
	// expiration :=  txNonce/1000 + 1000
	// bundle := GenBundle(items, expiration)
	//
	// signer01, err := goar.NewSigner([]byte(priv01))
	// assert.NoError(t, err)
	// signer02, err := goether.NewSigner(priv02)
	// assert.NoError(t, err)
	// sdk01, err  := New(signer01,"https://api-dev.everpay.io")
	// assert.NoError(t, err)
	// sdk02,err := New(signer02,"https://api-dev.everpay.io")
	// assert.NoError(t, err)
	//
	//
	// bundleData01, err := sdk01.SignBundleData(bundle)
	// assert.NoError(t, err)
	// bundleData02, err := sdk02.SignBundleData(bundle)
	// assert.NoError(t, err)
	//
	// bundleSigs := paySchema.BundleWithSigs{
	// 	Bundle: bundle,
	// 	Sigs: map[string]string{
	// 		sdk01.AccId: bundleData01.Sigs[sdk01.AccId],
	// 		sdk02.AccId: bundleData02.Sigs[sdk02.AccId],
	// 	},
	// }
	//
	// res, err := sdk01.Bundle("usdt",addr01,nil,bundleSigs)
	// assert.NoError(t, err)
	// t.Log(res.HexHash())
}

func Test_TransferArWallet(t *testing.T) {
	// payUrl := "https://api-dev.everpay.io"
	// signer, err := goar.NewSignerFromPath("../test-keyfile-watchmen.json")
	// if err != nil {
	// 	panic(err)
	// }
	// t.Log(signer.Address)
	// sdk , err := New(signer, payUrl)
	// assert.NoError(t, err)
	//
	// symbol := "AR"
	// to := "0x4002ED1a1410aF1b4930cF6c479ae373dEbD6223"
	// amount := big.NewInt(100000)
	// data := "sandy test sdk transfer rsa Sign"
	// res , err := sdk.Transfer(symbol, amount,to, data)
	// assert.NoError(t, err)
	// t.Log(res.HexHash())
}

func TestNew(t *testing.T) {
	// signer, err := goether.NewSigner("1a7ffbdae668acf43251ed8913596f7db0ce0f90bcd27d4aa85b2bd8a3d0c550")
	// if err != nil {
	// 	panic(err)
	// }
	// t.Log(signer.Address.String())
	// t.Log(951513707640286992 - 4550000000000000)
}

func TestBurnTx(t *testing.T) {
	// payUrl := "https://api-dev.everpay.io"
	// signer, err :=goether.NewSigner("1a7ffbdae668acf43251ed8913596f7db0ce0f90bcd27d4aa85b2bd8a3d0c550")
	// if err != nil {
	// 	panic(err)
	// }
	// t.Log(signer.Address)
	// sdk, err := New(signer, payUrl)
	// assert.NoError(t, err)
	//
	// tokenSymbol := "USDT"
	// targetChain := "ethereum"
	// to := "0x4002ED1a1410aF1b4930cF6c479ae373dEbD6223"
	// amount := big.NewInt(9000000)
	// res , err := sdk.Withdraw(tokenSymbol,amount,targetChain,to)
	// assert.NoError(t, err)
	// t.Log(res.HexHash())
}
