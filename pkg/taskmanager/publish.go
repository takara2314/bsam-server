package taskmanager

import (
	"context"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/samber/oops"
	repoFirestore "github.com/takara2314/bsam-server/pkg/infrastructure/repository/firestore"
)

func (m *Manager) Publish(
	ctx context.Context,
	targetID string,
	msgType string,
	payload []byte,
) error {
	taskID := targetID + "_" + ulid.Make().String()

	if err := repoFirestore.SetTask(
		ctx,
		m.FirestoreClient,
		taskID,
		msgType,
		payload,
		time.Now(),
	); err != nil {
		return oops.
			In("taskmanager.Publish").
			With("target_id", targetID).
			With("task_id", taskID).
			With("message_type", msgType).
			Wrapf(err, "failed to set task to firestore")
	}

	return nil
}
