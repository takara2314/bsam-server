package firestore

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/samber/oops"
)

type Task struct {
	ID        string    `firestore:"-"`
	Type      string    `firestore:"type"`
	Payload   []byte    `firestore:"payload"`
	UpdatedAt time.Time `firestore:"updatedAt"`
}

func SetTask(
	ctx context.Context,
	client *firestore.Client,
	id string,
	msgType string,
	payload []byte,
	updatedAt time.Time,
) error {
	_, err := client.Collection("tasks").Doc(id).Set(ctx, Task{
		ID:        id,
		Type:      msgType,
		Payload:   payload,
		UpdatedAt: updatedAt,
	})
	if err != nil {
		return oops.
			In("firestore.SetTask").
			Wrapf(err, "failed to set task")
	}

	return nil
}

func FetchTaskByID(ctx context.Context, client *firestore.Client, id string) (*Task, error) {
	doc, err := client.Collection("tasks").Doc(id).Get(ctx)
	if err != nil {
		return nil, oops.
			In("firestore.FetchTaskByID").
			Wrapf(err, "failed to fetch task")
	}

	var task Task
	err = doc.DataTo(&task)
	if err != nil {
		return nil, oops.
			In("firestore.FetchTaskByID").
			Wrapf(err, "failed to convert data to task")
	}

	task.ID = doc.Ref.ID

	return &task, err
}

func DeleteTaskByID(ctx context.Context, client *firestore.Client, id string) error {
	_, err := client.Collection("tasks").Doc(id).Delete(ctx)
	if err != nil {
		return oops.
			In("firestore.DeleteTaskByID").
			Wrapf(err, "failed to delete task")
	}

	return nil
}
