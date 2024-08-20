package domain

import "time"

type ContractType string

const (
	ThreeMonthContract ContractType = "3month"
	OneYearContract    ContractType = "1year"
	ThreeYearContract  ContractType = "3year"
	FiveYearContract   ContractType = "5year"
	FreeContract       ContractType = "free"
)

func (ct ContractType) Duration() time.Duration {
	switch ct {
	case ThreeMonthContract:
		return time.Duration(90 * 24 * time.Hour)
	case OneYearContract:
		return time.Duration(365 * 24 * time.Hour)
	case ThreeYearContract:
		return time.Duration(3 * 365 * 24 * time.Hour)
	case FiveYearContract:
		return time.Duration(5 * 365 * 24 * time.Hour)
	case FreeContract:
		return 0
	default:
		return 0
	}
}
