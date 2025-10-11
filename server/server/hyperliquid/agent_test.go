package hyperliquid_test

import (
	"context"
	"hyperliquid-server/hyperliquid"
	"testing"

	"github.com/ysmood/got"
)

func TestExtraAgents(t *testing.T) {
	g := got.New(t)
	got, gotErr := hyperliquid.ExtraAgents(context.Background(), "0x1755e9a4e305f8528f0b0705fd4d0e0c860b5fB8")
	if gotErr != nil {
		g.Error(gotErr)
		return
	}
	g.Log(g.ToJSONString(got))
}
