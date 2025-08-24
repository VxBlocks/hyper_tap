package handler

import (
	"context"
	"fmt"
	userwatchv1 "hyperliquid-server/gen/userwatch/v1"
	userwatchv1connect "hyperliquid-server/gen/userwatch/v1/userwatchv1connect"
	"timescale"

	"connectrpc.com/connect"
	"gorm.io/gorm"
)

type UserWatchHandler struct {
}

// GetUserWatchMsg implements userwatchv1connect.UserWatchServiceHandler.
func (u *UserWatchHandler) GetUserWatchMsg(ctx context.Context, req *connect.Request[userwatchv1.GetUserWatchMsgRequest]) (*connect.Response[userwatchv1.GetUserWatchMsgResponse], error) {
	res, err := gorm.G[userwatchv1.UserWatchMsgORM](timescale.GetPostgresGormDB()).
		Where("user_id_from_session(?) = user_id", req.Header().Get("Authorization")).
		Where("id = ?", req.Msg.Id).
		Take(ctx)

	if err != nil {
		return nil, err
	}

	pb, _ := res.ToPB(ctx)

	return connect.NewResponse(&userwatchv1.GetUserWatchMsgResponse{
		Results: &pb,
	}), nil
}

// ListUserWatchMsg implements userwatchv1connect.UserWatchServiceHandler.
func (u *UserWatchHandler) ListUserWatchMsg(ctx context.Context, req *connect.Request[userwatchv1.ListUserWatchMsgRequest]) (*connect.Response[userwatchv1.ListUserWatchMsgResponse], error) {
	res, err := gorm.G[userwatchv1.UserWatchMsgORM](timescale.GetPostgresGormDB()).
		// Where("user_id_from_session(?) = user_id OR user_id = 'all'", req.Header().Get("Authorization")).
		Where("user_id = 'all'").
		Order("id desc").
		Limit(200).
		Find(ctx)

	if err != nil {
		return nil, err
	}

	var pbs []*userwatchv1.UserWatchMsg

	for _, r := range res {
		pb, _ := r.ToPB(ctx)
		pbs = append(pbs, &pb)
	}

	return connect.NewResponse(&userwatchv1.ListUserWatchMsgResponse{
		Results: pbs,
	}), nil
}

// CreateUserWatch implements v1connect.UserWatchServiceHandler.
func (u *UserWatchHandler) CreateUserWatch(ctx context.Context, req *connect.Request[userwatchv1.CreateUserWatchRequest]) (*connect.Response[userwatchv1.CreateUserWatchResponse], error) {
	var user_id string

	err := timescale.GetPostgresGormDB(ctx).Raw("select user_id_from_session(?)", req.Header().Get("Authorization")).Scan(&user_id).Error
	if err != nil {
		return nil, fmt.Errorf("could not get user id from session: %w", err)
	}

	req.Msg.Payload.UserId = user_id

	res, err := userwatchv1.DefaultCreateUserWatch(ctx, req.Msg.Payload, timescale.GetPostgresGormDB(ctx))
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&userwatchv1.CreateUserWatchResponse{
		Result: res,
	}), nil
}

// DeleteUserWatch implements v1connect.UserWatchServiceHandler.
func (u *UserWatchHandler) DeleteUserWatch(ctx context.Context, req *connect.Request[userwatchv1.DeleteUserWatchRequest]) (*connect.Response[userwatchv1.DeleteUserWatchResponse], error) {
	rowsAffected, err := gorm.G[userwatchv1.UserWatchORM](timescale.GetPostgresGormDB()).
		Where("id = ? and user_id = user_id_from_session(?)", req.Msg.Id, req.Header().Get("Authorization")).
		Delete(ctx)
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("user watch not found"))
	}

	return &connect.Response[userwatchv1.DeleteUserWatchResponse]{}, nil
}

// ListUserWatch implements v1connect.UserWatchServiceHandler.
func (u *UserWatchHandler) ListUserWatch(ctx context.Context, req *connect.Request[userwatchv1.ListUserWatchRequest]) (*connect.Response[userwatchv1.ListUserWatchResponse], error) {

	res, err := gorm.G[userwatchv1.UserWatchORM](timescale.GetPostgresGormDB()).
		Where("user_id_from_session(?) = user_id", req.Header().Get("Authorization")).
		Order("id").
		Find(ctx)

	if err != nil {
		return nil, err
	}

	var pbRes []*userwatchv1.UserWatch
	for _, r := range res {
		pb, err := r.ToPB(ctx)
		if err != nil {
			return nil, err
		}
		pbRes = append(pbRes, &pb)
	}

	return &connect.Response[userwatchv1.ListUserWatchResponse]{
		Msg: &userwatchv1.ListUserWatchResponse{
			Results: pbRes,
		},
	}, nil
}

// UpUserWatch implements v1connect.UserWatchServiceHandler.
func (u *UserWatchHandler) UpUserWatch(ctx context.Context, req *connect.Request[userwatchv1.UpdateUserWatchRequest]) (*connect.Response[userwatchv1.UpdateUserWatchResponse], error) {
	db := timescale.GetPostgresGormDB(ctx).
		Where("user_id_from_session(?) = ?", req.Header().Get("Authorization"), req.Msg.UserWatch.UserId)

	old, err := userwatchv1.DefaultReadUserWatch(ctx, &userwatchv1.UserWatch{Id: req.Msg.UserWatch.GetId()}, db)
	if err != nil {
		return nil, err
	}

	updated, err := userwatchv1.DefaultApplyFieldMaskUserWatch(ctx, old, req.Msg.UserWatch, req.Msg.Masks, "", timescale.GetPostgresGormDB())
	if err != nil {
		return nil, err
	}

	updatedOrm, err := updated.ToORM(ctx)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := gorm.G[userwatchv1.UserWatchORM](db).Select("*").Updates(ctx, updatedOrm)
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("user watch not found"))
	}

	return &connect.Response[userwatchv1.UpdateUserWatchResponse]{
		Msg: &userwatchv1.UpdateUserWatchResponse{
			UserWatch: updated,
		},
	}, nil
}

var _ userwatchv1connect.UserWatchServiceHandler = (*UserWatchHandler)(nil)
