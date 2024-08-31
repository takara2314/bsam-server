package main

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/takara2314/bsam-server/pkg/auth"
	repoFirestore "github.com/takara2314/bsam-server/pkg/infrastructure/repository/firestore"
)

const projectIDFormat = "bsam-%s"

func main() {
	ctx := context.Background()

	// フラグをパース
	inputFlag, err := parseInputFlag()
	if err != nil {
		panic(err)
	}

	// Firestoreクライアントを作成
	projectID := fmt.Sprintf(projectIDFormat, inputFlag.Environment)
	firestoreClient, err := repoFirestore.NewClient(
		ctx,
		projectID,
	)
	if err != nil {
		panic(err)
	}
	defer firestoreClient.Close()

	// 協会情報をFirestoreに登録
	if err := setAssociation(
		ctx,
		firestoreClient,
		inputFlag,
	); err != nil {
		panic(err)
	}
}

func setAssociation(
	ctx context.Context,
	firestoreClient *firestore.Client,
	inputFlag *InputFlag,
) error {
	hashedPassword := auth.HashPassword(inputFlag.Password)
	expiresAt := inputFlag.ContractStartedAt.Add(
		inputFlag.ContractType.Duration(),
	)

	if err := repoFirestore.SetAssociation(
		ctx,
		firestoreClient,
		inputFlag.AssociationID,
		inputFlag.Name,
		hashedPassword,
		string(inputFlag.ContractType),
		expiresAt,
		time.Now(),
	); err != nil {
		return err
	}

	return nil
}
