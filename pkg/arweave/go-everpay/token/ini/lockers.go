package ini

import (
	"github.com/Hubmakerlabs/hoover/pkg/arweave/go-everpay/token/schema"

	"strconv"
)

func InitLockers(everChainId int) schema.Lockers {
	lockers := make(schema.Lockers, 0)
	// initArLockers
	lockers[schema.OracleArweaveChainType] = InitArLocker(everChainId)

	// initEThLockers, now ethereum only one HeightLocker
	lockers[schema.OracleEthChainType] = InitEthLocker(everChainId)

	// initMoonbeamLockers, only one HeightLocker
	lockers[schema.OracleMoonChainType] = InitMoonLocker(everChainId)

	// initConfluxLocker, only one HeightLocker
	lockers[schema.OracleCfxChainType] = InitCfxLocker(everChainId)

	// todo if add new chain, then need to add lockers here ...

	return lockers
}

func InitArLocker(everChainId int) schema.Locker {
	lk := schema.Locker{ChainType: schema.OracleArweaveChainType}
	switch strconv.Itoa(everChainId) {
	case everMainChainId:
		lk.ChainId = arweaveMainChainId
		lk.List = []schema.HeightLocker{
			{"uGx-QfBXSwABKxjha-00dI7vvfyqIYblY6Z5L6cyTFM", 0,
				1}, // address for recharge ar_locker fees, can not used minted
			{"8-GZSKB8VisgQltIiBmHiKYlTzVD46rCOlewOSOzbUA", 2, 726024},
			{"dH-_dwLlN86fitrFZzi86IVEEQFyYpTzWcqnFh460ys", 726025, 9999999999},
		}

	case everDevChainId:
		lk.ChainId = arweaveDevChainId
		lk.List = []schema.HeightLocker{
			{"xuCOkXMLtMMtMWwr8qJqIfp77nlLi8LdNbLLT75hVkU", 0, 717750},
			{"bX7sKd1s8L6PxUHxK-UPCfus7duyVFdf0J1lm90zehc", 717751, 725968},
			{"FyINHRSrHW0teUhvJzd6R33Tl50qxLnSj8LJCP5puiI", 725969, 9999999999},
		}

	default:
		panic("not support ethChainID")
	}

	return lk
}

// InitEthLocker notice: ethLockerList only one HeightLocker address
func InitEthLocker(everChainId int) schema.Locker {
	lk := schema.Locker{ChainType: schema.OracleEthChainType}

	switch strconv.Itoa(everChainId) {
	case everMainChainId:
		lk.ChainId = ethMainChainId
		lk.List = []schema.HeightLocker{
			{"0x38741a69785e84399fcf7c5ad61d572f7ecb1dab", 0, 999999999999},
		}

	case everDevChainId:
		lk.ChainId = ethKovanChainId
		lk.List = []schema.HeightLocker{
			{"0xa7ae99c13d82dd32fc6445ec09e38d197335f38a", 0, 999999999999},
		}

	default:
		panic("not support ethChainID")
	}
	// lockerList length must be 1
	if len(lk.List) != 1 {
		panic("ethLocker must only one address")
	}
	return lk
}

func InitMoonLocker(everChainId int) schema.Locker {
	lk := schema.Locker{ChainType: schema.OracleMoonChainType}

	switch strconv.Itoa(everChainId) {
	case everMainChainId:
		lk.ChainId = moonMainChainId
		lk.List = []schema.HeightLocker{
			{"0x93b2c8834264e9e88bf49467ae6cbe9ebee2a880", 0, 999999999999999},
		}

	case everDevChainId:
		lk.ChainId = moonDevChainId
		lk.List = []schema.HeightLocker{
			{"0xb3f2f559fe40c1f1ea1e941e982d9467208e17ae", 0, 999999999999999},
		}

	default:
		panic("not support ethChainID")
	}

	// lockerList length must be 1
	if len(lk.List) != 1 {
		panic("moonbeamLocker must only one address")
	}
	return lk
}

func InitCfxLocker(everChainId int) schema.Locker {
	lk := schema.Locker{ChainType: schema.OracleCfxChainType}

	switch strconv.Itoa(everChainId) {
	case everMainChainId:
		lk.ChainId = cfxMainChainId
		lk.List = []schema.HeightLocker{
			{"0xc68370c007cab6f0698ebb6a8da93d40c43ada5a", 0, 999999999999999},
		}

	case everDevChainId:
		lk.ChainId = cfxDevChainId
		lk.List = []schema.HeightLocker{
			{"0x7e6ef86b86141e82aed9c16eeda642ab41e647a9", 0, 999999999999999},
		}
	}

	// lockerList length must be 1
	if len(lk.List) != 1 {
		panic("conflux Locker must only one address")
	}
	return lk
}
