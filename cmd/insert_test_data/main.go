package main

import (
	"context"
	"log/slog"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/takara2314/bsam-server/pkg/auth"
	"github.com/takara2314/bsam-server/pkg/domain"
	"github.com/takara2314/bsam-server/pkg/environment"
	repoFirestore "github.com/takara2314/bsam-server/pkg/infrastructure/repository/firestore"
	"github.com/takara2314/bsam-server/pkg/logging"
)

var (
	sampleAssociations = []repoFirestore.Association{
		{
			ID:             "japan",
			Name:           "日本サンプルセーリング協会",
			HashedPassword: auth.HashPassword("nippon"),
			ContractType:   string(domain.OneYearContract),
			ExpiresAt:      time.Now().Add(domain.OneYearContract.Duration()),
			UpdatedAt:      time.Now(),
		},
		{
			ID:             "america",
			Name:           "アメリカサンプルセーリング協会",
			HashedPassword: auth.HashPassword("amerika"),
			ContractType:   string(domain.ThreeYearContract),
			ExpiresAt:      time.Now().Add(domain.ThreeYearContract.Duration()),
			UpdatedAt:      time.Now(),
		},
		{
			ID:             "china",
			Name:           "中国サンプルセーリング協会",
			HashedPassword: auth.HashPassword("chugoku"),
			ContractType:   string(domain.FiveYearContract),
			ExpiresAt:      time.Now().Add(domain.FiveYearContract.Duration()),
			UpdatedAt:      time.Now(),
		},
	}
)

func main() {
	var err error
	ctx := context.Background()

	logging.InitSlog()

	env, err := environment.LoadVariables(false)
	if err != nil {
		slog.Error(
			"failed to load env",
			"error", err,
		)
		panic(err)
	}

	firestoreClient, err := repoFirestore.NewClient(
		ctx,
		env.GoogleCloudProjectID,
	)
	if err != nil {
		panic(err)
	}
	defer firestoreClient.Close()

	if err := insertTestData(ctx, firestoreClient); err != nil {
		panic(err)
	}

	slog.Info("successfully inserted test data")
}

func insertTestData(ctx context.Context, client *firestore.Client) error {
	if err := insertTestAssociations(ctx, client); err != nil {
		return err
	}

	return nil
}

func insertTestAssociations(ctx context.Context, client *firestore.Client) error {
	for _, associations := range sampleAssociations {
		if err := repoFirestore.SetAssociation(
			ctx,
			client,
			associations.ID,
			associations.Name,
			associations.HashedPassword,
			associations.ContractType,
			associations.ExpiresAt,
			associations.UpdatedAt,
		); err != nil {
			return err
		}
	}

	return nil
}
