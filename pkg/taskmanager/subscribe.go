package taskmanager

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	repoFirestore "github.com/takara2314/bsam-server/pkg/infrastructure/repository/firestore"
	"github.com/takara2314/bsam-server/pkg/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SubscribeHandler func(string, []byte) error

func (m *Manager) SetSubscribeHandler(handler SubscribeHandler) {
	m.Mu.Lock()
	defer m.Mu.Unlock()
	m.subscribeHandler = handler
}

func (m *Manager) subscribeTasks(ctx context.Context, errCh chan error) {
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
			errCh <- err
			return
		}

		for _, change := range snap.Changes {
			// タスク追加のみ処理する
			if change.Kind != firestore.DocumentAdded {
				continue
			}

			// タスクIDのprefixによって処理を分岐
			taskIDPrefix := util.FindPrefixIfHasAnyPrefix(
				change.Doc.Ref.ID,
				[]string{PrefixManager, PrefixAssociation},
			)

			switch taskIDPrefix {
			case PrefixManager:
				m.callSubscribeHandlerIfMyManager(ctx, errCh, change)
			case PrefixAssociation:
				m.callSubscribeHandlerIfMyAssociation(ctx, errCh, change)
			default:
				slog.Warn(
					"unsupported prefix",
					"task_id", change.Doc.Ref.ID,
					"task_id_prefix", taskIDPrefix,
				)
			}
		}
	}
}

func (m *Manager) callSubscribeHandlerIfMyManager(
	ctx context.Context,
	errCh chan error,
	change firestore.DocumentChange,
) {
	// 自分のIDのものだけ処理する
	if checkIsMyTask(change.Doc.Ref.ID, m.ID) {
		return
	}

	task, err := repoFirestore.FetchTaskByID(
		ctx,
		m.FirestoreClient,
		change.Doc.Ref.ID,
	)
	if err != nil {
		errCh <- err
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
			errCh <- err
			return
		}
	}
}

func (m *Manager) callSubscribeHandlerIfMyAssociation(
	ctx context.Context,
	errCh chan error,
	change firestore.DocumentChange,
) {
	// 自分の協会IDのものだけ処理する
	if checkIsMyTask(change.Doc.Ref.ID, m.AssociationID) {
		return
	}

	task, err := repoFirestore.FetchTaskByID(
		ctx,
		m.FirestoreClient,
		change.Doc.Ref.ID,
	)
	if err != nil {
		errCh <- err
		return
	}

	// 登録されたハンドラを呼び出す
	if handlerErr := m.subscribeHandler(
		task.Type,
		task.Payload,
	); handlerErr == nil {
		// ハンドラーでエラーがなければタスクを完了としてマークする
		// 他のレースハブが未ウォッチかもしれないので、10秒後に削除する
		go func(ctx context.Context) {
			time.Sleep(10 * time.Second)
			if err := repoFirestore.DeleteTaskByID(
				ctx,
				m.FirestoreClient,
				task.ID,
			); err != nil {
				errCh <- err
				return
			}
		}(ctx)
	}
}

func checkIsMyTask(taskID string, myID string) bool {
	return fetchTargetManagerOrAssociationID(taskID) == myID
}

func fetchTargetManagerOrAssociationID(taskID string) string {
	elems := strings.Split(taskID, "_")
	if len(elems) != 3 {
		return ""
	}
	return elems[1]
}
