package handler

import (
	"context"
	"fmt"
	pricealertv1 "hyperliquid-server/gen/price_alert/v1"
	"hyperliquid-server/gen/price_alert/v1/pricealertv1connect"
	"hyperliquid-server/monitor"
	"timescale"

	"connectrpc.com/connect"
)

type PriceAlertHandler struct {
}

// ListPriceAlerts implements pricealertv1connect.PriceAlertServiceHandler.
func (p *PriceAlertHandler) ListPriceAlerts(ctx context.Context, req *connect.Request[pricealertv1.ListPriceAlertRequest]) (*connect.Response[pricealertv1.ListPriceAlertResponse], error) {

	var user_id string

	err := timescale.GetPostgresGormDB(ctx).Raw("select user_id_from_session(?)", req.Header().Get("Authorization")).Scan(&user_id).Error
	if err != nil {
		// return nil, fmt.Errorf("could not get user id from session: %w", err)
	}

	var msgs []*pricealertv1.PriceAlert
	err = timescale.G[monitor.PriceAlert]().
		Select("id, asset_id, symbol, change_percent, created_at, price_alert_is_read(?, id) as read", user_id).
		Where("asset_id in ?", req.Msg.AssetIds).
		Order("created_at DESC").Limit(50).
		Scan(ctx, &msgs)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&pricealertv1.ListPriceAlertResponse{
		Results: msgs,
	}), nil
}

// MarkPriceAlertsRead implements pricealertv1connect.PriceAlertServiceHandler.
func (p *PriceAlertHandler) MarkPriceAlertsRead(ctx context.Context, req *connect.Request[pricealertv1.MarkPriceAlertReadRequest]) (*connect.Response[pricealertv1.MarkPriceAlertReadResponse], error) {
	var user_id string

	err := timescale.GetPostgresGormDB(ctx).Raw("select user_id_from_session(?)", req.Header().Get("Authorization")).Scan(&user_id).Error
	if err != nil {
		return nil, fmt.Errorf("could not get user id from session: %w", err)
	}

	err = timescale.GetPostgresGormDB(ctx).Save(&PriceAlertReadOrm{
		AlertId: req.Msg.GetAlertId(),
		UserId:  user_id,
	}).Error

	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&pricealertv1.MarkPriceAlertReadResponse{}), nil
}

var _ pricealertv1connect.PriceAlertServiceHandler = (*PriceAlertHandler)(nil)

type PriceAlertReadOrm struct {
	ID      uint   `gorm:"primarykey"`
	AlertId uint32 `gorm:"index"`
	UserId  string `gorm:"index"`
}

func (PriceAlertReadOrm) TableName() string { return "price_alert_read" }
