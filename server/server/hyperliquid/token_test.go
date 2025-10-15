package hyperliquid

import (
	"testing"

	"github.com/ysmood/got"
)

func TestGetMetaAndAssetCtx(t *testing.T) {
	g := got.New(t)
	res, err := GetMetaAndAssetCtx(g.Context())
	g.Must().Nil(err)
	assetCtx := res.ToTokenMetaAndAssetCtx()
	g.Log(g.ToJSONString(assetCtx))
}
