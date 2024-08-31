package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/takara2314/bsam-server/pkg/domain"
)

const dateLayout = "2006-01-02"

type InputFlag struct {
	Environment       string
	AssociationID     string
	Name              string
	Password          string
	ContractStartedAt time.Time
	ContractType      domain.ContractType
}

func parseInputFlag() (*InputFlag, error) {
	environment := flag.String("environment", "", "デプロイ環境 ('stg' or 'prd')")
	associationID := flag.String("association_id", "", "協会ID")
	name := flag.String("name", "", "協会名")
	password := flag.String("password", "", "パスワード")
	contractStartedAtStr := flag.String("contract_started_at", "", "契約開始日 (YYYY-MM-DD)")
	contractType := flag.String("contract_type", "", "契約タイプ ('3month' or '1year' or '3year' or '5year' or 'free')")

	flag.Parse()

	if *associationID == "" ||
		*name == "" ||
		*password == "" ||
		*contractStartedAtStr == "" ||
		*contractType == "" {
		return nil, fmt.Errorf("All flags are required")
	}

	jst, _ := time.LoadLocation("Asia/Tokyo")
	contractStartedAt, err := time.ParseInLocation(dateLayout, *contractStartedAtStr, jst)
	if err != nil {
		return nil, err
	}

	return &InputFlag{
		Environment:       *environment,
		AssociationID:     *associationID,
		Name:              *name,
		Password:          *password,
		ContractStartedAt: contractStartedAt,
		ContractType:      domain.ContractType(*contractType),
	}, nil
}
