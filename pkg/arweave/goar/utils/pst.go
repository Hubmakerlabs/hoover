package utils

import "github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"

func PstTransferTags(contractId string, target string, qty int64,
	warp bool) ([]types.Tag, error) {
	input := types.Input{
		"function": "transfer",
		"target":   target,
		"qty":      qty,
	}

	inputStr, err := input.ToString()
	if err != nil {
		return nil, err
	}

	pstTags := []types.Tag{
		{Name: "App-Name", Value: "SmartWeaveAction"},
		{Name: "App-Version", Value: "0.3.0"},
		{Name: "Contract", Value: contractId},
		{Name: "Input", Value: inputStr},
	}

	if warp {
		pstTags = append(pstTags, types.Tag{Name: "SDK", Value: "Warp"})
	}

	return pstTags, nil
}
