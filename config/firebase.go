package config

import (
	"context"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

type FirebaseClient struct {
	Firestore *firestore.Client
	Auth      *auth.Client
	Ctx       context.Context
}

func InitFirebase(ctx context.Context) (*FirebaseClient, error) {
	projectID := os.Getenv("FIREBASE_PROJECT_ID")
	if projectID == "" {
		projectID = "chinese-learning-app-442609" // 預設專案ID
	}

	// 初始化 Firebase App
	var app *firebase.App
	var err error

	credentialsPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credentialsPath != "" {
		// 使用服務帳戶金鑰文件
		opt := option.WithCredentialsFile(credentialsPath)
		config := &firebase.Config{ProjectID: projectID}
		app, err = firebase.NewApp(ctx, config, opt)
	} else {
		// 使用預設憑證（適用於 Google Cloud 環境）
		config := &firebase.Config{ProjectID: projectID}
		app, err = firebase.NewApp(ctx, config)
	}

	if err != nil {
		return nil, err
	}

	// 初始化 Firestore 客戶端
	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}

	// 初始化 Auth 客戶端
	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}

	return &FirebaseClient{
		Firestore: firestoreClient,
		Auth:      authClient,
		Ctx:       ctx,
	}, nil
}

func (fc *FirebaseClient) Close() {
	if fc.Firestore != nil {
		fc.Firestore.Close()
	}
}
