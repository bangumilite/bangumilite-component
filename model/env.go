package model

import "os"

const RunningEnvironment = "RUNNING_ENVIRONMENT"

type Environment string

const (
	Local      Environment = "local"
	Production Environment = "production"
)

func GetRunningEnvironment() Environment {
	if os.Getenv(RunningEnvironment) == string(Production) {
		return Production
	}
	return Local
}
