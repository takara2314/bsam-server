package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/takara2314/bsam-server/internal/game/common"
	"github.com/takara2314/bsam-server/internal/game/presentation"
	"github.com/takara2314/bsam-server/pkg/environment"
	"github.com/takara2314/bsam-server/pkg/geolocationhub"
	repoFirestore "github.com/takara2314/bsam-server/pkg/infrastructure/repository/firestore"
	"github.com/takara2314/bsam-server/pkg/logging"
)

func main() {
	var err error
	ctx := context.Background()

	logging.InitSlog()

	common.Env, err = environment.LoadVariables(false)
	if err != nil {
		slog.Error(
			"failed to load env",
			"error", err,
		)
		panic(err)
	}

	common.FirestoreClient, err = repoFirestore.NewClient(
		ctx,
		common.Env.GoogleCloudProjectID,
	)
	if err != nil {
		panic(err)
	}
	defer common.FirestoreClient.Close()

	// TODO: Firestoreの位置情報データの変更を監視するコード。後で削除する
	go watchGeolocations(ctx)

	router := presentation.NewGin()
	presentation.RegisterRouter(router)

	slog.Info(
		"game server started",
		"is_production", os.Getenv("ENVIRONMENT") == "production",
	)

	if err := router.Run(":8080"); err != nil {
		slog.Error(
			"failed to run Gin router",
			"error", err,
		)
		panic(err)
	}
}

func watchGeolocations(ctx context.Context) {
	geohub := geolocationhub.NewGeolocationHub("ise", common.FirestoreClient)

	it := geohub.Snapshots(ctx)
	for {
		fmt.Println("watching...")

		item, err := it.Next()
		if err != nil {
			slog.Error(
				"failed to get next item",
				"error", err,
			)
			return
		}

		fmt.Println("geolocation item:", item)
	}
}
