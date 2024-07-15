package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGasStation(t *testing.T) {
	// fastestETH, fastETH, averageETH, err := GetGasStation("https://data-api.defipulse.com/api/v1/egs/api/ethgasAPI.json?api-key=bd99f150dcc3fb0117cc50e650a048373e73d77b49ef0b58b3f25ecff3ca")
	// assert.NoError(t, err)
	// assert.NotEqual(t, 0.0, fastestETH)
	// assert.NotEqual(t, 0.0, fastETH)
	// assert.NotEqual(t, 0.0, averageETH)
	// t.Log(fastestETH, fastETH, averageETH)
}

func TestGetGasEtherscan(t *testing.T) {
	// fastestETH, fastETH, averageETH, err := GetGasEtherscan("https://api.etherscan.io/api?module=gastracker&action=gasoracle&apikey=MYDQVP2XTAE7DCPACT9NMZ9CQX558M2WWS")
	// assert.NoError(t, err)
	// assert.NotEqual(t, 0.0, fastestETH)
	// assert.NotEqual(t, 0.0, fastETH)
	// assert.NotEqual(t, 0.0, averageETH)
	// t.Log(fastestETH, fastETH, averageETH)
}

func TestGetGasBlockNative(t *testing.T) {
	// fastestETH, fastETH, averageETH, err := GetGasBlockNative("https://api.blocknative.com/gasprices/blockprices")
	// assert.NoError(t, err)
	// assert.NotEqual(t, 0.0, fastestETH)
	// assert.NotEqual(t, 0.0, fastETH)
	// assert.NotEqual(t, 0.0, averageETH)
	// t.Log(fastestETH, fastETH, averageETH)
}

func TestNewGas(t *testing.T) {
	// gas := NewGas(map[string]func(url string) (fastestETH float64, fastETH float64, averageETH float64, err error){
	// 	"https://api.etherscan.io/api?module=gastracker&action=gasoracle&apikey=MYDQVP2XTAE7DCPACT9NMZ9CQX558M2WWS": GetGasEtherscan,
	// 	"https://api.blocknative.com/gasprices/blockprices":                                                         GetGasBlockNative,
	// })
	//
	// gas.Run()
	//
	// fastestETH, fastETH, averageETH, source, err := gas.GetGasPrice()
	// assert.NoError(t, err)
	// assert.NotEqual(t, 0.0, fastestETH)
	// assert.NotEqual(t, 0.0, fastETH)
	// assert.NotEqual(t, 0.0, averageETH)
	// t.Log(fastestETH, fastETH, averageETH, source)
}

func TestGas_Average(t *testing.T) {
	rpc := "https://rpc.api.moonbeam.network"
	a, _, _, err := GetGasNode(rpc)
	assert.NoError(t, err)
	t.Log(a)
}
