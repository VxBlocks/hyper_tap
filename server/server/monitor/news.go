package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"logger"
	"net/http"
	"time"
	"timescale"

	"gorm.io/gorm"
)

func StartNewsMonitor() {

	go func() {
		for {
			ctx := context.Background()
			logger.Info("Checking news")
			err := checkUpdateNews(ctx)
			if err != nil {
				logger.Error("Error checkUpdateNews", "error", err)
			}
			time.Sleep(time.Second * 300)
		}
	}()
}

type News struct {
	Entries []NewsOrm `json:"entries"`
}

type NewsOrm struct {
	UUID      string    `json:"uuid" gorm:"column:uuid;primaryKey"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	Preview   string    `json:"preview"`
	Hash      string    `json:"hash"`
	Category  string    `json:"category"`
	IsRead    bool      `json:"isRead" gorm:"-"`
}

func (NewsOrm) TableName() string {
	return "news"
}

func checkUpdateNews(ctx context.Context) error {
	news, err := getNews(ctx)
	if err != nil {
		return err
	}

	for _, entry := range news {
		e, err := gorm.G[NewsOrm](timescale.GetPostgresGormDB()).Where("uuid = ?", entry.UUID).Find(ctx)
		if err != nil {
			return err
		}
		if len(e) == 0 {
			err = gorm.G[NewsOrm](timescale.GetPostgresGormDB()).Create(ctx, &entry)
			if err != nil {
				return err
			}
			logger.Info("created news", "uuid", entry.UUID)
			sendToNewsTopic(ctx, &entry)
		}
	}

	return nil
}

func getNews(ctx context.Context) ([]NewsOrm, error) {
	var news News
	resp, err := http.Get("https://dzjnlsk4rxci0.cloudfront.net/mainnet/entries.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	err = json.NewDecoder(resp.Body).Decode(&news)
	if err != nil {
		return nil, err
	}
	return news.Entries, nil
}

func sendToNewsTopic(ctx context.Context, msg *NewsOrm) error {

	return sendToFcmTopic(ctx, "news", MessageOption{
		ChannelID: "news",
		Title:     msg.Title,
		Body:      msg.Preview,
		Data: map[string]string{
			"message": msg.Title,
			"type":    msg.Category,
			"id":      msg.UUID,
		},
	})
}
