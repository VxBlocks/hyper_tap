package handler

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	loginv1 "hyperliquid-server/gen/login/v1"
	"hyperliquid-server/gen/login/v1/loginv1connect"
	"hyperliquid-server/hyperliquid"
	"hyperliquid-server/models"
	"logger"
	"slices"
	"timescale"

	"connectrpc.com/connect"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type LoginHandler struct {
}

// RegisterFcm implements loginv1connect.LoginServiceHandler.
func (l *LoginHandler) RegisterFcm(ctx context.Context, req *connect.Request[loginv1.RegisterFcmRequest]) (*connect.Response[loginv1.RegisterFcmResponse], error) {

	logger.Info("session", "session", req.Header().Get("Authorization"))
	res, err := gorm.G[models.FcmToken](timescale.GetPostgresGormDB()).
		Where("user_id_from_session(?) = user_id", req.Header().Get("Authorization")).
		Take(ctx)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果记录不存在，创建新记录
			err = timescale.GetPostgresGormDB(ctx).Model(&models.FcmToken{}).Create(
				map[string]any{
					"UserId": clause.Expr{SQL: "user_id_from_session(?)", Vars: []any{req.Header().Get("Authorization")}},
					"Tokens": pq.StringArray{req.Msg.FcmToken},
				},
			).Error
			if err != nil {
				return nil, err
			}

			return connect.NewResponse(&loginv1.RegisterFcmResponse{}), nil
		} else {
			return nil, err
		}
	}

	if slices.Contains(res.Tokens, req.Msg.FcmToken) {
		return connect.NewResponse(&loginv1.RegisterFcmResponse{}), nil
	}

	// 如果记录存在，更新现有记录
	if len(res.Tokens) >= 3 {
		res.Tokens = res.Tokens[:2]
	}

	res.Tokens = append([]string{req.Msg.FcmToken}, res.Tokens...)

	err = timescale.GetPostgresGormDB(ctx).Model(&models.FcmToken{}).
		Where("user_id_from_session(?) = user_id", req.Header().Get("Authorization")).
		Updates(&res).Error
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&loginv1.RegisterFcmResponse{}), nil
}

// LoginIsValid implements loginv1connect.LoginServiceHandler.
func (l *LoginHandler) LoginIsValid(ctx context.Context, req *connect.Request[loginv1.LoginIsValidRequest]) (*connect.Response[loginv1.LoginIsValidResponse], error) {
	var authenticated bool
	row := timescale.GetPostgresGormDB(ctx).Raw("select is_authenticated(?,?)", req.Msg.Address, req.Msg.Session).Row()
	err := row.Scan(&authenticated)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&loginv1.LoginIsValidResponse{
		Valid: authenticated,
	}), nil
}

// Logout implements loginv1connect.LoginServiceHandler.
func (l *LoginHandler) Logout(ctx context.Context, req *connect.Request[loginv1.LogoutRequest]) (*connect.Response[loginv1.LogoutResponse], error) {
	err := timescale.GetPostgresGormDB(ctx).Delete(&models.SessionORMV1{}, "session = ? AND address = ?", req.Msg.Session, req.Msg.Address).Error
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&loginv1.LogoutResponse{}), nil
}

// Login implements loginv1connect.LoginServiceHandler.
func (l *LoginHandler) Login(ctx context.Context, req *connect.Request[loginv1.LoginRequest]) (*connect.Response[loginv1.LoginResponse], error) {
	agent, err := Authenticate(req.Msg.Msg, req.Msg.Signature)
	if err != nil {
		return nil, err
	}

	agents, err := hyperliquid.ExtraAgents(ctx, req.Msg.Address)
	if err != nil {
		return nil, err
	}

	if !slices.ContainsFunc(agents, func(v hyperliquid.ExtraAgent) bool {
		return common.HexToAddress(v.Address) == agent
	}) {
		return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("invalid agent %s is not an agent of %s", agent, req.Msg.Address))
	}

	// generate a random session
	var session_bytes [32]byte
	_, err = rand.Read(session_bytes[:])
	if err != nil {
		return nil, err
	}

	session := hexutil.Encode(session_bytes[:])

	sessionTable := models.SessionORMV1{
		Address: req.Msg.Address,
		Session: session,
		// CreatedAt: time.Now(),
	}

	err = timescale.GetPostgresGormDB().Create(&sessionTable).Error
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&loginv1.LoginResponse{
		Session: session,
	}), nil
}

// https://blog.gkomninos.com/metamask-login-using-golang-and-vuejs#heading-backend-web-server-andamp-endpoints
func Authenticate(nonce string, sigHex string) (common.Address, error) {

	// decode the provided signature into bytes
	sig, err := hexutil.Decode(sigHex)
	if err != nil {
		return common.Address{}, err
	}
	// https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L516
	// check here why I am subtracting 27 from the last byte
	// I spent a lot of time to figure out that
	if len(sig) < crypto.SignatureLength {
		return common.Address{}, fmt.Errorf("signature length is %d, want %d", len(sig), crypto.SignatureLength)
	}
	sig[crypto.RecoveryIDOffset] -= 27
	// now hash the nonce
	msg := accounts.TextHash([]byte(nonce))
	// recover the public key that signed that data
	recovered, err := crypto.SigToPub(msg, sig)
	if err != nil {
		return common.Address{}, err
	}
	// create an ethereum address from the extracted public key
	recoveredAddr := crypto.PubkeyToAddress(*recovered)

	return recoveredAddr, nil
}

var _ loginv1connect.LoginServiceHandler = (*LoginHandler)(nil)
