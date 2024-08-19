package taskmanager

import (
	"context"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/samber/oops"
	repoFirestore "github.com/takara2314/bsam-server/pkg/infrastructure/repository/firestore"
)

func (m *Manager) PublishToManager(
	ctx context.Context,
	targetID string,
	taskType string,
	payload []byte,
) error {
	return m.publish(
		ctx,
		PrefixManager,
		targetID,
		taskType,
		payload,
	)
}

func (m *Manager) PublishToAssociation(
	ctx context.Context,
	targetID string,
	taskType string,
	payload []byte,
) error {
	return m.publish(
		ctx,
		PrefixAssociation,
		targetID,
		taskType,
		payload,
	)
}

func (m *Manager) publish(
	ctx context.Context,
	prefix string,
	targetManagerOrAssociationID string,
	taskType string,
	payload []byte,
) error {
	taskID := createTaskID(prefix, targetManagerOrAssociationID)

	if err := repoFirestore.SetTask(
		ctx,
		m.FirestoreClient,
		taskID,
		taskType,
		payload,
		time.Now(),
	); err != nil {
		return oops.
			In("taskmanager.Publish").
			With("prefix", prefix).
			With("target_manager_or_association_id", targetManagerOrAssociationID).
			With("task_id", taskID).
			With("task_type", taskType).
			Wrapf(err, "failed to set task to firestore")
	}

	return nil
}

func createTaskID(prefix string, targetManagerOrAssociationID string) string {
	return prefix + targetManagerOrAssociationID + "_" + ulid.Make().String()
}
