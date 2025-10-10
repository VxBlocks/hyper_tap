package handler

import (
	"context"
	"fmt"
	favoritev1 "hyperliquid-server/gen/favorite/v1"
	favoritev1connect "hyperliquid-server/gen/favorite/v1/favoritev1connect"
	"timescale"

	"connectrpc.com/connect"
	"gorm.io/gorm"
)

type FavoriteHandler struct {
}

// CreateFavorite implements favoritev1connect.FavoriteServiceHandler.
func (f *FavoriteHandler) CreateFavorite(ctx context.Context, req *connect.Request[favoritev1.CreateFavoriteRequest]) (*connect.Response[favoritev1.CreateFavoriteResponse], error) {
	var user_id string

	err := timescale.GetPostgresGormDB(ctx).Raw("select user_id_from_session(?)", req.Header().Get("Authorization")).Scan(&user_id).Error
	if err != nil {
		return nil, fmt.Errorf("could not get user id from session: %w", err)
	}

	req.Msg.Payload.UserId = user_id

	res, err := favoritev1.DefaultCreateFavorite(ctx, req.Msg.Payload, timescale.GetPostgresGormDB(ctx))
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&favoritev1.CreateFavoriteResponse{
		Result: res,
	}), nil
}

// ListFavorite implements favoritev1connect.FavoriteServiceHandler.
func (f *FavoriteHandler) ListFavorite(ctx context.Context, req *connect.Request[favoritev1.ListFavoriteRequest]) (*connect.Response[favoritev1.ListFavoriteResponse], error) {
	res, err := gorm.G[favoritev1.FavoriteORM](timescale.GetPostgresGormDB()).
		Where("user_id_from_session(?) = user_id", req.Header().Get("Authorization")).
		Order("id").
		Find(ctx)

	if err != nil {
		return nil, err
	}

	var pbRes []*favoritev1.Favorite
	for _, r := range res {
		pb, err := r.ToPB(ctx)
		if err != nil {
			return nil, err
		}
		pbRes = append(pbRes, &pb)
	}

	return &connect.Response[favoritev1.ListFavoriteResponse]{
		Msg: &favoritev1.ListFavoriteResponse{
			Results: pbRes,
		},
	}, nil
}

// DeleteFavorite implements favoritev1connect.FavoriteServiceHandler.
func (f *FavoriteHandler) DeleteFavorite(ctx context.Context, req *connect.Request[favoritev1.DeleteFavoriteRequest]) (*connect.Response[favoritev1.DeleteFavoriteResponse], error) {
	rowsAffected, err := gorm.G[favoritev1.FavoriteORM](timescale.GetPostgresGormDB()).
		Where("id = ? and user_id = user_id_from_session(?)", req.Msg.Id, req.Header().Get("Authorization")).
		Delete(ctx)
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("favorite not found"))
	}

	return &connect.Response[favoritev1.DeleteFavoriteResponse]{}, nil
}

var _ favoritev1connect.FavoriteServiceHandler = (*FavoriteHandler)(nil)
