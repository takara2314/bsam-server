package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/bytedance/sonic"
)

func FetchTokenFromAPI(assocID string, password string) (string, error) {
	// リクエストのボディを作成
	requestBody, err := sonic.Marshal(map[string]string{
		"assoc_id": assocID,
		"password": password,
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	// HTTPクライアントを作成
	client := &http.Client{}

	// リクエストを作成
	req, err := http.NewRequest(
		"POST",
		"http://localhost:8082/verify/password",
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// ヘッダーを設定
	req.Header.Set("Content-Type", "application/json")

	// リクエストを送信
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// レスポンスボディを読み取る
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf(
			"failed to read response body: %w", err,
		)
	}

	// ステータスコードをチェック
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf(
			"unexpected status code: %d, body: %s",
			resp.StatusCode, string(body),
		)
	}

	// レスポンスからトークンを抽出 (実際のレスポンス形式に応じて調整が必要)
	var response struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response.Token, nil
}
