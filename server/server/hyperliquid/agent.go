package hyperliquid

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ExtraAgent struct {
	Name       string `json:"name"`
	Address    string `json:"address"`
	ValidUntil int64  `json:"validUntil"`
}

func ExtraAgents(ctx context.Context, address string) ([]ExtraAgent, error) {

	requestBody := map[string]string{
		"user": address,
		"type": "extraAgents",
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

	var response []ExtraAgent
	err = json.Unmarshal(body, &response)
	return response, err
}

type Announcement struct {
	Title     string    `json:"title"`
	Preview   string    `json:"preview"`
	Content   string    `json:"content"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"createdAt"`
	UUID      string    `json:"uuid"`
	Hash      string    `json:"hash"`
}

func GetAnnouncement(ctx context.Context, uuid string) (*Announcement, error) {
	request, err := http.NewRequestWithContext(ctx, "GET",
		fmt.Sprintf("https://dzjnlsk4rxci0.cloudfront.net/mainnet/entry-%s.json", uuid), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var announcement Announcement

	err = json.NewDecoder(resp.Body).Decode(&announcement)

	return &announcement, err
}
