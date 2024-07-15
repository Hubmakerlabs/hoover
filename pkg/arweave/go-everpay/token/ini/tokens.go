package ini

import (
	"fmt"
	"strconv"

	"github.com/Hubmakerlabs/hoover/pkg/arweave/go-everpay/token"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/go-everpay/token/schema"
)

var (
	ethMainChainId  = "1"
	ethKovanChainId = "42"

	everMainChainId = ethMainChainId
	everDevChainId  = ethKovanChainId

	arweaveMainChainId = "0"  // arweave mainnet
	arweaveDevChainId  = "99" // todo To prevent replay of dev environment transactions, pst needs to be distinguished by the chainId

	moonMainChainId = "1284" // moonbeam chainId
	moonDevChainId  = "1287" // moonbase chainId

	cfxMainChainId = "1030"
	cfxDevChainId  = "71"
)

// InitToken init supported token list
func InitToken(everChainId int,
	ethRpcUrl, moonRpcUrl, cfxRpcUrl, arNodeUrl, pstGwUrl string) (tokens map[string]*token.Token) {
	// lockers
	ethLocker := InitEthLocker(everChainId)
	arLocker := InitArLocker(everChainId)
	moonLocker := InitMoonLocker(everChainId)
	cfxLocker := InitCfxLocker(everChainId)

	tokenList := make([]*token.Token, 0)
	switch strconv.Itoa(everChainId) {

	case everMainChainId:
		// arweave tokens
		ar := token.New(
			schema.ArAddress+",0x4fadc7a98f2dc96510e42dd1a74141eeae0c1543",
			"AR", schema.ChainTypeCrossArEth,
			fmt.Sprintf("%s,%s", arweaveMainChainId, ethMainChainId), 12,
			[]schema.TargetChain{
				{
					ChainId:   ethMainChainId,
					ChainType: schema.OracleEthChainType,
					Decimals:  12,
					TokenID:   "0x4fadc7a98f2dc96510e42dd1a74141eeae0c1543",
					Locker:    &ethLocker,
					Rpc:       ethRpcUrl,
				},
				{
					ChainId:   arweaveMainChainId,
					ChainType: schema.OracleArweaveChainType,
					Decimals:  12,
					TokenID:   schema.ArAddress,
					Locker:    &arLocker,
					Rpc:       arNodeUrl,
				},
				{ // AR cross to moon target info
					ChainId:   moonMainChainId,
					ChainType: schema.OracleMoonChainType,
					Decimals:  12,
					TokenID:   "0xAc42091313105104AC5a884c3c7c7e5a7EF9Ea38",
					Locker:    &moonLocker,
					Rpc:       moonRpcUrl,
				},
			},
		)
		vrt := token.New(
			"usjm4PCxUd5mtaon7zc97-dt-3qf67yPyqgzLnLqk5A", "VRT",
			schema.ChainTypeArweave, arweaveMainChainId, 18,
			[]schema.TargetChain{
				{
					ChainId:   arweaveMainChainId,
					ChainType: schema.OracleArweaveChainType,
					Decimals:  0,
					TokenID:   "usjm4PCxUd5mtaon7zc97-dt-3qf67yPyqgzLnLqk5A",
					Locker:    &arLocker,
					Rpc:       arNodeUrl,
					PstGw:     pstGwUrl,
				},
			},
		)
		ardrive := token.New(
			"-8A6RexFkpfWwuyVO98wzSFZh0d6VJuI-buTJvlwOJQ", "ARDRIVE",
			schema.ChainTypeArweave, arweaveMainChainId, 18,
			[]schema.TargetChain{
				{
					ChainId:   arweaveMainChainId,
					ChainType: schema.OracleArweaveChainType,
					Decimals:  0,
					TokenID:   "-8A6RexFkpfWwuyVO98wzSFZh0d6VJuI-buTJvlwOJQ",
					Locker:    &arLocker,
					Rpc:       arNodeUrl,
					PstGw:     pstGwUrl,
				},
				{ // ARDRIVE cross to moon target info
					ChainId:   moonMainChainId,
					ChainType: schema.OracleMoonChainType,
					Decimals:  18,
					TokenID:   "0x826DB9e588217c1ca1166fd24A491537511b966b",
					Locker:    &moonLocker,
					Rpc:       moonRpcUrl,
				},
			},
		)

		tokenList = append(tokenList, ar, vrt, ardrive)

		// ethereum tokens
		eth := token.New(
			schema.EthAddress, "ETH", schema.ChainTypeEth, ethMainChainId, 18,
			[]schema.TargetChain{
				{
					ChainId:   ethMainChainId,
					ChainType: schema.OracleEthChainType,
					Decimals:  18,
					TokenID:   schema.EthAddress,
					Locker:    &ethLocker,
					Rpc:       ethRpcUrl,
				},
			},
		)

		wbtc := token.New(
			"0x2260fac5e5542a773aa44fbcfedf7c193bc2c599", "WBTC",
			schema.ChainTypeEth, ethMainChainId, 8,
			[]schema.TargetChain{
				{
					ChainId:   ethMainChainId,
					ChainType: schema.OracleEthChainType,
					Decimals:  8,
					TokenID:   "0x2260fac5e5542a773aa44fbcfedf7c193bc2c599",
					Locker:    &ethLocker,
					Rpc:       ethRpcUrl,
				},
			},
		)
		usdt := token.New(
			"0xdac17f958d2ee523a2206206994597c13d831ec7", "USDT",
			schema.ChainTypeEth, ethMainChainId, 6,
			[]schema.TargetChain{
				{
					ChainId:   ethMainChainId,
					ChainType: schema.OracleEthChainType,
					Decimals:  6,
					TokenID:   "0xdac17f958d2ee523a2206206994597c13d831ec7",
					Locker:    &ethLocker,
					Rpc:       ethRpcUrl,
				},
			},
		)
		usdc := token.New(
			"0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", "USDC",
			schema.ChainTypeEth, ethMainChainId, 6,
			[]schema.TargetChain{
				{
					ChainId:   ethMainChainId,
					ChainType: schema.OracleEthChainType,
					Decimals:  6,
					TokenID:   "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
					Locker:    &ethLocker,
					Rpc:       ethRpcUrl,
				},
			},
		)
		dai := token.New(
			"0x6b175474e89094c44da98b954eedeac495271d0f", "DAI",
			schema.ChainTypeEth, ethMainChainId, 18,
			[]schema.TargetChain{
				{
					ChainId:   ethMainChainId,
					ChainType: schema.OracleEthChainType,
					Decimals:  18,
					TokenID:   "0x6b175474e89094c44da98b954eedeac495271d0f",
					Locker:    &ethLocker,
					Rpc:       ethRpcUrl,
				},
			},
		)
		uni := token.New(
			"0x1f9840a85d5af5bf1d1762f925bdaddc4201f984", "UNI",
			schema.ChainTypeEth, ethMainChainId, 18,
			[]schema.TargetChain{
				{
					ChainId:   ethMainChainId,
					ChainType: schema.OracleEthChainType,
					Decimals:  18,
					TokenID:   "0x1f9840a85d5af5bf1d1762f925bdaddc4201f984",
					Locker:    &ethLocker,
					Rpc:       ethRpcUrl,
				},
			},
		)
		sos := token.New(
			"0x3b484b82567a09e2588a13d54d032153f0c0aee0", "SOS",
			schema.ChainTypeEth, ethMainChainId, 18,
			[]schema.TargetChain{
				{
					ChainId:   ethMainChainId,
					ChainType: schema.OracleEthChainType,
					Decimals:  18,
					TokenID:   "0x3b484b82567a09e2588a13d54d032153f0c0aee0",
					Locker:    &ethLocker,
					Rpc:       ethRpcUrl,
				},
			},
		)
		bank := token.New(
			"0x2d94aa3e47d9d5024503ca8491fce9a2fb4da198", "BANK",
			schema.ChainTypeEth, ethMainChainId, 18,
			[]schema.TargetChain{
				{
					ChainId:   ethMainChainId,
					ChainType: schema.OracleEthChainType,
					Decimals:  18,
					TokenID:   "0x2d94aa3e47d9d5024503ca8491fce9a2fb4da198",
					Locker:    &ethLocker,
					Rpc:       ethRpcUrl,
				},
			},
		)
		dodo := token.New(
			"0x43dfc4159d86f3a37a5a4b3d4580b888ad7d4ddd", "DODO",
			schema.ChainTypeEth, ethMainChainId, 18,
			[]schema.TargetChain{
				{
					ChainId:   ethMainChainId,
					ChainType: schema.OracleEthChainType,
					Decimals:  18,
					TokenID:   "0x43dfc4159d86f3a37a5a4b3d4580b888ad7d4ddd",
					Locker:    &ethLocker,
					Rpc:       ethRpcUrl,
				},
			},
		)
		mask := token.New(
			"0x69af81e73a73b40adf4f3d4223cd9b1ece623074", "MASK",
			schema.ChainTypeEth, ethMainChainId, 18,
			[]schema.TargetChain{
				{
					ChainId:   ethMainChainId,
					ChainType: schema.OracleEthChainType,
					Decimals:  18,
					TokenID:   "0x69af81e73a73b40adf4f3d4223cd9b1ece623074",
					Locker:    &ethLocker,
					Rpc:       ethRpcUrl,
				},
			},
		)
		t4ever := token.New(
			"0xeaba187306335dd773ca8042b3792c46e213636a", "T4EVER",
			schema.ChainTypeEth, ethMainChainId, 18,
			[]schema.TargetChain{
				{
					ChainId:   ethMainChainId,
					ChainType: schema.OracleEthChainType,
					Decimals:  18,
					TokenID:   "0xeaba187306335dd773ca8042b3792c46e213636a",
					Locker:    &ethLocker,
					Rpc:       ethRpcUrl,
				},
			},
		)

		tokenList = append(tokenList, eth, usdc, usdt, wbtc, dai, uni, sos,
			bank, dodo, mask, t4ever)

		// moonbeam tokens
		glmr := token.New(
			schema.MoonAddress, "GLMR", schema.ChainTypeMoonbeam,
			moonMainChainId, 18,
			[]schema.TargetChain{
				{
					ChainId:   moonMainChainId,
					ChainType: schema.OracleMoonChainType,
					Decimals:  18,
					TokenID:   schema.MoonAddress,
					Locker:    &moonLocker,
					Rpc:       moonRpcUrl,
				},
			},
		)
		zlk := token.New(
			"0x3fd9b6c9a24e09f67b7b706d72864aebb439100c", "ZLK",
			schema.ChainTypeMoonbeam, moonMainChainId, 18,
			[]schema.TargetChain{
				{
					ChainId:   moonMainChainId,
					ChainType: schema.OracleMoonChainType,
					Decimals:  18,
					TokenID:   "0x3fd9b6c9a24e09f67b7b706d72864aebb439100c",
					Locker:    &moonLocker,
					Rpc:       moonRpcUrl,
				},
			},
		)

		tokenList = append(tokenList, glmr, zlk)

		// conflux tokens
		cfx := token.New(
			schema.CfxAddress, "CFX", schema.ChainTypeCfx, cfxMainChainId, 18,
			[]schema.TargetChain{
				{
					ChainId:   cfxMainChainId,
					ChainType: schema.OracleCfxChainType,
					Decimals:  18,
					TokenID:   schema.CfxAddress,
					Locker:    &cfxLocker,
					Rpc:       cfxRpcUrl,
				},
			},
		)

		tokenList = append(tokenList, cfx)

	case everDevChainId:
		// arweave tokens
		ar := token.New(
			schema.ArAddress+",0xcc9141efa8c20c7df0778748255b1487957811be",
			"AR", schema.ChainTypeCrossArEth,
			fmt.Sprintf("%s,%s", arweaveMainChainId, ethKovanChainId), 12,
			[]schema.TargetChain{
				{ // AR cross to ethereum target info
					ChainId:   ethKovanChainId,
					ChainType: schema.OracleEthChainType,
					Decimals:  12,
					TokenID:   "0xcc9141efa8c20c7df0778748255b1487957811be",
					Locker:    &ethLocker,
					Rpc:       ethRpcUrl,
				},
				{
					ChainId:   arweaveMainChainId,
					ChainType: schema.OracleArweaveChainType,
					Decimals:  12,
					TokenID:   schema.ArAddress,
					Locker:    &arLocker,
					Rpc:       arNodeUrl,
				},
				{ // AR cross to moon target info
					ChainId:   moonDevChainId,
					ChainType: schema.OracleMoonChainType,
					Decimals:  12,
					TokenID:   "0xc8F0B30449fBB398C48231072863522C2eF36a05",
					Locker:    &moonLocker,
					Rpc:       moonRpcUrl,
				},
			},
		)
		vrt := token.New(
			"usjm4PCxUd5mtaon7zc97-dt-3qf67yPyqgzLnLqk5A", "VRT",
			schema.ChainTypeArweave, arweaveDevChainId, 18,
			[]schema.TargetChain{
				{
					ChainId:   arweaveMainChainId,
					ChainType: schema.OracleArweaveChainType,
					Decimals:  0,
					TokenID:   "usjm4PCxUd5mtaon7zc97-dt-3qf67yPyqgzLnLqk5A",
					Locker:    &arLocker,
					Rpc:       arNodeUrl,
					PstGw:     pstGwUrl,
				},
				{
					ChainId:   ethKovanChainId,
					ChainType: schema.OracleEthChainType,
					Decimals:  18,
					TokenID:   "0xde10c3040aDB1e3d63Dd0ce7965192610aE36712",
					Locker:    &ethLocker,
					Rpc:       ethRpcUrl,
				},
				{
					ChainId:   moonDevChainId,
					ChainType: schema.OracleMoonChainType,
					Decimals:  18,
					TokenID:   "0xb5EadFdbDB40257D1d24A1432faa2503A867C270",
					Locker:    &moonLocker,
					Rpc:       moonRpcUrl,
				},
			},
		)
		ardrive := token.New(
			"-8A6RexFkpfWwuyVO98wzSFZh0d6VJuI-buTJvlwOJQ", "ARDRIVE",
			schema.ChainTypeArweave, arweaveDevChainId, 18,
			[]schema.TargetChain{
				{
					ChainId:   arweaveMainChainId,
					ChainType: schema.OracleArweaveChainType,
					Decimals:  0,
					TokenID:   "-8A6RexFkpfWwuyVO98wzSFZh0d6VJuI-buTJvlwOJQ",
					Locker:    &arLocker,
					Rpc:       arNodeUrl,
					PstGw:     pstGwUrl,
				},
				{
					ChainId:   moonDevChainId,
					ChainType: schema.OracleMoonChainType,
					Decimals:  18,
					TokenID:   "0x2044b3b09E03C7398749d774Ae4C7771260F29b2",
					Locker:    &moonLocker,
					Rpc:       moonRpcUrl,
				},
			},
		)
		xyz := token.New(
			"mzvUgNc8YFk0w5K5H7c8pyT-FC5Y_ba0r7_8766Kx74", "XYZ",
			schema.ChainTypeArweave, arweaveDevChainId, 18,
			[]schema.TargetChain{
				{
					ChainId:   arweaveMainChainId,
					ChainType: schema.OracleArweaveChainType,
					Decimals:  0,
					TokenID:   "mzvUgNc8YFk0w5K5H7c8pyT-FC5Y_ba0r7_8766Kx74",
					Locker:    &arLocker,
					Rpc:       arNodeUrl,
					PstGw:     pstGwUrl,
				},
			},
		)
		pia := token.New(
			"n05LTiuWcAYjizXAu-ghegaWjL89anZ6VdvuHcU6dno", "PIA",
			schema.ChainTypeArweave, arweaveDevChainId, 18,
			[]schema.TargetChain{
				{
					ChainId:   arweaveMainChainId,
					ChainType: schema.OracleArweaveChainType,
					Decimals:  0,
					TokenID:   "n05LTiuWcAYjizXAu-ghegaWjL89anZ6VdvuHcU6dno",
					Locker:    &arLocker,
					Rpc:       arNodeUrl,
					PstGw:     pstGwUrl,
				},
			},
		)
		tokenList = append(tokenList, ar, vrt, ardrive, xyz, pia)

		// ethereum tokens
		eth := token.New(
			schema.EthAddress, "ETH", schema.ChainTypeEth, ethKovanChainId, 18,
			[]schema.TargetChain{
				{
					ChainId:   ethKovanChainId,
					ChainType: schema.OracleEthChainType,
					Decimals:  18,
					TokenID:   schema.EthAddress,
					Locker:    &ethLocker,
					Rpc:       ethRpcUrl,
				},
			},
		)
		usdt := token.New(
			"0xd85476c906b5301e8e9eb58d174a6f96b9dfc5ee", "USDT",
			schema.ChainTypeEth, ethKovanChainId, 6,
			[]schema.TargetChain{
				{
					ChainId:   ethKovanChainId,
					ChainType: schema.OracleEthChainType,
					Decimals:  6,
					TokenID:   "0xd85476c906b5301e8e9eb58d174a6f96b9dfc5ee",
					Locker:    &ethLocker,
					Rpc:       ethRpcUrl,
				},
			},
		)
		usdc := token.New(
			"0xb7a4f3e9097c08da09517b5ab877f7a917224ede", "USDC",
			schema.ChainTypeEth, ethKovanChainId, 6,
			[]schema.TargetChain{
				{
					ChainId:   ethKovanChainId,
					ChainType: schema.OracleEthChainType,
					Decimals:  6,
					TokenID:   "0xb7a4f3e9097c08da09517b5ab877f7a917224ede",
					Locker:    &ethLocker,
					Rpc:       ethRpcUrl,
				},
			},
		)
		dai := token.New(
			"0xc4375b7de8af5a38a93548eb8453a498222c4ff2", "DAI",
			schema.ChainTypeEth, ethKovanChainId, 18,
			[]schema.TargetChain{
				{
					ChainId:   ethKovanChainId,
					ChainType: schema.OracleEthChainType,
					Decimals:  18,
					TokenID:   "0xc4375b7de8af5a38a93548eb8453a498222c4ff2",
					Locker:    &ethLocker,
					Rpc:       ethRpcUrl,
				},
			},
		)
		tokenList = append(tokenList, eth, usdt, usdc, dai)

		// moonbase tokens
		glmr := token.New(
			schema.MoonAddress, "DEV", schema.ChainTypeMoonbase, moonDevChainId,
			18,
			[]schema.TargetChain{
				{
					ChainId:   moonDevChainId,
					ChainType: schema.OracleMoonChainType,
					Decimals:  18,
					TokenID:   schema.MoonAddress,
					Locker:    &moonLocker,
					Rpc:       moonRpcUrl,
				},
			},
		)
		zlk := token.New(
			"0x322f069e9b8b554f3fb43cefcb0c7b3222242f0e", "ZLK",
			schema.ChainTypeMoonbase, moonDevChainId, 18,
			[]schema.TargetChain{
				{
					ChainId:   moonDevChainId,
					ChainType: schema.OracleMoonChainType,
					Decimals:  18,
					TokenID:   "0x322f069e9b8b554f3fb43cefcb0c7b3222242f0e",
					Locker:    &moonLocker,
					Rpc:       moonRpcUrl,
				},
			})
		tokenList = append(tokenList, glmr, zlk)

		// conflux tokens
		cfx := token.New(
			schema.CfxAddress, "CFX", schema.ChainTypeCfx, cfxDevChainId, 18,
			[]schema.TargetChain{
				{
					ChainId:   cfxDevChainId,
					ChainType: schema.OracleCfxChainType,
					Decimals:  18,
					TokenID:   schema.CfxAddress,
					Locker:    &cfxLocker,
					Rpc:       cfxRpcUrl,
				},
			},
		)

		tokenList = append(tokenList, cfx)

	default:
		panic(fmt.Sprintf("Not support this chainId: %d", everChainId))
	}

	tokens = make(map[string]*token.Token, len(tokenList))
	for _, t := range tokenList {
		tokens[t.Tag()] = t
	}
	return
}

// InitTokenWithoutRpc only use get token info
func InitTokenWithoutRpc(everChainId int) (tokens map[string]*token.Token) {
	return InitToken(everChainId, "", "", "", "", "")
}
