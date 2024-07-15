package common

import (
	"fmt"
	"io/ioutil"
	"net/http"
	url2 "net/url"
	"strconv"
	"sync"
	"time"

	"github.com/Hubmakerlabs/hoover/pkg/arweave/ethrpc"
	tokSchema "github.com/Hubmakerlabs/hoover/pkg/arweave/go-everpay/token/schema"
	"github.com/getsentry/sentry-go"

	"github.com/tidwall/gjson"
)

var (
	EvmGas = func(chainType string) *Gas {
		switch chainType { // todo more chain
		case tokSchema.OracleEthChainType:
			return NewGas(map[string]func(url string) (fastestETH float64,
				fastETH float64, averageETH float64, err error){
				"https://api.etherscan.io/api?module=gastracker&action=gasoracle&apikey=MYDQVP2XTAE7DCPACT9NMZ9CQX558M2WWS": GetGasEtherscan,
				"https://api.blocknative.com/gasprices/blockprices":                                                         GetGasBlockNative,
			})
		case tokSchema.OracleMoonChainType:
			return NewGas(map[string]func(url string) (fastestETH float64,
				fastETH float64, averageETH float64, err error){
				"https://rpc.api.moonbeam.network": GetGasNode,
			})
		case tokSchema.OracleCfxChainType:
			return NewGas(map[string]func(url string) (fastestETH float64,
				fastETH float64, averageETH float64, err error){
				"https://evm.confluxrpc.com": GetGasNode,
			})
		}
		return nil
	}
)

type Gas struct {
	gasFuncMap map[string]func(url string) (fastestETH, fastETH, averageETH float64,
		err error) // key: url, val: get gasPrice func
	fastestGasPrice float64 // ETH uint
	fastGasPrice    float64
	averageGasPrice float64
	source          string

	gasPriceMux sync.RWMutex
}

func NewGas(gasFuncMap map[string]func(url string) (fastestETH, fastETH, averageETH float64,
	err error)) *Gas {
	return &Gas{
		gasFuncMap: gasFuncMap,
	}
}

func (g *Gas) Run() {
	if err := g.setGasPrice(); err != nil {
		panic(err)
	}

	go func() {
		for {
			time.Sleep(5 * time.Second)
			if err := g.setGasPrice(); err != nil {
				sentry.CaptureException(fmt.Errorf("get gas price err: %v",
					err))
			}
		}
	}()
}

func (g *Gas) Fastest() float64 {
	g.gasPriceMux.RLock()
	defer g.gasPriceMux.RUnlock()

	return g.fastestGasPrice
}

func (g *Gas) Fast() float64 {
	g.gasPriceMux.RLock()
	defer g.gasPriceMux.RUnlock()

	return g.fastGasPrice
}

func (g *Gas) Average() float64 {
	g.gasPriceMux.RLock()
	defer g.gasPriceMux.RUnlock()

	return g.averageGasPrice
}

func (g *Gas) Source() string {
	g.gasPriceMux.RLock()
	defer g.gasPriceMux.RUnlock()

	return g.source
}

func (g *Gas) GetGasPrice() (fastestETH, fastETH, averageETH float64,
	source string, err error) {
	for url, getFunc := range g.gasFuncMap {
		fastestETH, fastETH, averageETH, err = getFunc(url)
		if err == nil && fastestETH > 0 {
			u, errs := url2.Parse(url)
			if errs == nil {
				source = u.Host
			}
			return
		}
	}
	return
}

func (g *Gas) setGasPrice() error {
	fastestETH, fastETH, averageETH, source, err := g.GetGasPrice()
	if err != nil {
		return err
	}

	g.gasPriceMux.Lock()
	defer g.gasPriceMux.Unlock()

	g.fastestGasPrice = fastestETH
	g.fastGasPrice = fastETH
	g.averageGasPrice = averageETH
	g.source = source

	return nil
}

func GetGasEtherscan(url string) (fastestETH, fastETH, averageETH float64,
	err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	httpClient := http.DefaultClient
	res, err := httpClient.Do(req)
	if err != nil {
		return
	}
	if res.StatusCode != 200 {
		return
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	str := string(body)

	fastest := gjson.Get(str, "result.FastGasPrice").Float()
	fast := gjson.Get(str, "result.ProposeGasPrice").Float()
	average := gjson.Get(str, "result.ProposeGasPrice").Float()

	fastestETH = fastest / 1e9
	fastETH = fast / 1e9
	averageETH = average / 1e9
	return
}

func GetGasBlockNative(url string) (fastestETH, fastETH, averageETH float64,
	err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Add("Authorization",
		"73a999cc-b2c7-45c1-83e8-ee71a4605b61") // TODO blocknative key

	httpClient := http.DefaultClient
	res, err := httpClient.Do(req)
	if err != nil {
		return
	}
	if res.StatusCode != 200 {
		return
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	str := string(body)
	blockPrices := gjson.Get(str, "blockPrices").Array()
	if len(blockPrices) == 0 {
		return
	}
	estimatedPrices := blockPrices[0].Get("estimatedPrices").Array()
	if len(estimatedPrices) < 3 {
		return
	}

	fastest := estimatedPrices[0].Get("price").Float()
	fast := estimatedPrices[1].Get("price").Float()
	average := estimatedPrices[2].Get("price").Float()

	fastestETH = fastest / 1e9
	fastETH = fast / 1e9
	averageETH = average / 1e9
	return
}

func GetGasNode(rpc string) (fastestETH, fastETH, averageETH float64,
	err error) {
	price, err := ethrpc.New(rpc).EthGasPrice()
	if err != nil {
		return
	}
	fprice, err := strconv.ParseFloat(price.String(), 64)
	if err != nil {
		return
	}
	fastestETH = fprice / 1e18
	fastETH = fprice / 1e18
	averageETH = fprice / 1e18
	return
}
