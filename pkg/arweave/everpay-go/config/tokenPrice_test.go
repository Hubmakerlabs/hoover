package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetTokenPriceByRedstone(t *testing.T) {
	price, err := GetTokenPriceByRedstone("VRT", "AR")
	assert.NoError(t, err)
	t.Log(price)
	price, err = GetTokenPriceByRedstone("VRT", "ETH")
	assert.NoError(t, err)
	t.Log(price)

	price, err = GetTokenPriceByRedstone("xyz", "AR")
	assert.NoError(t, err)
	t.Log(price)
	price, err = GetTokenPriceByRedstone("ardrive", "AR")
	assert.NoError(t, err)
	t.Log(price)
	price, err = GetTokenPriceByRedstone("pia", "AR")
	assert.NoError(t, err)
	t.Log(price)
	price, err = GetTokenPriceByRedstone("dlt", "AR")
	assert.NoError(t, err)
	t.Log(price)
}
