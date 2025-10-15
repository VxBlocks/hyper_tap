package monitor

import (
	"context"
	"fmt"
	"hyperliquid-server/firebase"

	"firebase.google.com/go/v4/messaging"
)

type MessageOption struct {
	ChannelID string
	Title     string
	Body      string
	Data      map[string]string
}

func sendToFcmTopic(ctx context.Context, topic string, options MessageOption) error {
	client, err := firebase.GetFirebaseApp().Messaging(ctx)
	if err != nil {
		return fmt.Errorf("error initializing messaging client: %v", err)
	}

	fcmMsg := messaging.Message{
		Topic: topic,
		Data:  options.Data,
		Notification: &messaging.Notification{
			Title: options.Title,
			Body:  options.Body,
		},
		Android: &messaging.AndroidConfig{
			Priority:     "normal",
			DirectBootOK: true,
			Notification: &messaging.AndroidNotification{
				Title:     options.Title,
				Body:      options.Body,
				ChannelID: options.ChannelID,
			},
		},
		APNS: nil,
	}

	_, err = client.Send(ctx, &fcmMsg)
	if err != nil {
		return fmt.Errorf("error SendMulticast: %v", err)
	}

	return nil
}
