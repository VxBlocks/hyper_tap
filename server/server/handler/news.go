package handler

import (
	"context"
	"encoding/json"
	"fmt"
	newsv1 "hyperliquid-server/gen/news/v1"
	"hyperliquid-server/gen/news/v1/newsv1connect"
	"hyperliquid-server/hyperliquid"
	"hyperliquid-server/monitor"
	"timescale"

	"connectrpc.com/connect"
	"github.com/ejilay/draftjs"
)

type NewsHandler struct {
}

// GetNews implements newsv1connect.NewsServiceHandler.
func (n *NewsHandler) GetNews(ctx context.Context, req *connect.Request[newsv1.GetNewsRequest]) (*connect.Response[newsv1.GetNewsResponse], error) {
	state, err := hyperliquid.GetAnnouncement(ctx, req.Msg.Uuid)
	if err != nil {
		return nil, err
	}

	// prepare some config (HTML here)
	config := draftjs.NewDefaultConfig()

	var contentState draftjs.ContentState
	err = json.Unmarshal([]byte(state.Content), &contentState)
	if err != nil {
		return nil, err
	}

	// and just render content state to HTML-string
	s := draftjs.Render(&contentState, config)

	return connect.NewResponse(&newsv1.GetNewsResponse{
		Result: &newsv1.NewsContent{
			Uuid:    state.UUID,
			Title:   state.Title,
			Content: s,
		},
	}), nil
}

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

var _ newsv1connect.NewsServiceHandler = (*NewsHandler)(nil)

type NewsReadOrm struct {
	ID     uint   `gorm:"primarykey"`
	NewsId string `gorm:"index"`
	UserId string `gorm:"index"`
}

func (NewsReadOrm) TableName() string { return "news_read" }
