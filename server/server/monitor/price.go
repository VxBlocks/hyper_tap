package monitor

import (
	"context"
	"fmt"
	"hyperliquid-server/hyperliquid"
	"logger"
	"math"
	"strconv"
	"time"
	"timescale"
)

func StartPriceMonitor() {

	go func() {
		for {
			time.Sleep(time.Minute)
			ctx := context.Background()
			MonitorPriceChanges(ctx)
			time.Sleep(time.Minute * 20)
		}
	}()
}

type TokenPrice struct {
	AssetID      int
	Symbol       string
	Price        string
	PrevDayPrice string
}

type PriceAlert struct {
	ID            uint `gorm:"primaryKey"`
	AssetID       int  `gorm:"index"`
	Symbol        string
	CurrentPrice  float64
	PreviousPrice float64
	ChangePercent float64
	AlertType     string // "threshold_exceeded" or "additional_change"
	CreatedAt     time.Time
}

func (PriceAlert) TableName() string {
	return "price_alerts"
}

type PriceAlertState struct {
	AssetID        int `gorm:"primaryKey"`
	LastAlertPrice float64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (PriceAlertState) TableName() string {
	return "price_alert_states"
}

func getPrices(ctx context.Context) ([]TokenPrice, error) {
	metaAndAssetCtx, err := hyperliquid.GetMetaAndAssetCtx(ctx)
	if err != nil {
		return nil, err
	}
	tokenMetaAndAssetCtx := metaAndAssetCtx.ToTokenMetaAndAssetCtx()
	var tokenPrice []TokenPrice
	for index, assetCtx := range tokenMetaAndAssetCtx {
		tokenPrice = append(tokenPrice, TokenPrice{
			AssetID:      index,
			Symbol:       assetCtx.Universe.Name,
			Price:        assetCtx.AssetCtx.MarkPx,
			PrevDayPrice: assetCtx.AssetCtx.PrevDayPx,
		})
	}
	return tokenPrice, nil
}

// MonitorPriceChanges checks for significant price changes and sends alerts
func MonitorPriceChanges(ctx context.Context) {
	// Load last alert prices from database
	alertPrices, err := loadLastAlertPrices(ctx)
	if err != nil {
		logger.Error("Failed to load last alert prices", "error", err)
		return
	}

	var lastAlertPrices = make(map[int]PriceAlertState)
	for _, alert := range alertPrices {
		lastAlertPrices[alert.AssetID] = alert
	}

	prices, err := getPrices(ctx)
	if err != nil {
		logger.Error("Failed to get prices", "error", err)
		return
	}

	for _, price := range prices {
		currentPrice, err := strconv.ParseFloat(price.Price, 64)
		if err != nil {
			logger.Error("Failed to parse current price", "symbol", price.Symbol, "price", price.Price, "error", err)
			continue
		}

		prevDayPrice, err := strconv.ParseFloat(price.PrevDayPrice, 64)
		if err != nil {
			logger.Error("Failed to parse previous day price", "symbol", price.Symbol, "price", price.PrevDayPrice, "error", err)
			continue
		}

		// Calculate 24-hour percentage change
		if prevDayPrice == 0 {
			continue // Avoid division by zero
		}

		changePercent := (currentPrice - prevDayPrice) / prevDayPrice * 100

		// Check if 24-hour change exceeds 5%
		if math.Abs(changePercent) > 5.0 {
			// Check if we've sent an alert before
			lastAlertPrice, hasAlerted := lastAlertPrices[price.AssetID]

			// 没有发送过，或者过了12h
			if !hasAlerted || lastAlertPrice.CreatedAt.Add(time.Hour*12).Before(time.Now()) {
				alert := PriceAlert{
					AssetID:       price.AssetID,
					Symbol:        price.Symbol,
					CurrentPrice:  currentPrice,
					PreviousPrice: prevDayPrice,
					ChangePercent: changePercent,
					AlertType:     "threshold_exceeded",
					CreatedAt:     time.Now(),
				}

				sendPriceAlert(ctx, alert)
				saveLastAlertPrice(ctx, price.AssetID, currentPrice)
				continue
			} else {

				// If we have alerted before, check for additional 2% change from last alert

				additionalChangePercent := (currentPrice - lastAlertPrice.LastAlertPrice) / lastAlertPrice.LastAlertPrice * 100
				if math.Abs(additionalChangePercent) > 2.0 {
					alert := PriceAlert{
						AssetID:       price.AssetID,
						Symbol:        price.Symbol,
						CurrentPrice:  currentPrice,
						PreviousPrice: prevDayPrice,
						ChangePercent: changePercent,
						AlertType:     "additional_change",
						CreatedAt:     time.Now(),
					}

					sendPriceAlert(ctx, alert)
					saveLastAlertPrice(ctx, price.AssetID, currentPrice)
				}
			}
		}
	}
}

// loadLastAlertPrices loads the last alert prices from the database
func loadLastAlertPrices(ctx context.Context) ([]PriceAlertState, error) {
	var states []PriceAlertState
	if err := timescale.GetPostgresGormDB(ctx).Find(&states).Error; err != nil {
		return nil, err
	}

	return states, nil
}

// saveLastAlertPrice saves the last alert price to the database
func saveLastAlertPrice(ctx context.Context, assetID int, price float64) {
	state := PriceAlertState{
		AssetID:        assetID,
		LastAlertPrice: price,
		UpdatedAt:      time.Now(),
	}

	if err := timescale.GetPostgresGormDB(ctx).Where(PriceAlertState{AssetID: assetID}).Assign(state).FirstOrCreate(&state).Error; err != nil {
		logger.Error("Failed to save alert state", "asset_id", assetID, "error", err)
	}
}

// sendPriceAlert sends a price alert notification and stores it in the database
func sendPriceAlert(ctx context.Context, alert PriceAlert) {
	logger.Info("Price alert triggered",
		"asset_id", alert.AssetID,
		"symbol", alert.Symbol,
		"current_price", alert.CurrentPrice,
		"previous_price", alert.PreviousPrice,
		"change_percent", alert.ChangePercent,
		"alert_type", alert.AlertType)

	// Store alert in database
	if err := timescale.GetPostgresGormDB(ctx).Create(&alert).Error; err != nil {
		logger.Error("Failed to store price alert in database", "error", err)
		return
	}

	sendToFcmTopic(ctx, fmt.Sprintf("price-alert-%d", alert.AssetID), MessageOption{
		ChannelID: "monitor",
		Title:     "Price Alert",
		Body:      fmt.Sprintf("%s 24h price change is %.2f%%", alert.Symbol, alert.ChangePercent),
	})

}
