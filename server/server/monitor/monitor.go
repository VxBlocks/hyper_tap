package monitor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"hyperliquid-server/firebase"
	userwatchv1 "hyperliquid-server/gen/userwatch/v1"
	"hyperliquid-server/models"
	"io"
	"logger"
	"net/http"
	"strings"
	"time"
	"timescale"

	gofirebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"gorm.io/gorm"
)

func StartMonitor() {
	go func() {
		for {
			ctx := context.Background()
			doAllFillJob(ctx)
			time.Sleep(time.Second * 30)
		}
	}()
}

func doAllFillJob(ctx context.Context) {
	addresses, err := GetAddresses(ctx)
	if err != nil {
		logger.Error("Error getting addresses", "error", err)
		return
	}

	addresses = []string{"0x5263ABaa3dd77dDD0870CdbCCE78e85d82Ea4c0c", "0xbddfC16aE40Cb06Ceba041d2dB5F72B5175374ee"}

	logger.Info("Getting addresses", "addresses", addresses)

	for _, address := range addresses {
		err := DoUserFillJob(ctx, address)
		if err != nil {
			logger.Error("Error do fill job", "error", err)
		}
	}
}

func DoUserFillJob(ctx context.Context, userId string) error {
	fills, err := getAddressFills(ctx, userId)
	if err != nil {
		return err
	}

	progress, err := GetAddressFillProgress(ctx, userId)
	if err != nil {
		return err
	}

	var progressTime int64

	var newFills []UserFillResponse
	for _, fill := range fills {
		if fill.Time > progress {
			newFills = append(newFills, fill)
		}
		if fill.Time > progressTime {
			progressTime = fill.Time
		}
	}

	for _, fill := range newFills {
		msg, _ := json.Marshal(fill)
		logger.Info("New fill", "fill", fill)

		if strings.HasPrefix(fill.Coin, "@") {
			continue
		}
		err := AddUserWatchMsg(ctx, &userwatchv1.UserWatchMsgORM{
			Msg:         string(msg),
			Event:       fill.Dir,
			UserId:      "all", // all
			WatchUserId: userId,
			Time:        fill.Time,
			Size:        fill.Sz,
			Token:       fill.Coin,
		})
		if err != nil {
			logger.Error("AddUserWatchMsg failed", "err", err)
		}
	}

	logger.Info("saved progress")

	err = SaveAddressFillProgress(ctx, userId, progressTime)
	if err != nil {
		return fmt.Errorf("SaveAddressFillProgress failed for %s: %w", userId, err)
	}

	return nil
}

// UserFillResponse represents the response structure for a single fill
type UserFillResponse struct {
	Coin          string      `json:"coin"`
	Px            string      `json:"px"`
	Sz            string      `json:"sz"`
	Side          string      `json:"side"`
	Time          int64       `json:"time"`
	StartPosition string      `json:"startPosition"`
	Dir           string      `json:"dir"`
	ClosedPnl     string      `json:"closedPnl"`
	Hash          string      `json:"hash"`
	Oid           int64       `json:"oid"`
	Crossed       bool        `json:"crossed"`
	Fee           string      `json:"fee"`
	Tid           int64       `json:"tid"`
	FeeToken      string      `json:"feeToken"`
	TwapId        interface{} `json:"twapId"`
}

// RequestBody represents the expected request structure
type RequestBody struct {
	User string `json:"user"`
	Type string `json:"type"`
}

func GetAddresses(ctx context.Context) ([]string, error) {
	table := userwatchv1.UserWatchORM{}.TableName()
	var addresses []string
	err := timescale.GetPostgresGormDB().Raw("select distinct watch_user_id from " + table).Scan(&addresses).Error
	if err != nil {
		return nil, err
	}
	return addresses, nil
}

// FIXME: should call getUserFillsbytime
func getAddressFills(ctx context.Context, address string) ([]UserFillResponse, error) {
	var fills []UserFillResponse
	request := RequestBody{
		User: address,
		Type: "userFills",
	}
	jsonstr, _ := json.Marshal(request)
	bodyReader := bytes.NewReader(jsonstr)
	resp, err := http.Post("https://api.hyperliquid-testnet.xyz/info", "application/json", bodyReader)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &fills)
	return fills, err
}

func getMonitors(ctx context.Context, address string) ([]userwatchv1.UserWatchORM, error) {
	var monitors []userwatchv1.UserWatchORM
	err := timescale.GetPostgresGormDB(ctx).Model(&userwatchv1.UserWatchORM{}).
		Where("watch_user_id = ?", address).
		Scan(&monitors).Error
	return monitors, err
}

type AddressFillProgress struct {
	Address string `gorm:"primaryKey"`
	Time    int64
}

func (AddressFillProgress) TableName() string {
	return "user_fill_progress"
}

func GetAddressFillProgress(ctx context.Context, address string) (int64, error) {
	var progress AddressFillProgress
	err := timescale.GetPostgresGormDB().Model(&AddressFillProgress{}).
		Where("address = ?", address).
		Scan(&progress).Error
	if err != nil {
		return 0, err
	}

	if progress.Time == 0 {
		logger.Info("AddressFillProgress not found", "address", address)
		return 1754381950824, nil
	}

	return progress.Time, err
}

func SaveAddressFillProgress(ctx context.Context, address string, time int64) error {
	progress := AddressFillProgress{
		Address: address,
		Time:    time,
	}
	return timescale.GetPostgresGormDB().Model(&AddressFillProgress{}).
		Where("address = ?", address).
		Save(&progress).Error
}

func AddUserWatchMsg(ctx context.Context, msg *userwatchv1.UserWatchMsgORM) error {

	err := timescale.GetPostgresGormDB().Create(msg).Error
	if err != nil {
		return err
	}

	go func() {
		err = sendToTopic(ctx, "watch_address", msg, firebase.GetFirebaseApp())
		if err != nil {
			logger.Error("sendFirebaseMessage failed", "err", err)
		}
	}()

	return nil
}

func GetUserFcmTokens(ctx context.Context, userId string) ([]string, error) {
	t, err := gorm.G[models.FcmToken](timescale.GetPostgresGormDB()).Where("user_id = ?", userId).Take(ctx)
	return t.Tokens, err
}

func sendToTopic(ctx context.Context, topic string, msg *userwatchv1.UserWatchMsgORM, app *gofirebase.App) error {
	client, err := app.Messaging(ctx)
	if err != nil {
		return fmt.Errorf("error initializing messaging client: %v", err)
	}

	fcmMsg := messaging.Message{
		Topic: topic,
		Data: map[string]string{
			"message": msg.Msg,
			"type":    msg.Event,
			"address": msg.WatchUserId,
		},
		Notification: &messaging.Notification{
			Title: "Address Fill Monitor",
			Body:  fmt.Sprintf("%s has %s %s %s", msg.WatchUserId, msg.Event, msg.Size, msg.Token),
		},
		Android: &messaging.AndroidConfig{
			Priority:     "high",
			DirectBootOK: true,
		},
	}

	_, err = client.Send(ctx, &fcmMsg)
	if err != nil {
		return fmt.Errorf("error SendMulticast: %v", err)
	}

	return nil
}
