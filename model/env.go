package model

import "os"

const RunningEnvironment = "RUNNING_ENVIRONMENT"

type Environment string

const (
	Local      Environment = "local"
	Production Environment = "production"
)

func GetRunningEnvironment() Environment {
	var env Environment

	if os.Getenv(RunningEnvironment) == string(Production) {
		env = Production
	} else {
		env = Local
	}

	return env
}
