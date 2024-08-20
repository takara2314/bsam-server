package environment

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/samber/oops"
)

type Variables struct {
	JWTSecretKey         string `env:"JWT_SECRET_KEY,required"`
	GoogleCloudProjectID string `env:"GOOGLE_CLOUD_PROJECT_ID,required"`
}

func LoadVariables(fileLoadingForced bool) (*Variables, error) {
	err := godotenv.Load(".env")
	if err != nil && fileLoadingForced {
		return nil, oops.
			In("environment.LoadEnv").
			Wrapf(err, "failed to load .env file")
	}

	var envVar Variables
	if err := env.Parse(&envVar); err != nil {
		return nil, oops.
			In("environment.LoadEnv").
			Wrapf(err, "failed to parse environment variables")
	}

	return &envVar, nil
}
