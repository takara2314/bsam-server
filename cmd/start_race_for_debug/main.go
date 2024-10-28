package main

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/bytedance/sonic"
	"github.com/eiannone/keyboard"
	"github.com/takara2314/bsam-server/e2e/auth"
	"github.com/takara2314/bsam-server/e2e/raceclient"
	"github.com/takara2314/bsam-server/pkg/racehub"
)

const wsHost = "localhost:8081"
const deviceID = "manager1"

func main() {
	ctx := context.Background()
	var associationID string
	var password string

	fmt.Print("協会IDを入力してください: ")
	if _, err := fmt.Scan(&associationID); err != nil {
		fmt.Println("入力エラーが発生しました:", err)
		os.Exit(1)
	}

	fmt.Print("パスワードを入力してください: ")
	if _, err := fmt.Scan(&password); err != nil {
		fmt.Println("入力エラーが発生しました:", err)
		os.Exit(1)
	}

	// トークンを取得
	token, err := auth.FetchTokenFromAPI(associationID, password)
	if err != nil {
		fmt.Println("ログインに失敗しました")
		os.Exit(1)
	}

	var raceStarted bool
	raceStartedCh := make(chan bool)
	errorCh := make(chan error)

	client := raceclient.NewClient(url.URL{
		Scheme: "ws",
		Host:   wsHost,
		Path:   "/" + associationID,
	}, deviceID)

	fmt.Println("connecting...")

	go establishConnectionAndAuth(
		ctx,
		time.Second,
		client,
		deviceID,
		token,
		raceStartedCh,
		errorCh,
	)

	select {
	case err := <-errorCh:
		fmt.Println(err)
		os.Exit(1)
	case raceStarted = <-raceStartedCh:
	}

	fmt.Println("connected!")
	if raceStarted {
		fmt.Println("--- STARTED レースは開始されています | [s] で切り替え ---")
	} else {
		fmt.Println("--- FINISHED レースは停止されています | [s] で切り替え ---")
	}

	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer keyboard.Close()

	finishCh := make(chan bool)
	go handleKeyPress(client, raceStarted, finishCh)
	<-finishCh

	fmt.Println("bye!")
}

func establishConnectionAndAuth(
	ctx context.Context,
	timeout time.Duration,
	client *raceclient.Client,
	deviceID string,
	token string,
	raceStarted chan bool,
	errorCh chan error,
) {
	err := client.Connect(ctx, timeout)
	if err != nil {
		errorCh <- fmt.Errorf("接続に失敗しました: %v", err)
	}
	defer client.Close()

	// 認証メッセージを送信
	err = client.Send(racehub.AuthInput{
		MessageType: racehub.HandlerTypeAuth,
		Token:       token,
		DeviceID:    deviceID,
	})
	if err != nil {
		errorCh <- fmt.Errorf("メッセージの送信に失敗しました: %v", err)
	}

	it := client.ReceiveStream()
	for {
		payload, err := it.Read()
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				break
			}
			errorCh <- fmt.Errorf("メッセージの受信に失敗しました: %v", err)
		}

		var msg map[string]any
		err = sonic.Unmarshal(payload, &msg)
		if err != nil {
			errorCh <- fmt.Errorf("メッセージのパースに失敗しました: %v", err)
		}

		switch msg["type"] {
		case "manage_race_status":
			started, ok := msg["started"].(bool)
			if !ok {
				errorCh <- fmt.Errorf("startedの値が不正です")
			}
			raceStarted <- started
		}
	}
}

func handleKeyPress(
	client *raceclient.Client,
	raceStarted bool,
	finishCh chan bool,
) {
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		if char == 's' {
			raceStarted = !raceStarted

			err := client.Send(racehub.ManageRaceStatusInput{
				MessageType: racehub.HandlerTypeManageRaceStatus,
				Started:     raceStarted,
				StartedAt:   time.Now(),
			})
			if err != nil {
				fmt.Printf("メッセージの送信に失敗しました: %v\n", err)
				os.Exit(1)
			}

			if raceStarted {
				fmt.Println("START レースを開始します")
			} else {
				fmt.Println("FINISH レースを終了します")
			}
		}

		if key == keyboard.KeyCtrlC {
			finishCh <- true
			return
		}
	}
}
