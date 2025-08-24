package firebase

import (
	"context"
	_ "embed"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"

	"google.golang.org/api/option"
)

//go:embed serviceAccountKey.json
var serviceAccountKey []byte

var app *firebase.App

func GetFirebaseApp() *firebase.App {
	return app
}

func InitFirebaseApp(ctx context.Context) (*firebase.App, error) {
	opt := option.WithCredentialsJSON(serviceAccountKey)
	config := &firebase.Config{ProjectID: "hl-message"}
	var err error
	app, err = firebase.NewApp(ctx, config, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}
	return app, nil
}

func SendMessage(ctx context.Context, token string, app *firebase.App) error {
	client, err := app.Messaging(context.Background())
	if err != nil {
		return fmt.Errorf("error initializing messaging client: %v", err)
	}
	client.Send(ctx, &messaging.Message{
		Notification: &messaging.Notification{
			Title: "New Message",
			Body:  "You have a new messageaaaaa",
		},
		Token: token,
		Data: map[string]string{
			"message": "Hello World",
		},
		Android: &messaging.AndroidConfig{
			Priority: "high",
		},
	})
	return nil
}

func SendMessageTopic(ctx context.Context, topic string, app *firebase.App) error {
	client, err := app.Messaging(context.Background())
	if err != nil {
		return fmt.Errorf("error initializing messaging client: %v", err)
	}
	client.Send(ctx, &messaging.Message{
		Notification: &messaging.Notification{
			Title: "New Message 222",
			Body:  "You have a new messageaaaaa",
		},
		Topic: topic,
		Data: map[string]string{
			"message": "Hello World",
		},
		Android: &messaging.AndroidConfig{
			Priority:     "high",
			DirectBootOK: true,
		},
	})
	return nil
}
