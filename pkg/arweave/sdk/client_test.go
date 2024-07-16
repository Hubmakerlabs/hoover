package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew2(t *testing.T) {
	cli := New("https://seed-dev.everpay.io")
	bundler, err := cli.GetBundler()
	assert.NoError(t, err)
	t.Log(bundler)
}
