package taskmanager

import (
	"context"
	"strings"

	"cloud.google.com/go/firestore"
	repoFirestore "github.com/takara2314/bsam-server/pkg/infrastructure/repository/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SubscribeHandler func(string, []byte) error

func (m *Manager) SetSubscribeHandler(handler SubscribeHandler) {
	m.Mu.Lock()
	defer m.Mu.Unlock()
	m.subscribeHandler = handler
}

func (m *Manager) subscribeTasks(ctx context.Context, errPipe chan error) {
	if m.subscribeHandler == nil {
		panic("subscribe handler is not set")
	}

	// Firestore のタスクコレクションの変更を監視する
	it := m.FirestoreClient.Collection("tasks").Snapshots(ctx)
	for {
		snap, err := it.Next()

		if status.Code(err) == codes.DeadlineExceeded {
			return
		}
		if err != nil {
			errPipe <- err
			return
		}

		for _, change := range snap.Changes {
			// タスク追加のみ処理する
			if change.Kind != firestore.DocumentAdded {
				continue
			}

			// 自分のIDのものだけ処理する
			if !strings.HasPrefix(change.Doc.Ref.ID, m.ID) {
				continue
			}

			task, err := repoFirestore.FetchTaskByID(
				ctx,
				m.FirestoreClient,
				change.Doc.Ref.ID,
			)
			if err != nil {
				errPipe <- err
				return
			}

			// 登録されたハンドラを呼び出す
			if handlerErr := m.subscribeHandler(
				task.Type,
				task.Payload,
			); handlerErr == nil {
				// ハンドラーでエラーがなければタスクを完了としてマークする
				if err := repoFirestore.DeleteTaskByID(
					ctx,
					m.FirestoreClient,
					task.ID,
				); err != nil {
					errPipe <- err
					return
				}
			}
		}
	}
}
