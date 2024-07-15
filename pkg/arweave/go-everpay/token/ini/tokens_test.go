package ini

import (
	"testing"
)

func TestInitToken(t *testing.T) {
	InitToken(1, "", "", "", "", "")
	InitToken(42, "", "", "", "", "")
	InitTokenWithoutRpc(1)
	InitTokenWithoutRpc(42)
}
