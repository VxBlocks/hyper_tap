package handler

import (
	"context"
	"fmt"
	newsv1 "hyperliquid-server/gen/news/v1"
	"hyperliquid-server/gen/news/v1/newsv1connect"
	"hyperliquid-server/monitor"
	"timescale"

	"connectrpc.com/connect"
)

type NewsHandler struct {
}

var _ newsv1connect.NewsServiceHandler = (*NewsHandler)(nil)

func (n *NewsHandler) ListNews(ctx context.Context, req *connect.Request[newsv1.ListUserNewsRequest]) (*connect.Response[newsv1.ListUserNewsResponse], error) {
	var user_id string

	err := timescale.GetPostgresGormDB(ctx).Raw("select user_id_from_session(?)", req.Header().Get("Authorization")).Scan(&user_id).Error
	if err != nil {
		// return nil, fmt.Errorf("could not get user id from session: %w", err)
	}

	var msgs []*newsv1.News
	err = timescale.G[monitor.NewsOrm]().Select("*, news_is_read(?, uuid) as read", user_id).Order("created_at DESC").Limit(50).Scan(ctx, &msgs)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&newsv1.ListUserNewsResponse{
		Results: msgs,
	}), nil
}

func (n *NewsHandler) MarkNewsRead(ctx context.Context, req *connect.Request[newsv1.MarkNewsReadRequest]) (*connect.Response[newsv1.MarkNewsReadResponse], error) {
	var user_id string

	err := timescale.GetPostgresGormDB(ctx).Raw("select user_id_from_session(?)", req.Header().Get("Authorization")).Scan(&user_id).Error
	if err != nil {
		return nil, fmt.Errorf("could not get user id from session: %w", err)
	}

	err = timescale.GetPostgresGormDB(ctx).Save(&NewsReadOrm{
		NewsId: req.Msg.GetUuid(),
		UserId: user_id,
	}).Error

	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&newsv1.MarkNewsReadResponse{}), nil
}

type NewsReadOrm struct {
	NewsId string `gorm:"primaryKey"`
	UserId string `gorm:"index"`
}

func (NewsReadOrm) TableName() string { return "news_read" }
