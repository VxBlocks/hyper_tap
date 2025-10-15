package hyperliquid

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/tidwall/gjson"
)

type MetaAndAssetCtx struct {
	Meta     TokenMeta
	AssetCtx []AssetCtx
}

type AssetCtx struct {
	Funding      string   `json:"funding"`
	OpenInterest string   `json:"openInterest"`
	PrevDayPx    string   `json:"prevDayPx"`
	DayNtlVlm    string   `json:"dayNtlVlm"`
	Premium      *string  `json:"premium"`
	OraclePx     string   `json:"oraclePx"`
	MarkPx       string   `json:"markPx"`
	MidPx        *string  `json:"midPx"`
	ImpactPxs    []string `json:"impactPxs"`
	DayBaseVlm   string   `json:"dayBaseVlm"`
}

type TokenMeta struct {
	Universe        []Universe `json:"universe"`
	CollateralToken int64      `json:"collateralToken"`
}

type Universe struct {
	SzDecimals    int64  `json:"szDecimals"`
	Name          string `json:"name"`
	MaxLeverage   int64  `json:"maxLeverage"`
	MarginTableID int64  `json:"marginTableId"`
	IsDelisted    *bool  `json:"isDelisted,omitempty"`
	OnlyIsolated  *bool  `json:"onlyIsolated,omitempty"`
}

type MetaAndAssetCtxUnion struct {
	FluffyMetaAndAssetCtx      *TokenMeta
	PurpleMetaAndAssetCtxArray []AssetCtx
}

func GetMetaAndAssetCtx(ctx context.Context) (*MetaAndAssetCtx, error) {

	requestBody := map[string]string{
		"type": "metaAndAssetCtxs",
	}
	jsonstr, _ := json.Marshal(requestBody)
	bodyReader := bytes.NewReader(jsonstr)

	request, err := http.NewRequestWithContext(ctx, "POST",
		"https://api.hyperliquid.xyz/info", bodyReader)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: %s", string(body))
	}

	metaStr := gjson.Get(string(body), "0").Raw
	assetCtxStr := gjson.Get(string(body), "1").Raw
	var meta TokenMeta
	var assetCtx []AssetCtx
	err = json.Unmarshal([]byte(metaStr), &meta)
	if err != nil {
		return nil, fmt.Errorf("unmarshal meta %v", err)
	}
	err = json.Unmarshal([]byte(assetCtxStr), &assetCtx)
	if err != nil {
		return nil, fmt.Errorf("unmarshal assetCtx %v", err)
	}

	return &MetaAndAssetCtx{
		Meta:     meta,
		AssetCtx: assetCtx,
	}, err
}

type TokenMetaAndAssetCtx struct {
	Universe Universe
	AssetCtx AssetCtx
}

func (m MetaAndAssetCtx) ToTokenMetaAndAssetCtx() []TokenMetaAndAssetCtx {
	var res []TokenMetaAndAssetCtx
	for i := range m.Meta.Universe {
		res = append(res,
			TokenMetaAndAssetCtx{
				Universe: m.Meta.Universe[i],
				AssetCtx: m.AssetCtx[i],
			})
	}
	return res
}
