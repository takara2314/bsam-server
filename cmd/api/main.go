package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/takara2314/bsam-server/internal/api/common"
	"github.com/takara2314/bsam-server/internal/api/presentation"
	"github.com/takara2314/bsam-server/pkg/domain"
	"github.com/takara2314/bsam-server/pkg/infrastructure/repository"
	"github.com/takara2314/bsam-server/pkg/logging"
)

func main() {
	var err error
	ctx := context.Background()

	logging.InitSlog()

	common.FirestoreClient, err = repository.NewFirestore(
		ctx,
		"bsam-app",
	)
	if err != nil {
		panic(err)
	}
	defer common.FirestoreClient.Close()

	// TODO: 後で消す
	err = repository.CreateAssoc(
		ctx,
		common.FirestoreClient,
		"ise",
		"セーリング伊勢",
		"hogehoge",
		domain.ThreeYearContract,
	)
	if err != nil {
		panic(err)
	}

	user, err := repository.FetchAssocByID(
		ctx,
		common.FirestoreClient,
		"ise",
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(user)

	router := presentation.NewGin()
	presentation.RegisterRouter(router)

	slog.Info(
		"api server started",
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
